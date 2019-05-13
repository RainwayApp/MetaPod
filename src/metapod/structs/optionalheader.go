package structs


/// <summary>
/// Represents the IMAGE_OPTIONAL_HEADER32 structure format <see href="https://docs.microsoft.com/en-us/windows/desktop/api/winnt/ns-winnt-_image_optional_header">HERE</see>
/// </summary>
type OptionalHeader32 struct {
	/// <summary>
	/// The state of the image file. This member can be one of the following values.
	/// </summary>
	Magic                       uint16
	/// <summary>
	/// The major version number of the linker.
	/// </summary>
	MajorLinkerVersion          uint8
	/// <summary>
	/// The minor version number of the linker.
	/// </summary>
	MinorLinkerVersion          uint8
	/// <summary>
	/// The size of the code section, in bytes, or the sum of all such sections if there are multiple code sections.
	/// </summary>
	SizeOfCode                  uint32
	/// <summary>
	/// The size of the initialized data section, in bytes, or the sum of all such sections if there are multiple initialized data sections.
	/// </summary>
	SizeOfInitializedData       uint32
	/// <summary>
	/// The size of the uninitialized data section, in bytes, or the sum of all such sections if there are multiple uninitialized data sections.
	/// </summary>
	SizeOfUninitializedData     uint32
	/// <summary>
	/// A pointer to the entry point function, relative to the image base address.
	/// For executable files, this is the starting address. For device drivers, this is the address of the initialization function.
	/// The entry point function is optional for DLLs. When no entry point is present, this member is zero.
	/// </summary>
	AddressOfEntryPoint         uint32
	/// <summary>
	/// A pointer to the beginning of the code section, relative to the image base.
	/// </summary>
	BaseOfCode                  uint32
	/// <summary>
	/// A pointer to the beginning of the data section, relative to the image base.
	/// </summary>
	BaseOfData                  uint32
	/// <summary>
	/// The preferred address of the first byte of the image when it is loaded in memory.
	// This value is a multiple of 64K bytes. The default value for DLLs is 0x10000000.
	// The default value for applications is 0x00400000, except on Windows CE where it is 0x00010000.
	/// </summary>
	ImageBase                   uint32
	/// <summary>
	/// The alignment of sections loaded in memory, in bytes.
	// This value must be greater than or equal to the FileAlignment member.
	// The default value is the page size for the system.
	/// </summary>
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint32
	SizeOfStackCommit           uint32
	SizeOfHeapReserve           uint32
	SizeOfHeapCommit            uint32
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32

	ExportTable           DataDirectory
	ImportTable           DataDirectory
	ResourceTable         DataDirectory
	ExceptionTable        DataDirectory
	CertificateTable      DataDirectory
	BaseRelocationTable   DataDirectory
	Debug                 DataDirectory
	Architecture          DataDirectory
	GlobalPtr             DataDirectory
	TLSTable              DataDirectory
	LoadConfigTable       DataDirectory
	BoundImport           DataDirectory
	IAT                   DataDirectory
	DelayImportDescriptor DataDirectory
	CLRRuntimeHeader      DataDirectory
	Reserved              DataDirectory
}

/// <summary>
/// Represents the IMAGE_OPTIONAL_HEADER64 structure format <see href="https://docs.microsoft.com/en-us/windows/desktop/api/winnt/ns-winnt-image_optional_header64">HERE</see>
/// </summary>
type OptionalHeader64 struct {
	Magic                       uint16
	MajorLinkerVersion          uint8
	MinorLinkerVersion          uint8
	SizeOfCode                  uint32
	SizeOfInitializedData       uint32
	SizeOfUninitializedData     uint32
	AddressOfEntryPoint         uint32
	BaseOfCode                  uint32
	BaseOfData                  uint32
	ImageBase                   uint64
	SectionAlignment            uint32
	FileAlignment               uint32
	MajorOperatingSystemVersion uint16
	MinorOperatingSystemVersion uint16
	MajorImageVersion           uint16
	MinorImageVersion           uint16
	MajorSubsystemVersion       uint16
	MinorSubsystemVersion       uint16
	Win32VersionValue           uint32
	SizeOfImage                 uint32
	SizeOfHeaders               uint32
	CheckSum                    uint32
	Subsystem                   uint16
	DllCharacteristics          uint16
	SizeOfStackReserve          uint64
	SizeOfStackCommit           uint64
	SizeOfHeapReserve           uint64
	SizeOfHeapCommit            uint64
	LoaderFlags                 uint32
	NumberOfRvaAndSizes         uint32

	ExportTable           DataDirectory
	ImportTable           DataDirectory
	ResourceTable         DataDirectory
	ExceptionTable        DataDirectory
	CertificateTable      DataDirectory
	BaseRelocationTable   DataDirectory
	Debug                 DataDirectory
	Architecture          DataDirectory
	GlobalPtr             DataDirectory
	TLSTable              DataDirectory
	LoadConfigTable       DataDirectory
	BoundImport           DataDirectory
	IAT                   DataDirectory
	DelayImportDescriptor DataDirectory
	CLRRuntimeHeader      DataDirectory
	Reserved              DataDirectory
}
