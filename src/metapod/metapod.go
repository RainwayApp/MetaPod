package main

import "C"

import (
	"metapod/windows"
	"reflect"
	"unsafe"
)

//export Create
//Creates a new executable containing the payload, based upon a given stub template.
//Once complete, it returns a byte array containing the new Portable Executable.
//If the error code is greater than zero, an issue was encountered.
//Use GetErrorCodeMessage to retrieve the error message.
func Create(buffer unsafe.Pointer, count C.int, payload *C.char, output *unsafe.Pointer, outputCount *C.int) C.int {
	stub := C.GoBytes(buffer, count)
	portableExecutable, err := windows.GetPortableExecutable(stub)
	if err > 0 {
		return C.int(err)
	}
	var payloadContents = []byte(C.GoString(payload))
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	contents, err := targetExecutable.CreateFromTemplate(payloadContents)
	if err > 0 {
		return C.int(err)
	}
	*output = C.CBytes(contents)
	*outputCount = C.int(len(contents))
	return 0
}

//export Open
//Opens a portable executable file from a byte stream, seeking to find a payload certificate.
//If found, the payload will be returned as a string.
//If the error code is greater than zero, an issue was encountered.
//Use GetErrorCodeMessage to retrieve the error message.
func Open(pe unsafe.Pointer, count C.int, payload *C.char payloadCount *C.int) C.int {
	buffer := C.GoBytes(pe, count)
	portableExecutable, err := windows.GetPortableExecutable(buffer)
	if err > 0 {
		return C.int(err)
	}
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	_, err, payload := targetExecutable.GetPayload()
	if (err > 0) {
		return C.int(err)
	}
	if payload != nil {
		*payload = C.CString(string(payload))
		*payloadCount = C.int(len(payload))
		return C.Int(0)
	} else {
		return C.int(1050)
	}
}

//export GetErrorCodeMessage
//Returns the human readable error message for a given error code.
func GetErrorCodeMessage(code C.int) *C.char {
	switch code {
	case 1050:
		return C.CString("unable to locate payload within input file")
	case 1043:
		return C.CString("the input file contains no certificates")
	case 1042:
		return C.CString("failed to marshal ASN.1 structure")
	case 1041:
		return C.CString("failed to create X509Certificate from provided templates")
	case 1040:
		return C.CString("failed to generate RSA keypair")
	case 1033:
		return C.CString("internal error calculating certificate data offset")
	case 1032:
		return C.CString("certificate table entry is not at the end of the file")
	case 1031:
		return C.CString("reached EOF calculating end of certificate entry")
	case 1030:
		return C.CString("portable executable lacks valid certificate data entry")
	case 1029:
		return C.CString("unable to read IMAGE_SECTION_HEADER. EOF?")
	case 1028:
		return C.CString("unable to read IMAGE_OPTIONAL_HEADER32")
	case 1027:
		return C.CString("input file must be 32-bit. 64-bit support planned")
	case 1026:
		return C.CString("input file cannot be a DLL")
	case 1025:
		return C.CString("input file is not a valid portable executable")
	case 1024:
		return C.CString("unable to read IMAGE_FILE_HEADER")
	case 1023:
		return C.CString("unable to locate portable executable file header")
	case 1022:
		return C.CString("portable executable is malformed")
	case 1021:
		return C.CString("reached EOF searching for the portable executable signature")
	case 1004:
		return C.CString("incorrect number of bytes reading ASN.1 length")
	case 1005:
		return C.CString("ASN.1 structure is incorrect")
	case 1006:
		return C.CString("unknown certificate type")
	case 1007:
		return C.CString("unknown certificate revision")
	case 1008:
		return C.CString("multiple attribute certificates found. unable to proceed.")
	case 1009:
		return C.CString("attribute certificate seems malformed")
	default:
		return C.CString("unknown error code.")
	}
}

//TODO 64-bit support
func main() { }
