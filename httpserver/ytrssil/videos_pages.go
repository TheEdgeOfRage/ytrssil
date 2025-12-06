package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
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

	var req models.VideoURIRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		r.Component = pages.ErrorPage(err)
		c.Render(http.StatusBadRequest, r)
		return
	}

	err = srv.handler.MarkVideoAsWatched(c.Request.Context(), req.VideoID)
	if err != nil {
		r.Component = pages.ErrorPage(err)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	c.String(http.StatusOK, "")
}

func (srv server) SetVideoProgressPage(c *gin.Context) {
	r := pages.TemplRenderer{
		Ctx: c.Request.Context(),
	}
	var req struct {
		models.VideoURIRequest
		models.SetVideoProgressRequest
	}
	err := c.ShouldBindUri(&req.VideoURIRequest)
	if err != nil {
		r.Component = pages.ErrorPage(err)
		c.Render(http.StatusBadRequest, r)
		return
	}
	err = c.ShouldBind(&req.SetVideoProgressRequest)
	if err != nil {
		r.Component = pages.ErrorPage(err)
		c.Render(http.StatusBadRequest, r)
		return
	}

	video, err := srv.handler.SetVideoProgress(c.Request.Context(), req.VideoID, req.Progress)
	if err != nil {
		r.Component = pages.ErrorPage(err)
		c.Render(http.StatusInternalServerError, r)
		return
	}

	r.Component = pages.ProgressBar(*video)
	c.Render(http.StatusOK, r)
}
