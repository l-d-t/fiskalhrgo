package fiskalhrgo

import (
	"fmt"
	"os"
	"testing"
)

var certManager *CertManager

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

	certPath := os.Getenv("FISKALHRGO_TEST_CERT_PATH")
	certPassword := os.Getenv("FISKALHRGO_TEST_CERT_PASSWORD")

	if certPath == "" || certPassword == "" {
		fmt.Println("FISKALHRGO_TEST_CERT_PATH or FISKALHRGO_TEST_CERT_PASSWORD environment variables are not set")
		os.Exit(1)
	}

	fmt.Printf("Using certificate: %s\n", certPath)

	// Initialize the CertManager only once for all tests
	certManager = &CertManager{}

	// Load the certificate
	err := certManager.DecodeP12Cert(certPath, certPassword)
	if err != nil {
		fmt.Printf("Failed to load and decode certificate: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	// Run tests
	code := m.Run()

	// Clean up (if needed)

	// Exit with the code returned by m.Run()
	os.Exit(code)
}

// TestLoadCert checks if the certificate was loaded correctly.
func TestLoadCert(t *testing.T) {
	t.Logf("Testing certificate loading...")

	// Ensure that the certManager is not nil and contains a valid certificate
	if certManager == nil {
		t.Fatal("CertManager is nil")
	}

	if certManager.publicCert == nil {
		t.Fatalf("Certificate is not loaded correctly")
	}

	t.Logf("Successfully loaded certificate. Issuer: %s, Subject: %s", certManager.certInfo["issuer"], certManager.certInfo["subject"])
}

// TestGenerateZKI tests the ZKI generation using the previously loaded certificate
func TestGenerateZKI(t *testing.T) {
	t.Logf("Testing ZKI generation...")

	// Reuse the loaded certManager to generate ZKI
	zki, err := GenerateZKI("12345678901", "01.01.2024 12:34:00", "1", "LOC1", "1", "100.00", certManager)
	if err != nil {
		t.Fatalf("Failed to generate ZKI: %v", err)
	}

	if zki == "" {
		t.Fatalf("Expected non-empty ZKI, but got an empty string")
	}

	t.Logf("Generated ZKI: %s", zki)
}

// Please note that a ZKI is dependent on the private key used to sign the data.
// If you use a different private key, the ZKI will be different.
// The ZKI generated in this test is for the known certificate and data.
// When the certificate used for automated tests is changed, this ZKI will also change.
// The test need to be updated with the new ZKI value.
func TestKnownZKI(t *testing.T) {
	t.Logf("Testing Known LDT ZKI generation...")

	// Reuse the loaded certManager to generate ZKI
	zki, err := GenerateZKI("65049901548", "17.05.2024 16:00:38", "13", "TEST3", "1", "90.00", certManager)
	if err != nil {
		t.Fatalf("Failed to generate ZKI: %v", err)
	}

	if zki == "" {
		t.Fatalf("Expected non-empty ZKI, but got an empty string")
	}

	if zki != "0b173c6127809d4f0fff53e13222c819" {
		t.Fatalf("Expected ZKI: 0b173c6127809d4f0fff53e13222c819, got %s", zki)
	}
}
