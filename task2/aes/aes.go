package aes

import (
	"USST-CyberSecurity/task2/util"
)

type WorkMode int

const (
	ECB WorkMode = iota
	CBC
	CTR
)

type Cipher struct {
	keyExp *KeyExpansion
	mode   WorkMode
	iv     [16]byte
}

// NewCipher creates a new cipher
func NewCipher(key []byte, iv []byte, mode WorkMode, keyType KeyType) *Cipher {
	c := &Cipher{
		keyExp: NewKeyExpansion(keyType, key),
		mode:   mode,
	}

	// in mode CTR, the IV is used as a nonce
	if mode == CBC || mode == CTR {
		if len(iv) != 16 {
			panic("Invalid IV length")
		}
		copy(c.iv[:], iv)
	}

	return c
}

// prepareBlock prepares the block for encryption
func prepareBlock(block []byte, decrypt bool) [][]byte {
	if len(block)%16 != 0 && decrypt {
		panic("Invalid block size")
	}

	if !decrypt {
		// pad the block
		block = util.Pad(block, 16)
	}

	numBlocks := len(block) / 16
	blocks := make([][]byte, numBlocks)
	for i := 0; i < numBlocks; i++ {
		blocks[i] = block[i*16 : (i+1)*16]
	}
	return blocks
}

// Encrypt encrypts the input
func (c *Cipher) Encrypt(input []byte) []byte {
	switch c.mode {
	case ECB:
		blocks := prepareBlock(input, false)
		return c.encryptECB(blocks)
	case CBC:
		blocks := prepareBlock(input, false)
		return c.encryptCBC(blocks)
	case CTR:
		// CTR mode does not require padding
		return c.encryptCTR(input)
	default:
		panic("Unsupported mode")
	}
}

// Decrypt decrypts the input
func (c *Cipher) Decrypt(input []byte) []byte {
	var result []byte
	switch c.mode {
	case ECB:
		blocks := prepareBlock(input, true)
		result = c.decryptECB(blocks)
		break
	case CBC:
		blocks := prepareBlock(input, true)
		result = c.decryptCBC(blocks)
		break
	case CTR:
		// CTR mode does not require padding
		return c.encryptCTR(input)
	default:
		panic("Unsupported mode")
	}
	return util.Unpad(result)
}
