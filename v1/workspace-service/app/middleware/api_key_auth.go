package middleware

import (
	"crypto/subtle"
	"strings"

	"github.com/gin-gonic/gin"
)

const bearerPrefix = "Bearer "

// APIKeyAuth rejects requests whose Authorization header does not carry the expected bearer token.
func APIKeyAuth(apiKey string) gin.HandlerFunc {
	expected := []byte(apiKey)
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if len(expected) == 0 || !strings.HasPrefix(header, bearerPrefix) {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		provided := []byte(strings.TrimPrefix(header, bearerPrefix))
		if subtle.ConstantTimeCompare(provided, expected) != 1 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	}
}
