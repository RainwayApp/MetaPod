package windows

import (
	"bytes"
	"encoding/asn1"
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/RainwayApp/metapod/errors"
	"github.com/RainwayApp/metapod/structs"
)

// http://msdn.microsoft.com/en-us/library/ms920091.aspx.
const (
	certificateRevision   = 0x200
	certificateType       = 2
	certificateTableIndex = 4
)

//Takes a given input file and creates a Portable Executable wrapper.
//See the getAttributes documentation for more information.
func GetPortableExecutable(stub []byte) (*structs.PortableExecutable, error) {
	offset, size, certSizeOffset, err := getAttributes(stub)
	if err != nil {
		return nil, err
	}
	attributeCertificates := stub[offset : offset+size]
	asn1Data, appendedTag, err := processAttributeCertificates(attributeCertificates)
	if err != nil {
		return nil, err
	}

	var signedData structs.X509Certificate
	if _, err := asn1.Unmarshal(asn1Data, &signedData); err != nil {
		return nil, errors.NewError(1010)
	}

	der, errm := asn1.Marshal(signedData)
	if errm != nil {
		return nil, errors.NewError(1011)
	}

	if !bytes.Equal(der, asn1Data) {
		return nil, errors.NewError(1012)
	}

	return &structs.PortableExecutable{
		Contents:        stub,
		AttrCertOffset:  offset,
		CertSizeOffset:  certSizeOffset,
		Asn1Data:        asn1Data,
		AppendedTag:     appendedTag,
		X509Certificate: &signedData,
	}, nil
}

// Parses the certificates section of a portable executable, returning the ASN.1 data.
func processAttributeCertificates(certs []byte) (asn1, appendedTag []byte, err error) {
	if len(certs) < 8 {
		err = errors.NewError(1009)
		return
	}

	// This reads a WIN_CERTIFICATE structure from
	// TODO Struct this up
	// http://msdn.microsoft.com/en-us/library/ms920091.aspx.
	certificateLength := binary.LittleEndian.Uint32(certs[:4])
	revision := binary.LittleEndian.Uint16(certs[4:6])
	certificateType := binary.LittleEndian.Uint16(certs[6:8])

	if int(certificateLength) != len(certs) {
		err = errors.NewError(1008)
		return
	}

	if revision != certificateRevision {
		err = errors.NewError(1007)
		return
	}

	if certificateType != certificateType {
		err = errors.NewError(1006)
		return
	}

	asn1 = certs[8:]

	if len(asn1) < 2 {
		err = errors.NewError(1005)
		return
	}

	// Read the ASN.1 length of the object.
	var asn1Length int
	if asn1[1]&0x80 == 0 {
		// Short form length.
		asn1Length = int(asn1[1]) + 2
	} else {
		numBytes := int(asn1[1] & 0x7f)
		if numBytes == 0 || numBytes > 2 {
			err = errors.NewError(1004)
			return
		}
		if len(asn1) < numBytes+2 {
			err = errors.NewError(1005)
			return
		}
		asn1Length = int(asn1[2])
		if numBytes == 2 {
			asn1Length <<= 8
			asn1Length |= int(asn1[3])
		}
		asn1Length += 2 + numBytes
	}

	appendedTag = asn1[asn1Length:]
	asn1 = asn1[:asn1Length]

	return
}

