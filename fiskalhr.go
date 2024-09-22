package fiskalhrgo

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Some important constants
const production_url = "https://cis.porezna-uprava.hr:8449/FiskalizacijaService"
const demo_url = "https://cistest.apis-it.hr:8449/FiskalizacijaServiceTest"
const cistimeout = 10 //how long to wait at max for CIS response in seconds

// FiskalEntity represents an entity involved in the fiscalization process.
// It contains essential information and configurations required for generating
// and verifying fiscal invoices in compliance with Croatian fiscalization laws.
type FiskalEntity struct {
	// oib is the taxpayer's identification number in Croatia (oib) and must match the oib in the certificate.
	// This is a mandatory field for fiscalization.
	oib string

	// sustPDV indicates whether the entity is part of the VAT system.
	// If true, the entity will include VAT in the invoices.
	sustPDV bool

	// locationID is the unique identifier of the location where the fiscalization is taking place.
	// This identifier is alphanumeric and must be registered in the ePorezna system.
	locationID string

	// centralizedInvoiceNumber specifies whether invoice numbers are centralized per locationID.
	// If true, invoice numbers are centralized for the entire location.
	// If false, each register device within the location has its own sequence of invoice numbers.
	centralizedInvoiceNumber bool

	// cert holds the certificate and private key used to sign invoices.
	cert *certManager

	// ciscert holds the public key, issuer, subject, serial number, and validity dates of a CIS certificate.
	// It is used to check the signature on CIS responses and contains the SSL root CA pool for SSL verification.
	ciscert *signatureCheckCIScert

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
//     If using in demo mode the location don't have to be registered. And can use any alphanumeric value.
//   - centralizedInvoiceNumber: If true, invoice numbers are centralized for the entire location.
//     If not, each register device within the location has its own sequence of invoice numbers.
//   - demoMode: If true, the entity is in demo mode and will use the demo CIS certificate and endpoint.
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
//
// Best Practices:
//   - It is advisable to retain old certificates even after they expire, along with the ZKI, JIR, and the certificate's
//     serial number or fingerprint. This ensures traceability and proof of which certificate was used to sign each invoice.
//   - While expired certificates may be loaded to handle historical cases, it is mandatory to always use a valid,
//     non-expired certificate when generating and sending new invoices.
//
// Returns:
//   - (*FiskalEntity, error): A pointer to a new FiskalEntity instance with the provided values, or an error if the input is invalid.
func NewFiskalEntity(oib string, sustavPDV bool, locationID string, centralizedInvoiceNumber bool, demoMode bool, chk_expired bool, certPath string, certPassword string) (*FiskalEntity, error) {

	// Check if OIB is valid
	if !ValidateOIB(oib) {
		return nil, errors.New("invalid OIB")
	}

	//check if locationID is valid
	if !ValidateLocationID(locationID) {
		return nil, errors.New("invalid locationID")
	}

	//check path is valid
	if !IsFileReadable(certPath) {
		return nil, errors.New("invalid certificate path or file not readable")
	}

	var CIScert *signatureCheckCIScert
	var CIScerterror error

	if demoMode {
		CIScert, CIScerterror = getDemoPublicKey()
	} else {
		CIScert, CIScerterror = getProductionPublicKey()
	}

	if CIScerterror != nil {
		return nil, fmt.Errorf("failed to get CIS public key and CA pool: %v", CIScerterror)
	}

	cert := newCertManager()
	err := cert.decodeP12Cert(certPath, certPassword)
	if err != nil {
		return nil, fmt.Errorf("certificate decode fail: %v", err)
	}

	if !cert.init_ok {
		return nil, errors.New("failed to initialize the certificate manager")
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
		oib:                      oib,
		sustPDV:                  sustavPDV,
		locationID:               locationID,
		centralizedInvoiceNumber: centralizedInvoiceNumber,
		cert:                     cert,
		demoMode:                 demoMode,
		ciscert:                  CIScert,
		url:                      url,
	}, nil
}

// OIB returns the taxpayer's identification number.
func (fe *FiskalEntity) OIB() string {
	return fe.oib
}

// SustPDV indicates whether the entity is part of the VAT system.
func (fe *FiskalEntity) SustPDV() bool {
	return fe.sustPDV
}

// LocationID returns the unique identifier of the location where the fiscalization is taking place.
func (fe *FiskalEntity) LocationID() string {
	return fe.locationID
}

// CentralizedInvoiceNumber specifies whether invoice numbers are centralized per locationID. Or each register device within the location has its own sequence of invoice numbers.
func (fe *FiskalEntity) CentralizedInvoiceNumber() bool {
	return fe.centralizedInvoiceNumber
}

// DemoMode indicates whether the entity is in demo mode (Demo Fiskalizacija).
func (fe *FiskalEntity) DemoMode() bool {
	return fe.demoMode
}

func (fe *FiskalEntity) DisplayCertInfoText() string {
	return fe.cert.displayCertInfoText()
}

func (fe *FiskalEntity) DisplayCertInfoMarkdown() string {
	return fe.cert.displayCertInfoMarkdown()
}

func (fe *FiskalEntity) DisplayCertInfoHTML() string {

	return fe.cert.displayCertInfoHTML()
}

func (fe *FiskalEntity) DisplayCertInfoKeyPoints() [][2]string {

	return fe.cert.displayCertInfoKeyPoints()
}

// GetCertORG returns the organization name from the certificate.
// The organization name is typically included in the certificate's subject field.
func (fe *FiskalEntity) GetCertORG() string {
	return fe.cert.certORG
}

// GetCertSERIAL returns the serial number from the certificate.
// The serial number is a unique identifier assigned by the certificate issuer.
func (fe *FiskalEntity) GetCertSERIAL() string {
	return fe.cert.certSERIAL
}

// IsExpired returns whether the certificate is expired.
// This indicates if the certificate's validity period has ended.
func (fe *FiskalEntity) IsExpired() bool {
	return fe.cert.expired
}

// IsExpiringSoon returns whether the certificate is expiring soon.
// This indicates if the certificate is approaching its expiration date.
func (fe *FiskalEntity) IsExpiringSoon() bool {
	return fe.cert.expire_soon
}

// DaysUntilExpire returns the number of days until the certificate expires.
// This provides a countdown of days remaining before the certificate becomes invalid.
func (fe *FiskalEntity) DaysUntilExpire() uint16 {
	return fe.cert.expire_days
}

// GenerateZKI generates the ZKI (Zaštitni Kod Izdavatelja) based on the given data.
// The ZKI is a unique identifier for the invoice, generated by signing the concatenated
// invoice data with the taxpayer's private key and hashing the signature.
//
// Parameters:
//
//   - ssueDateTime time.Time: The date and time when the invoice was issued.
//   - invoiceNumber uint: The unique number of the invoice.
//   - deviceID uint: The unique identifier of the device issuing the invoice.
//   - totalAmount string: The total amount of the invoice, formatted as a string with 2 decimal places (e.g., "100.00").
//
// Returns:
//   - string: The generated ZKI as a hexadecimal string.
//   - error: An error if the ZKI generation fails, otherwise nil.
func (entity *FiskalEntity) GenerateZKI(issueDateTime time.Time, invoiceNumber uint, deviceID uint, totalAmount string) (string, error) {

	formattedTime := issueDateTime.Format("02.01.2006 15:04:05")

	// Ensure totalAmount is a valid decimal string with 2 decimal places
	if !IsValidCurrencyFormat(totalAmount) {
		return "", errors.New("invalid totalAmount format; expected a string with 2 decimal places (e.g., 100.00)")
	}

	// Convert invoiceNumber and deviceID from uint to string
	invoiceNumberStr := strconv.FormatUint(uint64(invoiceNumber), 10)
	deviceIDStr := strconv.FormatUint(uint64(deviceID), 10)

	// Concatenate the required data (oib, date, invoice number, location, device ID, total amount)
	guardCode := entity.oib + formattedTime + invoiceNumberStr + entity.locationID + deviceIDStr + totalAmount

	// Hash the concatenated data using SHA1
	hashed := sha1.Sum([]byte(guardCode))

	// Use the private key from the CertManager to sign the hashed data with RSA and SHA1
	var signature []byte
	signature, err := rsa.SignPKCS1v15(rand.Reader, entity.cert.privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Generate the MD5 hash of the signature
	md5Hash := md5.Sum(signature)

	// Return the ZKI as a hexadecimal string
	zki := fmt.Sprintf("%x", md5Hash[:])
	return zki, nil
}

// CISEcho sends an echo request to CIS and processes the response.
func (fe *FiskalEntity) CISEcho(text string) (string, error) {
	// Create an XML payload for the echo request
	echoRequest := &EchoRequest{
		Xmlns: DefaultNamespace,
		Text:  text,
	}

	xmlPayload, err := xml.Marshal(echoRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML payload: %w", err)
	}

	body, _, err := fe.GetResponse(xmlPayload)
	if err != nil {
		return "", err
	}

	// Process the XML response
	var echoResponse EchoResponse
	if err := xml.Unmarshal(body, &echoResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal XML response: %w", err)
	}

	return echoResponse.Text, nil
}

// PingCIS checks if connection and message exchange with CIS works using the CISEcho function.
// It sends a simple text message to CIS and expects the same message back.
// Returns:
//   - nil if the ping was successful
//   - error if the ping failed
func (fe *FiskalEntity) PingCIS() error {
	echoText := "Hello, CIS, from FiskalhrGo!"
	response, err := fe.CISEcho(echoText)
	if err != nil {
		return fmt.Errorf("CIS ping failed: %v", err)
	}
	if response != echoText {
		return fmt.Errorf("CIS ping failed: unexpected response")
	}
	return nil
}
