package main

import (
	"crypto/sha512"
	"encoding/hex"
)

func generateHash(data string) string {
	hash := sha512.New()
	hash.Write([]byte(data))
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
