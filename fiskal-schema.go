package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import "encoding/xml"

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
	JirPD     []string `xml:"tns:JirPD"`
	ZastKodPD []string `xml:"tns:ZastKodPD"`
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
