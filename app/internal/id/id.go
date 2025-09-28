package id

import (
	"crypto/rand"
	"encoding/hex"
)

// New returns a URL-safe 32 character identifier suitable for DynamoDB keys.
func New() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}
