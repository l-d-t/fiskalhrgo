package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

// certManager holds the private key, public certificate, and additional info
type certManager struct {
	privateKey  *rsa.PrivateKey
	publicCert  *x509.Certificate
	caCerts     []*x509.Certificate // This holds any CA certs
	certORG     string
	certOIB     string
	certSERIAL  string
	init_ok     bool
	expired     bool
	expire_soon bool
	expire_days uint16
}

func newCertManager() *certManager {
	return &certManager{
		privateKey:  nil,
		publicCert:  nil,
		caCerts:     []*x509.Certificate{},
		certORG:     "",
		certOIB:     "",
		certSERIAL:  "",
		init_ok:     false,
		expired:     false,
		expire_soon: false,
		expire_days: 0,
	}
}

// decodeP12Cert loads and decodes a P12 certificate, extracting the private key, public cert, and CA certificates
func (cm *certManager) decodeP12Cert(certPath string, password string) error {
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

	// Check if the certificate is expired
	now := time.Now()
	if now.Before(certificate.NotBefore) {
		return fmt.Errorf("certificate is not valid yet: valid from %v", certificate.NotBefore)
	}
	if now.After(certificate.NotAfter) {
		cm.expired = true
	}

	// Check if the certificate is expiring soon (within 30 days)
	daysUntilExpiration := certificate.NotAfter.Sub(now).Hours() / 24
	cm.expire_days = uint16(daysUntilExpiration)
	if daysUntilExpiration <= 30 {
		cm.expire_soon = true
	}

	// Extract the OIB
	oib, err := cm.getCertOIB()
	if err != nil {
		return fmt.Errorf("error extracting OIB: %v", err)
	}
	cm.certOIB = oib
	cm.certORG = certificate.Subject.Organization[0]

	cm.init_ok = true

	return nil
}

// getCertOIB extracts the OIB from the certificate's subject information
func (cm *certManager) getCertOIB() (string, error) {
	if cm.publicCert == nil {
		return "", fmt.Errorf("certificate not loaded")
	}

	// Extract the subject's organization (O) and country (C) fields
	organization := cm.publicCert.Subject.Organization
	country := cm.publicCert.Subject.Country

	if len(organization) == 0 || len(country) == 0 {
		return "", fmt.Errorf("organization or country fields missing in certificate")
	}

	// Try to extract the OIB by splitting the organization field at the country field
	ex := strings.Split(organization[0], country[0])
	if len(ex) < 2 {
		return "", fmt.Errorf("failed to extract OIB from certificate")
	}

	return ex[1], nil
}

func (cm *certManager) displayCertInfoText() string {
	if cm.publicCert == nil {
		return "No public certificate available."
	}

	result := "Certificate Information:\n"
	result += fmt.Sprintf("Issuer: %s\n", cm.publicCert.Issuer.String())
	result += fmt.Sprintf("Subject: %s\n", cm.publicCert.Subject.String())
	result += fmt.Sprintf("Serial Number: %s\n", cm.publicCert.SerialNumber.String())
	result += fmt.Sprintf("Valid From: %s\n", cm.publicCert.NotBefore.Format("02 Jan 2006 15:04:05 MST"))
	result += fmt.Sprintf("Valid Until: %s\n", cm.publicCert.NotAfter.Format("02 Jan 2006 15:04:05 MST"))

	// Display CA certificates if present
	if len(cm.caCerts) > 0 {
		result += "CA Certificates:\n"
		for i, caCert := range cm.caCerts {
			result += fmt.Sprintf("CA Cert %d: Issuer: %s, Subject: %s\n", i+1, caCert.Issuer.String(), caCert.Subject.String())
		}
	} else {
		result += "No CA certificates found.\n"
	}
	return result
}

