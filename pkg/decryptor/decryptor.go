package decryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"strings"
)

func KeyFromMasterKey(masterKey string) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(masterKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 master key value: %w", err)
	}

	return key, nil
}

func KeyFromPassword(password string) ([]byte, error) {
	salt := []byte{0x4f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x73}

	key := pbkdf2.Key([]byte(password), salt, 1000, 16, sha1.New)

	return key, nil
}

func DecryptString(key []byte, value string) (string, error) {
	// split the encrypted value into data and salt
	parts := strings.SplitN(value, "|", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid value format, expected: '<base64>|<base64>' but got '%v'", value)
	}

	// decode the values
	data, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data value '%v': %w", parts[0], err)
	}

	salt, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 salt value '%v': %w", parts[1], err)
	}

	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to initialize cipher: %w", err)
	}

	mode := cipher.NewCBCDecrypter(block, salt)

	// decrypt the data
	decrypted := make([]byte, len(data))
	mode.CryptBlocks(decrypted, data)

	// remove padding
	decryptedData, err := unpadPKCS7(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data value: %w", err)
	}

	// done
	return string(decryptedData), nil
}
