package views

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomID() string {
	idLength := 32

	randomBytes := make([]byte, idLength)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}
	id := hex.EncodeToString(randomBytes)

	return id
}
