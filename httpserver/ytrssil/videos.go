package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/pages"
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

func (srv *server) GetVideoFormatsJSON(c *gin.Context) {
	videoID := c.Param("video_id")

	formats, err := srv.handler.GetVideoFormats(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"formats": formats})
}

func (srv *server) DownloadVideoJSON(c *gin.Context) {
	videoID := c.Param("video_id")
	format := c.PostForm("format")

	err := srv.handler.DownloadVideoWithFormat(c.Request.Context(), videoID, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "download started"})
}

func (srv *server) GetResolutionModal(c *gin.Context) {
	videoID := c.Param("video_id")
	title := c.Query("title")

	formats, err := srv.handler.GetVideoFormats(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Render(http.StatusOK, pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.ResolutionModal(videoID, title, formats),
	})
}
