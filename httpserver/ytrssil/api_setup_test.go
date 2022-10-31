package ytrssil_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/lib/log"
)

var testConfig config.Config

func init() {
	// always use UTC
	time.Local = time.UTC
	testConfig = config.TestConfig()
}

func setupTestServer(t *testing.T) *http.Server {
	l := log.NewNopLogger()

	handler := handler.New(l, nil, nil)

	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(
		l,
		handler,
		auth.AuthMiddleware(nil),
	)
	assert.Nil(t, err)

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", testConfig.Gin.Port),
		Handler: router,
	}
}

func TestHealthz(t *testing.T) {
	server := setupTestServer(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "healthy", w.Body.String())
}
