package fiskalhrgo

// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors

import "testing"

// Expected serial number of the embedded CIS demo certificate currently in use
const expectedDemoSerial = "125647116929350878241146578378192994801"

// Expected serial number of the embedded CIS production certificate currently in use
const expectedProdSerial = "86885036496942243675781140548806350693"

// Test embedded CIS demo certificate
func TestParseAndVerifyEmbeddedCertsDemo(t *testing.T) {
	t.Logf("Testing embedded CIS demo certificate...")

	// Parse and verify the embedded CIS demo certificate
	cert, err := getDemoPublicKey()
	if err != nil {
		t.Fatalf("Failed to parse and verify embedded CIS demo certificate: %v", err)
	}

	t.Logf("Embedded CIS demo certificate parsed and verified successfully")
	t.Logf("Subject: %s", cert.Subject)
	t.Logf("Serial: %s", cert.Serial)
	t.Logf("Issuer: %s", cert.Issuer)
	t.Logf("Valid from: %v", cert.ValidFrom)
	t.Logf("Valid until: %v", cert.ValidUntil)

	// Check if the serial number matches the expected value
	if cert.Serial != expectedDemoSerial {
		t.Fatalf("Expected serial number %s, but got %s", expectedDemoSerial, cert.Serial)
	}
}

// Test embedded CIS production certificate
func TestParseAndVerifyEmbeddedCertsProd(t *testing.T) {
	t.Logf("Testing embedded CIS production certificate...")

	// Parse and verify the embedded CIS production certificate
	cert, err := getProductionPublicKey()
	if err != nil {
		t.Fatalf("Failed to parse and verify embedded CIS production certificate: %v", err)
	}

	t.Logf("Embedded CIS production certificate parsed and verified successfully")
	t.Logf("Subject: %s", cert.Subject)
	t.Logf("Serial: %s", cert.Serial)
	t.Logf("Issuer: %s", cert.Issuer)
	t.Logf("Valid from: %v", cert.ValidFrom)
	t.Logf("Valid until: %v", cert.ValidUntil)

	// Check if the serial number matches the expected value
	if cert.Serial != expectedProdSerial {
		t.Fatalf("Expected serial number %s, but got %s", expectedProdSerial, cert.Serial)
	}
}
