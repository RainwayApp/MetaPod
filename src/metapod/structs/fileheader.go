package structs

/// <summary>
/// Represents the IMAGE_FILE_HEADER structure from <see href="https://docs.microsoft.com/en-us/windows/desktop/api/winnt/ns-winnt-_image_file_header">here</see>
/// </summary>
type FileHeader struct {
	/// <summary>
	/// The architecture type of the computer. An image file can only be run on the specified computer or a system that emulates the specified computer. This member can be one of the following values.
	/// </summary>
	Machine               uint16
	/// <summary>
	/// The number of sections. This indicates the size of the section table, which immediately follows the headers. Note that the Windows loader limits the number of sections to 96.
	/// </summary>
	NumberOfSections      uint16
	/// <summary>
	/// The low 32 bits of the time stamp of the image. This represents the date and time the image was created by the linker.
	/// The value is represented in the number of seconds elapsed since midnight (00:00:00), January 1, 1970, Universal Coordinated Time, according to the system clock.
	/// </summary>
	TimeDateStamp         uint32
	/// <summary>
	/// The offset of the symbol table, in bytes, or zero if no COFF symbol table exists.
	/// </summary>
	PointerForSymbolTable uint32
	/// <summary>
	/// The number of symbols in the symbol table.
	/// </summary>
	NumberOfSymbols       uint32
	/// <summary>
	/// The size of the optional header, in bytes. This value should be 0 for object files.
	/// </summary>
	SizeOfOptionalHeader  uint16
	/// <summary>
	/// The characteristics of the image. This member can be one or more of the following values.
	/// </summary>
	Characteristics       uint16
}

//determines if a PE is based around 32x architecture.
func (fileHeader *FileHeader) Is32Bit() bool {
	const imageFile32BitMachine = uint16(0x0100)
	return (imageFile32BitMachine & fileHeader.Characteristics) == imageFile32BitMachine
}

//determines if a PE is a DLL.
func (fileHeader *FileHeader) IsDll() bool {
	return (fileHeader.Characteristics & uint16(0x2000)) > 0
}

//determines if a PE is a valid EXE
func (fileHeader *FileHeader) IsExe() bool {
	return (fileHeader.Characteristics & uint16(2)) > 0
}