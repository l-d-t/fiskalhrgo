//go:build linux && cgo && libxml2

package fiskalhrgo

import (
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/beevik/etree"
)

func TestLibxml2SignAndVerify(t *testing.T) {
	entity := newSignatureTestEntity(t)
	entity.useLibxml2 = true
	xmlData := []byte(`<tns:RacunZahtjev xmlns:tns="urn:test" Id="request-1"><tns:Value>42</tns:Value></tns:RacunZahtjev>`)
	signedXML, err := entity.SignXML(xmlData)
	if err != nil {
		t.Fatalf("SignXML failed: %v", err)
	}
	if !strings.Contains(string(signedXML), rsaSHA256Algorithm) {
		t.Fatalf("SignXML did not use RSA-SHA256: %s", signedXML)
	}
	if !strings.Contains(string(signedXML), digestSHA256Algorithm) {
		t.Fatalf("SignXML did not use SHA-256 digest: %s", signedXML)
	}
	if strings.Contains(string(signedXML), rsaSHA1Algorithm) || strings.Contains(string(signedXML), digestSHA1Algorithm) {
		t.Fatalf("SignXML still contains a SHA-1 XMLDSig algorithm: %s", signedXML)
	}
	verified, err := entity.VerifyXML(signedXML)
	if err != nil || !verified {
		t.Fatalf("VerifyXML failed: verified=%v err=%v", verified, err)
	}

	tamperedXML := []byte(strings.Replace(string(signedXML), ">42<", ">43<", 1))
	if verified, err := entity.VerifyXML(tamperedXML); err == nil || verified {
		t.Fatalf("VerifyXML accepted tampered XML: verified=%v err=%v", verified, err)
	}
}

func TestLibxml2B2BSignatures(t *testing.T) {
	entity := newSignatureTestEntity(t)
	for _, nativeSigner := range []bool{false, true} {
		t.Run(fmt.Sprintf("nativeSigner=%v", nativeSigner), func(t *testing.T) {
			entity.useLibxml2 = nativeSigner
			invoice, zki, err := entity.NewCISInvoice(time.Now(), 2, 1, nil, nil, nil,
				"0.00", "0.00", "0.00", nil, "10.00", CISCash, entity.OIB())
			if err != nil {
				t.Fatal(err)
			}
			if err := invoice.SetOibPrimateljaRacuna("65049901548"); err != nil {
				t.Fatal(err)
			}
			data, err := xml.Marshal(RacunZahtjev{Xmlns: DefaultNamespace, IdAttr: "b2b", Racun: invoice})
			if err != nil {
				t.Fatal(err)
			}
			signed, err := entity.SignXML(data)
			if err != nil {
				t.Fatal(err)
			}
			entity.useLibxml2 = true
			if ok, err := entity.VerifyXML(signed); err != nil || !ok {
				t.Fatalf("B2B signature failed: %v", err)
			}
			modified := strings.Replace(string(signed), "<tns:OibPrimateljaRacuna>65049901548", "<tns:OibPrimateljaRacuna>12345678900", 1)
			if ok, err := entity.VerifyXML([]byte(modified)); ok || err == nil {
				t.Fatal("buyer OIB not covered by signature")
			}
			if invoice.GetZKI() != zki {
				t.Fatal("buyer OIB changed ZKI")
			}
		})
	}
}

