package fiskalhrgo

import (
	"encoding/xml"
	"fmt"
)

// CISEcho sends an echo request to CIS and processes the response.
func (fe *FiskalEntity) CISEcho(text string) (string, error) {
	// Create an XML payload for the echo request
	echoRequest := &EchoRequest{
		Xmlns: DefaultNamespace,
		Text:  text,
	}

	xmlPayload, err := xml.Marshal(echoRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML payload: %w", err)
	}

	body, _, err := fe.GetResponse(xmlPayload)
	if err != nil {
		return "", err
	}

	// Process the XML response
	var echoResponse EchoResponse
	if err := xml.Unmarshal(body, &echoResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal XML response: %w", err)
	}

	return echoResponse.Text, nil
}
