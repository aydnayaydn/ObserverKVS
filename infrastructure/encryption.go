package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
	"os"
)

func ReadCipherKeyFromTextFile() (string, error) {
	file, err := os.Open("cipher.txt") // Must be 32 bytes
	if err != nil {
		return "", err
	}
	defer func() {
		if err = file.Close(); err != nil {
			return
		}
	}()

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(b), err
}

func Encrypt(plaintext string) (string, error) {
	cipherKey, err := ReadCipherKeyFromTextFile() // Must be 32 bytes
	if err != nil {
		return "", err
	}

	// Convert the plaintext and key to byte arrays
	plaintextBytes := []byte(plaintext)
	keyBytes := []byte(cipherKey)

	// Generate a new AES cipher block from the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Generate a new GCM cipher using the AES cipher block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Use a fixed nonce (IV) for the GCM cipher
	nonce := make([]byte, gcm.NonceSize())

	// Encrypt the plaintext using the GCM cipher
	ciphertext := gcm.Seal(nonce, nonce, plaintextBytes, nil)

	// Encode the ciphertext as a base64 string
	encryptedText := base64.StdEncoding.EncodeToString(ciphertext)

	return encryptedText, nil
}

func Decrypt(encryptedText string) (string, error) {
	cipherKey, err := ReadCipherKeyFromTextFile() // Must be 32 bytes
	if err != nil {
		return "", err
	}

	// Convert the encrypted text and key to byte arrays
	encryptedTextBytes, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	keyBytes := []byte(cipherKey)

	// Generate a new AES cipher block from the key
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// Generate a new GCM cipher using the AES cipher block
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract the nonce from the encrypted text
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := encryptedTextBytes[:nonceSize], encryptedTextBytes[nonceSize:]

	// Decrypt the ciphertext using the GCM cipher
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