func TestLibxml2RejectInvalidSignatures(t *testing.T) {
	entity := newSignatureTestEntity(t)
	entity.useLibxml2 = true
	signed, err := entity.SignXML([]byte(`<r Id="r" xmlns="urn:test"><value>42</value></r>`))
	if err != nil {
		t.Fatal(err)
	}
	mutate := func(change func(*etree.Element, *etree.Element)) []byte {
		t.Helper()
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(signed); err != nil {
			t.Fatal(err)
		}
		change(doc.Root(), findXMLDSigElements(doc.Root(), "Signature")[0])
		data, err := doc.WriteToBytes()
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	cases := map[string][]byte{
		"malformed":            []byte(`<broken>`),
		"empty":                nil,
		"unsigned":             []byte(`<r Id="r"/>`),
		"multiple roots":       append(append([]byte{}, signed...), []byte(`<other/>`)...),
		"DTD":                  append([]byte(`<!DOCTYPE r [<!ENTITY unused "value">]>`), signed...),
		"duplicate Id":         mutate(func(r, s *etree.Element) { r.CreateElement("other").CreateAttr("Id", "r") }),
		"qualified Id":         []byte(strings.Replace(string(signed), `Id="r"`, `xmlns:a="urn:a" a:Id="r"`, 1)),
		"duplicate attribute":  []byte(strings.Replace(string(signed), `Id="r"`, `Id="r" Id="r"`, 1)),
		"duplicate Signature":  mutate(func(r, s *etree.Element) { r.AddChild(s.Copy()) }),
		"duplicate SignedInfo": mutate(func(r, s *etree.Element) { s.AddChild(s.FindElement("SignedInfo").Copy()) }),
		"duplicate Reference": mutate(func(r, s *etree.Element) {
			si := s.FindElement("SignedInfo")
			si.AddChild(si.FindElement("Reference").Copy())
		}),
		"external reference":          []byte(strings.Replace(string(signed), `URI="#r"`, `URI="https://example.invalid/data"`, 1)),
		"wrong reference":             []byte(strings.Replace(string(signed), `URI="#r"`, `URI="#absent"`, 1)),
		"bad digest base64":           mutate(func(r, s *etree.Element) { s.FindElement("SignedInfo/Reference/DigestValue").SetText("!") }),
		"bad signature base64":        mutate(func(r, s *etree.Element) { s.FindElement("SignatureValue").SetText("!") }),
		"bad RSA signature":           mutate(func(r, s *etree.Element) { s.FindElement("SignatureValue").SetText("AAAA") }),
		"missing certificate":         mutate(func(r, s *etree.Element) { s.RemoveChild(s.FindElement("KeyInfo")) }),
		"bad certificate":             mutate(func(r, s *etree.Element) { s.FindElement("KeyInfo/X509Data/X509Certificate").SetText("AAAA") }),
		"unknown signature algorithm": []byte(strings.Replace(string(signed), rsaSHA256Algorithm, "urn:unsupported", 1)),
		"unknown digest algorithm":    []byte(strings.Replace(string(signed), digestSHA256Algorithm, "urn:unsupported", 1)),
		"unknown canonicalization": mutate(func(r, s *etree.Element) {
			s.FindElement("SignedInfo/CanonicalizationMethod").CreateAttr("Algorithm", "urn:unsupported")
		}),
		"unsupported prefix list": mutate(func(r, s *etree.Element) {
			s.FindElement("SignedInfo/CanonicalizationMethod").CreateElement("InclusiveNamespaces").CreateAttr("PrefixList", "test")
		}),
		"reversed transforms": mutate(func(r, s *etree.Element) {
			tr := s.FindElement("SignedInfo/Reference/Transforms")
			tr.AddChild(tr.ChildElements()[0])
		}),
		"duplicate transform": mutate(func(r, s *etree.Element) {
			tr := s.FindElement("SignedInfo/Reference/Transforms")
			tr.AddChild(tr.ChildElements()[0].Copy())
		}),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			if ok, err := entity.VerifyXML(data); ok || err == nil {
				t.Fatalf("accepted invalid signature: %v, %v", ok, err)
			}
		})
	}
	other := newSignatureTestEntity(t)
	other.useLibxml2 = true
	if ok, err := other.VerifyXML(signed); ok || err == nil {
		t.Fatal("accepted an untrusted signing certificate")
	}
}

