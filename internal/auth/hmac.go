package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHMAC(secret string, data []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}

func VerifyHMAC(secret string, data []byte, signature string) bool {
	expected := GenerateHMAC(secret, data)
	return hmac.Equal([]byte(expected), []byte(signature))
}
