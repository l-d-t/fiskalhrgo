package fiskalhrgo

import (
	"encoding/xml"
	"testing"
)

// Test for RacunZahtjev structure
func TestRacunZahtjevMarshal(t *testing.T) {
	racun := RacunZahtjev{
		Zaglavlje: &ZaglavljeType{
			IdPoruke:     "4ec0a35a-7db7-45c6-9260-7b25b7749d61",
			DatumVrijeme: "2024-09-19T10:00:00",
		},
		Racun: &RacunType{
			Oib:         "12345678901",
			USustPdv:    true,
			DatVrijeme:  "2024-09-19T10:00:00",
			OznSlijed:   "P",
			BrRac:       &BrojRacunaType{100, "POS1", 1},
			IznosUkupno: "150.50",
			Pdv:         &PdvType{[]*PorezType{{25, "120.40", "30.10"}}},
			NacinPlac:   "G",
			OibOper:     "98765432100",
			ZastKod:     "c3b2ecf807f56e294fbb3d536aad0f6c",
			NakDost:     false,
		},
		IdAttr: "646332",
	}

	// Marshal the RacunZahtjev to XML
	xmlData, err := xml.MarshalIndent(racun, "", "  ")
	if err != nil {
		t.Fatalf("Error marshalling RacunZahtjev: %v", err)
	}

	t.Log(string(xmlData))
}

// Test for RacunOdgovor
func TestRacunOdgovorUnmarshal(t *testing.T) {
	xmlData := `<RacunOdgovor Id="646332">
  <Zaglavlje>
    <IdPoruke>4ec0a35a-7db7-45c6-9260-7b25b7749d61</IdPoruke>
    <DatumVrijeme>2024-09-19T10:00:00</DatumVrijeme>
  </Zaglavlje>
  <Jir>12345678-1234-1234-1234-123456789012</Jir>
</RacunOdgovor>`

	var odgovor RacunOdgovor
	err := xml.Unmarshal([]byte(xmlData), &odgovor)
	if err != nil {
		t.Fatalf("Error unmarshalling RacunOdgovor: %v", err)
	}

	if odgovor.Zaglavlje.IdPoruke != "4ec0a35a-7db7-45c6-9260-7b25b7749d61" {
		t.Fatalf("Expected IdPoruke: 4ec0a35a-7db7-45c6-9260-7b25b7749d61, but got: %s", odgovor.Zaglavlje.IdPoruke)
	}

	if odgovor.Jir != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("Expected JIR: 12345678-1234-1234-1234-123456789012, but got: %s", odgovor.Jir)
	}

}
