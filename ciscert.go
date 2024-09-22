package fiskalhrgo

import (
	"crypto/rsa"
	"crypto/x509"
	"embed"
	"encoding/pem"
	"errors"
	"fmt"
	"path/filepath"
	"time"
)

// Embed demo certs
//
//go:embed certDemo/democis*.pem
var demoCISCert embed.FS

// Embed production certs
//
//go:embed certProd/fiskalcis*.pem
var prodCISCert embed.FS

// type signatureCheckCIScert holds the public key, issuer, subject, serial number, and validity dates
// of a CIS certificate to check signature on CIS responses. It also holds the SSL verify pool
type signatureCheckCIScert struct {
	PublicKey     *rsa.PublicKey
	Subject       string
	Serial        string
	Issuer        string
	ValidFrom     time.Time
	ValidUntil    time.Time
	SSLverifyPoll *x509.CertPool
}

// parseAndVerifyEmbeddedCerts parses the embedded certificates, verifies the chain, and returns the public key of the newest valid certificate
func parseAndVerifyEmbeddedCerts(certFS embed.FS, dir string, pattern string) (*signatureCheckCIScert, error) {
	var newestCert *x509.Certificate
	var roots *x509.CertPool

	// Read the embedded certificate files
	certFiles, err := certFS.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read embedded cert files: %w", err)
	}

	for _, certFile := range certFiles {

		if certFile.IsDir() {
			continue // Skip directories
		}

		if match, _ := filepath.Match(pattern, certFile.Name()); !match {
			continue
		}

		certData, err := certFS.ReadFile(dir + "/" + certFile.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file %s: %w", certFile.Name(), err)
		}

		// Parse the certificates
		var certs []*x509.Certificate
		for {
			block, rest := pem.Decode(certData)
			if block == nil {
				break
			}
			if block.Type != "CERTIFICATE" {
				return nil, errors.New("invalid PEM block type")
			}

			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse certificate: %w", err)
			}
			certs = append(certs, cert)
			certData = rest
		}

		// Verify the certificate chain
		roots = x509.NewCertPool()
		intermediates := x509.NewCertPool()

		// Add the root certificate to the roots pool
		roots.AddCert(certs[len(certs)-1])

		// Add intermediate certificates to the intermediates pool
		for i := 1; i < len(certs)-1; i++ {
			intermediates.AddCert(certs[i])
		}

		opts := x509.VerifyOptions{
			Roots:         roots,
			Intermediates: intermediates,
			CurrentTime:   time.Now(),
			KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		}

		leafCert := certs[0]
		if _, err := leafCert.Verify(opts); err != nil {
			continue // Skip invalid certificate chains
		}

		// Check if the certificate is valid and not expired
		now := time.Now()
		if now.Before(leafCert.NotBefore) || now.After(leafCert.NotAfter) {
			continue // Skip expired or not yet valid certificates
		}

		// Update the newest valid certificate
		if newestCert == nil || leafCert.NotBefore.After(newestCert.NotBefore) {
			newestCert = leafCert
		}
	}

	if newestCert == nil {
		return nil, errors.New("no suitable certificate found")
	}

	// Extract the public key from the newest valid certificate
	publicKey, ok := newestCert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key is not of type RSA")
	}

	return &signatureCheckCIScert{
		PublicKey:     publicKey,
		Subject:       newestCert.Subject.String(),
		Serial:        newestCert.SerialNumber.String(),
		Issuer:        newestCert.Issuer.String(),
		ValidFrom:     newestCert.NotBefore,
		ValidUntil:    newestCert.NotAfter,
		SSLverifyPoll: roots,
	}, nil
}

// Get demo public key
func getDemoPublicKey() (*signatureCheckCIScert, error) {
	return parseAndVerifyEmbeddedCerts(demoCISCert, "certDemo", "democis*.pem")
}

// Get production public key
func getProductionPublicKey() (*signatureCheckCIScert, error) {
	return parseAndVerifyEmbeddedCerts(prodCISCert, "certProd", "fiskalcis*.pem")
}
