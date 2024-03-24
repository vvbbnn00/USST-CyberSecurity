package test

import (
	"USST-CyberSecurity/task2/aes"
	"bufio"
	"os"
	"strings"
)

// Task represents a single encryption or decryption task.
type Task struct {
	Key         string
	IV          string
	Plaintext   string
	Ciphertext  string
	TaskType    string // "ENCRYPT" or "DECRYPT"
	EncryptMode string
}

// ParseRSPFile parses the given .rsp file and returns a list of tasks.
func ParseRSPFile(filename, encryptMode string) ([]Task, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tasks []Task
	var currentTask *Task
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "[ENCRYPT]") {
			currentTask = &Task{TaskType: "ENCRYPT", EncryptMode: encryptMode}
		} else if strings.HasPrefix(line, "[DECRYPT]") {
			currentTask = &Task{TaskType: "DECRYPT", EncryptMode: encryptMode}
		} else if strings.HasPrefix(line, "COUNT") {
			if currentTask != nil {
				tasks = append(tasks, *currentTask)
				currentTask = &Task{TaskType: currentTask.TaskType, EncryptMode: encryptMode}
			}
		} else if strings.HasPrefix(line, "KEY") {
			if currentTask != nil {
				currentTask.Key = strings.Split(line, " = ")[1]
			}
		} else if strings.HasPrefix(line, "IV") {
			if currentTask != nil {
				currentTask.IV = strings.Split(line, " = ")[1]
			}
		} else if strings.HasPrefix(line, "PLAINTEXT") {
			if currentTask != nil {
				currentTask.Plaintext = strings.Split(line, " = ")[1]
			}
		} else if strings.HasPrefix(line, "CIPHERTEXT") {
			if currentTask != nil {
				currentTask.Ciphertext = strings.Split(line, " = ")[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if currentTask != nil {
		tasks = append(tasks, *currentTask)
	}

	return tasks, nil
}

// GetKeyType returns the key type based on the length of the key.
func GetKeyType(key []byte) aes.KeyType {
	switch len(key) {
	case 16:
		return aes.Key128
	case 24:
		return aes.Key192
	case 32:
		return aes.Key256
	default:
		panic("Invalid key size")
	}
}
