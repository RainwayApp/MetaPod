package windows

import (
	"bytes"
	"encoding/asn1"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"metapod/structs"
	"unsafe"
)


// http://msdn.microsoft.com/en-us/library/ms920091.aspx.
const (
	certificateRevision = 0x200
	certificateType     = 2
	certificateTableIndex = 4
)

//Takes a given input file and creates a Portable Executable wrapper.
//See the getAttributes documentation for more information.
func GetPortableExecutable(stub []byte) (*structs.PortableExecutable, error) {
	offset, size, certSizeOffset, err := getAttributes(stub)
	if err != nil {
		return nil, errors.New("Error parsing headers: " + err.Error())
	}
	attributeCertificates := stub[offset : offset+size]
	asn1Data, appendedTag, err := processAttributeCertificates(attributeCertificates)
	if err != nil {
		return nil, errors.New("Error parsing attribute certificate section: " + err.Error())
	}

	var signedData structs.X509Certificate
	if _, err := asn1.Unmarshal(asn1Data, &signedData); err != nil {
		return nil, errors.New("Error while parsing X509Certificate structure: " + err.Error())
	}

	der, err := asn1.Marshal(signedData)
	if err != nil {
		return nil, errors.New("Error while marshaling X509Certificate structure: " + err.Error())
	}

	if !bytes.Equal(der, asn1Data) {
		return nil, errors.New("ASN.1 marshaling test failed: " + err.Error())
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
		err = errors.New("attribute certificate truncated")
		return
	}

	// This reads a WIN_CERTIFICATE structure from
	// TODO Struct this up
	// http://msdn.microsoft.com/en-us/library/ms920091.aspx.
	certificateLength := binary.LittleEndian.Uint32(certs[:4])
	revision := binary.LittleEndian.Uint16(certs[4:6])
	certificateType := binary.LittleEndian.Uint16(certs[6:8])

	if int(certificateLength) != len(certs) {
		err = errors.New("multiple attribute certificates found")
		return
	}

	if revision != certificateRevision {
		err = fmt.Errorf("unknown attribute certificate revision: %x", revision)
		return
	}

	if certificateType != certificateType {
		err = fmt.Errorf("unknown attribute certificate type: %d", certificateType)
		return
	}

	asn1 = certs[8:]

	if len(asn1) < 2 {
		err = errors.New("ASN.1 structure truncated")
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
			err = fmt.Errorf("bad number of bytes in ASN.1 length: %d", numBytes)
			return
		}
		if len(asn1) < numBytes+2 {
			err = errors.New("ASN.1 structure truncated")
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
		err = errors.New("Portable Executable malformed.")
		return
	}

	peOffset := int(binary.LittleEndian.Uint32(stub[offsetOfPEHeaderOffset:]))
	if peOffset < 0 || peOffset+4 < peOffset {
		err = errors.New("Overflowed searching for PE signature")
		return
	}
	if len(stub) < peOffset+4 {
		err = errors.New("Portable Executable malformed.")
		return
	}
	pe := stub[peOffset:]
	if !bytes.Equal(pe[:4], []byte{'P', 'E', 0, 0}) {
		err = errors.New("PE header not found. Is this a Portable Executable?")
		return
	}

	r := io.Reader(bytes.NewReader(pe[4:]))
	var fileHeader structs.FileHeader
	if err = binary.Read(r, binary.LittleEndian, &fileHeader); err != nil {
		return
	}

	if !fileHeader.IsExe() {
		err = errors.New("Input file is not a valid Portable Executable.")
		return
	}

	if fileHeader.IsDll() {
		err = errors.New("Template file cannot be a DLL")
		return
	}

	if !fileHeader.Is32Bit() {
		err = errors.New("Template executable must be 32x.")
		return
	}

	var press = int64(fileHeader.SizeOfOptionalHeader) + (int64(unsafe.Sizeof(structs.SectionHeader{})) * int64(fileHeader.NumberOfSections))

	r = io.LimitReader(r, press)
	var optionalHeader structs.OptionalHeader32
	if err = binary.Read(r, binary.LittleEndian, &optionalHeader); err != nil {
		fmt.Println(unsafe.Sizeof(optionalHeader))
		return
	}

	var sectionHeaders = make([]structs.SectionHeader, fileHeader.NumberOfSections)


	for headerNumber := 0; headerNumber < len(sectionHeaders); headerNumber++ {

		var sectionHeader structs.SectionHeader
		if err = binary.Read(r, binary.LittleEndian, &sectionHeader); err != nil {
			return
		}
		sectionHeaders[headerNumber] = sectionHeader
	}

	if optionalHeader.CertificateTable.VirtualAddress == 0 {
		err = fmt.Errorf("Portable Executable does not have certificate data")
		return
	}

	var certEntryEnd = optionalHeader.CertificateTable.VirtualAddress + optionalHeader.CertificateTable.Size
	if certEntryEnd < optionalHeader.CertificateTable.VirtualAddress {

		err = fmt.Errorf("Overflow while calculating end of certificate entry")
		return
	}
	if certEntryEnd != uint32(len(stub)) {
		err = fmt.Errorf("Certificiate entry is not at end of file: {certEntryEnd} vs { fileReader.BaseStream.Length}")
		return
	}


	 offset = int(optionalHeader.CertificateTable.VirtualAddress)
	size = int(optionalHeader.CertificateTable.Size)
	sizeOffset = int(peOffset) + 4 + int(unsafe.Sizeof(structs.FileHeader{})) + int(fileHeader.SizeOfOptionalHeader) - 8 * (int(optionalHeader.NumberOfRvaAndSizes) - certificateTableIndex) + 4

	 if binary.LittleEndian.Uint32(stub[sizeOffset:]) != optionalHeader.CertificateTable.Size {
		err = errors.New("Internal error when calculating certificate data size offset.")
		return
	}

	return
}
