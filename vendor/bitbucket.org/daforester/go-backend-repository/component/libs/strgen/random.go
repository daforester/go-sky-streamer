package strgen

import (
	"crypto/rand"
	"encoding/base32"
)

func RandString(length int64) (string, error) {
	if length == 0 {
		length = 32
	}

	randomBytes, err := RandBytes(length)

	if err != nil {
		return "", err
	}

	return base32.StdEncoding.EncodeToString(randomBytes)[:length], err
}

func RandBytes(length int64) ([]byte, error) {
	if length == 0 {
		length = 32
	}

	randomBytes := make([]byte, length)

	_, err := rand.Read(randomBytes)

	return randomBytes, err
}
