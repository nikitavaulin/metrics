package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Hash(value, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}
