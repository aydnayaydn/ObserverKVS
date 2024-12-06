package infrastructure

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateAPIKey() (string, error) {
	// Generate a random byte slice with 32 bytes
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	// Encode the byte slice to base64 string
	apiKey := base64.StdEncoding.EncodeToString(key)

	return apiKey, nil
}
