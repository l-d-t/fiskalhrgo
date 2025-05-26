package fiskalhrgo

import (
	"errors"
	"strings"
)

// extractSignatureValue returns the base64 encoded SignatureValue from the
// XML document.
func extractSignatureValue(xmlData []byte) ([]byte, error) {
	s := string(xmlData)
	startSig := strings.Index(s, "<SignatureValue>")
	if startSig == -1 {
		return nil, errors.New("SignatureValue not found")
	}
	endVal := strings.Index(s[startSig:], "</SignatureValue>")
	if endVal == -1 {
		return nil, errors.New("SignatureValue not closed")
	}
	value := strings.TrimSpace(s[startSig+len("<SignatureValue>") : startSig+endVal])
	return []byte(value), nil
}
