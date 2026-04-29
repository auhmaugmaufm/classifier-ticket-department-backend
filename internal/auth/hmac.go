package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func ComputeHMAC(secret string, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}
