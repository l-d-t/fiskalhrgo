package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

// PaymentMethod defines a custom type for means of payment
type PaymentMethod string

// Constants representing allowed values for PaymentMethod
const (
	CISCash         PaymentMethod = "G" // Cash
	CISCard         PaymentMethod = "K" // Card
	CISMixOther     PaymentMethod = "O" // Mix/Other
	CISBankTransfer PaymentMethod = "T" // Bank Transfer (usually not sent to CIS, not mandatory)
	CISCheck        PaymentMethod = "C" // Check [deprecated]
)

// IsValid checks if PaymentMethod is one of the allowed values
func (p PaymentMethod) IsValid() error {
	switch p {
	case CISCash, CISCard, CISMixOther, CISBankTransfer, CISCheck:
		return nil
	default:
		return errors.New("PaymentMethod must be one of the following values: G - Cash, K - Card, O - Mix/Other, T - Bank Transfer, C - Check (deprecated)")
	}
}

// NewCISInvoice initializes and returns a RacunType instance
//
// This method creates a new instance of RacunType, which represents an invoice with all necessary fields.
// The instance can be marshaled to XML and sent to the CIS for fiscalization.
// ALWAYS use the provided methods to set or modify the values of the RacunType instance.
// Using the provided methods is safe end ensure the correct format of the data and no discrepancies with the ZKI.
// DO NOT MODIFY the returned RacunType instance directly, as it may lead to invalid XML output or ZKI problems.
// Unsafe and unexpected results may happen if you modify the data manually and use it to send later.
// It can result in wrong data sent to CIS or discrepancy between sent data and ZKI, witch can lead to severe consequences.
//
// Parameters:
//
//   - dateTime (time.Time): The date and time of the invoice.
//   - centralized (bool): Indicates whether the sequence mark is centralized.
//   - invoiceNumber (uint): The unique number of the invoice.
//   - locationIdentifier (string): The identifier for the business location where the invoice was issued.
//   - registerDeviceID (uint): The identifier for the cash register device used to issue the invoice.
//   - pdvValues ([][]interface{}): A 2D array for VAT details (nullable).
//   - pnpValues ([][]interface{}): A 2D array for consumption tax details (nullable).
//   - ostaliPorValues ([][]interface{}): A 2D array for other tax details (nullable).
//   - iznosOslobPdv (string): The amount exempt from VAT.
//   - iznosMarza (string): The margin amount.
//   - iznosNePodlOpor (string): The amount not subject to taxation.
//   - naknadeValues ([][]string): A 2D array for fees details (nullable).
//   - iznosUkupno (string): The total amount.
//   - paymentMethod (string): The payment method.
//   - oibOper (string): The OIB of the operator.
//   - attachedDocumentJIRorZKI (string): The JIR or ZKI of the attached document (empty if no attached document).
//
// Returns:
//
//	(*RacunType, string, error): A pointer to a new RacunType instance with the provided values, generated zki or an error if the input is invalid.
func (fe *FiskalEntity) NewCISInvoice(
	dateTime time.Time,
	invoiceNumber uint,
	registerDeviceID uint,
	pdvValues [][]interface{},
	pnpValues [][]interface{},
	ostaliPorValues [][]interface{},
	iznosOslobPdv string,
	iznosMarza string,
	iznosNePodlOpor string,
	naknadeValues [][]string,
	iznosUkupno string,
	paymentMethod PaymentMethod,
	oibOper string,
) (*RacunType, string, error) {
	// Format the date and time
	formattedDate := dateTime.Format("02.01.2006T15:04:05")

	// Determine the sequence mark
	oznSlijed := "N"
	if fe.centralizedInvoiceNumber {
		oznSlijed = "P"
	}

	if !IsValidCurrencyFormat(iznosUkupno) {
		return nil, "", errors.New("the total amount must be a valid currency format")
	}

	if !IsValidCurrencyFormat(iznosOslobPdv) {
		return nil, "", errors.New("the amount exempt from VAT must be a valid currency format")
	}

	if !IsValidCurrencyFormat(iznosMarza) {
		return nil, "", errors.New("the margin amount must be a valid currency format")
	}

	if !IsValidCurrencyFormat(iznosNePodlOpor) {
		return nil, "", errors.New("the amount not subject to taxation must be a valid currency format")
	}

	if iznosOslobPdv == "0.00" {
		iznosOslobPdv = ""
	}
	if iznosMarza == "0.00" {
		iznosMarza = ""
	}
	if iznosNePodlOpor == "0.00" {
		iznosNePodlOpor = ""
	}

	// Use helper functions to create the necessary types
	var pdv *PdvType
	var err error
	if pdvValues != nil {
		pdv, err = newPdv(pdvValues)
		if err != nil {
			return nil, "", err
		}
	}

	var pnp *PorezNaPotrosnjuType
	if pnpValues != nil {
		pnp, err = newPNP(pnpValues)
		if err != nil {
			return nil, "", err
		}
	}

	var ostaliPor *OstaliPoreziType
	if ostaliPorValues != nil {
		ostaliPor, err = otherTaxes(ostaliPorValues)
		if err != nil {
			return nil, "", err
		}
	}

	var naknade *NaknadeType
	if naknadeValues != nil {
		naknade, err = genNaknade(naknadeValues)
		if err != nil {
			return nil, "", err
		}
	}

	// Create the BrojRacunaType instance
	brRac := &BrojRacunaType{
		BrOznRac: invoiceNumber,
		OznPosPr: fe.locationID,
		OznNapUr: registerDeviceID,
	}

	//check means of payment can be:  G - Cash, K - Card, O - Mix/other
	//								, T - Bank transfer (usually not sent to CIS not mandatory)
	//                              , C - Check [deprecated]
	err = paymentMethod.IsValid()
	if err != nil {
		return nil, "", err
	}

	zki, err := fe.GenerateZKI(dateTime, invoiceNumber, registerDeviceID, iznosUkupno)

	if err != nil {
		return nil, "", err
	}

	return &RacunType{
		Oib:                fe.oib,
		USustPdv:           fe.sustPDV,
		DatVrijeme:         formattedDate,
		OznSlijed:          oznSlijed,
		BrRac:              brRac,
		Pdv:                pdv,
		Pnp:                pnp,
		OstaliPor:          ostaliPor,
		IznosOslobPdv:      iznosOslobPdv,
		IznosMarza:         iznosMarza,
		IznosNePodlOpor:    iznosNePodlOpor,
		Naknade:            naknade,
		IznosUkupno:        iznosUkupno,
		NacinPlac:          string(paymentMethod),
		OibOper:            oibOper,
		ZastKod:            zki,
		NakDost:            false,
		pointerToEntity:    fe,
		oldEntityForOldZKI: nil,
	}, zki, nil
}

