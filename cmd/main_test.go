package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/config"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/db"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/handler"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/auth"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/httpserver/ytrssil"
	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
)

var testConfig config.Config

func init() {
	testConfig = config.TestConfig()
}

func setupTestServer(t *testing.T) (*http.Server, db.DB) {
	l := slog.New(slog.NewTextHandler(io.Discard, nil))

	db, err := db.NewPostgresDB(l, testConfig.DB)
	if !assert.NoError(t, err) {
		return nil, nil
	}
	parser := feedparser.NewParser(l)
	handler := handler.New(l, db, parser)
	gin.SetMode(gin.TestMode)
	router, err := ytrssil.SetupGinRouter(
		l,
		handler,
		auth.AuthMiddleware(db),
	)
	if !assert.NoError(t, err) {
		return nil, nil
	}

	return &http.Server{
		Addr:    fmt.Sprintf(":%v", testConfig.Gin.Port),
		Handler: router,
	}, db
}

func TestHealthz(t *testing.T) {
	server, _ := setupTestServer(t)
	if !assert.NotNil(t, server) {
		return
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthz", nil)
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "healthy", w.Body.String())
}

func TestCreateUser(t *testing.T) {
	server, db := setupTestServer(t)
	if !assert.NotNil(t, server) {
		return
	}

	jsonData, err := json.Marshal(models.User{Username: "test", Password: "test"})
	if !assert.Nil(t, err) {
		return
	}
	data := bytes.NewBuffer(jsonData)
	req, _ := http.NewRequest("POST", "/register", data)
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"msg":"user created"}`, w.Body.String())

	db.DeleteUser(context.TODO(), "test")
}
