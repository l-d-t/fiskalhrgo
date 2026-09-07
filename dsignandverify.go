package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
)

// doc14n applies Exclusive Canonical XML (http://www.w3.org/2001/10/xml-exc-c14n#) to the input XML data
func doc14n(xmlData []byte) ([]byte, error) {
	// Parse the input XML string into an etree.Document
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	canonicalizer := MakeC14N10ExclusiveCanonicalizerWithPrefixList("") // No prefix list
	canonicalizedXML, err := canonicalizer.Canonicalize(doc.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize the XML: %v", err)
	}

	return canonicalizedXML, nil
}

func createSignedInfoElement(referenceURI, digestValue string) *etree.Element {
	signedInfo := etree.NewElement("SignedInfo")
	signedInfo.CreateAttr("xmlns", "http://www.w3.org/2000/09/xmldsig#")

	canonicalizationMethod := signedInfo.CreateElement("CanonicalizationMethod")
	canonicalizationMethod.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")

	signatureMethod := signedInfo.CreateElement("SignatureMethod")
	signatureMethod.CreateAttr("Algorithm", rsaSHA256Algorithm)

	reference := signedInfo.CreateElement("Reference")
	reference.CreateAttr("URI", "#"+referenceURI)

	transforms := reference.CreateElement("Transforms")

	transform1 := transforms.CreateElement("Transform")
	transform1.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")

	transform2 := transforms.CreateElement("Transform")
	transform2.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")

	digestMethod := reference.CreateElement("DigestMethod")
	digestMethod.CreateAttr("Algorithm", digestSHA256Algorithm)

	digestValueElement := reference.CreateElement("DigestValue")
	digestValueElement.SetText(digestValue)

	return signedInfo
}

func createSignatureElement(signedInfoElement *etree.Element, signatureValue string, cert *x509.Certificate) *etree.Element {
	signatureElement := etree.NewElement("Signature")
	signatureElement.CreateAttr("xmlns", "http://www.w3.org/2000/09/xmldsig#")

	// Add the canonicalized SignedInfo element
	signatureElement.AddChild(signedInfoElement)

	// Add the SignatureValue
	signatureValueElement := signatureElement.CreateElement("SignatureValue")
	signatureValueElement.SetText(signatureValue)

	// Add the KeyInfo
	keyInfoElement := signatureElement.CreateElement("KeyInfo")
	x509DataElement := keyInfoElement.CreateElement("X509Data")

	// Add the X509Certificate
	x509CertificateElement := x509DataElement.CreateElement("X509Certificate")
	x509CertificateElement.SetText(base64.StdEncoding.EncodeToString(cert.Raw))

	// Add the X509IssuerSerial
	x509IssuerSerialElement := x509DataElement.CreateElement("X509IssuerSerial")

	x509IssuerNameElement := x509IssuerSerialElement.CreateElement("X509IssuerName")
	x509IssuerNameElement.SetText(cert.Issuer.String())

	x509SerialNumberElement := x509IssuerSerialElement.CreateElement("X509SerialNumber")
	x509SerialNumberElement.SetText(cert.SerialNumber.String())

	return signatureElement
}

func (fe *FiskalEntity) signXML(xmlRequest []byte) ([]byte, error) {
	// Step 1: Parse and Canonicalize the XML document using etree
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlRequest); err != nil {
		return nil, fmt.Errorf("failed to parse XML document: %v", err)
	}

	// Step 6: Insert the Signature block before the closing tag of the root element
	root := doc.Root()
	if root == nil {
		return nil, fmt.Errorf("invalid XML: root element not found")
	}

	referenceID := root.SelectAttrValue("Id", "")
	if referenceID == "" {
		return nil, fmt.Errorf("no Id attribute found in the root element")
	}

	// Canonicalize the XML document
	xmlCanonical, err := fe.canonicalizeForSigning(xmlRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize XML document: %v", err)
	}

	// DigestValue calculation using SHA-256
	digest := sha256.New()
	if _, err := digest.Write([]byte(xmlCanonical)); err != nil {
		return nil, fmt.Errorf("failed to calculate digest: %v", err)
	}
	digestValue := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	// Step 2: Create SignedInfo block with DigestValue using etree
	signedInfoElement := createSignedInfoElement(referenceID, digestValue)

	// Convert the SignedInfo element to a string
	signedInfoDocument := etree.NewDocument()
	signedInfoDocument.SetRoot(signedInfoElement)
	signedInfoString, err := signedInfoDocument.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize SignedInfo: %v", err)
	}

	// Canonicalize the SignedInfo block
	canonicalizedSignedInfo, err := fe.canonicalizeForSigning(signedInfoString)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize SignedInfo: %v", err)
	}

	// Step 3: Compute hash of canonicalized SignedInfo
	hashedSignedInfo := sha256.Sum256(canonicalizedSignedInfo)

	// Step 4: Generate the SignatureValue using the private key
	signature, err := rsa.SignPKCS1v15(nil, fe.cert.privateKey, crypto.SHA256, hashedSignedInfo[:])
	if err != nil {
		return nil, fmt.Errorf("failed to generate signature: %v", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	// Step 5: Build the Signature block with certificate details using etree
	signatureBlock := createSignatureElement(
		signedInfoElement,
		signatureValue,
		fe.cert.publicCert,
	)

	root.AddChild(signatureBlock)

	// Serialize the updated document back to bytes
	output, err := doc.WriteToBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed XML: %v", err)
	}

	return output, nil
}