func (invoice *RacunType) GetZKI() string {
	return invoice.ZastKod
}

func (invoice *RacunType) GetOib() string {
	return invoice.Oib
}

// Set late delivery to true, and set the ZKI you pass from saved data when you issued the invoice to customer
// Don't worry the ZKI you set will be validated with the current certificate before sending unless to set
// IhaveZKIwithExpiredCertificateEdgeCase method then the old certificate provided will be used to validate the ZKI
//
// So just set the ZKI you got from the invoice you issued to the customer,
// and the system will validate it with the current certificate
func (invoice *RacunType) SetLateDelivery(ZKI string) error {
	invoice.ZastKod = ZKI
	invoice.NakDost = true

	invoiceTime, err := time.Parse("02.01.2006T15:04:05", invoice.DatVrijeme)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	// Validate the ZKI with the current certificate
	calculatedZKI, err := invoice.pointerToEntity.GenerateZKI(invoiceTime, uint(invoice.BrRac.BrOznRac), uint(invoice.BrRac.OznNapUr), invoice.IznosUkupno)

	if err != nil {
		return fmt.Errorf("failed to generate ZKI: %w", err)
	}

	if calculatedZKI != invoice.ZastKod {
		return errors.New("ZKI is not valid")
	}

	return nil
}

// IhaveZKIwithExpiredCertificateEdgeCase sets the old FiskalEntity instance for the old ZKI verification
// This is used in the edge case that the ZKI was generated with one certificate and the fiscalization failed
// But the certificate expired or had to be changed and now fiscalization have to be repeated with new certificate
// If we replace the original ZKI its a problem we already gave the invoice with old ZKI out
// So we have to keep the old ZKI and validate it with the old certificate before signing and sending with new one
func (invoice *RacunType) IhaveZKIwithExpiredCertificateEdgeCase(oldZKI string, oldCertPath string, oldCertPassword string) error {
	invoice.ZastKod = oldZKI
	invoice.NakDost = true

	// Create a new old FiskalEntity
	oldFiskalEntity, err := NewFiskalEntity(
		invoice.pointerToEntity.oib,
		invoice.pointerToEntity.sustPDV,
		invoice.pointerToEntity.locationID,
		invoice.pointerToEntity.centralizedInvoiceNumber,
		invoice.pointerToEntity.demoMode,
		false,
		oldCertPath,
		oldCertPassword,
	)
	if err != nil {
		return fmt.Errorf("failed to create FiskalEntity: %v", err)
	}

	invoice.oldEntityForOldZKI = oldFiskalEntity

	invoiceTime, err := time.Parse("02.01.2006T15:04:05", invoice.DatVrijeme)
	if err != nil {
		return fmt.Errorf("failed to parse date: %w", err)
	}

	// Validate the ZKI with the current certificate
	calculatedZKI, err := invoice.oldEntityForOldZKI.GenerateZKI(invoiceTime, uint(invoice.BrRac.BrOznRac), uint(invoice.BrRac.OznNapUr), invoice.IznosUkupno)

	if err != nil {
		return fmt.Errorf("failed to generate ZKI: %w", err)
	}

	if calculatedZKI != invoice.ZastKod {
		return errors.New("ZKI is not valid")
	}

	return nil
}

