package structs

/// <summary>
/// Represents the data directory <see href="https://docs.microsoft.com/en-us/windows/desktop/api/winnt/ns-winnt-_image_data_directory">here</see>
/// </summary>
type DataDirectory struct {
	/// <summary>
	/// The relative virtual address of the table.
	/// </summary>
	VirtualAddress uint32
	/// <summary>
	/// The size of the table, in bytes.
	/// </summary>
	Size           uint32
}