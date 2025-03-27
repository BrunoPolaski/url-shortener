package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(c *gin.Context) {
	u, p, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		c.Abort()
		return
	}

	if len(u) == 0 || len(p) == 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Unauthorized"},
		)
		c.Abort()
		return
	}
}
