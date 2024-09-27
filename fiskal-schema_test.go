package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import (
	"encoding/xml"
	"testing"
)

func TestEchoRequestMarshal(t *testing.T) {
	echoReq := EchoRequest{
		Xmlns: DefaultNamespace,
		Text:  "Hello, world!",
	}

	// Marshal the EchoRequest to XML
	xmlData, err := xml.MarshalIndent(echoReq, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling EchoRequest: %v", err)
	}

	t.Log(string(xmlData))
}

// Test for RacunZahtjev structure
func TestRacunZahtjevMarshal(t *testing.T) {
	racun := RacunZahtjev{
		Zaglavlje: newFiskalHeader(),
		Racun: &RacunType{
			Oib:         "12345678901",
			USustPdv:    true,
			DatVrijeme:  "2024-09-19T10:00:00",
			OznSlijed:   "P",
			BrRac:       &BrojRacunaType{100, "POS1", 1},
			IznosUkupno: "150.50",
			Pdv:         &PdvType{[]*PorezType{{"25.00", "120.40", "30.10"}}},
			NacinPlac:   "G",
			OibOper:     "98765432100",
			ZastKod:     "c3b2ecf807f56e294fbb3d536aad0f6c",
			NakDost:     false,
		},
		Xmlns:  DefaultNamespace,
		IdAttr: "646332",
	}

	// Marshal the RacunZahtjev to XML
	xmlData, err := xml.MarshalIndent(racun, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling RacunZahtjev: %v", err)
	}

	t.Log(string(xmlData))
}
