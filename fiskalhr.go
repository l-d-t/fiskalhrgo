package fiskalhrgo

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
)

// GenerateZKI generates the ZKI based on the given data, similar to your PHP version
func GenerateZKI(oib, issueDateTime, invoiceNumber, location, deviceID, totalAmount string, cm *CertManager) (string, error) {
	// Concatenate the required data (oib, date, invoice number, location, device ID, total amount)
	guardCode := oib + issueDateTime + invoiceNumber + location + deviceID + totalAmount

	// Hash the concatenated data using SHA1
	hashed := sha1.Sum([]byte(guardCode))

	// Use the private key from the CertManager to sign the hashed data with RSA and SHA1
	var signature []byte
	signature, err := rsa.SignPKCS1v15(rand.Reader, cm.privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Generate the MD5 hash of the signature
	md5Hash := md5.Sum(signature)

	// Return the ZKI as a hexadecimal string
	zki := fmt.Sprintf("%x", md5Hash[:])
	return zki, nil
}
