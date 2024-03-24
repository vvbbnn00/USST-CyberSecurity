package aes

// Counter is a counter for CTR mode
type Counter struct {
	iv    []byte
	count uint64
}

// NewCounter creates a new counter
func NewCounter(iv []byte) *Counter {
	return &Counter{
		iv:    iv,
		count: 0,
	}
}

// Next returns the next counter value
func (c *Counter) Next() []byte {
	result := make([]byte, 16)
	copy(result, c.iv)
	result[15] += byte(c.count)
	c.count++
	return result
}

func (c *Cipher) encryptCTR(input []byte) []byte {
	result := make([]byte, len(input))
	counter := NewCounter(c.iv[:])

	for i := 0; i < len(input); i += 16 {
		var block []byte
		if i+16 > len(input) {
			block = input[i:]
		} else {
			block = input[i : i+16]
		}

		// encrypt the counter, then XOR it with the block
		encrypted := EncryptBlock(MatrixFromBytes(counter.Next()), c.keyExp)

		for j := 0; j < 16; j++ {
			if i+j >= len(input) {
				break
			}
			result[i+j] = block[j] ^ encrypted.Bytes()[j]
		}
	}

	return result
}

func (c *Cipher) decryptCTR(input []byte) []byte {
	return c.encryptCTR(input)
}
