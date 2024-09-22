package fiskalhrgo

import "testing"

func TestCheckCurrency(t *testing.T) {
	t.Logf("Testing currency validation...")

	// Test a valid currency
	if !IsValidCurrencyFormat("100.00") {
		t.Fatalf("Expected currency 100.00 to be valid")
	}

	// Test a valid currency
	if !IsValidCurrencyFormat("13.12") {
		t.Fatalf("Expected currency 13.12 to be valid")
	}

	// Test a valid currency
	if !IsValidCurrencyFormat("1.12") {
		t.Fatalf("Expected currency 1.12 to be valid")
	}

	// Test a valid currency
	if !IsValidCurrencyFormat("134876348653847632687.99") {
		t.Fatalf("Expected currency 134876348653847632687.99 to be valid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("100.001") {
		t.Fatalf("Expected currency 100.001 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("100,00") {
		t.Fatalf("Expected currency 100,00 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("100") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("abc") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("abc.fg") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("abc.23") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test an invalid currency
	if IsValidCurrencyFormat("100.ab") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	// Test negative currency
	if IsValidCurrencyFormat("-100.00") {
		t.Fatalf("Expected currency 100 to be invalid")
	}

	//Test zero
	if !IsValidCurrencyFormat("0.00") {
		t.Fatalf("Expected currency 0.00 to be valid")
	}

	//Test zero
	if IsValidCurrencyFormat("0.0") {
		t.Fatalf("Expected currency 0.0 to be invalid")
	}

	//Test zero
	if IsValidCurrencyFormat("0") {
		t.Fatalf("Expected currency 0 to be invalid")
	}
}

func TestValidateLocationID(t *testing.T) {
	t.Logf("Testing location ID validation...")

	if !ValidateLocationID("12345678") {
		t.Fatalf("Expected location ID 12345678 to be valid")
	}

	if !ValidateLocationID("TEST3") {
		t.Fatalf("Expected location ID TEST3 to be valid")
	}

	if !ValidateLocationID("POS1") {
		t.Fatalf("Expected location ID POS1 to be valid")
	}

	if !ValidateLocationID("1") {
		t.Fatalf("Expected location ID 1 to be valid")
	}

	if !ValidateLocationID("1234567a") {
		t.Fatalf("Expected location ID 1234567a to be valid")
	}

	if ValidateLocationID("1234567!") {
		t.Fatalf("Expected location ID 1234567! to be invalid")
	}

	if ValidateLocationID("1234567.") {
		t.Fatalf("Expected location ID 1234567. to be invalid")
	}

	if ValidateLocationID("POS 1") {
		t.Fatalf("Expected location ID POS 1 to be invalid")
	}

	if ValidateLocationID("fdkjhfdhjfdshfshfdkhfd87549549875kjfhhfdshfjhjdshkjdfsk7554875kjgfkjgfsssssssssssssssssssssssss") {
		t.Fatalf("Expected location ID fdkjhfdhjfdshfshfdkhfd87549549875kjfhhfdshfjhjdshkjdfsk7554875kjgfkjgfsssssssssssssssssssssssss to be invalid")
	}
}

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
