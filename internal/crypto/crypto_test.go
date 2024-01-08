package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	// key := []byte("my-secret-key-32-characters-long")
	key := []byte("1234567812345678")

	plaintext := []byte("Hello, AES-GCM with Random Salt!")

	ciphertextHex, err := encryptWithRandomSalt(plaintext, key)
	assert.Nil(t, err)

	t.Logf("ciphertext: %s", ciphertextHex)

	decryptedText, err := decryptWithRandomSalt(ciphertextHex, key)
	assert.Nil(t, err)

	t.Logf("plaintexttext: %s", decryptedText)
}
