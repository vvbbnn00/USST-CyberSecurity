package main

import (
	"USST-CyberSecurity/task2/aes"
	"USST-CyberSecurity/task2/util"
	"fmt"
)

func main() {
	testPlaintext := util.BytesFromHex("1fd4ee65603e6130cfc2a82ab3d56c241fd4ee65603e6130cfc2a82ab3d56c24")
	testKey := util.BytesFromHex("8809e7dd3a959ee5d8dbb13f501f2274")
	testIV := util.BytesFromHex("ffffffffffffffffffffffffffffffff")

	// create a new cipher
	cipher := aes.NewCipher(testKey, testIV, aes.CTR, aes.Key128)
	data := cipher.Encrypt(testPlaintext)
	fmt.Printf("Encrypted: %x\n", data)

	// decrypt the data
	decrypted := cipher.Decrypt(data)
	fmt.Printf("Decrypted: %x\n", decrypted)
}
