package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	db_mock "github.com/TheEdgeOfRage/ytrssil-api/mocks/db"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func setupTestServer() *http.Server {
	db := &db_mock.DBMock{
		AuthenticateUserFunc: func(ctx context.Context, user models.User) (bool, error) {
			return user.Username == "username" && user.Password == "password", nil
		},
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Middlewares are executed top to bottom in a stack-like manner
	router.Use(
		gin.Recovery(), // Recovery needs to go before other middlewares to catch panics
		AuthMiddleware(db),
	)
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	return &http.Server{Handler: router}
}

func TestSuccessfulAuthentication(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.SetBasicAuth("username", "password") // Valid credentials
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestMissingAuthorizationHeader(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"invalid basic auth header"}`, w.Body.String())
}

func TestWrongCredentials(t *testing.T) {
	server := setupTestServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.SetBasicAuth("test", "test") // Invalid credentials
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, `{"error":"invalid username or password"}`, w.Body.String())
}
