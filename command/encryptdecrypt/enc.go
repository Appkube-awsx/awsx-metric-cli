package encryptdecrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var EncryptDecrypt = &cobra.Command{
	Use:   "encryptDecrypt",
	Short: "Encrypt and decrypt using AES cipher",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'encrypt' or 'decrypt' subcommands.")
	},
}

var encryptCmd = &cobra.Command{
	Use:   "encrypt [accessKey]",
	Short: "Encrypts an access key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accessKey := args[0]
		key := []byte("qwertyuioplkjhgfdsa1234!@MNB?>P)")

		encryptedAccessKey, err := encrypt(key, accessKey)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}

		fmt.Println("Encrypted Access Key:", encryptedAccessKey)
	},
}

var decryptCmd = &cobra.Command{
	Use:   "decrypt [encryptedAccessKey]",
	Short: "Decrypts an encrypted access key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encryptedAccessKey := args[0]
		key := []byte("qwertyuioplkjhgfdsa1234!@MNB?>P)")

		decryptedAccessKey, err := decrypt(key, encryptedAccessKey)
		if err != nil {
			fmt.Println("Decryption error:", err)
			return
		}

		fmt.Println("Decrypted Access Key:", decryptedAccessKey)
	},
}

func encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

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

func init() {
	EncryptDecrypt.AddCommand(encryptCmd)
	EncryptDecrypt.AddCommand(decryptCmd)
}

func main() {
	if err := EncryptDecrypt.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
