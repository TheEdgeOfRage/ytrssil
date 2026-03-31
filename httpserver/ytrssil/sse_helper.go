package ytrssil

import (
	"github.com/gin-gonic/gin"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func newSSE(c *gin.Context) *datastar.ServerSentEventGenerator {
	return datastar.NewSSE(c.Writer, c.Request)
}
