//go:build linux && cgo && libxml2

package fiskalhrgo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"strings"
	"testing"
	"time"
)

func TestLibxml2SignAndVerify(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "XML signature test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certificateDER, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	certificate, err := x509.ParseCertificate(certificateDER)
	if err != nil {
		t.Fatal(err)
	}

	entity := &FiskalEntity{
		cert:       &certManager{privateKey: privateKey, publicCert: certificate},
		ciscert:    &signatureCheckCIScert{PublicCert: certificate},
		useLibxml2: true,
	}
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
