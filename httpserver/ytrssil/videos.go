package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *server) GetNewVideosJSON(c *gin.Context) {
	videos, err := srv.handler.GetNewVideos(c.Request.Context(), false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

func (srv *server) GetWatchedVideosJSON(c *gin.Context) {
	videos, err := srv.handler.GetWatchedVideos(c.Request.Context(), false, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"videos": videos})
}

func (srv *server) FetchVideosJSON(c *gin.Context) {
	err := srv.handler.FetchVideos(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "videos fetched successfully"})
}

func (srv *server) MarkVideoAsWatchedJSON(c *gin.Context) {
	err := srv.handler.MarkVideoAsWatched(c.Request.Context(), c.Param("video_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "marked video as watched"})
}

func (srv *server) MarkVideoAsUnwatchedJSON(c *gin.Context) {
	err := srv.handler.MarkVideoAsUnwatched(c.Request.Context(), c.Param("video_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "cleared video from watch history"})
}
