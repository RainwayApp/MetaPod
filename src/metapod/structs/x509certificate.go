package structs

import "encoding/asn1"

//Represents the structure of an X509Certificate.
//We are able to use this struct to easily marshal to and from DER encoding.
type X509Certificate struct {
	Type  asn1.ObjectIdentifier
	PKCS7 struct {
		Version      int
		Digests      asn1.RawValue
		ContentInfo  asn1.RawValue
		Certificates []asn1.RawValue `asn1:"tag:0, optional, set"`
		SignerInfos  asn1.RawValue
	} `asn1:"explicit, tag:0"`
}
