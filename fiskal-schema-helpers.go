package fiskalhrgo

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// NewCISInvoice initializes and returns a RacunType instance
//
// This method creates a new instance of RacunType, which represents an invoice with all necessary fields.
//
// Parameters:
//
//	dateTime (time.Time): The date and time of the invoice.
//	centralized (bool): Indicates whether the sequence mark is centralized.
//	invoiceNumber (uint): The unique number of the invoice.
//	locationIdentifier (string): The identifier for the business location where the invoice was issued.
//	registerDeviceID (uint): The identifier for the cash register device used to issue the invoice.
//	pdvValues ([][]interface{}): A 2D array for VAT details (nullable).
//	pnpValues ([][]interface{}): A 2D array for consumption tax details (nullable).
//	ostaliPorValues ([][]interface{}): A 2D array for other tax details (nullable).
//	iznosOslobPdv (string): The amount exempt from VAT (optional).
//	iznosMarza (string): The margin amount (optional).
//	iznosNePodlOpor (string): The amount not subject to taxation (optional).
//	naknadeValues ([][]string): A 2D array for fees details (nullable).
//	iznosUkupno (string): The total amount.
//	nacinPlac (string): The payment method.
//	oibOper (string): The OIB of the operator.
//	nakDost (bool): Indicates whether the invoice is delivered.
//	paragonBrRac (string): The paragon invoice number (optional).
//	specNamj (string): Special purpose (optional).
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
	nacinPlac string,
	oibOper string,
	nakDost bool,
	paragonBrRac string,
	specNamj string,
) (*RacunType, string, error) {
	// Format the date and time
	formattedDate := dateTime.Format("2006-01-02T15:04:05")

	// Determine the sequence mark
	oznSlijed := "N"
	if fe.centralizedInvoiceNumber {
		oznSlijed = "P"
	}

	// Use helper functions to create the necessary types
	var pdv *PdvType
	var err error
	if pdvValues != nil {
		pdv, err = NewPdv(pdvValues)
		if err != nil {
			return nil, "", err
		}
	}

	var pnp *PorezNaPotrosnjuType
	if pnpValues != nil {
		pnp, err = NewPNP(pnpValues)
		if err != nil {
			return nil, "", err
		}
	}

	var ostaliPor *OstaliPoreziType
	if ostaliPorValues != nil {
		ostaliPor, err = OtherTaxes(ostaliPorValues)
		if err != nil {
			return nil, "", err
		}
	}

	var naknade *NaknadeType
	if naknadeValues != nil {
		naknade, err = Naknade(naknadeValues)
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
	if nacinPlac != "G" && nacinPlac != "K" && nacinPlac != "O" && nacinPlac != "T" && nacinPlac != "C" {
		return nil, "", errors.New("NacinPlac must be one of the following values: G, K, O, T, C (deprecated)")
	}

	zki, err := fe.GenerateZKI(dateTime, invoiceNumber, registerDeviceID, iznosUkupno)

	if err != nil {
		return nil, "", err
	}

	return &RacunType{
		Oib:             fe.oib,
		USustPdv:        fe.sustPDV,
		DatVrijeme:      formattedDate,
		OznSlijed:       oznSlijed,
		BrRac:           brRac,
		Pdv:             pdv,
		Pnp:             pnp,
		OstaliPor:       ostaliPor,
		IznosOslobPdv:   iznosOslobPdv,
		IznosMarza:      iznosMarza,
		IznosNePodlOpor: iznosNePodlOpor,
		Naknade:         naknade,
		IznosUkupno:     iznosUkupno,
		NacinPlac:       nacinPlac,
		OibOper:         oibOper,
		ZastKod:         zki,
		NakDost:         nakDost,
		ParagonBrRac:    paragonBrRac,
		SpecNamj:        specNamj,
	}, zki, nil
}

// NewFiskalHeader creates a new instance of ZaglavljeType with a unique message ID and the current timestamp
//
// This function generates a new UUIDv4 for the IdPoruke field to ensure that each message has a unique identifier.
// It also sets the DatumVrijeme field to the current time formatted as "2006-01-02T15:04:05" to indicate when the message was created.
//
// Returns:
//
//	*ZaglavljeType: A pointer to a new ZaglavljeType instance with the IdPoruke and DatumVrijeme fields populated.
func NewFiskalHeader() *ZaglavljeType {
	return &ZaglavljeType{
		IdPoruke:     uuid.New().String(),
		DatumVrijeme: time.Now().Format("2006-01-02T15:04:05"),
	}
}

// Naknade initializes and returns a NaknadeType instance
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
//	naknade, err := Naknade(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("Naknade: %+v\n", naknade)
//	}
func Naknade(values [][]string) (*NaknadeType, error) {
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

// OtherTaxes initializes and returns an OstaliPoreziType instance
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
//	ostaliPorezi, err := OtherTaxes(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("OstaliPorezi: %+v\n", ostaliPorezi)
//	}
func OtherTaxes(values [][]interface{}) (*OstaliPoreziType, error) {
	porezi := make([]*PorezOstaloType, len(values))
	for i, v := range values {
		if len(v) != 4 {
			return nil, errors.New("each inner array must contain exactly four elements")
		}
		name, ok := v[0].(string)
		if !ok {
			return nil, errors.New("the first element of each inner array must be a string (name)")
		}
		rate, ok := v[1].(int)
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

// NewPNP initializes and returns a PorezNaPotrosnjuType instance
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
//	pnp, err := NewPNP(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("PorezNaPotrosnju: %+v\n", pnp)
//	}
func NewPNP(values [][]interface{}) (*PorezNaPotrosnjuType, error) {
	porezi := make([]*PorezType, len(values))
	for i, v := range values {
		if len(v) != 3 {
			return nil, errors.New("each inner array must contain exactly three elements")
		}
		taxRate, ok := v[0].(int)
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

// NewPdv initializes and returns a PdvType instance
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
//	pdv, err := NewPdv(values)
//	if err != nil {
//	  fmt.Printf("Error: %v\n", err)
//	} else {
//	  fmt.Printf("Pdv: %+v\n", pdv)
//	}
func NewPdv(values [][]interface{}) (*PdvType, error) {
	porezi := make([]*PorezType, len(values))
	for i, v := range values {
		if len(v) != 3 {
			return nil, errors.New("each inner array must contain exactly three elements")
		}
		taxRate, ok := v[0].(int)
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

// NewInvoiceNumber initializes and returns a BrojRacunaType instance
//
// This function creates a new instance of BrojRacunaType, which represents the structure for an invoice number.
// It takes three parameters: the invoice number, the location identifier, and the register device ID.
//
// Parameters:
//
//	InvoiceNumber (int): The unique number of the invoice.
//	LocationIdentifier (string): The identifier for the business location where the invoice was issued.
//	RegisterDeviceID (int): The identifier for the cash register device used to issue the invoice.
//
// Returns:
//
//	*BrojRacunaType: A pointer to a new BrojRacunaType instance with the provided values.
func NewInvoiceNumber(InvoiceNumber uint, LocationIdentifier string, RegisterDeviceID uint) *BrojRacunaType {
	return &BrojRacunaType{
		BrOznRac: InvoiceNumber,
		OznPosPr: LocationIdentifier,
		OznNapUr: RegisterDeviceID,
	}
}
