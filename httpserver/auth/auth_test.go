package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	"github.com/TheEdgeOfRage/ytrssil-api/config"
)

type AuthTestSuite struct {
	suite.Suite

	cfg    config.Config
	server *http.Server
	engine *gin.Engine
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) SetupSuite() {
	s.cfg = config.TestConfig()

	gin.SetMode(gin.TestMode)
	s.engine = gin.New()
	s.engine.Use(
		gin.Recovery(),
		APIAuthMiddleware(s.cfg.AuthToken),
	)
	s.engine.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	s.server = &http.Server{Handler: s.engine}
}

func (s *AuthTestSuite) TestSuccessfulAuthentication() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header["Authorization"] = []string{"foo"}
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Equal("OK", w.Body.String())
}

func (s *AuthTestSuite) TestMissingAuthorizationHeader() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Equal(`{"error":"missing Authorization header"}`, w.Body.String())
}

func (s *AuthTestSuite) TestWrongCredentials() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header["Authorization"] = []string{"bar"}
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.Equal(`{"error":"invalid auth token"}`, w.Body.String())
}