func (cm *certManager) displayCertInfoMarkdown() string {
	if cm.publicCert == nil {
		return "No public certificate available."
	}

	result := "# Certificate Information\n"
	result += fmt.Sprintf("**Issuer**: %s\n\n", cm.publicCert.Issuer.String())
	result += fmt.Sprintf("**Subject**: %s\n\n", cm.publicCert.Subject.String())
	result += fmt.Sprintf("**Serial Number**: %s\n\n", cm.publicCert.SerialNumber.String())
	result += fmt.Sprintf("**Valid From**: %s\n\n", cm.publicCert.NotBefore.Format("02 Jan 2006 15:04:05 MST"))
	result += fmt.Sprintf("**Valid Until**: %s\n\n", cm.publicCert.NotAfter.Format("02 Jan 2006 15:04:05 MST"))

	// Display CA certificates if present
	if len(cm.caCerts) > 0 {
		result += "## CA Certificates:\n"
		for i, caCert := range cm.caCerts {
			result += fmt.Sprintf("- CA Cert %d: **Issuer**: %s, **Subject**: %s\n", i+1, caCert.Issuer.String(), caCert.Subject.String())
		}
	} else {
		result += "No CA certificates found.\n"
	}
	return result
}

func (cm *certManager) displayCertInfoHTML() string {
	if cm.publicCert == nil {
		return "<p>No public certificate available.</p>"
	}

	result := "<h3>Certificate Information</h3>"
	result += fmt.Sprintf("<p><strong>Issuer:</strong> %s</p>", cm.publicCert.Issuer.String())
	result += fmt.Sprintf("<p><strong>Subject:</strong> %s</p>", cm.publicCert.Subject.String())
	result += fmt.Sprintf("<p><strong>Serial Number:</strong> %s</p>", cm.publicCert.SerialNumber.String())
	result += fmt.Sprintf("<p><strong>Valid From:</strong> %s</p>", cm.publicCert.NotBefore.Format("02 Jan 2006 15:04:05 MST"))
	result += fmt.Sprintf("<p><strong>Valid Until:</strong> %s</p>", cm.publicCert.NotAfter.Format("02 Jan 2006 15:04:05 MST"))

	// Display CA certificates if present
	if len(cm.caCerts) > 0 {
		result += "<h4>CA Certificates</h4><ul>"
		for i, caCert := range cm.caCerts {
			result += fmt.Sprintf("<li>CA Cert %d: <strong>Issuer:</strong> %s, <strong>Subject:</strong> %s</li>", i+1, caCert.Issuer.String(), caCert.Subject.String())
		}
		result += "</ul>"
	} else {
		result += "<p>No CA certificates found.</p>"
	}
	return result
}

func (cm *certManager) displayCertInfoKeyPoints() [][2]string {
	var result [][2]string

	if cm.publicCert == nil {
		return append(result, [2]string{"Error", "No public certificate available."})
	}

	result = append(result, [2]string{"Title", "Certificate Information"})
	result = append(result, [2]string{"Issuer", cm.publicCert.Issuer.String()})
	result = append(result, [2]string{"Subject", cm.publicCert.Subject.String()})
	result = append(result, [2]string{"Serial Number", cm.publicCert.SerialNumber.String()})
	result = append(result, [2]string{"Valid From", cm.publicCert.NotBefore.Format("02 Jan 2006 15:04:05 MST")})
	result = append(result, [2]string{"Valid Until", cm.publicCert.NotAfter.Format("02 Jan 2006 15:04:05 MST")})

	// Display CA certificates if present
	if len(cm.caCerts) > 0 {
		result = append(result, [2]string{"Title", "CA Certificates"})
		for i, caCert := range cm.caCerts {
			result = append(result, [2]string{fmt.Sprintf("CA Cert %d Issuer", i+1), caCert.Issuer.String()})
			result = append(result, [2]string{fmt.Sprintf("CA Cert %d Subject", i+1), caCert.Subject.String()})
		}
	} else {
		result = append(result, [2]string{"Info", "No CA certificates found."})
	}

	return result
}