// xmlsec1 is an independent signer, avoiding a round-trip test in which our
// own signer and verifier might share the same canonicalization mistake.
func xmlsecSOAPFixture(t *testing.T, entity *FiskalEntity, signedInfoAlgorithm, referenceAlgorithm string) []byte {
	t.Helper()
	tool, err := exec.LookPath("xmlsec1")
	if err != nil {
		t.Skip("xmlsec1 is required for independent interoperability tests (installed by Linux CI)")
	}
	dir := t.TempDir()
	keyPath, certPath := filepath.Join(dir, "key.pem"), filepath.Join(dir, "cert.pem")
	files := map[string][]byte{
		keyPath:  pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(entity.cert.privateKey)}),
		certPath: pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: entity.cert.publicCert.Raw}),
	}
	template := fmt.Sprintf(`<s:Envelope xmlns:s="%s" xmlns:tns="%s" xmlns:unused="urn:unused"><s:Body><tns:RacunOdgovor Id="response-1"><tns:Jir>c7310c17-6245-4b36-bf20-b908ded86823</tns:Jir><Signature xmlns="%s"><SignedInfo><CanonicalizationMethod Algorithm="%s"/><SignatureMethod Algorithm="%s"/><Reference URI="#response-1"><Transforms><Transform Algorithm="%senveloped-signature"/><Transform Algorithm="%s"/></Transforms><DigestMethod Algorithm="%s"/><DigestValue/></Reference></SignedInfo><SignatureValue/><KeyInfo><X509Data><X509Certificate/></X509Data></KeyInfo></Signature></tns:RacunOdgovor></s:Body></s:Envelope>`, soapNamespace, DefaultNamespace, xmlDSigNamespace, signedInfoAlgorithm, rsaSHA256Algorithm, xmlDSigNamespace, referenceAlgorithm, digestSHA256Algorithm)
	input, output := filepath.Join(dir, "template.xml"), filepath.Join(dir, "signed.xml")
	files[input] = []byte(strings.Replace(template, "<KeyInfo>", "<KeyInfo><KeyName>fixture-key</KeyName>", 1))
	for name, data := range files {
		if err := os.WriteFile(name, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(tool, "--sign", "--privkey-pem:fixture-key", keyPath+","+certPath,
		"--id-attr:Id", DefaultNamespace+":RacunOdgovor", "--output", output, input)
	if result, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("xmlsec1 signing failed: %v\n%s", err, result)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestLibxml2IndependentSOAPSignatures(t *testing.T) {
	entity := newSignatureTestEntity(t)
	entity.useLibxml2 = true
	for _, si := range []string{c14n10Algorithm, exclusiveC14NAlgorithm} {
		for _, ref := range []string{c14n10Algorithm, exclusiveC14NAlgorithm} {
			t.Run(si+"/"+ref, func(t *testing.T) {
				data := xmlsecSOAPFixture(t, entity, si, ref)
				if ok, err := entity.VerifyXML(data); !ok || err != nil {
					t.Fatalf("independent signature rejected: %v", err)
				}
			})
		}
	}
}

func TestLibxml2SOAPResponseBinding(t *testing.T) {
	entity := newSignatureTestEntity(t)
	entity.useLibxml2 = true
	signed := xmlsecSOAPFixture(t, entity, exclusiveC14NAlgorithm, exclusiveC14NAlgorithm)
	modify := func(change func(*etree.Element, *etree.Element)) []byte {
		t.Helper()
		doc := etree.NewDocument()
		if err := doc.ReadFromBytes(signed); err != nil {
			t.Fatal(err)
		}
		change(doc.Root(), doc.Root().FindElement("s:Body"))
		data, err := doc.WriteToBytes()
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	cases := map[string][]byte{
		"valid":                   signed,
		"changed JIR":             []byte(strings.Replace(string(signed), "c7310c17", "c7310c18", 1)),
		"duplicate Body":          modify(func(root, body *etree.Element) { root.AddChild(body.Copy()) }),
		"multiple payloads":       modify(func(root, body *etree.Element) { body.AddChild(body.ChildElements()[0].Copy()) }),
		"wrong SOAP namespace":    []byte(strings.ReplaceAll(string(signed), soapNamespace, "urn:wrong-soap")),
		"wrong payload namespace": []byte(strings.ReplaceAll(string(signed), DefaultNamespace, "urn:wrong-cis")),
		"signed payload in Header": modify(func(root, body *etree.Element) {
			payload := body.ChildElements()[0]
			body.RemoveChild(payload)
			root.RemoveChild(body)
			root.CreateElement("s:Header").AddChild(payload)
			root.AddChild(body)
			body.CreateElement("tns:RacunOdgovor").CreateAttr("Id", "unsigned")
		}),
		"signed nested payload": modify(func(root, body *etree.Element) {
			payload := body.ChildElements()[0]
			body.RemoveChild(payload)
			wrapper := body.CreateElement("tns:RacunOdgovor")
			wrapper.CreateAttr("Id", "unsigned")
			wrapper.AddChild(payload)
		}),
	}
	for name, data := range cases {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/xml")
				_, _ = w.Write(data)
			}))
			defer server.Close()
			entity.url = server.URL
			entity.ciscert.SSLverifyPoll = x509.NewCertPool()
			entity.ciscert.SSLverifyPoll.AddCert(server.Certificate())
			body, _, err := entity.GetResponse([]byte(`<r Id="request"/>`), true)
			if name == "valid" {
				if err != nil || !strings.Contains(string(body), "c7310c17-6245-4b36-bf20-b908ded86823") {
					t.Fatalf("valid signed SOAP response failed: %v", err)
				}
			} else if err == nil {
				t.Fatal("GetResponse accepted an invalid or unsigned response payload")
			}
		})
	}
}
