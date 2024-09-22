package fiskalhrgo

import (
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// GenerateZKI generates the ZKI based on the given data
func (entity *FiskalEntity) GenerateZKI(issueDateTime time.Time, invoiceNumber uint, deviceID uint, totalAmount string) (string, error) {

	formattedTime := issueDateTime.Format("02.01.2006 15:04:05")

	// Ensure totalAmount is a valid decimal string with 2 decimal places
	if !isValidCurrencyFormat(totalAmount) {
		return "", errors.New("invalid totalAmount format; expected a string with 2 decimal places (e.g., 100.00)")
	}

	// Convert invoiceNumber and deviceID from uint to string
	invoiceNumberStr := strconv.FormatUint(uint64(invoiceNumber), 10)
	deviceIDStr := strconv.FormatUint(uint64(deviceID), 10)

	// Concatenate the required data (oib, date, invoice number, location, device ID, total amount)
	guardCode := entity.oib + formattedTime + invoiceNumberStr + entity.locationID + deviceIDStr + totalAmount

	// Hash the concatenated data using SHA1
	hashed := sha1.Sum([]byte(guardCode))

	// Use the private key from the CertManager to sign the hashed data with RSA and SHA1
	var signature []byte
	signature, err := rsa.SignPKCS1v15(rand.Reader, entity.cert.privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %v", err)
	}

	// Generate the MD5 hash of the signature
	md5Hash := md5.Sum(signature)

	// Return the ZKI as a hexadecimal string
	zki := fmt.Sprintf("%x", md5Hash[:])
	return zki, nil
}
