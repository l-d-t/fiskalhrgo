package fiskalhrgo

import "testing"

func TestCISSpecificationMetadata(t *testing.T) {
	if ReleaseVersion != "v1.27.1" {
		t.Fatalf("ReleaseVersion = %q, want %q", ReleaseVersion, "v1.27.1")
	}
	if CISSpecificationVersion != "2.7" {
		t.Fatalf("CISSpecificationVersion = %q, want %q", CISSpecificationVersion, "2.7")
	}
	if CISSpecificationDate != "2026-07-21" {
		t.Fatalf("CISSpecificationDate = %q, want %q", CISSpecificationDate, "2026-07-21")
	}
}
