//go:build libxml2 && linux

package fiskalhrgo

import (
	"testing"
	"time"
)

// TestPingLibxml2 ensures PingCIS succeeds when libxml2 validation
// is enabled on Linux.
func TestPingLibxml2(t *testing.T) {
	testEntity.SetUseLibxml2(true)
	defer testEntity.SetUseLibxml2(false)

	if err := testEntity.PingCIS(); err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

// TestEchoLibxml2 ensures EchoRequest succeeds when libxml2
// validation of the response is enabled on Linux.
func TestEchoLibxml2(t *testing.T) {
	testEntity.SetUseLibxml2(true)
	defer testEntity.SetUseLibxml2(false)

	msg := "Hello, libxml2"
	resp, err := testEntity.EchoRequest(msg)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resp != msg {
		t.Fatalf("Expected the sent message returned!")
	}
}

// TestSimpleInvoiceFromReadmeVerifyLibXml2 exercises a simple invoice
// request with libxml2 based signature validation enabled.
func TestSimpleInvoiceFromReadmeVerifyLibXml2(t *testing.T) {
	testEntity.SetUseLibxml2(true)
	defer testEntity.SetUseLibxml2(false)

	invoice, _, err := testEntity.NewCISInvoice(
		time.Now(),
		uint(1236),
		uint(1),
		[][]interface{}{{"25.00", "1000.00", "250.00"}},
		nil,
		nil,
		"0.00",
		"0.00",
		"0.00",
		nil,
		"1250.00",
		CISCash,
		"12345678901",
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	jir, zkiR, err := invoice.InvoiceRequest()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	t.Logf("We got a JIR!: %v, ZKI: %v", jir, zkiR)
}
