package main

import "C"

import (
	"unsafe"

	"github.com/RainwayApp/metapod/src/metapod/windows"
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
	return C.int(0)
}

//export Open
//Opens a portable executable file from a byte stream, seeking to find a payload certificate.
//If found, the payload will be returned as a string.
//If the error code is greater than zero, an issue was encountered.
//Use GetErrorCodeMessage to retrieve the error message.
func Open(pe unsafe.Pointer, count C.int, payload **C.char, payloadCount *C.int) C.int {
	buffer := C.GoBytes(pe, count)
	portableExecutable, err := windows.GetPortableExecutable(buffer)
	if err > 0 {
		return C.int(err)
	}
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	_, err, rawPayload := targetExecutable.GetPayload()
	if err > 0 {
		return C.int(err)
	}
	if payload != nil {
		*payload = C.CString(string(rawPayload))
		*payloadCount = C.int(len(rawPayload))
		return C.int(0)
	} else {
		return C.int(1050)
	}
}

//export GetErrorCodeMessage
//Returns the human readable error message for a given error code.
func GetErrorCodeMessage(code C.int, text **C.char) C.int {
	var errorText string
	switch code {
	case 1050:
		errorText = "unable to locate payload within input file"
	case 1043:
		errorText = "the input file contains no certificates"
	case 1042:
		errorText = "failed to marshal ASN.1 structure"
	case 1041:
		errorText = "failed to create X509Certificate from provided templates"
	case 1040:
		errorText = "failed to generate RSA keypair"
	case 1033:
		errorText = "internal error calculating certificate data offset"
	case 1032:
		errorText = "certificate table entry is not at the end of the file"
	case 1031:
		errorText = "reached EOF calculating end of certificate entry"
	case 1030:
		errorText = "portable executable lacks valid certificate data entry"
	case 1029:
		errorText = "unable to read IMAGE_SECTION_HEADER. EOF?"
	case 1028:
		errorText = "unable to read IMAGE_OPTIONAL_HEADER32"
	case 1027:
		errorText = "input file must be 32-bit. 64-bit support planned"
	case 1026:
		errorText = "input file cannot be a DLL"
	case 1025:
		errorText = "input file is not a valid portable executable"
	case 1024:
		errorText = "unable to read IMAGE_FILE_HEADER"
	case 1023:
		errorText = "unable to locate portable executable file header"
	case 1022:
		errorText = "portable executable is malformed"
	case 1021:
		errorText = "reached EOF searching for the portable executable signature"
	case 1020:
		errorText = "the length of the input file is less than the PE offset"
	case 1004:
		errorText = "incorrect number of bytes reading ASN.1 length"
	case 1005:
		errorText = "ASN.1 structure is incorrect"
	case 1006:
		errorText = "unknown certificate type"
	case 1007:
		errorText = "unknown certificate revision"
	case 1008:
		errorText = "multiple attribute certificates found. unable to proceed."
	case 1009:
		errorText = "attribute certificate seems malformed"
	default:
		errorText = "unknown error code."
	}

	*text = C.CString(errorText)
	return C.int(len(errorText))
}

//TODO 64-bit support
func main() {}
