package fiskalhrgo

import (
	"errors"
	"fmt"
)

// Some important constants
const production_url = "https://cis.porezna-uprava.hr:8449/FiskalizacijaService"
const demo_url = "https://cistest.apis-it.hr:8449/FiskalizacijaServiceTest"
const CIStimeout = 10 //how long to wait at max for CIS response in seconds

// FiskalEntity represents an entity involved in the fiscalization process.
// It contains essential information and configurations required for generating
// and verifying fiscal invoices in compliance with Croatian fiscalization laws.
type FiskalEntity struct {
	// OIB is the taxpayer's identification number in Croatia (OIB) and must match the OIB in the certificate.
	// This is a mandatory field for fiscalization.
	OIB string

	// SustPDV indicates whether the entity is part of the VAT system.
	// If true, the entity will include VAT in the invoices.
	SustPDV bool

	// locationID is the unique identifier of the location where the fiscalization is taking place.
	// This identifier is alphanumeric and must be registered in the ePorezna system.
	locationID string

	// centralizedInvoiceNumber specifies whether invoice numbers are centralized per locationID.
	// If true, invoice numbers are centralized for the entire location.
	// If false, each register device within the location has its own sequence of invoice numbers.
	centralizedInvoiceNumber bool

	// Cert holds the certificate and private key used to sign invoices.
	Cert *CertManager

	// CIScert holds the public key, issuer, subject, serial number, and validity dates of a CIS certificate.
	// It is used to check the signature on CIS responses and contains the SSL root CA pool for SSL verification.
	CIScert *signatureCheckCIScert

	// demoMode indicates whether the entity is in demo mode.
	// If true, the entity will use the demo CIS certificate and endpoint for testing purposes.
	demoMode bool

	// url is the endpoint URL for the CIS service.
	// This URL is used to send fiscalization requests to the CIS system.
	url string
}

// NewFiskalEntity creates a new FiskalEntity with provided values, validates certificates and input before returning an entity.
//
// Parameters:
//   - oib: The taxpayer's OIB, which will be validated against the OIB in the certificate.
//   - sustavPDV: If true, the entity is part of the VAT system and will include VAT in the invoices.
//   - locationID: The unique identifier of the location where fiscalization is taking place. This identifier must be
//     registered with ePorezna, is case-sensitive, and must be identical to the one registered there.
//	   If using in demo mode the location don't have to be registered. And can use any alphanumeric value.
//   - centralizedInvoiceNumber: If true, invoice numbers are centralized for the entire location.
//	   If not, each register device within the location has its own sequence of invoice numbers.
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
//     recalculating the ZKI for older invoices. This is sometimes required by law during an inspection to prove that the invoice was
//     not modified after the original fiscalization took place. It is recommended to save for each invoice a pointer or identifier of
//     certificate used to generate the ZKI at the time not only the ZKI itself. And to keep in store old certificates. Normally fiskal certificates
//     are valid for 5 years.

// Best Practices:
//   - It is advisable to retain old certificates even after they expire, along with the ZKI, JIR, and the certificate's
//     serial number or fingerprint. This ensures traceability and proof of which certificate was used to sign each invoice.
//   - While expired certificates may be loaded to handle historical cases, it is mandatory to always use a valid,
//     non-expired certificate when generating and sending new invoices.
//
// Returns:
//   - (*FiskalEntity, error): A pointer to a new FiskalEntity instance with the provided values, or an error if the input is invalid.
func NewFiskalEntity(oib string, sustavPDV bool, locationID string, centralizedInvoiceNumber bool, demoMode bool, cert *CertManager, chk_expired bool, cert_config ...string) (*FiskalEntity, error) {

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
		OIB:                      oib,
		SustPDV:                  sustavPDV,
		locationID:               locationID,
		centralizedInvoiceNumber: centralizedInvoiceNumber,
		Cert:                     cert,
		demoMode:                 demoMode,
		CIScert:                  CIScert,
		url:                      url,
	}, nil
}
