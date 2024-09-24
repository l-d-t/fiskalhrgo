// SPDX-License-Identifier: Apache-2.0
// This file is adapted from the github.com/russellhaering/goxmldsig project.
package fiskalhrgo

import (
	"crypto"
	"crypto/x509"
	"sort"

	"github.com/beevik/etree"
	"github.com/l-d-t/fiskalhrgo/etreeutils" // Import the local etreeutils package
)

const (
	DefaultPrefix = "ds"
	Namespace     = "http://www.w3.org/2000/09/xmldsig#"
)

// Tags
const (
	SignatureTag              = "Signature"
	SignedInfoTag             = "SignedInfo"
	CanonicalizationMethodTag = "CanonicalizationMethod"
	SignatureMethodTag        = "SignatureMethod"
	ReferenceTag              = "Reference"
	TransformsTag             = "Transforms"
	TransformTag              = "Transform"
	DigestMethodTag           = "DigestMethod"
	DigestValueTag            = "DigestValue"
	SignatureValueTag         = "SignatureValue"
	KeyInfoTag                = "KeyInfo"
	X509DataTag               = "X509Data"
	X509CertificateTag        = "X509Certificate"
	InclusiveNamespacesTag    = "InclusiveNamespaces"
)

const (
	AlgorithmAttr  = "Algorithm"
	URIAttr        = "URI"
	DefaultIdAttr  = "Id"
	PrefixListAttr = "PrefixList"
)

type AlgorithmID string

func (id AlgorithmID) String() string {
	return string(id)
}

const (
	RSASHA1SignatureMethod     = "http://www.w3.org/2000/09/xmldsig#rsa-sha1"
	RSASHA256SignatureMethod   = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha256"
	RSASHA384SignatureMethod   = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha384"
	RSASHA512SignatureMethod   = "http://www.w3.org/2001/04/xmldsig-more#rsa-sha512"
	ECDSASHA1SignatureMethod   = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha1"
	ECDSASHA256SignatureMethod = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha256"
	ECDSASHA384SignatureMethod = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha384"
	ECDSASHA512SignatureMethod = "http://www.w3.org/2001/04/xmldsig-more#ecdsa-sha512"
)

// Well-known signature algorithms
const (
	// Supported canonicalization algorithms
	CanonicalXML10ExclusiveAlgorithmId             AlgorithmID = "http://www.w3.org/2001/10/xml-exc-c14n#"
	CanonicalXML10ExclusiveWithCommentsAlgorithmId AlgorithmID = "http://www.w3.org/2001/10/xml-exc-c14n#WithComments"

	CanonicalXML11AlgorithmId             AlgorithmID = "http://www.w3.org/2006/12/xml-c14n11"
	CanonicalXML11WithCommentsAlgorithmId AlgorithmID = "http://www.w3.org/2006/12/xml-c14n11#WithComments"

	CanonicalXML10RecAlgorithmId          AlgorithmID = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315"
	CanonicalXML10WithCommentsAlgorithmId AlgorithmID = "http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments"

	EnvelopedSignatureAltorithmId AlgorithmID = "http://www.w3.org/2000/09/xmldsig#enveloped-signature"
)

var digestAlgorithmIdentifiers = map[crypto.Hash]string{
	crypto.SHA1:   "http://www.w3.org/2000/09/xmldsig#sha1",
	crypto.SHA256: "http://www.w3.org/2001/04/xmlenc#sha256",
	crypto.SHA384: "http://www.w3.org/2001/04/xmldsig-more#sha384",
	crypto.SHA512: "http://www.w3.org/2001/04/xmlenc#sha512",
}

type signatureMethodInfo struct {
	PublicKeyAlgorithm x509.PublicKeyAlgorithm
	Hash               crypto.Hash
}

var digestAlgorithmsByIdentifier = map[string]crypto.Hash{}
var signatureMethodsByIdentifier = map[string]signatureMethodInfo{}

func init() {
	for hash, id := range digestAlgorithmIdentifiers {
		digestAlgorithmsByIdentifier[id] = hash
	}
	for algo, hashToMethod := range signatureMethodIdentifiers {
		for hash, method := range hashToMethod {
			signatureMethodsByIdentifier[method] = signatureMethodInfo{
				PublicKeyAlgorithm: algo,
				Hash:               hash,
			}
		}
	}
}