func (fe *FiskalEntity) canonicalizeForSigning(xmlData []byte) ([]byte, error) {
	if fe.useLibxml2 {
		return canonicalizeNativeDocument(xmlData, true)
	}
	return doc14n(xmlData)
}

// SignXML signs an XML document whose root element has an Id attribute.
func (fe *FiskalEntity) SignXML(xmlData []byte) ([]byte, error) {
	return fe.signXML(xmlData)
}

// verifyXML preserves the intentional legacy bypass when native mode is off.
// In that mode true means verification was skipped, NOT a valid XML signature.
// Native mode checks the signed payload, reference digest and pinned certificate.
func (fe *FiskalEntity) verifyXML(xmlData []byte) (bool, error) {
	if !fe.useLibxml2 {
		return true, nil
	}
	return fe.verifyXMLLibxml2(xmlData)
}

// VerifyXML verifies an XML signature and pins its certificate to the CIS
// certificate embedded in this FiskalEntity. The signature must cover the root
// or the sole CIS response in a SOAP Body. Native support must be enabled.
func (fe *FiskalEntity) VerifyXML(xmlData []byte) (bool, error) {
	if !fe.useLibxml2 {
		return false, errors.New("XML signature verification requires libxml2; call SetUseLibxml2(true)")
	}
	return fe.verifyXMLLibxml2(xmlData)
}

const (
	xmlDSigNamespace       = "http://www.w3.org/2000/09/xmldsig#"
	c14n10Algorithm        = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	exclusiveC14NAlgorithm = "http://www.w3.org/2001/10/xml-exc-c14n#"
	rsaSHA1Algorithm       = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	rsaSHA256Algorithm     = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	digestSHA1Algorithm    = "http://www.w3.org/2000/09/xmldsig#sha1"
	digestSHA256Algorithm  = "http://www.w3.org/2001/04/xmlenc#sha256"
)

