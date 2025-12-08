package ytrssil

import (
	"github.com/gin-gonic/gin"
)

func returnErr(c *gin.Context, status int, err error) {
	c.Data(status, "text/html", []byte(err.Error()))
}

func returnMsg(c *gin.Context, status int, msg string) {
	c.Data(status, "text/html", []byte(msg))
}
