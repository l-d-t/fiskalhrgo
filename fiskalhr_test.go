package fiskalhrgo

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var testCert *CertManager
var testEntity *FiskalEntity

var certPath string
var certPassword string
var testOIB string

// Function to clean up the global variables
func cleanupGlobals() {
	// Reset or clean up any global variables
	testEntity = nil
	testCert = nil
	certPassword = ""

	fmt.Println("Cleaned up global variables.")
}

// TestMain is run before any other tests. It sets up the shared CertManager instance.
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

	certPath = os.Getenv("FISKALHRGO_TEST_CERT_PATH")
	certPassword = os.Getenv("FISKALHRGO_TEST_CERT_PASSWORD")
	testOIB = os.Getenv("FISKALHRGO_TEST_CERT_OIB")

	if certPath == "" || certPassword == "" || testOIB == "" {
		fmt.Println("FISKALHRGO_TEST_CERT_PATH or FISKALHRGO_TEST_CERT_PASSWORD or FISKALHRGO_TEST_CERT_OIB environment variables are not set")
		os.Exit(1)
	}

	fmt.Printf("Using certificate: %s\n", certPath)
	fmt.Printf("Test OIB: %s\n", testOIB)

	fmt.Println("Running tests...")
	// Run tests
	code := m.Run()

	// Clean up (if needed)
	cleanupGlobals()
	// Exit with the code returned by m.Run()
	os.Exit(code)
}

// SECTION 1 - TEST NOT DEPENDING ON THE CERTIFICATE

// Test check OIB
func TestCheckOIB(t *testing.T) {
	t.Logf("Testing OIB validation...")

	// Test a valid OIB
	if !ValidateOIB("65049901548") {
		t.Fatalf("Expected OIB 65049901548 to be valid")
	}

	// Test an invalid OIB
	if ValidateOIB("12345678900") {
		t.Fatalf("Expected OIB 12345678900 to be invalid")
	}
}

// Test valid currency format
func TestCheckCurrency(t *testing.T) {
	t.Logf("Testing currency validation...")

	// Test a valid currency
	if !isValidCurrencyFormat("100.00") {
		t.Fatalf("Expected currency 100.00 to be valid")
	}

	// Test a valid currency
	if !isValidCurrencyFormat("13.12") {
		t.Fatalf("Expected currency 13.12 to be valid")
	}

	// Test a valid currency
	if !isValidCurrencyFormat("1.12") {
		t.Fatalf("Expected currency 1.12 to be valid")
	}

	// Test a valid currency
	if !isValidCurrencyFormat("134876348653847632687.99") {
		t.Fatalf("Expected currency 134876348653847632687.99 to be valid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("100.001") {
		t.Fatalf("Expected currency 100.001 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("100,00") {
		t.Fatalf("Expected currency 100,00 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("100") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("abc") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("abc.fg") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("abc.23") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if isValidCurrencyFormat("100.ab") {
		t.Fatalf("Expected currency 100 to be invalid")
	}
}

// SECTION 2 - TEST DEPENDING ON THE CERTIFICATE

// TestLoadCert checks if the certificate was loaded correctly.
func TestLoadCert(t *testing.T) {
	t.Logf("Testing certificate loading...")

	testCert = NewCertManager()
	// Load the certificate
	err := testCert.DecodeP12Cert(certPath, certPassword)

	if err != nil {
		t.Fatalf("Failed to load certificate: %v", err)
	}

	// Ensure that the certManager is not nil and contains a valid certificate
	if testCert == nil {
		t.Fatal("CertManager is nil")
	}

	if testCert.publicCert == nil {
		t.Fatalf("Certificate is not loaded correctly")
	}

	if !testCert.init_ok {
		t.Fatalf("CertManager initialization failed")
	}

	// Log issuer and subject
	t.Logf("Certificate Subject: %s", testCert.publicCert.Subject)
	t.Logf("Certificate Issuer: %s", testCert.publicCert.Issuer)

	// Log the certificate's serial number
	t.Logf("Certificate Serial Number: %s", testCert.publicCert.SerialNumber)

	// Log the certificate's validity dates
	t.Logf("Certificate Valid From: %v", testCert.publicCert.NotBefore)
	t.Logf("Certificate Valid Until: %v", testCert.publicCert.NotAfter)
}

func TestDisplayCertInfo(t *testing.T) {
	t.Logf("Testing certificate display...")

	t.Log("Cert Text:")
	// Display the certificate information
	fmt.Print(testCert.DisplayCertInfoText())

	t.Log("Cert Markdown:")
	t.Log(testCert.DisplayCertInfoMarkdown())
	t.Log("Cert HTML:")
	t.Log(testCert.DisplayCertInfoHTML())
}

func TestExtractOIB(t *testing.T) {
	t.Logf("Testing OIB extraction...")

	// Reuse the loaded certManager to extract the OIB
	oib, err := testCert.GetCertOIB()
	if err != nil {
		t.Fatalf("Failed to extract OIB: %v", err)
	}

	if oib == "" {
		t.Fatalf("Expected non-empty OIB, but got an empty string")
	}

	if !ValidateOIB(oib) {
		t.Fatalf("Extracted OIB is not valid")
	}

	t.Logf("Extracted OIB: %s", oib)
}

// SECTION 3 - TEST DEPENDING ON FiskalEntity and all working parts

// Test without passing a loaded certificate
func TestNewFiskalEntityStandalone(t *testing.T) {
	t.Logf("Testing FiskalEntity with cert init creation...")
	testEntityStandalone, err := NewFiskalEntity(testOIB, nil, true, certPath, certPassword)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if !testEntityStandalone.Cert.init_ok {
		t.Fatalf("Failed to initialize CertManager")
	}
}

// Test with passing the loaded certificate
// Also init global testEntity for other tests
func TestNewFiskalEntity(t *testing.T) {
	t.Logf("Testing FiskalEntity with passing cert menager creation...")
	var err error
	testEntity, err = NewFiskalEntity(testOIB, testCert, true)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}
}

// TestGenerateZKI tests the ZKI generation using the previously loaded certificate
func TestGenerateZKI(t *testing.T) {
	t.Logf("Testing ZKI generation...")

	// Reuse the loaded certManager to generate ZKI
	zki, err := testEntity.GenerateZKI(time.Now(), 1, "LOC1", 1, "100.00")
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
	zki, err := testEntity.GenerateZKI(parsedTime, 13, "TEST3", 1, "90.00")
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
