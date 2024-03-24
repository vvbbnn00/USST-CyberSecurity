package util

// fromHexChar converts a hex character to a byte
func fromHexChar(c byte) byte {
	switch {
	case '0' <= c && c <= '9':
		return c - '0'
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10
	default:
		panic("invalid hex character")
	}
}

// BytesFromHex converts a hex string to a byte slice
func BytesFromHex(hex string) []byte {
	if len(hex)%2 != 0 {
		panic("invalid hex length")
	}

	bytes := make([]byte, len(hex)/2)
	for i := 0; i < len(hex); i += 2 {
		bytes[i/2] = byte((fromHexChar(hex[i]) << 4) | fromHexChar(hex[i+1]))
	}
	return bytes
}
