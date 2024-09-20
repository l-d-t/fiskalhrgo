package fiskalhrgo

import (
	"errors"
	"fmt"
)

// Some important settings
const production_url = "https://cis.porezna-uprava.hr:8449/FiskalizacijaService"
const demo_url = "https://cistest.apis-it.hr:8449/FiskalizacijaServiceTest"
const CIStimeout = 10 //how long to wait for CIS response in seconds

// FiskalEntity represents an entity involved in the fiscalization process.
// It contains essential information and configurations required for generating
// and verifying fiscal invoices in compliance with Croatian fiscalization laws.
type FiskalEntity struct {
	OIB     string
	SustPDV bool
	Cert    *CertManager
	// true if the entity is in demo mode and will use the demo CIS certificate and endpoint
	demoMode bool
	// holds the public key, issuer, subject, serial number, and validity dates of a CIS certificate to check signature on CIS responses
	// also contains SSL root CA pool for SSL verification
	CIScert *signatureCheckCIScert
	url     string
}

// NewFiskalEntity creates a new FiskalEntity with default values and an optional CertManager.
//
// Parameters:
//   - oib: The taxpayer's OIB, which will be validated against the OIB in the certificate.
//   - sustavPDV: If true, the entity is part of the VAT system and will include VAT in the invoices.
//   - demoMode: If true, the entity is in demo mode and will use the demo CIS certificate and endpoint.
//   - certManager: If nil, a new CertManager is initialized using the provided certificate path and password.
//     Otherwise, the existing CertManager is used as is.
//   - chk_expired: If true, the entity creation will fail if the certificate is expired (recommended).
//   - certPath, certPassword: These are required if certManager is nil and are used to load the certificate.
//
// Certificate Handling and Expiry:
//   - If the certificate is expired and the `chk_expired` flag is set to true, the entity creation will fail.
//     This is recommended as invoices signed with an expired certificate will be rejected by the Croatian CIS,
//     and no JIR (unique invoice identifier) will be returned.
//   - The `chk_expired` flag exists for situations where an expired certificate must be loaded, such as when
//     recalculating the ZKI for older invoices. This is sometimes required by law to prove that the invoice was
//     originally signed with a valid certificate at the time of issuance.
//
// Best Practices:
//   - It is advisable to retain old certificates even after they expire, along with the ZKI, JIR, and the certificate's
//     serial number or fingerprint. This ensures traceability and proof of which certificate was used to sign each invoice.
//   - While expired certificates may be loaded to handle historical cases, it is recommended to always use a valid,
//     non-expired certificate when generating and sending new invoices.
func NewFiskalEntity(oib string, sustavPDV bool, demoMode bool, cert *CertManager, chk_expired bool, cert_config ...string) (*FiskalEntity, error) {

	// Check if OIB is valid
	if !ValidateOIB(oib) {
		return nil, errors.New("invalid OIB")
	}

	var CIScert *signatureCheckCIScert
	var CIScerterror error

	if demoMode {
		CIScert, CIScerterror = GetDemoPublicKey()
	} else {
		CIScert, CIScerterror = GetProductionPublicKey()
	}

	if CIScerterror != nil {
		return nil, fmt.Errorf("failed to get CIS public key and CA pool: %v", CIScerterror)
	}

	// Initialize a new CertManager if it's nil, otherwise use the provided one
	if cert == nil {

		// Check if the certificate path and password are provided
		if len(cert_config) < 2 {
			return nil, errors.New("certificate path and password are required")
		}

		cert = NewCertManager()
		err := cert.DecodeP12Cert(cert_config[0], cert_config[1])
		if err != nil {
			return nil, fmt.Errorf("cert decode fail: %v", err)
		}
	}

	if !cert.init_ok {
		return nil, errors.New("failed to initialize CertManager")
	}
	if cert.certOIB != oib {
		return nil, errors.New("OIB does not match the certificate")
	}
	if chk_expired && cert.expired {
		return nil, errors.New("certificate expired")
	}

	var url string
	if demoMode {
		url = demo_url
	} else {
		url = production_url
	}

	return &FiskalEntity{
		OIB:      oib,
		SustPDV:  sustavPDV,
		Cert:     cert,
		demoMode: demoMode,
		CIScert:  CIScert,
		url:      url,
	}, nil
}
