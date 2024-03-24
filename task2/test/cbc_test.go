package test

import (
	"USST-CyberSecurity/task2/aes"
	"USST-CyberSecurity/task2/util"
	"bytes"
	"fmt"
	"path"
	"testing"
)

var cbcFiles = []string{
	"CBCMMT128.rsp",
	"CBCMMT192.rsp",
	"CBCMMT256.rsp",
}

func TestCBC(t *testing.T) {
	for _, file := range cbcFiles {
		t.Run(file, func(t *testing.T) {
			fmt.Printf("Running test for file %s\n", file)

			p := path.Join("rsp", file)
			tasks, err := ParseRSPFile(p, "CBC")
			if err != nil {
				t.Errorf("Error parsing file %s: %v", file, err)
				return
			}

			for _, task := range tasks {
				if task.Key == "" || task.IV == "" || task.Plaintext == "" || task.Ciphertext == "" {
					continue
				}
				key := util.BytesFromHex(task.Key)
				iv := util.BytesFromHex(task.IV)
				plaintext := util.BytesFromHex(task.Plaintext)
				ciphertext := util.BytesFromHex(task.Ciphertext)
				taskType := task.TaskType

				if taskType == "ENCRYPT" {
					cipher := aes.NewCipher(key, iv, aes.CBC, GetKeyType(key))
					encrypted := cipher.Encrypt(plaintext)

					// Truncate the encrypted data to the length of the ciphertext
					encrypted = encrypted[:len(ciphertext)]

					if !bytes.Equal(encrypted, ciphertext) {
						t.Errorf("Encryption failed for task %s", task.Plaintext)
					}
				} else {
					cipher := aes.NewCipher(key, iv, aes.CBC, GetKeyType(key))
					decrypted := cipher.Decrypt(ciphertext)
					if !bytes.Equal(decrypted, plaintext) {
						t.Errorf("Decryption failed for task %s", task.Plaintext)
					}
				}
			}
		})
	}
}
