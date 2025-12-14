package ytrssil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
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
	channel, err := srv.handler.SubscribeToChannel(c.Request.Context(), c.PostForm("channel_id"))
	if err != nil {
		if errors.Is(err, db.ErrAlreadySubscribed) {
			returnErr(c, http.StatusConflict, err)
			return
		}
		if errors.Is(err, feedparser.ErrInvalidChannelID) {
			returnErr(c, http.StatusBadRequest, err)
			return
		}

		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	r := pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.ChannelCard(*channel),
	}
	c.Render(http.StatusOK, r)
}

func (srv *server) UnsubscribeFromChannelPage(c *gin.Context) {
	err := srv.handler.UnsubscribeFromChannel(c.Request.Context(), c.Param("channel_id"))
	if err != nil {
		if errors.Is(err, db.ErrChannelNotFound) {
			returnErr(c, http.StatusNotFound, err)
			return
		}

		returnErr(c, http.StatusInternalServerError, err)
		return
	}

	returnMsg(c, http.StatusOK, "")
}