// InvoiceRequest sends an invoice request to the CIS (Croatian Fiscalization System) and processes the response.
//
// This function performs the following steps:
//  1. Minimally validates the provided invoice for required fields
//     (any business logic and math is the responsibility of the invoicing application using the library)
//     PLEASE NOTE: the CIS also don't do any extensive validation of the invoice, only basic checks.
//     so you could get a JIR back even if the invoice is not correct.
//     But if you do that you can have problems later with inspections or periodic CIS checks of the data.
//     The library will send the data as is to the CIS.
//     So please validate and chek the invoice data according to you business logic
//     before sending it to the CIS.
//  2. Sends the XML request to the CIS and receives the response.
//  3. Unmarshals the response XML to extract the response data.
//  4. Checks for errors in the response and aggregates them if any are found.
//  5. Returns the JIR (Unique Invoice Identifier) if the request was successful.
//
// Parameters:
// - invoice: A pointer to a RacunType struct representing the invoice to be sent.
//
// Returns:
// - A string representing the JIR (Unique Invoice Identifier) if the request was successful.
// - A string representing the ZKI (Protection Code of the Issuer) from the invoice.
// - An error if any issues occurred during the process.
//
// Possible errors:
// - If the invoice is nil or something is invalid (only basic checks).
// - If the SpecNamj field of the invoice is not empty.
// - If the ZastKod field of the invoice is empty.
// - If there is an error marshalling the request to XML.
// - If there is an error making the request to the CIS.
// - If there is an error unmarshalling the response XML.
// - If the IdPoruke in the response does not match the request.
// - If the response status is not 200 and there are errors in the response.
// - If the JIR in the response is empty.
// - If an unexpected error occurs.
func (invoice *RacunType) InvoiceRequest() (string, string, error) {

	//some basic tests for invoice
	if invoice == nil {
		return "", "", errors.New("invoice is nil")
	}

	if invoice.SpecNamj != "" {
		return "", "", errors.New("invoice SpecNamj must be empty")
	}

	if invoice.ZastKod == "" {
		return "", "", errors.New("invoice ZKI (Zastitni Kod Izdavatelja) must be set")
	}

	//check ZKI
	invoiceTime, err := time.Parse("02.01.2006T15:04:05", invoice.DatVrijeme)
	if err != nil {
		return "", invoice.ZastKod, fmt.Errorf("failed to parse date: %w", err)
	}

	var chkEntity *FiskalEntity
	if invoice.oldEntityForOldZKI != nil {
		chkEntity = invoice.oldEntityForOldZKI
	} else {
		chkEntity = invoice.pointerToEntity
	}

	// Validate the ZKI with the old certificate
	calculatedZKI, err := chkEntity.GenerateZKI(invoiceTime, uint(invoice.BrRac.BrOznRac), uint(invoice.BrRac.OznNapUr), invoice.IznosUkupno)

	if err != nil {
		return "", invoice.ZastKod, fmt.Errorf("failed to check ZKI: %w", err)
	}

	if calculatedZKI != invoice.ZastKod {
		return "", invoice.ZastKod, errors.New("ZKI is not valid")
	}

	//Combine with zahtjev for final XML
	zahtjev := RacunZahtjev{
		Zaglavlje: newFiskalHeader(),
		Racun:     invoice,
		Xmlns:     DefaultNamespace,
		IdAttr:    generateUniqueID(),
	}

	// Marshal the RacunZahtjev to XML
	xmlData, err := xml.MarshalIndent(zahtjev, "", " ")
	if err != nil {
		return "", invoice.ZastKod, fmt.Errorf("error marshalling RacunZahtjev: %w", err)
	}

	// Let's send it to CIS
	body, status, errComm := invoice.pointerToEntity.GetResponse(xmlData, true)

	if errComm != nil {
		return "", invoice.ZastKod, fmt.Errorf("failed to make request: %w", errComm)
	}

	//unmarshad body to get Racun Odgovor
	var racunOdgovor RacunOdgovor
	if err := xml.Unmarshal(body, &racunOdgovor); err != nil {
		return "", invoice.ZastKod, fmt.Errorf("failed to unmarshal XML response: %w", err)
	}

	if zahtjev.Zaglavlje.IdPoruke != racunOdgovor.Zaglavlje.IdPoruke {
		return "", invoice.ZastKod, errors.New("IdPoruke mismatch")
	}

	if status != 200 {

		// Aggregate all errors into a single error message
		var errorMessages []string
		for _, greska := range racunOdgovor.Greske.Greska {
			errorMessages = append(errorMessages, fmt.Sprintf("%s: %s", greska.SifraGreske, greska.PorukaGreske))
		}
		if len(errorMessages) > 0 {
			return "", invoice.ZastKod, fmt.Errorf("errors in response: %s", strings.Join(errorMessages, "; "))
		}

	} else {
		if ValidateJIR(racunOdgovor.Jir) {
			return racunOdgovor.Jir, invoice.ZastKod, nil
		} else {
			return "", invoice.ZastKod, errors.New("JIR is not valid")
		}
	}

	// Add a default return statement to handle unexpected cases
	return "", invoice.ZastKod, errors.New("unexpected error")
}

