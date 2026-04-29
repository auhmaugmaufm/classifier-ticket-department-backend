package middleware

import (
	"crypto/hmac"
	"io"
	"net/http"
	"strings"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func HMACMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-HMAC-Signature")
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		message := ""
		if len(bodyBytes) > 0 {
			message = string(bodyBytes)
		}
		expected := auth.ComputeHMAC(secret, message)
		incoming := strings.TrimPrefix(signature, "sha256=")
		if !hmac.Equal([]byte(expected), []byte(incoming)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}
		c.Next()
	}
}
