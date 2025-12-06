package ytrssil

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/db"
	"github.com/TheEdgeOfRage/ytrssil-api/feedparser"
	"github.com/TheEdgeOfRage/ytrssil-api/models"
)

func (srv *server) SubscribeToChannelJSON(c *gin.Context) {
	var channel models.Channel
	err := c.ShouldBindUri(&channel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.GetString("username")

	err = srv.handler.SubscribeToChannel(c.Request.Context(), username, channel.ID)
	if err != nil {
		if errors.Is(err, db.ErrAlreadySubscribed) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, feedparser.ErrInvalidChannelID) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "subscribed to channel successfully"})
}

func (srv *server) UnsubscribeFromChannelJSON(c *gin.Context) {
	var channel models.Channel
	err := c.ShouldBindUri(&channel)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	username := c.GetString("username")

	err = srv.handler.UnsubscribeFromChannel(c.Request.Context(), username, channel.ID)
	if err != nil {
		if errors.Is(err, db.ErrChannelNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "unsubscribed from channel successfully"})
}
