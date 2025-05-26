//go:build !libxml2 || !linux

package fiskalhrgo

import "errors"

func (fe *FiskalEntity) verifyXMLLibxml(xmlData []byte) (bool, error) {
	return false, errors.New("libxml2 validation not supported")
}
