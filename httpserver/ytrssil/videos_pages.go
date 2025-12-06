package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/pages"
)

func (srv server) NewVideosPage(c *gin.Context) {
	videos, err := srv.handler.GetNewVideos(c.Request.Context(), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	r := pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.NewVideosPage(videos),
	}
	c.Render(http.StatusOK, r)
}

func (srv server) MarkVideoAsWatchedPage(c *gin.Context) {
	r := pages.TemplRenderer{
		Ctx: c.Request.Context(),
	}

	err := srv.handler.MarkVideoAsWatched(c.Request.Context(), c.Param("video_id"))
	if err != nil {
		r.Component = pages.ErrorPage(err.Error())
		c.Render(http.StatusInternalServerError, r)
		return
	}

	c.String(http.StatusOK, "")
}

func (srv server) SetVideoProgressPage(c *gin.Context) {
	r := pages.TemplRenderer{
		Ctx: c.Request.Context(),
	}
	progress := c.PostForm("progress")
	if progress == "" {
		r.Component = pages.ErrorPage("missing progress")
		c.Render(http.StatusBadRequest, r)
		return
	}

	video, err := srv.handler.SetVideoProgress(c.Request.Context(), c.Param("video_id"), progress)
	if err != nil {
		r.Component = pages.ErrorPage(err.Error())
		c.Render(http.StatusInternalServerError, r)
		return
	}

	r.Component = pages.ProgressBar(*video)
	c.Render(http.StatusOK, r)
}
