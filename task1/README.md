## Task1 Vigenere Cipher Implementation

- [x] Implement Vigenere Cipher
- [x] Implement Encryption
- [x] Implement Decryption
- [x] Implement Decryption without key

### Usage

#### Windows

```shell
del vigenere.exe
go build -o vigenere.exe
# Remove .txt files in output directory
del output\*.txt
# Encryption
type ./test/plain.txt | .\vigenere.exe -e -k "key" > output/test1_encrypted.txt
# Decryption Without Key
type ./output/test1_encrypted.txt | .\vigenere.exe -d
# Decryption With Key
type ./output/test1_encrypted.txt | .\vigenere.exe -d -k "key"

# Other tests
type ./test/cipher1.txt | .\vigenere.exe -d
type ./test/cipher2.txt | .\vigenere.exe -d
```
