package fiskalhrgo

import "testing"

// Test embeded CIS demo certificate
func TestParseAndVerifyEmbeddedCertsDemo(t *testing.T) {
	t.Logf("Testing embedded CIS demo certificate...")

	// Parse and verify the embedded CIS demo certificate
	_, err := GetDemoPublicKey()
	if err != nil {
		t.Fatalf("Failed to parse and verify embedded CIS demo certificate: %v", err)
	}

	t.Logf("Embedded CIS demo certificate parsed and verified successfully")
}

// Test embeded CIS production certificate
func TestParseAndVerifyEmbeddedCertsProd(t *testing.T) {
	t.Logf("Testing embedded CIS production certificate...")

	// Parse and verify the embedded CIS production certificate
	_, err := GetProductionPublicKey()
	if err != nil {
		t.Fatalf("Failed to parse and verify embedded CIS production certificate: %v", err)
	}

	t.Logf("Embedded CIS production certificate parsed and verified successfully")
}