func (fe *FiskalEntity) verifyXMLLibxml2(xmlData []byte) (bool, error) {
	if err := validateSignedXMLDocument(xmlData); err != nil {
		return false, err
	}
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(xmlData); err != nil {
		return false, fmt.Errorf("failed to parse signed XML: %w", err)
	}
	signatures := findXMLDSigElements(doc.Root(), "Signature")
	if len(signatures) != 1 {
		return false, fmt.Errorf("expected one XML signature, found %d", len(signatures))
	}
	signature := signatures[0]
	signedInfo, err := requiredDSigChild(signature, "SignedInfo")
	if err != nil {
		return false, err
	}
	canonicalizationMethod, err := requiredDSigChild(signedInfo, "CanonicalizationMethod")
	if err != nil {
		return false, err
	}
	signedInfoExclusive, err := canonicalizationMode(canonicalizationMethod.SelectAttrValue("Algorithm", ""))
	if err != nil {
		return false, err
	}

	reference, err := requiredDSigChild(signedInfo, "Reference")
	if err != nil {
		return false, err
	}
	referenceURI := reference.SelectAttrValue("URI", "")
	if !strings.HasPrefix(referenceURI, "#") || len(referenceURI) == 1 {
		return false, errors.New("XML signature Reference URI must identify a local Id")
	}
	if err := validateSignatureTarget(doc.Root(), signature, referenceURI); err != nil {
		return false, err
	}
	if len(canonicalizationMethod.ChildElements()) != 0 {
		return false, errors.New("canonicalization parameters are not supported")
	}
	referenceExclusive, err := referenceCanonicalizationMode(reference)
	if err != nil {
		return false, err
	}
	canonicalReference, err := canonicalizeNativeReference(xmlData, referenceURI[1:], referenceExclusive)
	if err != nil {
		return false, fmt.Errorf("failed to canonicalize XML signature reference: %w", err)
	}
	digestMethod, err := requiredDSigChild(reference, "DigestMethod")
	if err != nil {
		return false, err
	}
	digestValue, err := requiredDSigChild(reference, "DigestValue")
	if err != nil {
		return false, err
	}
	calculatedDigest, err := digestBytes(digestMethod.SelectAttrValue("Algorithm", ""), canonicalReference)
	if err != nil {
		return false, err
	}
	expectedDigest, err := decodeBase64Text(digestValue.Text())
	if err != nil {
		return false, fmt.Errorf("invalid DigestValue: %w", err)
	}
	if !bytesEqual(calculatedDigest, expectedDigest) {
		return false, errors.New("XML signature reference digest mismatch")
	}

	certificateElements := findXMLDSigElements(signature, "X509Certificate")
	if len(certificateElements) != 1 {
		return false, fmt.Errorf("expected one signing certificate, found %d", len(certificateElements))
	}
	certificateDER, err := decodeBase64Text(certificateElements[0].Text())
	if err != nil {
		return false, fmt.Errorf("invalid signing certificate: %w", err)
	}
	certificate, err := x509.ParseCertificate(certificateDER)
	if err != nil {
		return false, fmt.Errorf("invalid signing certificate: %w", err)
	}
	if fe.ciscert == nil || fe.ciscert.PublicCert == nil || !certificate.Equal(fe.ciscert.PublicCert) {
		return false, errors.New("XML signing certificate is not the trusted CIS certificate")
	}

	canonicalSignedInfo, err := canonicalizeNativeSignedInfo(xmlData, signedInfoExclusive)
	if err != nil {
		return false, fmt.Errorf("failed to canonicalize SignedInfo: %w", err)
	}
	signatureMethod, err := requiredDSigChild(signedInfo, "SignatureMethod")
	if err != nil {
		return false, err
	}
	hash, digest, err := signatureDigest(signatureMethod.SelectAttrValue("Algorithm", ""), canonicalSignedInfo)
	if err != nil {
		return false, err
	}
	signatureValue, err := requiredDSigChild(signature, "SignatureValue")
	if err != nil {
		return false, err
	}
	signatureBytes, err := decodeBase64Text(signatureValue.Text())
	if err != nil {
		return false, fmt.Errorf("invalid SignatureValue: %w", err)
	}
	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("CIS signing certificate does not contain an RSA public key")
	}
	if err := rsa.VerifyPKCS1v15(publicKey, hash, digest, signatureBytes); err != nil {
		return false, fmt.Errorf("invalid XML signature: %w", err)
	}
	return true, nil
}

func findXMLDSigElements(root *etree.Element, tag string) []*etree.Element {
	if root == nil {
		return nil
	}
	var result []*etree.Element
	if root.Tag == tag && root.NamespaceURI() == xmlDSigNamespace {
		result = append(result, root)
	}
	for _, child := range root.ChildElements() {
		result = append(result, findXMLDSigElements(child, tag)...)
	}
	return result
}

func requiredDSigChild(parent *etree.Element, tag string) (*etree.Element, error) {
	var matches []*etree.Element
	for _, child := range parent.ChildElements() {
		if child.Tag == tag && child.NamespaceURI() == xmlDSigNamespace {
			matches = append(matches, child)
		}
	}
	if len(matches) != 1 {
		return nil, fmt.Errorf("expected one %s element, found %d", tag, len(matches))
	}
	return matches[0], nil
}

func canonicalizationMode(algorithm string) (bool, error) {
	switch algorithm {
	case c14n10Algorithm:
		return false, nil
	case exclusiveC14NAlgorithm:
		return true, nil
	default:
		return false, fmt.Errorf("unsupported canonicalization algorithm %q", algorithm)
	}
}

