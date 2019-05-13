package main

import "C"

import (
	"fmt"
	"metapod/windows"
	"os"
	"reflect"
	"unsafe"
)

//export Create
//Creates a new executable containing the payload, based upon a given stub template.
//Once complete, it returns a byte array containing the new Portable Executable.
//TODO raise exceptions to the higher level langauge.
func Create(buffer unsafe.Pointer, count C.int, payload *C.char, output *unsafe.Pointer) C.int {
	slice := &reflect.SliceHeader{Data: uintptr(buffer), Len: int(count), Cap: int(count)}
	stub := *(*[]byte)(unsafe.Pointer(slice))
	portableExecutable, err := windows.GetPortableExecutable(stub)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 0
	}
	var payloadContents = []byte(C.GoString(payload))
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	contents, err := targetExecutable.CreateFromTemplate(payloadContents)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while setting superfluous certificate tag: %s\n", err)
		return 0
	}
	*output = C.CBytes(contents)

	return C.int(len(contents))
}

//export Open
//Opens a portable executable file from a byte stream, seeking to find a payload certificate.
//If found, the payload will be returned as a string -- otherwise an empty string is returned.
func Open( pe unsafe.Pointer, count C.int) *C.char {
	slice := &reflect.SliceHeader{Data: uintptr(pe), Len: int(count), Cap: int(count)}
	buffer := *(*[]byte)(unsafe.Pointer(slice))
	portableExecutable, err := windows.GetPortableExecutable(buffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return C.CString("")
	}
	var targetExecutable = windows.TargetExecutable{*portableExecutable}
	_, err, payload := targetExecutable.GetPayload()
	if err == nil && payload != nil {
		return C.CString(string(payload))
	}
	return C.CString("")
}



//TODO
//proper error handling
//refactor code some more
//comment code
func main() { }
