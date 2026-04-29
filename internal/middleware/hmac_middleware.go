package middleware

import (
	"crypto/hmac"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/auhmaugmaufm/predict-ticket-department-backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func HMACMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp := c.GetHeader("X-Timestamp")
		signature := c.GetHeader("X-Signature")
		if timestamp == "" || signature == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing signature headers"})
			return
		}
		ts, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil || math.Abs(float64(time.Now().Unix()-ts)) > 300 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired timestamp"})
			return
		}
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
			return
		}
		c.Request.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))
		message := timestamp + "." + string(bodyBytes)
		expected := auth.ComputeHMAC(secret, message)
		incoming := strings.TrimPrefix(signature, "sha256=")
		if !hmac.Equal([]byte(expected), []byte(incoming)) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid signature"})
			return
		}
		c.Next()
	}
}
