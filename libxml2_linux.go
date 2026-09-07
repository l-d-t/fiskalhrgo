//go:build linux && cgo && libxml2

package fiskalhrgo

/*
#cgo pkg-config: libxml-2.0
#include <libxml/parser.h>
#include <libxml/c14n.h>
#include <libxml/xpath.h>
#include <libxml/xpathInternals.h>
#include <stdlib.h>

static const char *dsig_ns = "http://www.w3.org/2000/09/xmldsig#";

static int is_element(xmlNodePtr node, const char *name, const char *ns) {
	return node != NULL && node->type == XML_ELEMENT_NODE &&
		xmlStrEqual(node->name, BAD_CAST name) && node->ns != NULL &&
		xmlStrEqual(node->ns->href, BAD_CAST ns);
}

static void find_signature_parts(xmlNodePtr node, xmlNodePtr *signature,
		xmlNodePtr *signed_info, int *signature_count) {
	for (; node != NULL; node = node->next) {
		if (is_element(node, "Signature", dsig_ns)) {
			(*signature_count)++;
			if (*signature == NULL) *signature = node;
		}
		if (is_element(node, "SignedInfo", dsig_ns) && node->parent == *signature) {
			*signed_info = node;
		}
		find_signature_parts(node->children, signature, signed_info, signature_count);
	}
}

static void find_id(xmlNodePtr node, const xmlChar *id, xmlNodePtr *target, int *count) {
	for (; node != NULL; node = node->next) {
		if (node->type == XML_ELEMENT_NODE) {
			xmlChar *value = xmlGetProp(node, BAD_CAST "Id");
			if (value != NULL) {
				if (xmlStrEqual(value, id)) {
					(*count)++;
					if (*target == NULL) *target = node;
				}
				xmlFree(value);
			}
		}
		find_id(node->children, id, target, count);
	}
}

static int add_subtree(xmlNodeSetPtr set, xmlNodePtr node, xmlNodePtr excluded) {
	if (node == NULL || node == excluded) return 0;
	if (xmlXPathNodeSetAddUnique(set, node) < 0) return -1;
	if (node->type == XML_ELEMENT_NODE) {
		for (xmlAttrPtr attr = node->properties; attr != NULL; attr = attr->next) {
			if (xmlXPathNodeSetAddUnique(set, (xmlNodePtr)attr) < 0) return -1;
		}
		for (xmlNsPtr ns = node->nsDef; ns != NULL; ns = ns->next) {
			if (xmlXPathNodeSetAddNs(set, node, ns) < 0) return -1;
		}
	}
	for (xmlNodePtr child = node->children; child != NULL; child = child->next) {
		if (add_subtree(set, child, excluded) < 0) return -1;
	}
	return 0;
}

static int canonicalize_nodes(xmlDocPtr doc, xmlNodePtr target,
		xmlNodePtr excluded, int exclusive, xmlChar **out) {
	xmlNodeSetPtr set = xmlXPathNodeSetCreate(NULL);
	if (set == NULL) return -10;
	if (add_subtree(set, target, excluded) < 0) {
		xmlXPathFreeNodeSet(set);
		return -11;
	}
	int mode = exclusive ? XML_C14N_EXCLUSIVE_1_0 : XML_C14N_1_0;
	int result = xmlC14NDocDumpMemory(doc, set, mode, NULL, 0, out);
	xmlXPathFreeNodeSet(set);
	return result;
}

int canonicalize_document(char *xml, int xml_len, int exclusive, xmlChar **out) {
	xmlDocPtr doc = xmlReadMemory(xml, xml_len, NULL, NULL, XML_PARSE_NONET);
	if (doc == NULL) return -1;
	int mode = exclusive ? XML_C14N_EXCLUSIVE_1_0 : XML_C14N_1_0;
	int result = xmlC14NDocDumpMemory(doc, NULL, mode, NULL, 0, out);
	xmlFreeDoc(doc);
	return result;
}

int canonicalize_signed_info(char *xml, int xml_len, int exclusive, xmlChar **out) {
	xmlDocPtr doc = xmlReadMemory(xml, xml_len, NULL, NULL, XML_PARSE_NONET);
	if (doc == NULL) return -1;
	xmlNodePtr signature = NULL;
	xmlNodePtr signed_info = NULL;
	int count = 0;
	find_signature_parts(xmlDocGetRootElement(doc), &signature, &signed_info, &count);
	int result = (count == 1 && signed_info != NULL)
		? canonicalize_nodes(doc, signed_info, NULL, exclusive, out) : -2;
	xmlFreeDoc(doc);
	return result;
}

int canonicalize_reference(char *xml, int xml_len, char *id, int exclusive, xmlChar **out) {
	xmlDocPtr doc = xmlReadMemory(xml, xml_len, NULL, NULL, XML_PARSE_NONET);
	if (doc == NULL) return -1;
	xmlNodePtr signature = NULL;
	xmlNodePtr signed_info = NULL;
	int signature_count = 0;
	find_signature_parts(xmlDocGetRootElement(doc), &signature, &signed_info, &signature_count);
	xmlNodePtr target = NULL;
	int target_count = 0;
	find_id(xmlDocGetRootElement(doc), BAD_CAST id, &target, &target_count);
	int result = (signature_count == 1 && target_count == 1)
		? canonicalize_nodes(doc, target, signature, exclusive, out) : -2;
	xmlFreeDoc(doc);
	return result;
}

void free_canonical(xmlChar *value) { xmlFree(value); }
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

func canonicalizeNativeDocument(data []byte, exclusive bool) ([]byte, error) {
	return runNativeCanonicalizer(data, "", exclusive, 0)
}

func canonicalizeNativeSignedInfo(data []byte, exclusive bool) ([]byte, error) {
	return runNativeCanonicalizer(data, "", exclusive, 1)
}

func canonicalizeNativeReference(data []byte, id string, exclusive bool) ([]byte, error) {
	if id == "" {
		return nil, errors.New("XML signature reference ID is empty")
	}
	return runNativeCanonicalizer(data, id, exclusive, 2)
}

func runNativeCanonicalizer(data []byte, id string, exclusive bool, operation int) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("cannot canonicalize empty XML")
	}
	cXML := (*C.char)(unsafe.Pointer(&data[0]))
	cExclusive := C.int(0)
	if exclusive {
		cExclusive = 1
	}
	var output *C.xmlChar
	var result C.int
	switch operation {
	case 0:
		result = C.canonicalize_document(cXML, C.int(len(data)), cExclusive, &output)
	case 1:
		result = C.canonicalize_signed_info(cXML, C.int(len(data)), cExclusive, &output)
	case 2:
		cID := C.CString(id)
		defer C.free(unsafe.Pointer(cID))
		result = C.canonicalize_reference(cXML, C.int(len(data)), cID, cExclusive, &output)
	default:
		return nil, errors.New("invalid native canonicalization operation")
	}
	if result < 0 {
		return nil, fmt.Errorf("libxml2 canonicalization failed with code %d", int(result))
	}
	defer C.free_canonical(output)
	return C.GoBytes(unsafe.Pointer(output), result), nil
}

const libxml2Available = true
