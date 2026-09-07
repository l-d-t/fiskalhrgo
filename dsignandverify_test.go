package fiskalhrgo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"testing"
	"time"
)

func newSignatureTestEntity(t *testing.T) *FiskalEntity {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: "Local XML signature test"},
		NotBefore: time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &FiskalEntity{
		oib: "65049901548", locationID: "TEST", centralizedInvoiceNumber: true,
		cert:    &certManager{privateKey: key, publicCert: cert},
		ciscert: &signatureCheckCIScert{PublicCert: cert},
	}
}

func TestLegacyVerificationBypass(t *testing.T) {
	entity := &FiskalEntity{}
	if ok, err := entity.verifyXML([]byte("not even XML")); !ok || err != nil {
		t.Fatalf("intentional legacy bypass changed: %v, %v", ok, err)
	}
	if ok, err := entity.VerifyXML([]byte("not even XML")); ok || err == nil {
		t.Fatal("public strict verification must not claim success without native mode")
	}
}
