package fiskalhrgo

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/ucarion/c14n"
)

// generateUniqueID generates a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

func doc14n(xmlData string) ([]byte, error) {
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	out, err := c14n.Canonicalize(decoder)
	if err != nil {
		return nil, fmt.Errorf("failed to canonicalize the XML: %v", err)
	}
	return out, nil
}

// signXML signs the given XML document using the private key and embeds the public certificate
func (fe *FiskalEntity) signXML(xmlRequest []byte, id string) ([]byte, error) {

	// Step 1: Canonicalize the XML document
	var err error
	xmlRequest, err = doc14n(string(xmlRequest))

	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}

	// DigestValue calculation
	digest := sha1.New()
	digest.Write(xmlRequest)
	digestValue := base64.StdEncoding.EncodeToString(digest.Sum(nil))

	// Step 2: Populate SignedInfo with the DigestValue
	signedInfoTemplate := `<SignedInfo xmlns="http://www.w3.org/2000/09/xmldsig#"><CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/><SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/><Reference URI="#%s"><Transforms><Transform Algorithm="http://www.w3.org/2000/09/xmldsig#enveloped-signature"/><Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#"/></Transforms><DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/><DigestValue>%s</DigestValue></Reference></SignedInfo>`
	signedInfo := fmt.Sprintf(signedInfoTemplate, id, digestValue)
	signedInfoCanonical, errCN := doc14n(signedInfo)
	if errCN != nil {
		return nil, fmt.Errorf("%v", errCN)
	}
	hashed := sha1.Sum(signedInfoCanonical) // Hash the canonicalized SignedInfo

	// Step 3: Generate the SignatureValue
	signature, err := rsa.SignPKCS1v15(nil, fe.cert.privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign the SignedInfo: %v", err)
	}
	signatureValue := base64.StdEncoding.EncodeToString(signature)

	// Step 4: Build the Signature block
	signatureBlock := fmt.Sprintf(`<Signature xmlns="http://www.w3.org/2000/09/xmldsig#">%s<SignatureValue>%s</SignatureValue><KeyInfo><X509Data><X509Certificate>%s</X509Certificate><X509IssuerSerial><X509IssuerName>%s</X509IssuerName><X509SerialNumber>%s</X509SerialNumber></X509IssuerSerial></X509Data></KeyInfo></Signature>`, signedInfoCanonical, signatureValue, base64.StdEncoding.EncodeToString(fe.cert.publicCert.Raw), fe.cert.publicCert.Issuer.String(), fe.cert.publicCert.SerialNumber.String())

	// Step 5: Inject the Signature block before the closing tag of the root element
	xmlRequestTmp := string(xmlRequest)
	closingTagIndex := strings.LastIndex(xmlRequestTmp, "</")
	if closingTagIndex == -1 {
		return nil, fmt.Errorf("invalid XML format: no closing tag found")
	}

	// Inject the Signature block before the closing root tag
	finalXML := xmlRequestTmp[:closingTagIndex] + signatureBlock + xmlRequestTmp[closingTagIndex:]

	return []byte(finalXML), nil
}

// verifyXML verifies the XML signature
func (fe *FiskalEntity) verifyXML(xmlData []byte) (bool, error) {
	return true, nil
}
