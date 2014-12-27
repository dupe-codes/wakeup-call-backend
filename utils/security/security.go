package security

import (
    "crypto/sha256"
    "encoding/hex"
    "math/rand"
)

var (
    numSaltDigits = 12
    randRange = 10000000
)

// RunSHA2 implements the SHA-2 hashing function to hash the given string
func RunSHA2(str string) string {
    strBytes := []byte(str)
    hasher := sha256.New()
    hasher.Write(strBytes)
    return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateSalt returns a randomized hash string to use as a password salt
func GenerateSalt() string {
    return RunSHA2(string(rand.Intn(randRange)))[numSaltDigits:]
}
