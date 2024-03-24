package aes

// encryptECB encrypts the blocks using ECB mode
func (c *Cipher) encryptECB(blocks [][]byte) []byte {
	result := make([]byte, len(blocks)*16)
	for i, block := range blocks {
		encrypted := EncryptBlock(MatrixFromBytes(block), c.keyExp)
		copy(result[i*16:], encrypted.Bytes())
	}
	return result
}

// decryptECB decrypts the blocks using ECB mode
func (c *Cipher) decryptECB(blocks [][]byte) []byte {
	result := make([]byte, len(blocks)*16)
	for i, block := range blocks {
		decrypted := DecryptBlock(MatrixFromBytes(block), c.keyExp)
		copy(result[i*16:], decrypted.Bytes())
	}
	return result
}
