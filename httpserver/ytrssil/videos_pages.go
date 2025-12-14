package ytrssil

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/pages"
)

func (srv server) NewVideosPage(c *gin.Context) {
	videos, err := srv.handler.GetNewVideos(c.Request.Context(), true)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	c.Render(http.StatusOK, pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.NewVideosPage(videos),
	})
}

func (srv server) WatchedVideosPage(c *gin.Context) {
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	videos, err := srv.handler.GetWatchedVideos(c.Request.Context(), true, page)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	c.Render(http.StatusOK, pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.WatchedVideosPage(videos, page),
	})
}

func (srv server) MarkVideoAsWatchedPage(c *gin.Context) {
	err := srv.handler.MarkVideoAsWatched(c.Request.Context(), c.Param("video_id"))
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	returnMsg(c, http.StatusOK, "")
}

func (srv server) MarkVideoAsUnwatchedPage(c *gin.Context) {
	err := srv.handler.MarkVideoAsUnwatched(c.Request.Context(), c.Param("video_id"))
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	returnMsg(c, http.StatusOK, "")
}

func (srv server) SetVideoProgressPage(c *gin.Context) {
	progress := c.PostForm("progress")
	if progress == "" {
		returnMsg(c, http.StatusBadRequest, "missing progress")
		return
	}

	video, err := srv.handler.SetVideoProgress(c.Request.Context(), c.Param("video_id"), progress)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	c.Render(http.StatusOK, pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.ProgressBar(*video),
	})
}

func (srv server) AddVideoPage(c *gin.Context) {
	videoID := c.PostForm("video_id")
	if videoID == "" {
		returnMsg(c, http.StatusBadRequest, "missing video ID")
		return
	}

	err := srv.handler.AddCustomVideo(c.Request.Context(), videoID)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	returnMsg(c, http.StatusAccepted, "")
}
