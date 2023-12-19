package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func main() {
	// Replace "yourSecretKey" with your actual secret key
	accessKey := ""
	key := []byte("0123456789012345")

	encryptedAccessKey, err := encrypt(key, accessKey)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}

	fmt.Println("Encrypted Access Key:", encryptedAccessKey)

	decryptedAccessKey, err := decrypt(key, encryptedAccessKey)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}

	fmt.Println("Decrypted Access Key:", decryptedAccessKey)
}

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	// Use a fixed IV instead of rand.Reader
	copy(iv, key[:aes.BlockSize])

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decrypt(key []byte, ciphertext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := decodedCiphertext[:aes.BlockSize]
	decodedCiphertext = decodedCiphertext[aes.BlockSize:]

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(decodedCiphertext, decodedCiphertext)

	return string(decodedCiphertext), nil
}
