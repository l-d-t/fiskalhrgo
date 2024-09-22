package fiskalhrgo

import "testing"

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
