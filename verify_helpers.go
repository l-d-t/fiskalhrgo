package fiskalhrgo

import (
	"errors"
	"strings"
)

// extractSignatureParts extracts the SignedInfo element and SignatureValue
// from the provided XML data. It returns the canonical SignedInfo XML and
// the base64 encoded SignatureValue.
func extractSignatureParts(xmlData []byte) ([]byte, []byte, error) {
	s := string(xmlData)
	startSig := strings.Index(s, "<ds:Signature")
	closeSig := "</ds:Signature>"
	if startSig == -1 {
		startSig = strings.Index(s, "<Signature")
		closeSig = "</Signature>"
	}
	if startSig == -1 {
		return nil, nil, errors.New("signature not found")
	}
	endSig := strings.Index(s[startSig:], closeSig)
	if endSig == -1 {
		return nil, nil, errors.New("signature not closed")
	}
	endSig += startSig + len(closeSig)
	sigBlock := s[startSig:endSig]

	startInfo := strings.Index(sigBlock, "<ds:SignedInfo")
	closeInfo := "</ds:SignedInfo>"
	if startInfo == -1 {
		startInfo = strings.Index(sigBlock, "<SignedInfo")
		closeInfo = "</SignedInfo>"
	}
	if startInfo == -1 {
		return nil, nil, errors.New("SignedInfo not found")
	}
	endInfo := strings.Index(sigBlock[startInfo:], closeInfo)
	if endInfo == -1 {
		return nil, nil, errors.New("SignedInfo not closed")
	}
	endInfo += startInfo + len(closeInfo)
	signedInfoStr := sigBlock[startInfo:endInfo]

	startVal := strings.Index(sigBlock, "<SignatureValue>")
	if startVal == -1 {
		return nil, nil, errors.New("SignatureValue not found")
	}
	endVal := strings.Index(sigBlock[startVal:], "</SignatureValue>")
	if endVal == -1 {
		return nil, nil, errors.New("SignatureValue not closed")
	}
	value := strings.TrimSpace(sigBlock[startVal+len("<SignatureValue>") : startVal+endVal])
	return []byte(signedInfoStr), []byte(value), nil
}
