package structs

//For Enum like formatting.
type DataSectionFlags uint32

/// <summary>
/// Represents the IMAGE_SECTION_HEADER PE section
/// </summary>
type SectionHeader struct {
	Name                 [8]byte
	VirtualSize          uint32
	VirtualAddress       uint32
	SizeOfRawData        uint32
	PointerToRawData     uint32
	PointerToRelocations uint32
	PointerToLinenumbers uint32
	NumberOfRelocations  uint16
	NumberOfLinenumbers  uint16
	Characteristics      DataSectionFlags
}


const (
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeReg DataSectionFlags = 0x00000000
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeDsect DataSectionFlags = 0x00000001
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeNoLoad DataSectionFlags = 0x00000002
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeGroup DataSectionFlags = 0x00000004
	/// <summary>
	/// The section should not be padded to the next boundary. This flag is obsolete and is replaced by IMAGE_SCN_ALIGN_1BYTES. This is valid only for object files.
	/// </summary>
	TypeNoPadded DataSectionFlags = 0x00000008
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeCopy DataSectionFlags = 0x00000010
	/// <summary>
	/// The section contains executable code.
	/// </summary>
	ContentCode DataSectionFlags = 0x00000020
	/// <summary>
	/// The section contains initialized data.
	/// </summary>
	ContentInitializedData DataSectionFlags = 0x00000040
	/// <summary>
	/// The section contains uninitialized data.
	/// </summary>
	ContentUninitializedData DataSectionFlags = 0x00000080
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	LinkOther DataSectionFlags = 0x00000100
	/// <summary>
	/// The section contains comments or other information. The .drectve section has this type. This is valid for object files only.
	/// </summary>
	LinkInfo DataSectionFlags = 0x00000200
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	TypeOver DataSectionFlags = 0x00000400
	/// <summary>
	/// The section will not become part of the image. This is valid only for object files.
	/// </summary>
	LinkRemove DataSectionFlags = 0x00000800
	/// <summary>
	/// The section contains COMDAT data. For more information see section 5.5.6 COMDAT Sections (Object Only). This is valid only for object files.
	/// </summary>
	LinkComDat DataSectionFlags = 0x00001000
	/// <summary>
	/// Reset speculative exceptions handling bits in the TLB entries for this section.
	/// </summary>
	NoDeferSpecExceptions DataSectionFlags = 0x00004000
	/// <summary>
	/// The section contains data referenced through the global pointer (GP).
	/// </summary>
	RelativeGP DataSectionFlags = 0x00008000
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	MemPurgeable DataSectionFlags = 0x00020000
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	Memory16Bit DataSectionFlags = 0x00020000
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	MemoryLocked DataSectionFlags = 0x00040000
	/// <summary>
	/// Reserved for future use.
	/// </summary>
	MemoryPreload DataSectionFlags = 0x00080000
	/// <summary>
	/// Align data on a 1-byte boundary. Valid only for object files.
	/// </summary>
	Align1Bytes DataSectionFlags = 0x00100000
	/// <summary>
	/// Align data on a 2-byte boundary. Valid only for object files.
	/// </summary>
	Align2Bytes DataSectionFlags = 0x00200000
	/// <summary>
	/// Align data on a 4-byte boundary. Valid only for object files.
	/// </summary>
	Align4Bytes DataSectionFlags = 0x00300000
	/// <summary>
	/// Align data on an 8-byte boundary. Valid only for object files.
	/// </summary>
	Align8Bytes DataSectionFlags = 0x00400000
	/// <summary>
	/// Align data on a 16-byte boundary. Valid only for object files.
	/// </summary>
	Align16Bytes DataSectionFlags = 0x00500000
	/// <summary>
	/// Align data on a 32-byte boundary. Valid only for object files.
	/// </summary>
	Align32Bytes DataSectionFlags = 0x00600000
	/// <summary>
	/// Align data on a 64-byte boundary. Valid only for object files.
	/// </summary>
	Align64Bytes DataSectionFlags = 0x00700000
	/// <summary>
	/// Align data on a 128-byte boundary. Valid only for object files.
	/// </summary>
	Align128Bytes DataSectionFlags = 0x00800000
	/// <summary>
	/// Align data on a 256-byte boundary. Valid only for object files.
	/// </summary>
	Align256Bytes DataSectionFlags = 0x00900000
	/// <summary>
	/// Align data on a 512-byte boundary. Valid only for object files.
	/// </summary>
	Align512Bytes DataSectionFlags = 0x00A00000
	/// <summary>
	/// Align data on a 1024-byte boundary. Valid only for object files.
	/// </summary>
	Align1024Bytes DataSectionFlags = 0x00B00000
	/// <summary>
	/// Align data on a 2048-byte boundary. Valid only for object files.
	/// </summary>
	Align2048Bytes DataSectionFlags = 0x00C00000
	/// <summary>
	/// Align data on a 4096-byte boundary. Valid only for object files.
	/// </summary>
	Align4096Bytes DataSectionFlags = 0x00D00000
	/// <summary>
	/// Align data on an 8192-byte boundary. Valid only for object files.
	/// </summary>
	Align8192Bytes DataSectionFlags = 0x00E00000
	/// <summary>
	/// The section contains extended relocations.
	/// </summary>
	LinkExtendedRelocationOverflow DataSectionFlags = 0x01000000
	/// <summary>
	/// The section can be discarded as needed.
	/// </summary>
	MemoryDiscardable DataSectionFlags = 0x02000000
	/// <summary>
	/// The section cannot be cached.
	/// </summary>
	MemoryNotCached DataSectionFlags = 0x04000000
	/// <summary>
	/// The section is not pageable.
	/// </summary>
	MemoryNotPaged DataSectionFlags = 0x08000000
	/// <summary>
	/// The section can be shared in memory.
	/// </summary>
	MemoryShared DataSectionFlags = 0x10000000
	/// <summary>
	/// The section can be executed as code.
	/// </summary>
	MemoryExecute DataSectionFlags = 0x20000000
	/// <summary>
	/// The section can be read.
	/// </summary>
	MemoryRead DataSectionFlags = 0x40000000
	/// <summary>
	/// The section can be written to.
	/// </summary>
	MemoryWrite DataSectionFlags = 0x80000000
)