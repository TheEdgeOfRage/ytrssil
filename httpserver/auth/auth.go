package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

// AuthMiddleware will authenticate against a static API key
func AuthMiddleware(db db.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid basic auth header"})
			return
		}
		user := models.User{
			Username: username,
			Password: password,
		}
		authenticated, err := db.AuthenticateUser(c.Request.Context(), user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		if !authenticated {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
			return
		}

		c.Set("username", username)

		// handle request
		c.Next()
	}
}
