package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateRandomSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func encryptWithRandomSalt(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate a random salt
	salt, err := generateRandomSalt(aesGCM.NonceSize())
	if err != nil {
		return "", err
	}

	// Append the salt to the ciphertext
	ciphertext := aesGCM.Seal(nil, salt, plaintext, nil)

	// Combine salt and ciphertext into a hex-encoded string
	result := hex.EncodeToString(append(salt, ciphertext...))
	return result, nil
}

func decryptWithRandomSalt(ciphertextHex string, key []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Extract the salt from the ciphertext
	saltSize := aesGCM.NonceSize()
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt, ciphertext := ciphertext[:saltSize], ciphertext[saltSize:]

	// Decrypt the data using the salt
	plaintext, err := aesGCM.Open(nil, salt, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
