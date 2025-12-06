package ytrssil

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/TheEdgeOfRage/ytrssil-api/pages"
)

func (srv server) AuthPage(c *gin.Context) {
	r := pages.TemplRenderer{
		Ctx:       c.Request.Context(),
		Component: pages.AuthPage(""),
	}
	c.Render(http.StatusOK, r)
}

func (srv server) HandleAuth(c *gin.Context) {
	token := c.PostForm("token")

	if token != srv.cfg.AuthToken {
		r := pages.TemplRenderer{
			Ctx:       c.Request.Context(),
			Component: pages.AuthPage("Invalid token"),
		}
		c.Render(http.StatusBadRequest, r)
		return
	}

	c.SetCookie(
		"token",
		token,
		0,            // maxAge (0 = session cookie)
		"/",          // path
		"",           // domain (empty = current domain)
		!srv.cfg.Dev, // secure (false = works on HTTP)
		true,         // httpOnly (true = not accessible via JavaScript)
	)
	c.Redirect(http.StatusFound, "/")
}
