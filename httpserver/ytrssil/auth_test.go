package ytrssil_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	EndpointsTestSuite
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestAuthPage() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth", nil)
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.Contains(w.Body.String(), "Authentication")
}

func (s *AuthTestSuite) TestHandleAuthSuccess() {
	form := url.Values{}
	form.Add("token", s.cfg.AuthToken)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusFound, w.Code)
	cookies := w.Result().Cookies()
	s.Require().Len(cookies, 1)
	s.Equal("token", cookies[0].Name)
	s.Equal(s.cfg.AuthToken, cookies[0].Value)
	s.Equal("/", w.Header().Get("Location"))
}

func (s *AuthTestSuite) TestHandleAuthInvalidToken() {
	form := url.Values{}
	form.Add("token", "invalid-token")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
	s.Contains(w.Body.String(), "Invalid token")
}

func (s *AuthTestSuite) TestAuthMiddleware() {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: s.cfg.AuthToken})
	s.server.Handler.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
}
