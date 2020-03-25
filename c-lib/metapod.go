package clib

import "C"

import (
	"unsafe"

	"github.com/RainwayApp/metapod"
	"github.com/RainwayApp/metapod/windows"
)

//export Create
//Creates a new executable containing the payload, based upon a given stub template.
//Once complete, it returns a byte array containing the new Portable Executable.
//If the error code is greater than zero, an issue was encountered.
//Use GetErrorCodeMessage to retrieve the error message.
func Create(buffer unsafe.Pointer, count C.int, payload *C.char, output *unsafe.Pointer, outputCount *C.int) C.int {
	stub := C.GoBytes(buffer, count)
	var payloadContents = []byte(C.GoString(payload))

	result, err := metapod.Create(stub, payloadContents)

	if err != nil {
		return C.int(err.(MetapodError).code)
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
func Open(pe unsafe.Pointer, count C.int, payload **C.char, payloadCount *C.int) C.int {
	buffer := C.GoBytes(pe, count)
	var targetExecutable = windows.TargetExecutable{*portableExecutable}

	rawPayload, err := metapod.Open(buffer)

	if err != nil {
		return C.int(err.(MetapodError).code)
	} else if rawPayload == nil {
		return C.int(1050)
	}

	*payload = C.CString(string(rawPayload))
	*payloadCount = C.int(len(rawPayload))
	return C.int(0)
}

//export GetErrorCodeMessage
//Returns the human readable error message for a given error code.
func GetErrorCodeMessage(code C.int, text **C.char) C.int {
	err := metapod.MetapodError{code}

	errorText := err.Error()

	*text = C.CString(errorText)
	return C.int(len(errorText))
}
