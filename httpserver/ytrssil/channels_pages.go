package ytrssil

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	datastar "github.com/starfederation/datastar-go/datastar"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/pages"
)

func (srv *server) ChannelsPage(c *gin.Context) {
	channels, err := srv.handler.ListChannels(c.Request.Context())
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	r := pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.ChannelsPage(channels),
	}
	c.Render(http.StatusOK, r)
}

func (srv *server) SubscribeToChannelPage(c *gin.Context) {
	var signals struct {
		ChannelID string `json:"channelID"`
	}
	if err := datastar.ReadSignals(c.Request, &signals); err != nil {
		sse := newSSE(c)
		sse.ExecuteScript(fmt.Sprintf(`showFormError("subscription-modal", %q)`, err.Error()))
		return
	}

	_, err := srv.handler.SubscribeToChannel(c.Request.Context(), signals.ChannelID)
	if err != nil {
		sse := newSSE(c)
		sse.ExecuteScript(fmt.Sprintf(`showFormError("subscription-modal", %q)`, err.Error()))
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(`
		bootstrap.Modal.getInstance(document.getElementById("subscription-modal")).hide();
		location.reload()
	`)
}

func (srv *server) UnsubscribeFromChannelPage(c *gin.Context) {
	channelID := c.Param("channel_id")
	err := srv.handler.UnsubscribeFromChannel(c.Request.Context(), channelID)
	if err != nil {
		if errors.Is(err, db.ErrChannelNotFound) {
			returnErr(c, http.StatusNotFound, err)
			return
		}

		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.ExecuteScript(fmt.Sprintf(`animateRemove("#channel-card-%s")`, channelID))
}

func (srv *server) ToggleChannelShortsPage(c *gin.Context) {
	var signals struct {
		Enable bool `json:"enable"`
	}
	if err := datastar.ReadSignals(c.Request, &signals); err != nil {
		returnErr(c, http.StatusBadRequest, err)
		return
	}

	channelID := c.Param("channel_id")
	err := srv.handler.ToggleChannelShorts(c.Request.Context(), channelID, signals.Enable)
	if err != nil {
		if errors.Is(err, db.ErrChannelNotFound) {
			returnErr(c, http.StatusNotFound, err)
			return
		}

		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	channel, err := srv.handler.GetChannelByID(c.Request.Context(), channelID)
	if err != nil {
		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	sse := newSSE(c)
	sse.PatchElementTempl(pages.ChannelCard(*channel))
}