//Validates an input file is a Portable Executable, meeting the baseline requirements.
//The input file must be 32x, not a DLL, and a valid EXE.
//Return  offset information about the certificate table.
func getAttributes(stub []byte) (offset, size, sizeOffset int, err error) {
	// offsetOfPEHeaderOffset is the offset in the binary where the PE header is found.
	const offsetOfPEHeaderOffset = 0x3c
	if len(stub) < offsetOfPEHeaderOffset+4 {
		err = errors.NewError(1020)
		return
	}

	peOffset := int(binary.LittleEndian.Uint32(stub[offsetOfPEHeaderOffset:]))
	if peOffset < 0 || peOffset+4 < peOffset {
		err = errors.NewError(1021)
		return
	}
	if len(stub) < peOffset+4 {
		err = errors.NewError(1020)
		return
	}
	pe := stub[peOffset:]
	if !bytes.Equal(pe[:4], []byte{'P', 'E', 0, 0}) {
		err = errors.NewError(1023)
		return
	}

	reader := io.Reader(bytes.NewReader(pe[4:]))
	var fileHeader structs.FileHeader
	if readError := binary.Read(reader, binary.LittleEndian, &fileHeader); readError != nil {
		err = errors.NewError(1024)
		return
	}

	if !fileHeader.IsExe() {
		err = errors.NewError(1025)
		return
	}

	if fileHeader.IsDll() {
		err = errors.NewError(1026)
		return
	}

	is32Bit := fileHeader.Is32Bit()

	press := int64(fileHeader.SizeOfOptionalHeader) + (int64(unsafe.Sizeof(structs.SectionHeader{})) * int64(fileHeader.NumberOfSections))

	reader = io.LimitReader(reader, press)

	optionalHeaderAttributesFn := getOptionalHeaderAttributes64

	if is32Bit {
		optionalHeaderAttributesFn = getOptionalHeaderAttributes64
	}

	certificateAddress, certificateSize, numRvasAndSizes, err := optionalHeaderAttributesFn(reader)

	if err != nil {
		return
	}

	var sectionHeaders = make([]structs.SectionHeader, fileHeader.NumberOfSections)

	for headerNumber := 0; headerNumber < len(sectionHeaders); headerNumber++ {

		var sectionHeader structs.SectionHeader
		if readError := binary.Read(reader, binary.LittleEndian, &sectionHeader); readError != nil {
			err = errors.NewError(1029)
			return
		}
		sectionHeaders[headerNumber] = sectionHeader
	}

	if certificateAddress == 0 {
		err = errors.NewError(1030)
		return
	}

	certEntryEnd := certificateAddress + certificateSize
	if certEntryEnd < certificateAddress {
		err = errors.NewError(1031)
		return
	}
	if certEntryEnd != uint32(len(stub)) {
		err = errors.NewError(1032)
		return
	}

	offset = int(certificateAddress)
	size = int(certificateSize)
	sizeOffset = int(peOffset) + 4 + int(unsafe.Sizeof(structs.FileHeader{})) + int(fileHeader.SizeOfOptionalHeader) - 8*(int(numRvasAndSizes)-certificateTableIndex) + 4

	if binary.LittleEndian.Uint32(stub[sizeOffset:]) != certificateSize {
		err = errors.NewError(1033)
		return
	}

	return
}

func getOptionalHeaderAttributes32(reader io.Reader) (certificateAddress, certificateSize, numRvasAndSizes uint32, err error) {
	var optionalHeader structs.OptionalHeader32
	if readError := binary.Read(reader, binary.LittleEndian, &optionalHeader); readError != nil {
		err = errors.NewError(1028)
		return
	}

	certificateAddress = optionalHeader.CertificateTable.VirtualAddress
	certificateSize = optionalHeader.CertificateTable.Size
	numRvasAndSizes = optionalHeader.NumberOfRvaAndSizes

	return
}

func getOptionalHeaderAttributes64(reader io.Reader) (certificateAddress, certificateSize, numRvasAndSizes uint32, err error) {
	var optionalHeader structs.OptionalHeader64
	if readError := binary.Read(reader, binary.LittleEndian, &optionalHeader); readError != nil {
		err = errors.NewError(1028)
		return
	}

	certificateAddress = optionalHeader.CertificateTable.VirtualAddress
	certificateSize = optionalHeader.CertificateTable.Size
	numRvasAndSizes = optionalHeader.NumberOfRvaAndSizes

	return
}
