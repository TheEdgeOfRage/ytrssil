package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// PageAuthMiddleware will authenticate against a static auth token
func PageAuthMiddleware(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenCookie, err := c.Request.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth token cookie"})
			return
		}
		if tokenCookie.Value != authToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth token"})
			return
		}

		// handle request
		c.Next()
	}
}

// APIAuthMiddleware will authenticate API endpoints against a static auth token
func APIAuthMiddleware(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := c.Request.Header["Authorization"]
		if !ok || len(token) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			return
		}
		if token[0] != authToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth token"})
			return
		}

		// handle request
		c.Next()
	}
}
