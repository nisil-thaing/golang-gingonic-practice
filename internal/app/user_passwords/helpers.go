package userpasswords

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)

	if err != nil {
		return "", err
	}

	encodedSalt := base64.StdEncoding.EncodeToString(salt)

	return encodedSalt, nil
}
