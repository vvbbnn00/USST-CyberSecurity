package aes

type KeyType int

const (
	Key128 KeyType = iota
	Key192
	Key256
)

// KeyExpansion represents the key expansion
type KeyExpansion struct {
	keyType KeyType
	word    []uint32
}

// NewKeyExpansion creates a new key expansion
func NewKeyExpansion(keyType KeyType, key []byte) *KeyExpansion {
	keyLength := 16 + int(keyType)*8
	if len(key) != keyLength {
		panic("Invalid key length")
	}

	rounds := 10 + 2*int(keyType)
	totalWords := (rounds + 1) * 4

	ke := &KeyExpansion{
		keyType: keyType,
		word:    make([]uint32, totalWords),
	}

	// initialize the first words based on the key type
	num := keyLength / 4
	for i := 0; i < num; i++ {
		ke.word[i] = uint32(key[4*i])<<24 | uint32(key[4*i+1])<<16 | uint32(key[4*i+2])<<8 | uint32(key[4*i+3])
	}

	ke.expandKey()

	return ke
}

func (ke *KeyExpansion) expandKey() {
	num := 4 + 2*int(ke.keyType)
	total := (10 + 2*int(ke.keyType) + 1) * 4

	for i := num; i < total; i++ {
		temp := ke.word[i-1]
		//fmt.Printf("%x\n", temp)

		if i%num == 0 {
			temp = subWord(rotWord(temp)) ^ rcon[i/num-1]
		} else if ke.keyType == Key256 && i%num == 4 {
			temp = subWord(temp)
		}

		ke.word[i] = ke.word[i-num] ^ temp
	}

	//fmt.Printf("Expanded key: %v\n", ke.word)
}

// GetKey returns the key at the specified index
func (ke *KeyExpansion) GetKey(index int) Matrix {
	key := new(Matrix)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			key[j][i] = byte(ke.word[index*4+i] >> uint(24-8*j))
		}
	}
	return *key
}

func rotWord(word uint32) uint32 {
	return word<<8 | word>>24
}

func subWord(word uint32) uint32 {
	b1 := byte(word >> 24)
	b2 := byte(word >> 16 & 0xff)
	b3 := byte(word >> 8 & 0xff)
	b4 := byte(word & 0xff)

	s1 := sbox[b1>>4][b1&0x0f]
	s2 := sbox[b2>>4][b2&0x0f]
	s3 := sbox[b3>>4][b3&0x0f]
	s4 := sbox[b4>>4][b4&0x0f]

	return uint32(s1)<<24 | uint32(s2)<<16 | uint32(s3)<<8 | uint32(s4)
}
