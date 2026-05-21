package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateCacheKey(sourceCode string) string {
	hasher := sha256.New()

	hasher.Write([]byte(sourceCode))

	return hex.EncodeToString(hasher.Sum(nil))
}
