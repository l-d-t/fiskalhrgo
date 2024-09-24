// SPDX-License-Identifier: MIT
// Copyright (c) 2024 L. D. T. d.o.o.
// Copyright (c) contributors for their respective contributions. See https://github.com/l-d-t/fiskalhrgo/graphs/contributors
package fiskalhrgo

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// iSOAPEnvelope represents a SOAP envelope
type iSOAPEnvelope struct {
	XMLName xml.Name  `xml:"soapenv:Envelope"`
	XmlnsT  string    `xml:"xmlns:tns,attr"` // Declare the tns namespace
	Xmlns   string    `xml:"xmlns:soapenv,attr"`
	Body    iSOAPBody `xml:"soapenv:Body"`
}

// iSOAPBody represents the body of a SOAP envelope
type iSOAPBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Content []byte   `xml:",innerxml"`
}

// iSOAPEnvelopeNoNamespace represents a SOAP envelope without namespace (for CIS responses)
// This to be more flexible and permissive on unmarhaling responses.
type iSOAPEnvelopeNoNamespace struct {
	XMLName xml.Name             `xml:"Envelope"`
	Body    iSOAPBodyNoNamespace `xml:"Body"`
}

// iSOAPBodyNoNamespace represents the body of a SOAP envelope without namespace (for CIS responses)
type iSOAPBodyNoNamespace struct {
	XMLName xml.Name `xml:"Body"`
	Content []byte   `xml:",innerxml"`
}

// GetResponse wraps the XML payload in a SOAP envelope, makes an HTTPS request, and returns the extracted response body.
// - Input: XML payload
// - Output: Response body, error, HTTP status code
func (fe *FiskalEntity) GetResponse(xmlPayload []byte, sign bool) ([]byte, int, error) {
	if fe.ciscert == nil || fe.ciscert.SSLverifyPoll == nil {
		return nil, 0, errors.New("CIScert or SSLverifyPoll is not initialized")
	}

	// Create a custom TLS configuration using TLS 1.3 and the CA pool
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    fe.ciscert.SSLverifyPoll,
	}

	// Create a custom HTTP client with the custom TLS configuration
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: cistimeout * time.Second, // Set a timeout for the request
	}

	if sign {
		// Sign the XML payload
		signedXML, err := fe.signXML(xmlPayload)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to sign XML: %w", err)
		}
		xmlPayload = signedXML
	}

	// Prepare the SOAP envelope with the payload
	soapEnvelope := iSOAPEnvelope{
		XmlnsT: DefaultNamespace,
		Xmlns:  "http://schemas.xmlsoap.org/soap/envelope/",
		Body:   iSOAPBody{Content: xmlPayload},
	}
	// Marshal the SOAP envelope to XML
	marshaledEnvelope, err := xml.Marshal(soapEnvelope)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to marshal SOAP envelope: %w", err)
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", fe.url, bytes.NewBuffer([]byte(marshaledEnvelope)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/xml")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse the SOAP response
	var soapResp iSOAPEnvelopeNoNamespace
	err = xml.Unmarshal(body, &soapResp)
	if err != nil {
		return body, resp.StatusCode, fmt.Errorf("failed to unmarshal SOAP response: %w", err)
	}

	if sign {
		// Verify the signature
		_, err := fe.verifyXML(soapResp.Body.Content)
		if err != nil {
			return soapResp.Body.Content, resp.StatusCode, fmt.Errorf("failed to verify CIS signature: %w", err)
		}
	}

	// Return the inner content of the SOAP Body (the actual response)
	if resp.StatusCode == http.StatusOK {
		return soapResp.Body.Content, resp.StatusCode, nil
	} else {
		return soapResp.Body.Content, resp.StatusCode, fmt.Errorf("CIS returned an error: %v", resp.Status)
	}
}
