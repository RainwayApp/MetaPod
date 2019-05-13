package structs

// PortableExecutable represents a PE binary.
type PortableExecutable struct {
	//The contents of the loaded file.
	Contents        []byte
	//The offset to the certificates table.
	AttrCertOffset  int
	//The offset to the size of the certificates.
	CertSizeOffset  int
	//The embedded X509Certificate (DER).
	Asn1Data        []byte
	//Any extra data, if any.
	AppendedTag     []byte
	//The decoded DER data.
	X509Certificate *X509Certificate
}
