package vigenere

import (
	"USST-CyberSecurity/task1/util"
	"strings"
)

func encrypt(cipher, key string, factor int) string {
	keyLength := len(key)
	key = strings.ToLower(key)
	key = util.FilterWords(key)

	keyDeltas := make([]int, keyLength)
	for i, char := range key {
		keyDeltas[i] = int(char - 'a')
	}

	var text strings.Builder
	cnt := 0
	for _, char := range cipher {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')) {
			text.WriteRune(char)
			continue
		}

		delta := keyDeltas[cnt%keyLength] * factor // factor determines whether to encrypt or decrypt
		if delta < 0 {
			delta += 26
		}

		if char >= 'a' && char <= 'z' {
			text.WriteRune(rune((int(char-'a')+delta)%26 + 'a'))
		} else {
			text.WriteRune(rune((int(char-'A')+delta)%26 + 'A'))
		}

		cnt++
	}

	return text.String()
}

// Encrypt the plain text using the Vigenere cipher
func Encrypt(plain, key string) string {
	return encrypt(plain, key, 1)
}

// Decrypt the cipher text using the Vigenere cipher
func Decrypt(cipher, key string) string {
	return encrypt(cipher, key, -1)
}
