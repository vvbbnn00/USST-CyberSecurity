package aes

// encryptCBC encrypts the blocks using CBC mode
func (c *Cipher) encryptCBC(blocks [][]byte) []byte {
	result := make([]byte, len(blocks)*16)
	prev := c.iv
	for i, block := range blocks {
		// XOR the block with the previous block
		for j := 0; j < 16; j++ {
			block[j] ^= prev[j]
		}

		// encrypt the block
		encrypted := EncryptBlock(MatrixFromBytes(block), c.keyExp)
		copy(result[i*16:], encrypted.Bytes())

		// save the encrypted block for the next iteration
		copy(prev[:], encrypted.Bytes())
	}
	return result
}

// decryptCBC decrypts the blocks using CBC mode
func (c *Cipher) decryptCBC(blocks [][]byte) []byte {
	result := make([]byte, len(blocks)*16)
	prev := c.iv
	for i, block := range blocks {
		// save the block for the next iteration
		next := make([]byte, 16)
		copy(next, block)

		// decrypt the block
		decrypted := DecryptBlock(MatrixFromBytes(block), c.keyExp)
		decryptedBytes := decrypted.Bytes()

		// XOR the decrypted block with the previous block
		for j := 0; j < 16; j++ {
			decryptedBytes[j] ^= prev[j]
		}
		copy(result[i*16:], decryptedBytes)

		// save the encrypted block for the next iteration
		copy(prev[:], next)
	}
	return result
}
