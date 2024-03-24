package test

import (
	"USST-CyberSecurity/task2/aes"
	"USST-CyberSecurity/task2/util"
	"bytes"
	"fmt"
	"path"
	"testing"
)

var ecbFiles = []string{
	"ECBMMT128.rsp",
	"ECBMMT192.rsp",
	"ECBMMT256.rsp",
}

func TestECB(t *testing.T) {
	for _, file := range ecbFiles {
		t.Run(file, func(t *testing.T) {
			fmt.Printf("Running test for file %s\n", file)

			p := path.Join("rsp", file)
			tasks, err := ParseRSPFile(p, "ECB")
			if err != nil {
				t.Errorf("Error parsing file %s: %v", file, err)
				return
			}

			for _, task := range tasks {
				if task.Key == "" || task.Plaintext == "" || task.Ciphertext == "" {
					continue
				}
				key := util.BytesFromHex(task.Key)
				plaintext := util.BytesFromHex(task.Plaintext)
				ciphertext := util.BytesFromHex(task.Ciphertext)
				taskType := task.TaskType

				if taskType == "ENCRYPT" {
					cipher := aes.NewCipher(key, make([]byte, 0), aes.ECB, GetKeyType(key))
					encrypted := cipher.Encrypt(plaintext)

					// Truncate the encrypted data to the length of the ciphertext
					encrypted = encrypted[:len(ciphertext)]
					if !bytes.Equal(encrypted, ciphertext) {
						t.Errorf("Encryption failed for task %s", task.Plaintext)
					}
				} else {
					cipher := aes.NewCipher(key, make([]byte, 0), aes.ECB, GetKeyType(key))
					decrypted := cipher.Decrypt(ciphertext)
					if !bytes.Equal(decrypted, plaintext) {
						t.Errorf("Decryption failed for task %s", task.Plaintext)
					}
				}
			}
		})
	}
}
