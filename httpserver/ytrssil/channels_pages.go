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

	returnMsg(c, http.StatusAccepted, "")
}
