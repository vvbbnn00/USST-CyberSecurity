package main

import (
	"USST-CyberSecurity/task1/vigenere"
	"flag"
	"fmt"
	"os"
	"strings"
)

// doDecryptWithoutKey decrypts the cipher without knowing the key
func doDecryptWithoutKey(cipher string) {
	// decrypt the cipher without knowing the key
	resultList := vigenere.MultiThreadDecrypt(cipher)

	// print the result (only the first 100 characters)
	cipher = strings.ReplaceAll(cipher, "\n", " ")
	if len(cipher) > 100 {
		cipher = cipher[:100] + "..."
	}

	for _, result := range resultList {
		fmt.Printf("Key: %s, Confidence: %.10f\n%s\n",
			result.Key,
			result.Possibility,
			vigenere.Decrypt(cipher, result.Key))
	}
}

// doDecrypt decrypts the cipher with the given key
func doDecrypt(cipher string, key string) {
	fmt.Println(vigenere.Decrypt(cipher, key))
}

// doEncrypt encrypts the plain text with the given key
func doEncrypt(plain string, key string) {
	fmt.Println(vigenere.Encrypt(plain, key))
}

// readText reads text from stdin until EOF
func readText() string {
	var text string
	for {
		var s string
		_, err := fmt.Scan(&s)
		if err != nil {
			break
		}
		text += s + " "
	}
	return text
}

func main() {
	var encryptMode bool
	var decryptMode bool
	var key string

	const usage = "Usage: cat cipher.txt | %s -d/-e -k [key]\n"

	flag.BoolVar(&encryptMode, "e", false, "Enable encryption mode")
	flag.BoolVar(&decryptMode, "d", false, "Enable decryption mode")
	flag.StringVar(&key, "k", "", "Encryption/Decryption key")

	flag.Parse()

	if encryptMode && decryptMode {
		fmt.Println("Error: Cannot specify both encryption and decryption modes")
		os.Exit(1)
	}

	if encryptMode {
		if key == "" {
			fmt.Println("Error: Encryption mode requires a key (-k)")
			os.Exit(1)
		}
		plainText := readText()
		doEncrypt(plainText, key)
	} else if decryptMode {
		cipherText := readText()
		if key != "" {
			doDecrypt(cipherText, key)
		} else {
			doDecryptWithoutKey(cipherText)
		}
	} else {
		fmt.Printf(usage, os.Args[0])
		os.Exit(1)
	}
}
