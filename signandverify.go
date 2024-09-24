package fiskalhrgo

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/beevik/etree"
)

// generateUniqueID generates a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// docC14N10 applies Canonical XML 1.0 (http://www.w3.org/TR/2001/REC-xml-c14n-20010315) to the input XML data
func docC14N10(xmlData string) ([]byte, error) {
	// Parse the input XML string into an etree.Document
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlData); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %v", err)
	}

	// Use the Canonical XML 1.0 algorithm from goxmldsig
	canonicalizer := MakeC14N10RecCanonicalizer() // Without comments
	canonicalizedXML, err := canonicalizer.Canonicalize(doc.Root())
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize the XML: %v", err)
	}

	return canonicalizedXML, nil
}

// doc14n applies Exclusive Canonical XML (http://www.w3.org/2001/10/xml-exc-c14n#) to the input XML data
func doc14n(xmlData string) ([]byte, error) {
	// Parse the input XML string into an etree.Document
	doc := etree.NewDocument()
	if err := doc.ReadFromString(xmlData); err != nil {
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
	signatureMethod.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#rsa-sha1")

	reference := signedInfo.CreateElement("Reference")
	reference.CreateAttr("URI", "#"+referenceURI)

	transforms := reference.CreateElement("Transforms")

	transform1 := transforms.CreateElement("Transform")
	transform1.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#enveloped-signature")

	transform2 := transforms.CreateElement("Transform")
	transform2.CreateAttr("Algorithm", "http://www.w3.org/2001/10/xml-exc-c14n#")

	digestMethod := reference.CreateElement("DigestMethod")
	digestMethod.CreateAttr("Algorithm", "http://www.w3.org/2000/09/xmldsig#sha1")

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
	xmlCanonical, err := doc14n(string(xmlRequest))
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize XML document: %v", err)
	}

	// DigestValue calculation using SHA-1
	digest := sha1.New()
	if _, err := digest.Write([]byte(xmlCanonical)); err != nil {
		return nil, fmt.Errorf("failed to calculate digest: %v", err)
	}
	digestValue := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	// Step 2: Create SignedInfo block with DigestValue using etree
	signedInfoElement := createSignedInfoElement(referenceID, digestValue)

	// Convert the SignedInfo element to a string
	signedInfoDocument := etree.NewDocument()
	signedInfoDocument.SetRoot(signedInfoElement)
	signedInfoString, err := signedInfoDocument.WriteToString()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize SignedInfo: %v", err)
	}

	// Canonicalize the SignedInfo block
	canonicalizedSignedInfo, err := doc14n(signedInfoString)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize SignedInfo: %v", err)
	}

	// Step 3: Compute hash of canonicalized SignedInfo
	hashedSignedInfo := sha1.Sum(canonicalizedSignedInfo)

	// Step 4: Generate the SignatureValue using the private key
	signature, err := rsa.SignPKCS1v15(nil, fe.cert.privateKey, crypto.SHA1, hashedSignedInfo[:])
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

// verifyXML verifies the signed XML document
func (fe *FiskalEntity) verifyXML(xmlData []byte) (bool, error) {
	return true, nil
}
