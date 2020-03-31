package errors

func errorText(code int) string {
	switch code {
	case 1050:
		return "unable to locate payload within input file"
	case 1043:
		return "the input file contains no certificates"
	case 1042:
		return "failed to marshal ASN.1 structure"
	case 1041:
		return "failed to create X509Certificate from provided templates"
	case 1040:
		return "failed to generate RSA keypair"
	case 1033:
		return "internal error calculating certificate data offset"
	case 1032:
		return "certificate table entry is not at the end of the file"
	case 1031:
		return "reached EOF calculating end of certificate entry"
	case 1030:
		return "portable executable lacks valid certificate data entry"
	case 1029:
		return "unable to read IMAGE_SECTION_HEADER. EOF?"
	case 1028:
		return "unable to read IMAGE_OPTIONAL_HEADER32"
	case 1027:
		return "input file must be 32-bit. 64-bit support planned"
	case 1026:
		return "input file cannot be a DLL"
	case 1025:
		return "input file is not a valid portable executable"
	case 1024:
		return "unable to read IMAGE_FILE_HEADER"
	case 1023:
		return "unable to locate portable executable file header"
	case 1022:
		return "portable executable is malformed"
	case 1021:
		return "reached EOF searching for the portable executable signature"
	case 1020:
		return "the length of the input file is less than the PE offset"
	case 1004:
		return "incorrect number of bytes reading ASN.1 length"
	case 1005:
		return "ASN.1 structure is incorrect"
	case 1006:
		return "unknown certificate type"
	case 1007:
		return "unknown certificate revision"
	case 1008:
		return "multiple attribute certificates found. unable to proceed."
	case 1009:
		return "attribute certificate seems malformed"
	}
	return "unknown error code."
}

type MetapodError struct {
	err int
}

func NewError(code int) MetapodError {
	return MetapodError{code}
}

func (m MetapodError) Error() string {
	return errorText(m.err)
}

func (m MetapodError) ErrCode() int {
	return m.err
}
