package fiskalhrgo

import "testing"

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
