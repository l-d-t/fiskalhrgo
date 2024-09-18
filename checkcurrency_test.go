package fiskalhrgo

import "testing"

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

	// Test negative currency
	if isValidCurrencyFormat("-100.00") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	//Test zero
	if !isValidCurrencyFormat("0.00") {
		t.Fatalf("Expected currency 0.00 to be valid")
	}

	//Test zero
	if isValidCurrencyFormat("0.0") {
		t.Fatalf("Expected currency 0.0 to be invalid")
	}

	//Test zero
	if isValidCurrencyFormat("0") {
		t.Fatalf("Expected currency 0 to be invalid")
	}
}
