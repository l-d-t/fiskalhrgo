package fiskalhrgo

import (
	"encoding/xml"
	"fmt"
	"os"
	"testing"
	"time"

	"math/rand"
)

var testEntity *FiskalEntity

// TestMain is run before any other tests. It sets up the shared instances and read env variables.
func TestMain(m *testing.M) {

	fmt.Println(`
___________.__        __           .__    ___ _____________    ________        
\_   _____/|__| _____|  | _______  |  |  /   |   \______   \  /  _____/  ____  
 |    __)  |  |/  ___/  |/ /\__  \ |  | /    ~    \       _/ /   \  ___ /  _ \ 
 |     \   |  |\___ \|    <  / __ \|  |_\    Y    /    |   \ \    \_\  (  <_> )
 \___  /   |__/____  >__|_ \(____  /____/\___|_  /|____|_  /  \______  /\____/ 
     \/            \/     \/     \/            \/        \/          \/        
	`)

	fmt.Println("Setting up...")

	certPath := os.Getenv("FISKALHRGO_TEST_CERT_PATH")
	certPassword := os.Getenv("FISKALHRGO_TEST_CERT_PASSWORD")
	testOIB := os.Getenv("FISKALHRGO_TEST_CERT_OIB")

	if certPath == "" || certPassword == "" || testOIB == "" {
		fmt.Println("FISKALHRGO_TEST_CERT_PATH or FISKALHRGO_TEST_CERT_PASSWORD or FISKALHRGO_TEST_CERT_OIB environment variables are not set")
		os.Exit(1)
	}

	fmt.Printf("Using certificate: %s\n", certPath)
	fmt.Printf("Test OIB: %s\n", testOIB)

	var err error
	testEntity, err = NewFiskalEntity(testOIB, true, "TEST3", true, true, true, certPath, certPassword)
	if err != nil {
		fmt.Printf("Failed to create FiskalEntity: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	// Run tests
	code := m.Run()

	// Exit with the code returned by m.Run()
	os.Exit(code)
}

func TestCertOutput(t *testing.T) {
	t.Logf("Testing certificate output...")

	fmt.Println(testEntity.DisplayCertInfoText())

	testEntity.DisplayCertInfoKeyPoints()
	testEntity.DisplayCertInfoMarkdown()
	testEntity.DisplayCertInfoHTML()

	if testEntity.IsExpiringSoon() {
		fmt.Println("WARNING: Certificate is expiring soon!")
	}

	fmt.Printf("Test certificate expires in %d days\n", testEntity.DaysUntilExpire())

}

// TestGenerateZKI tests the ZKI generation using the previously loaded certificate
func TestGenerateZKI(t *testing.T) {
	t.Logf("Testing ZKI generation...")

	// Reuse the loaded certManager to generate ZKI
	zki, err := testEntity.GenerateZKI(time.Now(), 1, 1, "100.00")
	if err != nil {
		t.Fatalf("Failed to generate ZKI: %v", err)
	}

	if zki == "" {
		t.Fatalf("Expected non-empty ZKI, but got an empty string")
	}

	t.Logf("Generated ZKI: %s", zki)
}

// Note: The ZKI (Protection Code) is dependent on the private key used during its generation.
// Any change in the private key will result in a different ZKI, even with the same input data.
// The ZKI generated in this test is based on a specific known certificate and corresponding data.
//
// For external contributors: It is normal for this test to fail if you do not have the required test certificate.
// You can specify a different ZKI to check against using the FISKALHRGO_TEST_KNOWN_ZKI environment variable.
// Alternatively, you can allow this test to fail locally. This test should pass successfully in the official
// GitHub repository and CI environment, where the correct test certificate is available.

func TestKnownZKI(t *testing.T) {
	t.Logf("Testing Known LDT ZKI generation...")

	expectedZKI := "0b173c6127809d4f0fff53e13222c819" // This is the ZKI for the test certificate

	// Check if the external contributor has set their own known ZKI via environment variable
	if envZKI := os.Getenv("FISKALHRGO_TEST_KNOWN_ZKI"); envZKI != "" {
		expectedZKI = envZKI
	}

	timeString := "17.05.2024 16:00:38"
	// Define the layout (format) for parsing
	layout := "02.01.2006 15:04:05"

	// Parse the time string into a time.Time object
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		t.Fatalf("Error parsing time: %v", err)
	}

	// Reuse the loaded certManager to generate ZKI
	zki, err := testEntity.GenerateZKI(parsedTime, 13, 1, "90.00")
	if err != nil {
		t.Fatalf("Failed to generate ZKI: %v", err)
	}

	if zki == "" {
		t.Fatalf("Expected non-empty ZKI, but got an empty string")
	}

	if zki != expectedZKI {
		t.Error("Note:")
		t.Error("- The ZKI (Protection Code) is dependent on the private key used.")
		t.Error("- If you're an external contributor and don't have the correct test certificate, this test is expected to fail.")
		t.Error("- You can use the FISKALHRGO_TEST_KNOWN_ZKI environment variable to specify your own expected ZKI for testing.")
		t.Error("- This test should pass on the official GitHub and CI setup where the correct certificate is available.")
		t.Fatalf("ERROR - Expected ZKI: %s, got %s", expectedZKI, zki)
	} else {
		t.Log("ZKI matched the expected value.")
	}
}

// Test CISEcho
func TestCISEcho(t *testing.T) {
	t.Logf("Testing CISEcho...")
	msg := "Hello, CIS, from FiskalhrGo!"

	t.Logf("Sending message to CIS: %s", msg)

	// Reuse the loaded certManager to generate ZKI
	resp, err := testEntity.EchoRequest(msg)
	if err != nil {
		t.Fatalf("Failed to make CISEcho request: %v", err)
	}

	t.Logf("CISEcho Response: %v", resp)

	if resp != msg {
		t.Fatalf("Expected the sent message returned!")
	}
}

func TestPing(t *testing.T) {
	t.Log("Testing Ping...")
	err := testEntity.PingCIS()
	if err != nil {
		t.Fatalf("Failed to make Ping request: %v", err)
	}
	t.Log("Ping OK!")
}

// Test CIS invoice with helper functions
func TestNewCISInvoice(t *testing.T) {
	pdvValues := [][]interface{}{
		{"25.00", "1000.00", "250.00"},
	}

	pnpValues := [][]interface{}{
		{"3.00", "1000.00", "30.00"},
	}

	ostaliPorValues := [][]interface{}{
		{"Other Tax", "5.00", "1000.00", "50.00"},
	}

	naknadeValues := [][]string{
		{"Povratna", "0.50"},
	}

	dateTime := time.Now()
	brOznRac := uint(rand.Intn(6901) + 100)
	oznNapUr := uint(1)
	iznosUkupno := "1330.50"
	nacinPlac := "G"
	oibOper := "12345678901"
	nakDost := true
	paragonBrRac := "12345"
	specNamj := "Special Purpose"

	invoice, zki, err := testEntity.NewCISInvoice(
		dateTime,
		brOznRac,
		oznNapUr,
		pdvValues,
		pnpValues,
		ostaliPorValues,
		"0.00",
		"0.00",
		"0.00",
		naknadeValues,
		iznosUkupno,
		nacinPlac,
		oibOper,
		nakDost,
		paragonBrRac,
		specNamj,
	)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if invoice.Oib != testEntity.oib {
		t.Errorf("Expected Oib %v, got %v", testEntity.oib, invoice.Oib)
	}

	if invoice.USustPdv != true {
		t.Errorf("Expected USustPdv true, got %v", invoice.USustPdv)
	}

	if invoice.DatVrijeme != dateTime.Format("02.01.2006T15:04:05") {
		t.Errorf("Expected DatVrijeme %v, got %v", dateTime.Format("02.01.2006T15:04:05"), invoice.DatVrijeme)
	}

	if invoice.OznSlijed != "P" {
		t.Errorf("Expected OznSlijed %v, got %v", "P", invoice.OznSlijed)
	}

	if invoice.BrRac.BrOznRac != brOznRac {
		t.Errorf("Expected BrOznRac %v, got %v", brOznRac, invoice.BrRac.BrOznRac)
	}

	if invoice.BrRac.OznPosPr != "TEST3" {
		t.Errorf("Expected OznPosPr %v, got %v", "TEST3", invoice.BrRac.OznPosPr)
	}

	if invoice.BrRac.OznNapUr != oznNapUr {
		t.Errorf("Expected OznNapUr %v, got %v", oznNapUr, invoice.BrRac.OznNapUr)
	}

	if invoice.IznosUkupno != iznosUkupno {
		t.Errorf("Expected IznosUkupno %v, got %v", iznosUkupno, invoice.IznosUkupno)
	}

	if invoice.NacinPlac != nacinPlac {
		t.Errorf("Expected NacinPlac %v, got %v", nacinPlac, invoice.NacinPlac)
	}

	if invoice.OibOper != oibOper {
		t.Errorf("Expected OibOper %v, got %v", oibOper, invoice.OibOper)
	}

	if invoice.ZastKod != zki {
		t.Errorf("Expected ZastKod %v, got %v", zki, invoice.ZastKod)
	}

	if invoice.NakDost != nakDost {
		t.Errorf("Expected NakDost %v, got %v", nakDost, invoice.NakDost)
	}

	if invoice.ParagonBrRac != paragonBrRac {
		t.Errorf("Expected ParagonBrRac %v, got %v", paragonBrRac, invoice.ParagonBrRac)
	}

	if invoice.SpecNamj != specNamj {
		t.Errorf("Expected SpecNamj %v, got %v", specNamj, invoice.SpecNamj)
	}

	// Additional checks for nullable fields
	if invoice.Pdv == nil {
		t.Errorf("Expected Pdv to be non-nil")
	}

	if invoice.Pnp == nil {
		t.Errorf("Expected Pnp to be non-nil")
	}

	if invoice.OstaliPor == nil {
		t.Errorf("Expected OstaliPor to be non-nil")
	}

	if invoice.Naknade == nil {
		t.Errorf("Expected Naknade to be non-nil")
	}

	//Combine with zahtjev for final XML
	zahtjev := RacunZahtjev{
		Zaglavlje: NewFiskalHeader(),
		Racun:     invoice,
		Xmlns:     DefaultNamespace,
		IdAttr:    generateUniqueID(),
	}

	t.Logf("Zahtijev UUID: %s", zahtjev.Zaglavlje.IdPoruke)
	t.Logf("Zahtijev Timestamp: %s", zahtjev.Zaglavlje.DatumVrijeme)

	// Marshal the RacunZahtjev to XML
	xmlData, err := xml.MarshalIndent(zahtjev, "", " ")
	if err != nil {
		t.Fatalf("Error marshalling RacunZahtjev: %v", err)
	}

	t.Log(string(xmlData))

	// Lets send it to CIS and see if we get a response
	body, status, errComm := testEntity.GetResponse(xmlData, true)

	if errComm != nil {
		t.Fatalf("Failed to make request: %v", errComm)
	}

	//unmarshad bodyu to get Racun Odgovor
	var racunOdgovor RacunOdgovor
	if err := xml.Unmarshal(body, &racunOdgovor); err != nil {
		t.Fatalf("failed to unmarshal XML response: %v\n%v", err, string(body))
	}

	//output zaglavlje first all elements
	t.Logf("Racun Odgovor IdPoruke: %s", racunOdgovor.Zaglavlje.IdPoruke)
	t.Logf("Racun Odgovor DatumVrijeme: %s", racunOdgovor.Zaglavlje.DatumVrijeme)

	if status != 200 {

		//all errors one by one
		for _, greska := range racunOdgovor.Greske.Greska {
			t.Logf("Racun Odgovor Greska: %s: %s", greska.SifraGreske, greska.PorukaGreske)
		}

	} else {
		//output JIR: Jedinicni identifikator racuna
		t.Logf("Racun Odgovor JIR: %s", racunOdgovor.Jir)
	}
}
