package util

// Pad adds padding to the input slice according to PKCS7 standard.
func Pad(input []byte, blockSize int) []byte {
	padding := blockSize - (len(input) % blockSize)
	for i := 0; i < padding; i++ {
		input = append(input, byte(padding))
	}
	return input
}

// Unpad removes the PKCS7 padding from the input slice.
func Unpad(input []byte) []byte {
	if len(input) == 0 {
		return input
	}

	padding := int(input[len(input)-1])
	if padding > len(input) || padding == 0 {
		return input
	}

	for i := len(input) - padding; i < len(input); i++ {
		if input[i] != byte(padding) {
			return input
		}
	}

	return input[:len(input)-padding]
}
