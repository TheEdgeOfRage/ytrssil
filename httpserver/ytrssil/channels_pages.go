package ytrssil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
)

func (srv *server) SubscribeToChannelPage(c *gin.Context) {
	err := srv.handler.SubscribeToChannel(c.Request.Context(), c.PostForm("channel_id"))
	if err != nil {
		if errors.Is(err, db.ErrAlreadySubscribed) {
			c.String(http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, feedparser.ErrInvalidChannelID) {
			c.Data(http.StatusNotFound, "text/html", []byte(err.Error()))
			return
		}

		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, "")
}
