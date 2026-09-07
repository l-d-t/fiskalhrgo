package fiskalhrgo

import (
	"encoding/xml"
	"strings"
	"testing"
	"time"

	"github.com/beevik/etree"
)

func TestBuyerOIBXML(t *testing.T) {
	entity := newSignatureTestEntity(t)
	for _, payment := range []PaymentMethod{CISCash, CISCard, CISMixOther} {
		t.Run(string(payment), func(t *testing.T) {
			invoice, zki, err := entity.NewCISInvoice(time.Now(), 1, 1, nil, nil, nil,
				"0.00", "0.00", "0.00", nil, "10.00", payment, entity.OIB())
			if err != nil {
				t.Fatal(err)
			}
			marshal := func() []byte {
				t.Helper()
				data, err := xml.Marshal(RacunZahtjev{Xmlns: DefaultNamespace, IdAttr: "b2b-1", Racun: invoice})
				if err != nil {
					t.Fatal(err)
				}
				return data
			}
			before := string(marshal())
			if strings.Contains(before, "OibPrimateljaRacuna") {
				t.Fatal("consumer invoice includes buyer OIB")
			}
			if err := invoice.SetOibPrimateljaRacuna("65049901548"); err != nil {
				t.Fatal(err)
			}
			if invoice.GetZKI() != zki {
				t.Fatal("setting buyer OIB changed ZKI")
			}
			doc := etree.NewDocument()
			if err := doc.ReadFromBytes(marshal()); err != nil {
				t.Fatal(err)
			}
			racun := doc.Root().FindElement("tns:Racun")
			children := racun.ChildElements()
			buyer := children[len(children)-1]
			if buyer.Tag != "OibPrimateljaRacuna" || buyer.NamespaceURI() != DefaultNamespace || buyer.Text() != "65049901548" {
				t.Fatalf("incorrect buyer OIB element or XSD sequence: %v", buyer)
			}
			if err := invoice.SetOibPrimateljaRacuna(""); err != nil {
				t.Fatal(err)
			}
			if string(marshal()) != before || invoice.GetZKI() != zki {
				t.Fatal("clearing buyer OIB did not restore original invoice")
			}
		})
	}
}

func TestBuyerOIBValidation(t *testing.T) {
	for _, oib := range []string{"12345678900", "6504990154", "650499015480", "6504990154x", " 65049901548", "65049901548 "} {
		t.Run(oib, func(t *testing.T) {
			invoice := &RacunType{NacinPlac: "G", OibPrimateljaRacuna: "65049901548"}
			if err := invoice.SetOibPrimateljaRacuna(oib); err == nil {
				t.Fatal("accepted invalid buyer OIB")
			}
			if invoice.OibPrimateljaRacuna != "65049901548" {
				t.Fatal("failed setter changed buyer OIB")
			}
			invoice.OibPrimateljaRacuna = oib
			if _, _, err := invoice.InvoiceRequest(); err == nil || !strings.Contains(err.Error(), "OibPrimateljaRacuna") {
				t.Fatalf("send did not reject directly modified buyer OIB: %v", err)
			}
		})
	}
	for _, payment := range []string{"T", "C", "", "invalid"} {
		invoice := &RacunType{NacinPlac: payment}
		if err := invoice.SetOibPrimateljaRacuna("65049901548"); err == nil {
			t.Fatalf("accepted buyer OIB with payment %q", payment)
		}
		invoice.OibPrimateljaRacuna = "65049901548"
		if _, _, err := invoice.InvoiceRequest(); err == nil || !strings.Contains(err.Error(), "OibPrimateljaRacuna") {
			t.Fatalf("send accepted buyer OIB with payment %q: %v", payment, err)
		}
		if err := invoice.SetOibPrimateljaRacuna(""); err != nil {
			t.Fatal(err)
		}
	}
	var invoice *RacunType
	if err := invoice.SetOibPrimateljaRacuna("65049901548"); err == nil {
		t.Fatal("nil invoice must return an error")
	}
}
