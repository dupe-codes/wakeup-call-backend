package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

var (
	saltLength = 64
)

// RunSHA2 implements the SHA-2 hashing function to hash the given string
func RunSHA2(str string) string {
	strBytes := []byte(str)
	hasher := sha256.New()
	hasher.Write(strBytes)
	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSalt returns a randomized hash string to use as a password salt
// TODO: Continue looking in to best way to do this
func GenerateSalt() string {
	randomBytes := make([]byte, saltLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err) // TODO: Better way of handling, maybe pass error up?
	}

	return base64.URLEncoding.EncodeToString(randomBytes)
}
