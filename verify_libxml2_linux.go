//go:build libxml2 && linux

package fiskalhrgo

/*
#cgo pkg-config: libxml-2.0
#include <libxml/parser.h>
#include <libxml/c14n.h>
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <stdlib.h>

int canonicalize_signed_info(char* xml, int xmlLen, char **out) {
    xmlDocPtr doc = xmlReadMemory(xml, xmlLen, NULL, NULL, 0);
    if (doc == NULL) return -1;

    xmlNodePtr root = xmlDocGetRootElement(doc);
    if (root == NULL) {
        xmlFreeDoc(doc);
        return -2;
    }

    xmlNodeSetPtr nodeset = xmlXPathNodeSetCreate(root);
    if (nodeset == NULL) {
        xmlFreeDoc(doc);
        return -3;
    }

    int ret = xmlC14NDocDumpMemory(doc, nodeset, XML_C14N_1_0, NULL, 0, (xmlChar**)out);

    xmlXPathFreeNodeSet(nodeset);
    xmlFreeDoc(doc);
    return ret;
}

void freeCanon(char* p) { xmlFree(p); }
*/
import "C"

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"unsafe"
)

func canonicalizeSignedInfo(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}
	cxml := (*C.char)(unsafe.Pointer(&data[0]))
	var out *C.char
	res := C.canonicalize_signed_info(cxml, C.int(len(data)), &out)
	if res < 0 {
		return nil, errors.New("libxml2 canonicalization failed")
	}
	defer C.freeCanon(out)
	return C.GoBytes(unsafe.Pointer(out), res), nil
}

func (fe *FiskalEntity) verifyXMLLibxml(xmlData []byte) (bool, error) {
	signedInfo, sigValue, err := extractSignatureParts(xmlData)
	if err != nil {
		return false, err
	}
	canonical, err := canonicalizeSignedInfo(signedInfo)
	if err != nil {
		return false, err
	}
	hash := sha1.Sum(canonical)
	sigBytes, err := base64.StdEncoding.DecodeString(string(sigValue))
	if err != nil {
		return false, err
	}
	pub, ok := fe.ciscert.PublicCert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return false, errors.New("invalid public key type")
	}
	if err := rsa.VerifyPKCS1v15(pub, crypto.SHA1, hash[:], sigBytes); err != nil {
		return false, err
	}
	return true, nil
}
