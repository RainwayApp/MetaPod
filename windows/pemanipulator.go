package windows

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/binary"
	"math/big"

	"github.com/RainwayApp/metapod/structs"
	"github.com/RainwayApp/metapod/utils"
)

//You cannot extend types defined in other packages, so we need to wrap it here.
//I can't help but to feel the design of this project is flawed if we've ended up here. Oh well, first Go project.
type TargetExecutable struct {
	structs.PortableExecutable
}

//this OID is not official and is used purely as a way to identify our custom certificate
var metaPodOID = asn1.ObjectIdentifier([]int{2, 4, 6, 8, 5, 1, 94659, 2, 1, 9000})

//issuerSubjectName will be set on the CA certificate,
// payloadSubjectName will be set on the certificate containing the OID and payload.
const (
	issuerSubjectName  = "Johto"
	payloadSubjectName = "Metapod Cert"
	metaPodSerial      = int64(102946554)
)

const (
	//The certificate validity period must be expired for this to work correctly.
	notBeforeTime = "Mon Jan 1 1:00:00 UTC 2018"
	notAfterTime  = "Mon Apr 1 1:00:00 UTC 2018"
)

//CreateFromTemplate adds a superfluous certificate to a portable executable.
//The appended data is "unverified" and does affect the PE's digital signature.
//This means metadata of any kind can be added to a base executable.
func (portableExecutable *TargetExecutable) CreateFromTemplate(payload []byte) (contents []byte, err int) {
	cert, err, _ := portableExecutable.GetPayload()

	//remove the previous payload if it already existed within the template
	//should we throw here because the template is technically already processed?
	if cert != nil {
		pkcs7 := &portableExecutable.X509Certificate.PKCS7
		pkcs7.Certificates = pkcs7.Certificates[:len(pkcs7.Certificates)-1]
	}

	notBefore := utils.ParseUnixTimeOrDie(notBeforeTime)
	notAfter := utils.ParseUnixTimeOrDie(notAfterTime)
	privateKey, rsaError := rsa.GenerateKey(rand.Reader, 2048)
	if rsaError != nil {
		return nil, 1040
	}

	//this certificate acts as our CA
	issuerCertificate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(metaPodSerial),
		Subject: pkix.Name{
			CommonName: issuerSubjectName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	//this certificate contains our payload, which is added to the extensions.
	payloadCertificate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(metaPodSerial),
		Subject: pkix.Name{
			CommonName: payloadSubjectName,
		},
		Issuer: pkix.Name{
			CommonName: issuerSubjectName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		BasicConstraintsValid: true,
		IsCA:                  false,
		ExtraExtensions: []pkix.Extension{
			{
				Id:    metaPodOID,
				Value: payload,
			},
		},
	}

	//creates a single X509Certificate (DER encoded).
	derCert, certError := x509.CreateCertificate(rand.Reader, &payloadCertificate, &issuerCertificate, &privateKey.PublicKey, privateKey)
	if certError != nil {
		return nil, 1041
	}

	portableExecutable.X509Certificate.PKCS7.Certificates = append(portableExecutable.X509Certificate.PKCS7.Certificates, asn1.RawValue{
		FullBytes: derCert,
	})

	asn1Bytes, asnError := asn1.Marshal(*portableExecutable.X509Certificate)
	if asnError != nil {
		return nil, 1042
	}

	return portableExecutable.restructure(asn1Bytes, portableExecutable.AppendedTag), 0
}

//This function takes the newly appended certificate (that has been serialized into an ASN.1 object)
//and restructure the PE as to replace the previous data -- creating an entirely new executable.
func (portableExecutable *TargetExecutable) restructure(asn1Data, tag []byte) (contents []byte) {
	contents = append(contents, portableExecutable.Contents[:portableExecutable.CertSizeOffset]...)
	for (len(asn1Data)+len(tag))&7 > 0 {
		tag = append(tag, 0)
	}
	attrCertSectionLen := uint32(8 + len(asn1Data) + len(tag))
	var lengthBytes [4]byte
	binary.LittleEndian.PutUint32(lengthBytes[:], attrCertSectionLen)
	contents = append(contents, lengthBytes[:4]...)
	contents = append(contents, portableExecutable.Contents[portableExecutable.CertSizeOffset+4:portableExecutable.AttrCertOffset]...)

	var header [8]byte
	binary.LittleEndian.PutUint32(header[:], attrCertSectionLen)
	binary.LittleEndian.PutUint16(header[4:], certificateRevision)
	binary.LittleEndian.PutUint16(header[6:], certificateType)
	contents = append(contents, header[:]...)
	contents = append(contents, asn1Data...)
	return append(contents, tag...)
}

//Searches a portable executable for the Metapod OID.
//If found, it will return the []value which can then be converted into a string.
//The string is arbitrary, as any format can be included. So it is up to the host program to parse it.
func (portableExecutable *TargetExecutable) GetPayload() (cert *x509.Certificate, err int, payload []byte) {
	n := len(portableExecutable.X509Certificate.PKCS7.Certificates)
	if n == 0 {
		return nil, 1043, nil
	}
	//A Metapod cert should always be the last one on the stack, however I've seen other languages flip the order.
	//So because I "don't trust like that" we are going to loop and find it ourselves.
	for _, der := range portableExecutable.X509Certificate.PKCS7.Certificates {
		if cert, certError := x509.ParseCertificate(der.FullBytes); certError == nil {
			for _, ext := range cert.Extensions {
				if !ext.Critical && ext.Id.Equal(metaPodOID) {
					return cert, 0, ext.Value
				}
			}
		}
	}
	return nil, 0, nil
}
