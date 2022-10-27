package ytrssil

import (
	"net/http"

	"gitea.theedgeofrage.com/TheEdgeOfRage/ytrssil-api/models"
	"github.com/gin-gonic/gin"
)

func (s *server) GetNewVideos(c *gin.Context) {
	username := c.MustGet("username").(string)
	videos, err := s.handler.GetNewVideos(c.Request.Context(), username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.GetNewVideosResponse{
		Videos: videos,
	})
}