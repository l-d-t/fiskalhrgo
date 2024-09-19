package fiskalhrgo

import "encoding/xml"

// RacunZahtjev ...
type RacunZahtjev struct {
	IdAttr      string         `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeType `xml:"Zaglavlje"`
	Racun       *RacunType     `xml:"Racun"`
	DsSignature *SignatureType `xml:"ds:Signature"`
}

// RacunOdgovor ...
type RacunOdgovor struct {
	IdAttr      string                `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir         string                `xml:"Jir"`
	Greske      *GreskeType           `xml:"Greske"`
	DsSignature *SignatureType        `xml:"ds:Signature"`
}

// PrateciDokumentiZahtjev ...
type PrateciDokumentiZahtjev struct {
	IdAttr          string               `xml:"Id,attr,omitempty"`
	Zaglavlje       *ZaglavljeType       `xml:"Zaglavlje"`
	PrateciDokument *PrateciDokumentType `xml:"PrateciDokument"`
	DsSignature     *SignatureType       `xml:"ds:Signature"`
}

// PrateciDokumentiOdgovor ...
type PrateciDokumentiOdgovor struct {
	IdAttr      string                `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir         string                `xml:"Jir"`
	Greske      *GreskeType           `xml:"Greske"`
	DsSignature *SignatureType        `xml:"ds:Signature"`
}

// RacunPDZahtjev ...
type RacunPDZahtjev struct {
	IdAttr      string         `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeType `xml:"Zaglavlje"`
	Racun       *RacunType     `xml:"Racun"`
	DsSignature *SignatureType `xml:"ds:Signature"`
}

// RacunPDOdgovor ...
type RacunPDOdgovor struct {
	IdAttr      string                `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir         string                `xml:"Jir"`
	Greske      *GreskeType           `xml:"Greske"`
	DsSignature *SignatureType        `xml:"ds:Signature"`
}

// PromijeniNacPlacZahtjev ...
type PromijeniNacPlacZahtjev struct {
	IdAttr      string         `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeType `xml:"Zaglavlje"`
	Racun       *RacunType     `xml:"Racun"`
	DsSignature *SignatureType `xml:"ds:Signature"`
}

// PromijeniNacPlacOdgovor ...
type PromijeniNacPlacOdgovor struct {
	IdAttr         string                `xml:"Id,attr,omitempty"`
	Zaglavlje      *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	PorukaOdgovora *PorukaOdgovoraType   `xml:"PorukaOdgovora"`
	Greske         *GreskeType           `xml:"Greske"`
	DsSignature    *SignatureType        `xml:"ds:Signature"`
}

// NapojnicaZahtjev ...
type NapojnicaZahtjev struct {
	IdAttr      string         `xml:"Id,attr,omitempty"`
	Zaglavlje   *ZaglavljeType `xml:"Zaglavlje"`
	Racun       *RacunType     `xml:"Racun"`
	DsSignature *SignatureType `xml:"ds:Signature"`
}

// NapojnicaOdgovor ...
type NapojnicaOdgovor struct {
	IdAttr         string                `xml:"Id,attr,omitempty"`
	Zaglavlje      *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	PorukaOdgovora *PorukaOdgovoraType   `xml:"PorukaOdgovora"`
	Greske         *GreskeType           `xml:"Greske"`
	DsSignature    *SignatureType        `xml:"ds:Signature"`
}

// EchoRequest represents a simple request with a text body
type EchoRequest struct {
	XMLName xml.Name `xml:"tns:EchoRequest"`
	Text    string   `xml:",chardata"`
}

// EchoResponse represents a simple response with a text body
type EchoResponse struct {
	XMLName xml.Name `xml:"tns:EchoResponse"`
	Text    string   `xml:",chardata"`
}

// PorukaOdgovoraType ...
type PorukaOdgovoraType struct {
	SifraPoruke string `xml:"SifraPoruke"`
	Poruka      string `xml:"Poruka"`
}

// ZaglavljeType is Datum i vrijeme slanja poruke.
type ZaglavljeType struct {
	IdPoruke     string `xml:"IdPoruke"`
	DatumVrijeme string `xml:"DatumVrijeme"`
}

// ZaglavljeOdgovorType ...
type ZaglavljeOdgovorType struct {
	IdPoruke     string `xml:"IdPoruke"`
	DatumVrijeme string `xml:"DatumVrijeme"`
}