var signatureMethodIdentifiers = map[x509.PublicKeyAlgorithm]map[crypto.Hash]string{
	x509.RSA: {
		crypto.SHA1:   RSASHA1SignatureMethod,
		crypto.SHA256: RSASHA256SignatureMethod,
		crypto.SHA384: RSASHA384SignatureMethod,
		crypto.SHA512: RSASHA512SignatureMethod,
	},
	x509.ECDSA: {
		crypto.SHA1:   ECDSASHA1SignatureMethod,
		crypto.SHA256: ECDSASHA256SignatureMethod,
		crypto.SHA384: ECDSASHA384SignatureMethod,
		crypto.SHA512: ECDSASHA512SignatureMethod,
	},
}

// Canonicalizer is an implementation of a canonicalization algorithm.
type Canonicalizer interface {
	Canonicalize(el *etree.Element) ([]byte, error)
	Algorithm() AlgorithmID
}

type NullCanonicalizer struct {
}

func MakeNullCanonicalizer() Canonicalizer {
	return &NullCanonicalizer{}
}

func (c *NullCanonicalizer) Algorithm() AlgorithmID {
	return AlgorithmID("NULL")
}

func (c *NullCanonicalizer) Canonicalize(el *etree.Element) ([]byte, error) {
	return canonicalSerialize(canonicalPrep(el, false, true))
}

type c14N10ExclusiveCanonicalizer struct {
	prefixList string
	comments   bool
}

// MakeC14N10ExclusiveCanonicalizerWithPrefixList constructs an exclusive Canonicalizer
// from a PrefixList in NMTOKENS format (a white space separated list).
func MakeC14N10ExclusiveCanonicalizerWithPrefixList(prefixList string) Canonicalizer {
	return &c14N10ExclusiveCanonicalizer{
		prefixList: prefixList,
		comments:   false,
	}
}

// MakeC14N10ExclusiveWithCommentsCanonicalizerWithPrefixList constructs an exclusive Canonicalizer
// from a PrefixList in NMTOKENS format (a white space separated list).
func MakeC14N10ExclusiveWithCommentsCanonicalizerWithPrefixList(prefixList string) Canonicalizer {
	return &c14N10ExclusiveCanonicalizer{
		prefixList: prefixList,
		comments:   true,
	}
}

// Canonicalize transforms the input Element into a serialized XML document in canonical form.
func (c *c14N10ExclusiveCanonicalizer) Canonicalize(el *etree.Element) ([]byte, error) {
	err := etreeutils.TransformExcC14n(el, c.prefixList, c.comments)
	if err != nil {
		return nil, err
	}

	return canonicalSerialize(el)
}

func (c *c14N10ExclusiveCanonicalizer) Algorithm() AlgorithmID {
	if c.comments {
		return CanonicalXML10ExclusiveWithCommentsAlgorithmId
	}
	return CanonicalXML10ExclusiveAlgorithmId
}

type c14N11Canonicalizer struct {
	comments bool
}

// MakeC14N11Canonicalizer constructs an inclusive canonicalizer.
func MakeC14N11Canonicalizer() Canonicalizer {
	return &c14N11Canonicalizer{
		comments: false,
	}
}

// MakeC14N11WithCommentsCanonicalizer constructs an inclusive canonicalizer.
func MakeC14N11WithCommentsCanonicalizer() Canonicalizer {
	return &c14N11Canonicalizer{
		comments: true,
	}
}

// Canonicalize transforms the input Element into a serialized XML document in canonical form.
func (c *c14N11Canonicalizer) Canonicalize(el *etree.Element) ([]byte, error) {
	return canonicalSerialize(canonicalPrep(el, true, c.comments))
}

func (c *c14N11Canonicalizer) Algorithm() AlgorithmID {
	if c.comments {
		return CanonicalXML11WithCommentsAlgorithmId
	}
	return CanonicalXML11AlgorithmId
}

type c14N10RecCanonicalizer struct {
	comments bool
}

// MakeC14N10RecCanonicalizer constructs an inclusive canonicalizer.
func MakeC14N10RecCanonicalizer() Canonicalizer {
	return &c14N10RecCanonicalizer{
		comments: false,
	}
}

// MakeC14N10WithCommentsCanonicalizer constructs an inclusive canonicalizer.
func MakeC14N10WithCommentsCanonicalizer() Canonicalizer {
	return &c14N10RecCanonicalizer{
		comments: true,
	}
}

// Canonicalize transforms the input Element into a serialized XML document in canonical form.
func (c *c14N10RecCanonicalizer) Canonicalize(inputXML *etree.Element) ([]byte, error) {
	parentNamespaceAttributes, parentXmlAttributes := getParentNamespaceAndXmlAttributes(inputXML)
	inputXMLCopy := inputXML.Copy()
	enhanceNamespaceAttributes(inputXMLCopy, parentNamespaceAttributes, parentXmlAttributes)
	return canonicalSerialize(canonicalPrep(inputXMLCopy, true, c.comments))
}