func referenceCanonicalizationMode(reference *etree.Element) (bool, error) {
	transforms, err := requiredDSigChild(reference, "Transforms")
	if err != nil {
		return false, err
	}
	children := transforms.ChildElements()
	if len(children) != 2 {
		return false, errors.New("expected enveloped-signature followed by canonicalization")
	}
	for _, transform := range children {
		if transform.Tag != "Transform" || transform.NamespaceURI() != xmlDSigNamespace {
			return false, errors.New("unexpected transform element")
		}
		if len(transform.ChildElements()) != 0 {
			return false, errors.New("transform parameters are not supported")
		}
	}
	if children[0].SelectAttrValue("Algorithm", "") != xmlDSigNamespace+"enveloped-signature" {
		return false, errors.New("first transform must be enveloped-signature")
	}
	return canonicalizationMode(children[1].SelectAttrValue("Algorithm", ""))
}

// Reject constructs which could be interpreted differently by the Go and C
// parsers. External resources and DTD-defined entities are never needed by CIS.
func validateSignedXMLDocument(data []byte) error {
	decoder := xml.NewDecoder(strings.NewReader(string(data)))
	depth, roots := 0, 0
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("invalid signed XML: %w", err)
		}
		switch token := token.(type) {
		case xml.Directive:
			return errors.New("DTD and XML directives are not permitted")
		case xml.StartElement:
			if depth == 0 {
				roots++
			}
			depth++
			attrs := make(map[xml.Name]bool)
			for _, attr := range token.Attr {
				if attrs[attr.Name] {
					return errors.New("duplicate XML attribute")
				}
				attrs[attr.Name] = true
				if attr.Name.Local == "Id" && attr.Name.Space != "" {
					return errors.New("signature Id must be unqualified")
				}
			}
		case xml.EndElement:
			depth--
		case xml.CharData:
			if depth == 0 && strings.TrimSpace(string(token)) != "" {
				return errors.New("text outside XML root")
			}
		}
	}
	if roots != 1 {
		return errors.New("expected exactly one XML root")
	}
	return nil
}

const soapNamespace = "http://schemas.xmlsoap.org/soap/envelope/"

func validateSignatureTarget(root, signature *etree.Element, referenceURI string) error {
	target := root
	if root.Tag == "Envelope" || root.NamespaceURI() == soapNamespace {
		if root.Tag != "Envelope" || root.NamespaceURI() != soapNamespace {
			return errors.New("invalid SOAP Envelope namespace")
		}
		var body *etree.Element
		headers := 0
		for _, child := range root.ChildElements() {
			if child.NamespaceURI() != soapNamespace {
				return errors.New("unexpected SOAP Envelope child namespace")
			}
			switch child.Tag {
			case "Body":
				if body != nil {
					return errors.New("multiple SOAP Body elements")
				}
				body = child
			case "Header":
				headers++
				if headers > 1 || body != nil {
					return errors.New("invalid SOAP Header placement")
				}
			default:
				return errors.New("unexpected SOAP Envelope child")
			}
		}
		if body == nil || len(body.ChildElements()) != 1 || strings.TrimSpace(body.Text()) != "" {
			return errors.New("expected one CIS response in SOAP Body")
		}
		target = body.ChildElements()[0]
		if target.NamespaceURI() != DefaultNamespace || !strings.HasSuffix(target.Tag, "Odgovor") {
			return errors.New("SOAP Body does not contain a CIS response")
		}
	}
	if target.SelectAttrValue("Id", "") == "" || referenceURI != "#"+target.SelectAttrValue("Id", "") {
		return errors.New("XML signature must reference the consumed response root")
	}
	if signature.Parent() != target {
		return errors.New("XML signature must be a direct child of the signed response")
	}
	return nil
}

func digestBytes(algorithm string, data []byte) ([]byte, error) {
	switch algorithm {
	case digestSHA1Algorithm:
		digest := sha1.Sum(data)
		return digest[:], nil
	case digestSHA256Algorithm:
		digest := sha256.Sum256(data)
		return digest[:], nil
	default:
		return nil, fmt.Errorf("unsupported digest algorithm %q", algorithm)
	}
}

func signatureDigest(algorithm string, data []byte) (crypto.Hash, []byte, error) {
	switch algorithm {
	case rsaSHA1Algorithm:
		digest := sha1.Sum(data)
		return crypto.SHA1, digest[:], nil
	case rsaSHA256Algorithm:
		digest := sha256.Sum256(data)
		return crypto.SHA256, digest[:], nil
	default:
		return 0, nil, fmt.Errorf("unsupported signature algorithm %q", algorithm)
	}
}

func decodeBase64Text(value string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(strings.Join(strings.Fields(value), ""))
}

func bytesEqual(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	var difference byte
	for index := range left {
		difference |= left[index] ^ right[index]
	}
	return difference == 0
}