// RacunType represents the invoice type with various details required for fiscalization.
type RacunType struct {
	Oib                   string                `xml:"Oib"`                             // Required
	USustPdv              bool                  `xml:"USustPdv"`                        // Required, Boolean
	DatVrijeme            string                `xml:"DatVrijeme"`                      // Required, DateTime in format dd.mm.yyyyThh:mm:ss
	OznSlijed             string                `xml:"OznSlijed"`                       // Required, 'P' or 'N' (Char(1))
	BrRac                 *BrojRacunaType       `xml:"BrRac"`                           // Required, custom format for "number/business unit/pos device"
	Pdv                   *PdvType              `xml:"Pdv,omitempty"`                   // Optional, list of VAT (PDV) taxes
	Pnp                   *PorezNaPotrosnjuType `xml:"Pnp,omitempty"`                   // Optional, tax on consumption (Porez na potrošnju)
	OstaliPor             *OstaliPoreziType     `xml:"OstaliPor,omitempty"`             // Optional, list of other taxes
	IznosOslobPdv         string                `xml:"IznosOslobPdv,omitempty"`         // Optional, amount of VAT exemption
	IznosMarza            string                `xml:"IznosMarza,omitempty"`            // Optional, amount subject to special margin taxation
	IznosNePodlOpor       string                `xml:"IznosNePodlOpor,omitempty"`       // Optional, amount not subject to tax
	Naknade               *NaknadeType          `xml:"Naknade,omitempty"`               // Optional, list of fees (e.g., packaging fee)
	IznosUkupno           string                `xml:"IznosUkupno"`                     // Required, total amount (Decimal)
	NacinPlac             string                `xml:"NacinPlac"`                       // Required, 'G' (cash), 'K' (card), 'C' (check), 'T' (transaction), 'O' (other) (Char(1))
	OibOper               string                `xml:"OibOper"`                         // Required, OIB of the operator issuing the invoice
	ZastKod               string                `xml:"ZastKod"`                         // Required, 32-character alphanumeric protection code
	NakDost               bool                  `xml:"NakDost"`                         // Required, whether the invoice is delivered later (Boolean)
	ParagonBrRac          string                `xml:"ParagonBrRac,omitempty"`          // Optional, number of a paragon invoice (in case of full device failure)
	SpecNamj              string                `xml:"SpecNamj,omitempty"`              // Optional, specific purpose for additional data
	PrateciDokument       *PrateciDokument      `xml:"PrateciDokument,omitempty"`       // Optional, additional document (prateci dokument)
	PromijenjeniNacinPlac string                `xml:"PromijenjeniNacinPlac,omitempty"` // Optional, changed payment method
	Napojnica             *NapojnicaType        `xml:"Napojnica,omitempty"`             // Optional, tip information
}

// PrateciDokumentType ...
type PrateciDokumentType struct {
	Oib                 string      `xml:"Oib"`
	DatVrijeme          string      `xml:"DatVrijeme"`
	BrPratecegDokumenta *BrojPDType `xml:"BrPratecegDokumenta"`
	IznosUkupno         string      `xml:"IznosUkupno"`
	ZastKodPD           string      `xml:"ZastKodPD"`
	NakDost             bool        `xml:"NakDost"`
}

// PrateciDokument ...
type PrateciDokument struct {
	JirPD     []string `xml:"JirPD"`
	ZastKodPD []string `xml:"ZastKodPD"`
}

// NapojnicaType ...
type NapojnicaType struct {
	IznosNapojnice         string `xml:"iznosNapojnice"`
	NacinPlacanjaNapojnice string `xml:"nacinPlacanjaNapojnice"`
}

// GreskeType ...
type GreskeType struct {
	Greska []*GreskaType `xml:"Greska"`
}

// GreskaType ...
type GreskaType struct {
	SifraGreske  string `xml:"SifraGreske"`
	PorukaGreske string `xml:"PorukaGreske"`
}

// NaknadeType ...
type NaknadeType struct {
	Naknada []*NaknadaType `xml:"Naknada"`
}

// NaknadaType ...
type NaknadaType struct {
	NazivN string `xml:"NazivN"`
	IznosN string `xml:"IznosN"`
}

// OstaliPoreziType ...
type OstaliPoreziType struct {
	Porez []*PorezOstaloType `xml:"Porez"`
}

// PorezNaPotrosnjuType ...
type PorezNaPotrosnjuType struct {
	Porez []*PorezType `xml:"Porez"`
}

// PdvType ...
type PdvType struct {
	Porez []*PorezType `xml:"Porez"`
}

// PorezOstaloType ...
type PorezOstaloType struct {
	Naziv    string `xml:"Naziv"`
	Stopa    int    `xml:"Stopa"`
	Osnovica string `xml:"Osnovica"`
	Iznos    string `xml:"Iznos"`
}

// PorezType ...
type PorezType struct {
	Stopa    int    `xml:"Stopa"`
	Osnovica string `xml:"Osnovica"`
	Iznos    string `xml:"Iznos"`
}

// BrojRacunaType ...
type BrojRacunaType struct {
	BrOznRac int    `xml:"BrOznRac"`
	OznPosPr string `xml:"OznPosPr"`
	OznNapUr int    `xml:"OznNapUr"`
}

// BrojPDType ...
type BrojPDType struct {
	BrOznPD  int    `xml:"BrOznPD"`
	OznPosPr string `xml:"OznPosPr"`
	OznNapUr int    `xml:"OznNapUr"`
}
