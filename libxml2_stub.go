//go:build !linux || !cgo || !libxml2

package fiskalhrgo

import "errors"

const libxml2Available = false

func canonicalizeNativeDocument([]byte, bool) ([]byte, error) {
	return nil, errors.New("libxml2 support is unavailable")
}

func canonicalizeNativeSignedInfo([]byte, bool) ([]byte, error) {
	return nil, errors.New("libxml2 support is unavailable")
}

func canonicalizeNativeReference([]byte, string, bool) ([]byte, error) {
	return nil, errors.New("libxml2 support is unavailable")
}
