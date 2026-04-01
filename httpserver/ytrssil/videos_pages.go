package ytrssil

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	datastar "github.com/starfederation/datastar-go/datastar"

	"github.com/TheEdgeOfRage/ytrssil-api/models"
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

func (srv server) FetchVideosPage(c *gin.Context) {
	err := srv.handler.FetchVideos(c.Request.Context())
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(`window.location.reload()`)
}

func (srv server) MarkVideoAsWatchedPage(c *gin.Context) {
	videoID := c.Param("video_id")
	err := srv.handler.MarkVideoAsWatched(c.Request.Context(), videoID)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(fmt.Sprintf(`animateRemove("#video-card-%s")`, videoID))
}

func (srv server) MarkVideoAsUnwatchedPage(c *gin.Context) {
	videoID := c.Param("video_id")
	err := srv.handler.MarkVideoAsUnwatched(c.Request.Context(), videoID)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(fmt.Sprintf(`animateRemove("#video-card-%s")`, videoID))
}

func (srv server) SetVideoProgressPage(c *gin.Context) {
	var signals struct {
		Progress string `json:"progress"`
	}
	if err := datastar.ReadSignals(c.Request, &signals); err != nil {
		returnErr(c, http.StatusBadRequest, err)
		return
	}
	if signals.Progress == "" {
		returnMsg(c, http.StatusBadRequest, "missing progress")
		return
	}

	video, err := srv.handler.SetVideoProgress(c.Request.Context(), c.Param("video_id"), signals.Progress)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.PatchElementTempl(pages.VideoCard(*video))
}

func (srv server) AddVideoPage(c *gin.Context) {
	var signals struct {
		VideoID string `json:"videoID"`
	}
	if err := datastar.ReadSignals(c.Request, &signals); err != nil {
		sse := newSSE(c)
		sse.ExecuteScript(fmt.Sprintf(`showFormError("add-video-modal", %q)`, err.Error()))
		return
	}
	if signals.VideoID == "" {
		sse := newSSE(c)
		sse.ExecuteScript(`showFormError("add-video-modal", "missing video ID")`)
		return
	}

	err := srv.handler.AddCustomVideo(c.Request.Context(), signals.VideoID)
	if err != nil {
		sse := newSSE(c)
		sse.ExecuteScript(fmt.Sprintf(`showFormError("add-video-modal", %q)`, err.Error()))
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(`bootstrap.Modal.getInstance(document.getElementById("add-video-modal")).hide(); location.reload()`)
}

func (srv server) DownloadVideoPage(c *gin.Context) {
	videoID := c.Param("video_id")
	resolution := c.PostForm("format")

	err := srv.handler.DownloadVideo(c.Request.Context(), videoID, resolution)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	videos, err := srv.handler.GetNewVideos(c.Request.Context(), true)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	var targetVideo *models.Video
	for _, v := range videos {
		if v.ID == videoID {
			targetVideo = &v
			break
		}
	}

	if targetVideo == nil {
		returnErr(c, http.StatusNotFound, fmt.Errorf("video not found"))
		return
	}

	c.Render(http.StatusOK, pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.VideoCard(*targetVideo),
	})
}

func (srv server) GetVideoCardPage(c *gin.Context) {
	videoID := c.Param("video_id")

	videos, err := srv.handler.GetNewVideos(c.Request.Context(), true)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	var targetVideo *models.Video
	for _, v := range videos {
		if v.ID == videoID {
			targetVideo = &v
			break
		}
	}

	if targetVideo == nil {
		returnErr(c, http.StatusNotFound, fmt.Errorf("video not found"))
		return
	}

	sse := newSSE(c)
	sse.PatchElementTempl(pages.VideoCard(*targetVideo))
}

func (srv server) ServeVideoFilePage(c *gin.Context) {
	videoID := c.Param("video_id")

	filePath, filename, err := srv.handler.ServeVideoFile(c.Request.Context(), videoID)
	if err != nil {
		returnErr(c, http.StatusNotFound, err)
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.File(filePath)
}