func (c *c14N10RecCanonicalizer) Algorithm() AlgorithmID {
	if c.comments {
		return CanonicalXML10WithCommentsAlgorithmId
	}
	return CanonicalXML10RecAlgorithmId

}

const nsSpace = "xmlns"

// canonicalPrep accepts an *etree.Element and transforms it into one which is ready
// for serialization into inclusive canonical form. Specifically this
// entails:
//
// 1. Stripping re-declarations of namespaces
// 2. Sorting attributes into canonical order
//
// Inclusive canonicalization does not strip unused namespaces.
//
// TODO(russell_h): This is very similar to excCanonicalPrep - perhaps they should
// be unified into one parameterized function?
func canonicalPrep(el *etree.Element, strip bool, comments bool) *etree.Element {
	return canonicalPrepInner(el, make(map[string]string), strip, comments)
}

func canonicalPrepInner(el *etree.Element, seenSoFar map[string]string, strip bool, comments bool) *etree.Element {
	_seenSoFar := make(map[string]string)
	for k, v := range seenSoFar {
		_seenSoFar[k] = v
	}

	ne := el.Copy()
	sort.Sort(etreeutils.SortedAttrs(ne.Attr))
	n := 0
	for _, attr := range ne.Attr {
		if attr.Space != nsSpace && !(attr.Space == "" && attr.Key == nsSpace) {
			ne.Attr[n] = attr
			n++
			continue
		}

		if attr.Space == nsSpace {
			key := attr.Space + ":" + attr.Key
			if uri, seen := _seenSoFar[key]; !seen || attr.Value != uri {
				ne.Attr[n] = attr
				n++
				_seenSoFar[key] = attr.Value
			}
		} else {
			if uri, seen := _seenSoFar[nsSpace]; (!seen && attr.Value != "") || attr.Value != uri {
				ne.Attr[n] = attr
				n++
				_seenSoFar[nsSpace] = attr.Value
			}
		}
	}
	ne.Attr = ne.Attr[:n]

	if !comments {
		c := 0
		for c < len(ne.Child) {
			if _, ok := ne.Child[c].(*etree.Comment); ok {
				ne.RemoveChildAt(c)
			} else {
				c++
			}
		}
	}

	for i, token := range ne.Child {
		childElement, ok := token.(*etree.Element)
		if ok {
			ne.Child[i] = canonicalPrepInner(childElement, _seenSoFar, strip, comments)
		}
	}

	return ne
}

func canonicalSerialize(el *etree.Element) ([]byte, error) {
	doc := etree.NewDocument()
	doc.SetRoot(el.Copy())

	doc.WriteSettings = etree.WriteSettings{
		CanonicalAttrVal: true,
		CanonicalEndTags: true,
		CanonicalText:    true,
	}

	return doc.WriteToBytes()
}

func getParentNamespaceAndXmlAttributes(el *etree.Element) (map[string]string, map[string]string) {
	namespaceMap := make(map[string]string, 23)
	xmlMap := make(map[string]string, 5)
	parents := make([]*etree.Element, 0, 23)
	n1 := el.Parent()
	if n1 == nil {
		return namespaceMap, xmlMap
	}
	parent := n1
	for parent != nil {
		parents = append(parents, parent)
		parent = parent.Parent()
	}
	for i := len(parents) - 1; i > -1; i-- {
		elementPos := parents[i]
		for _, attr := range elementPos.Attr {
			if attr.Space == "xmlns" && (attr.Key != "xml" || attr.Value != "http://www.w3.org/XML/1998/namespace") {
				namespaceMap[attr.Key] = attr.Value
			} else if attr.Space == "" && attr.Key == "xmlns" {
				namespaceMap[attr.Key] = attr.Value
			} else if attr.Space == "xml" {
				xmlMap[attr.Key] = attr.Value
			}
		}
	}
	return namespaceMap, xmlMap
}

func enhanceNamespaceAttributes(el *etree.Element, parentNamespaces map[string]string, parentXmlAttributes map[string]string) {
	for prefix, uri := range parentNamespaces {
		if prefix == "xmlns" {
			el.CreateAttr("xmlns", uri)
		} else {
			el.CreateAttr("xmlns:"+prefix, uri)
		}
	}
	for attr, value := range parentXmlAttributes {
		el.CreateAttr("xml:"+attr, value)
	}
}