// genNaknade initializes and returns a NaknadeType instance
//
// This function creates a new instance of NaknadeType, which represents a collection of fees (NaknadaType) entries.
// It takes a 2D array of values where each inner array represents a single fee entry with the name and amount.
//
// Parameters:
//
//	values ([][]string): A 2D array where each inner array contains two elements:
//	  - string: The name of the fee (NazivN)
//	  - string: The amount of the fee (IznosN)
//
// Returns:
//
//	(*NaknadeType, error): A pointer to a new NaknadeType instance with the provided fee entries, or an error if the input is invalid.
//
// Example:
//
//	values := [][]string{
//	  {"Service Fee", "100"},
//	  {"Delivery Fee", "50"},
//	}
//	naknade, err := genNaknade(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("genNaknade: %+v\n", naknade)
//	}
func genNaknade(values [][]string) (*NaknadeType, error) {
	naknade := make([]*NaknadaType, len(values))
	for i, v := range values {
		if len(v) != 2 {
			return nil, errors.New("each inner array must contain exactly two elements")
		}
		feeName := v[0]
		feeAmount := v[1]
		if !IsValidCurrencyFormat(feeAmount) {
			return nil, errors.New("the second element of each inner array must be a valid currency format (fee amount)")
		}
		naknade[i] = &NaknadaType{NazivN: feeName, IznosN: feeAmount}
	}
	return &NaknadeType{Naknada: naknade}, nil
}

// otherTaxes initializes and returns an OstaliPoreziType instance
//
// This function creates a new instance of OstaliPoreziType, which represents a collection of other taxes (PorezOstaloType) entries.
// It takes a 2D array of values where each inner array represents a single other tax entry with the name, rate, base, and amount.
//
// Parameters:
//
//	values ([][]interface{}): A 2D array where each inner array contains four elements:
//	  - string: The name of the tax (Naziv)
//	  - int: The tax rate (Stopa)
//	  - string: The tax base (Osnovica)
//	  - string: The tax amount (Iznos)
//
// Returns:
//
//	(*OstaliPoreziType, error): A pointer to a new OstaliPoreziType instance with the provided other tax entries, or an error if the input is invalid.
//
// Example:
//
//	values := [][]interface{}{
//	  {"Other Tax", 5, "1000", "50"},
//	}
//	ostaliPorezi, err := otherTaxes(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("OstaliPorezi: %+v\n", ostaliPorezi)
//	}
func otherTaxes(values [][]interface{}) (*OstaliPoreziType, error) {
	porezi := make([]*PorezOstaloType, len(values))
	for i, v := range values {
		if len(v) != 4 {
			return nil, errors.New("each inner array must contain exactly four elements")
		}
		name, ok := v[0].(string)
		if !ok {
			return nil, errors.New("the first element of each inner array must be a string (name)")
		}
		rate, ok := v[1].(string)
		if !ok {
			return nil, errors.New("the second element of each inner array must be an int (rate)")
		}
		base, ok := v[2].(string)
		if !ok {
			return nil, errors.New("the third element of each inner array must be a string (base)")
		}
		if !IsValidCurrencyFormat(base) {
			return nil, errors.New("the third element of each inner array must be a valid currency format (base)")
		}
		amount, ok := v[3].(string)
		if !ok {
			return nil, errors.New("the fourth element of each inner array must be a string (amount)")
		}
		if !IsValidCurrencyFormat(amount) {
			return nil, errors.New("the fourth element of each inner array must be a valid currency format (amount)")
		}
		porezi[i] = &PorezOstaloType{Naziv: name, Stopa: rate, Osnovica: base, Iznos: amount}
	}
	return &OstaliPoreziType{Porez: porezi}, nil
}

