package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const DefaultNamespace = "http://www.apis-it.hr/fin/2012/types/f73"

// RacunZahtjev ...
type RacunZahtjev struct {
	XMLName   xml.Name       `xml:"tns:RacunZahtjev"`
	Xmlns     string         `xml:"xmlns:tns,attr"` // Declare the tns namespace
	IdAttr    string         `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeType `xml:"tns:Zaglavlje"`
	Racun     *RacunType     `xml:"tns:Racun"`
}

// RacunOdgovor ...
type RacunOdgovor struct {
	XMLName   xml.Name              `xml:"RacunOdgovor"`
	IdAttr    string                `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir       string                `xml:"Jir"`
	Greske    *GreskeType           `xml:"Greske"`
}

// PrateciDokumentiZahtjev ...
type PrateciDokumentiZahtjev struct {
	XMLName         xml.Name             `xml:"tns:PrateciDokumentiZahtjev"`
	Xmlns           string               `xml:"xmlns:tns,attr"` // Declare the tns namespace
	IdAttr          string               `xml:"Id,attr,omitempty"`
	Zaglavlje       *ZaglavljeType       `xml:"tns:Zaglavlje"`
	PrateciDokument *PrateciDokumentType `xml:"tns:PrateciDokument"`
}

// PrateciDokumentiOdgovor ...
type PrateciDokumentiOdgovor struct {
	XMLName   xml.Name              `xml:"PrateciDokumentiOdgovor"`
	IdAttr    string                `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir       string                `xml:"Jir"`
	Greske    *GreskeType           `xml:"Greske"`
}

// RacunPDZahtjev ...
type RacunPDZahtjev struct {
	XMLName   xml.Name       `xml:"tns:RacunPDZahtjev"`
	Xmlns     string         `xml:"xmlns:tns,attr"` // Declare the tns namespace
	IdAttr    string         `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeType `xml:"tns:Zaglavlje"`
	Racun     *RacunType     `xml:"tns:Racun"`
}

// RacunPDOdgovor ...
type RacunPDOdgovor struct {
	XMLName   xml.Name              `xml:"RacunPDOdgovor"`
	IdAttr    string                `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	Jir       string                `xml:"Jir"`
	Greske    *GreskeType           `xml:"Greske"`
}

// PromijeniNacPlacZahtjev ...
type PromijeniNacPlacZahtjev struct {
	XMLName   xml.Name       `xml:"tns:PromijeniNacPlacZahtjev"`
	Xmlns     string         `xml:"xmlns:tns,attr"` // Declare the tns namespace
	IdAttr    string         `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeType `xml:"tns:Zaglavlje"`
	Racun     *RacunType     `xml:"tns:Racun"`
}

// PromijeniNacPlacOdgovor ...
type PromijeniNacPlacOdgovor struct {
	XMLName        xml.Name              `xml:"PromijeniNacPlacOdgovor"`
	IdAttr         string                `xml:"Id,attr,omitempty"`
	Zaglavlje      *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	PorukaOdgovora *PorukaOdgovoraType   `xml:"PorukaOdgovora"`
	Greske         *GreskeType           `xml:"Greske"`
}

// NapojnicaZahtjev ...
type NapojnicaZahtjev struct {
	XMLName   xml.Name       `xml:"tns:NapojnicaZahtjev"`
	Xmlns     string         `xml:"xmlns:tns,attr"` // Declare the tns namespace
	IdAttr    string         `xml:"Id,attr,omitempty"`
	Zaglavlje *ZaglavljeType `xml:"tns:Zaglavlje"`
	Racun     *RacunType     `xml:"tns:Racun"`
}

// NapojnicaOdgovor ...
type NapojnicaOdgovor struct {
	XMLName        xml.Name              `xml:"NapojnicaOdgovor"`
	IdAttr         string                `xml:"Id,attr,omitempty"`
	Zaglavlje      *ZaglavljeOdgovorType `xml:"Zaglavlje"`
	PorukaOdgovora *PorukaOdgovoraType   `xml:"PorukaOdgovora"`
	Greske         *GreskeType           `xml:"Greske"`
}

// EchoRequest represents a simple request with a text body
type EchoRequest struct {
	XMLName xml.Name `xml:"tns:EchoRequest"`
	Xmlns   string   `xml:"xmlns:tns,attr"` // Declare the tns namespace
	Text    string   `xml:",chardata"`
}

// EchoResponse represents a simple response with a text body
type EchoResponse struct {
	XMLName xml.Name `xml:"EchoResponse"`
	Text    string   `xml:",chardata"`
}

// PorukaOdgovoraType ...
type PorukaOdgovoraType struct {
	SifraPoruke string `xml:"SifraPoruke"`
	Poruka      string `xml:"Poruka"`
}

// ZaglavljeType is Datum i vrijeme slanja poruke.
type ZaglavljeType struct {
	IdPoruke     string `xml:"tns:IdPoruke"`
	DatumVrijeme string `xml:"tns:DatumVrijeme"`
}

// ZaglavljeOdgovorType ...
type ZaglavljeOdgovorType struct {
	IdPoruke     string `xml:"IdPoruke"`
	DatumVrijeme string `xml:"DatumVrijeme"`
}

