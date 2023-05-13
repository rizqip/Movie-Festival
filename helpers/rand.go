package helpers

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomByte(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}
