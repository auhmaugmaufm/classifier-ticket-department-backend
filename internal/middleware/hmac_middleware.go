package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func HMACMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signature := c.GetHeader("X-HMAC-Signature")
		body, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if !auth.VerifyHMAC(secret, body, signature) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid HMAC"})
			return
		}
		c.Next()
	}
}