// newPNP initializes and returns a PorezNaPotrosnjuType instance
//
// This function creates a new instance of PorezNaPotrosnjuType, which represents a collection of consumption tax (PorezType) entries.
// It takes a 2D array of values where each inner array represents a single consumption tax entry with the tax rate, tax base, and tax amount.
//
// Parameters:
//
//	values ([][]interface{}): A 2D array where each inner array contains three elements:
//	  - int: The tax rate (Stopa)
//	  - string: The tax base (Osnovica)
//	  - string: The tax amount (Iznos)
//
// Returns:
//
//	(*PorezNaPotrosnjuType, error): A pointer to a new PorezNaPotrosnjuType instance with the provided consumption tax entries, or an error if the input is invalid.
//
// Example:
//
//	values := [][]interface{}{
//	  {3, "1000", "30"},
//	  {5, "2000", "100"},
//	}
//	pnp, err := newPNP(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("PorezNaPotrosnju: %+v\n", pnp)
//	}
func newPNP(values [][]interface{}) (*PorezNaPotrosnjuType, error) {
	porezi := make([]*PorezType, len(values))
	for i, v := range values {
		if len(v) != 3 {
			return nil, errors.New("each inner array must contain exactly three elements")
		}
		taxRate, ok := v[0].(string)
		if !ok {
			return nil, errors.New("the first element of each inner array must be an int (tax rate)")
		}
		taxBase, ok := v[1].(string)
		if !ok {
			return nil, errors.New("the second element of each inner array must be a string (tax base)")
		}
		if !IsValidCurrencyFormat(taxBase) {
			return nil, errors.New("the second element of each inner array must be a valid currency format (tax base)")
		}
		taxAmount, ok := v[2].(string)
		if !ok {
			return nil, errors.New("the third element of each inner array must be a string (tax amount)")
		}
		if !IsValidCurrencyFormat(taxAmount) {
			return nil, errors.New("the third element of each inner array must be a valid currency format (tax amount)")
		}
		porezi[i] = &PorezType{Stopa: taxRate, Osnovica: taxBase, Iznos: taxAmount}
	}
	return &PorezNaPotrosnjuType{Porez: porezi}, nil
}

// newPdv initializes and returns a PdvType instance
//
// This function creates a new instance of PdvType, which represents a collection of VAT (PorezType) entries.
// It takes a 2D array of values where each inner array represents a single VAT entry with the tax rate, tax base, and tax amount.
//
// Parameters:
//
//	values ([][]interface{}): A 2D array where each inner array contains three elements:
//	  - int: The tax rate (Stopa)
//	  - string: The tax base (Osnovica)
//	  - string: The tax amount (Iznos)
//
// Returns:
//
//	(*PdvType, error): A pointer to a new PdvType instance with the provided VAT entries, or an error if the input is invalid.
//
// Example:
//
//	values := [][]interface{}{
//	  {25, "1000", "250"},
//	  {13, "500", "65"},
//	}
//	pdv, err := newPdv(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("Pdv: %+v\n", pdv)
//	}
func newPdv(values [][]interface{}) (*PdvType, error) {
	porezi := make([]*PorezType, len(values))
	for i, v := range values {
		if len(v) != 3 {
			return nil, errors.New("each inner array must contain exactly three elements")
		}
		taxRate, ok := v[0].(string)
		if !ok {
			return nil, errors.New("the first element of each inner array must be an int (tax rate)")
		}
		taxBase, ok := v[1].(string)
		if !ok {
			return nil, errors.New("the second element of each inner array must be a string (tax base)")
		}
		if !IsValidCurrencyFormat(taxBase) {
			return nil, errors.New("the second element of each inner array must be a valid currency format (tax base)")
		}
		taxAmount, ok := v[2].(string)
		if !ok {
			return nil, errors.New("the third element of each inner array must be a string (tax amount)")
		}
		if !IsValidCurrencyFormat(taxAmount) {
			return nil, errors.New("the third element of each inner array must be a valid currency format (tax amount)")
		}
		porezi[i] = &PorezType{Stopa: taxRate, Osnovica: taxBase, Iznos: taxAmount}
	}
	return &PdvType{Porez: porezi}, nil
}
