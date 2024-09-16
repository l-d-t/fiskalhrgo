package fiskalhrgo

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"

	"golang.org/x/crypto/pkcs12"
)

// CertManager holds the private key, public certificate, and additional info
type CertManager struct {
	privateKey *rsa.PrivateKey
	publicCert *x509.Certificate
	caCerts    []*x509.Certificate // This holds any CA certs
	certInfo   map[string]interface{}
}

// DecodeP12Cert loads and decodes a P12 certificate, extracting the private key, public cert, and CA certificates
func (cm *CertManager) DecodeP12Cert(certPath string, password string) error {
	// Read the P12 file
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		return fmt.Errorf("failed to read certificate: %v", err)
	}

	// Convert the P12 file to PEM blocks using the password
	pemBlocks, err := pkcs12.ToPEM(certBytes, password)
	if err != nil {
		return fmt.Errorf("failed to convert P12 to PEM: %v", err)
	}

	var privateKey *rsa.PrivateKey
	var certificate *x509.Certificate
	var caCerts []*x509.Certificate

	// Iterate over the PEM blocks to extract the private key, certificate, and CA certificates
	for _, block := range pemBlocks {
		switch block.Type {
		case "PRIVATE KEY":
			// Try parsing the key as PKCS8 first
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				// If PKCS8 parsing fails, try PKCS1
				key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
				if err != nil {
					return fmt.Errorf("failed to parse private key (tried PKCS8 and PKCS1): %v", err)
				}
			}
			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return fmt.Errorf("private key is not of RSA type")
			}
			privateKey = rsaKey
		case "CERTIFICATE":
			// Parse the certificate
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse certificate: %v", err)
			}
			// Check if it's a CA cert (we assume it's a CA cert if it's not self-issued)
			if cert.IsCA {
				caCerts = append(caCerts, cert)
			} else {
				certificate = cert
			}
		}
	}

	if privateKey == nil {
		return fmt.Errorf("private key not found in P12 file")
	}
	if certificate == nil {
		return fmt.Errorf("certificate not found in P12 file")
	}

	// Store the parsed certificate information
	cm.privateKey = privateKey
	cm.publicCert = certificate
	cm.caCerts = caCerts
	cm.certInfo = map[string]interface{}{
		"issuer":      certificate.Issuer.String(),
		"subject":     certificate.Subject.String(),
		"valid_from":  certificate.NotBefore,
		"valid_until": certificate.NotAfter,
	}

	return nil
}

// DisplayCertInfo prints the certificate details (Issuer, Subject, Validity, and CA certs)
func (cm *CertManager) DisplayCertInfo() {
	fmt.Println("Certificate Information:")
	fmt.Printf("Issuer: %s\n", cm.certInfo["issuer"])
	fmt.Printf("Subject: %s\n", cm.certInfo["subject"])
	fmt.Printf("Valid From: %s\n", cm.certInfo["valid_from"])
	fmt.Printf("Valid Until: %s\n", cm.certInfo["valid_until"])

	// Display CA certificates if present
	if len(cm.caCerts) > 0 {
		fmt.Println("CA Certificates:")
		for i, caCert := range cm.caCerts {
			fmt.Printf("CA Cert %d: Issuer: %s, Subject: %s\n", i+1, caCert.Issuer.String(), caCert.Subject.String())
		}
	} else {
		fmt.Println("No CA certificates found.")
	}
}