// RacunType represents the invoice type with various details required for fiscalization.
type RacunType struct {
	XMLName               xml.Name              `xml:"tns:Racun"`
	Oib                   string                `xml:"tns:Oib"`
	USustPdv              bool                  `xml:"tns:USustPdv"`
	DatVrijeme            string                `xml:"tns:DatVrijeme"`
	OznSlijed             string                `xml:"tns:OznSlijed"`
	BrRac                 *BrojRacunaType       `xml:"tns:BrRac"`
	Pdv                   *PdvType              `xml:"tns:Pdv,omitempty"`
	Pnp                   *PorezNaPotrosnjuType `xml:"tns:Pnp,omitempty"`
	OstaliPor             *OstaliPoreziType     `xml:"tns:OstaliPor,omitempty"`
	IznosOslobPdv         string                `xml:"tns:IznosOslobPdv,omitempty"`
	IznosMarza            string                `xml:"tns:IznosMarza,omitempty"`
	IznosNePodlOpor       string                `xml:"tns:IznosNePodlOpor,omitempty"`
	Naknade               *NaknadeType          `xml:"tns:Naknade,omitempty"`
	IznosUkupno           string                `xml:"tns:IznosUkupno"`
	NacinPlac             string                `xml:"tns:NacinPlac"`
	OibOper               string                `xml:"tns:OibOper"`
	ZastKod               string                `xml:"tns:ZastKod"`
	NakDost               bool                  `xml:"tns:NakDost"`
	ParagonBrRac          string                `xml:"tns:ParagonBrRac,omitempty"`
	SpecNamj              string                `xml:"tns:SpecNamj,omitempty"`
	PrateciDokument       *PrateciDokument      `xml:"tns:PrateciDokument,omitempty"`
	PromijenjeniNacinPlac string                `xml:"tns:PromijenjeniNacinPlac,omitempty"`
	Napojnica             *NapojnicaType        `xml:"tns:Napojnica,omitempty"`

	// Additional functional non XML fields
	pointerToEntity    *FiskalEntity // Pointer to the FiskalEntity
	oldEntityForOldZKI *FiskalEntity // Pointer to the old FiskalEntity for the old ZKI
	// This is used in the edge case that the ZKI was generated with one certificate and the fiscalization failed
	// But the certificate expired or had to be changed and now fiscalization have to be repeated with new certificate
	// If we replace the original ZKI its a problem we already gave the invoice with old ZKI out
	// So we have to keep the old ZKI and validate it with the old certificate before signing and sending with new one
	// In any case this is set by IhaveZKIwithExpiredCertificateEdgeCase(EntityWithOldCertLoaded *FiskalEntity) method
}

// PrateciDokumentType ...
type PrateciDokumentType struct {
	Oib                 string      `xml:"tns:Oib"`
	DatVrijeme          string      `xml:"tns:DatVrijeme"`
	BrPratecegDokumenta *BrojPDType `xml:"tns:BrPratecegDokumenta"`
	IznosUkupno         string      `xml:"tns:IznosUkupno"`
	ZastKodPD           string      `xml:"tns:ZastKodPD"`
	NakDost             bool        `xml:"tns:NakDost"`
}

// PrateciDokument ...
type PrateciDokument struct {
	JirPD     string `xml:"tns:JirPD"`
	ZastKodPD string `xml:"tns:ZastKodPD"`
}

// NapojnicaType ...
type NapojnicaType struct {
	IznosNapojnice         string `xml:"tns:iznosNapojnice"`
	NacinPlacanjaNapojnice string `xml:"tns:nacinPlacanjaNapojnice"`
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
	Naknada []*NaknadaType `xml:"tns:Naknada"`
}

// NaknadaType ...
type NaknadaType struct {
	NazivN string `xml:"tns:NazivN"`
	IznosN string `xml:"tns:IznosN"`
}

// OstaliPoreziType ...
type OstaliPoreziType struct {
	Porez []*PorezOstaloType `xml:"tns:Porez"`
}

// PorezNaPotrosnjuType ...
type PorezNaPotrosnjuType struct {
	Porez []*PorezType `xml:"tns:Porez"`
}

// PdvType ...
type PdvType struct {
	Porez []*PorezType `xml:"tns:Porez"`
}

// PorezOstaloType ...
type PorezOstaloType struct {
	Naziv    string `xml:"tns:Naziv"`
	Stopa    string `xml:"tns:Stopa"`
	Osnovica string `xml:"tns:Osnovica"`
	Iznos    string `xml:"tns:Iznos"`
}

// PorezType ...
type PorezType struct {
	Stopa    string `xml:"tns:Stopa"`
	Osnovica string `xml:"tns:Osnovica"`
	Iznos    string `xml:"tns:Iznos"`
}

// BrojRacunaType ...
type BrojRacunaType struct {
	BrOznRac uint   `xml:"tns:BrOznRac"`
	OznPosPr string `xml:"tns:OznPosPr"`
	OznNapUr uint   `xml:"tns:OznNapUr"`
}

// BrojPDType ...
type BrojPDType struct {
	BrOznPD  int    `xml:"tns:BrOznPD"`
	OznPosPr string `xml:"tns:OznPosPr"`
	OznNapUr int    `xml:"tns:OznNapUr"`
}

// generateUniqueID generates a unique ID
func generateUniqueID() string {
	return fmt.Sprintf("%x", time.Now().UnixNano())
}

// newFiskalHeader creates a new instance of ZaglavljeType with a unique message ID and the current timestamp
//
// This function generates a new UUIDv4 for the IdPoruke field to ensure that each message has a unique identifier.
// It also sets the DatumVrijeme field to the current time formatted as "2006-01-02T15:04:05" to indicate when the message was created.
//
// Returns:
//
//	*ZaglavljeType: A pointer to a new ZaglavljeType instance with the IdPoruke and DatumVrijeme fields populated.
func newFiskalHeader() *ZaglavljeType {
	return &ZaglavljeType{
		IdPoruke:     uuid.New().String(),
		DatumVrijeme: time.Now().Format("02.01.2006T15:04:05"),
	}
}
