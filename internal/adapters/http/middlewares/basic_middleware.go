package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(c *gin.Context) {
	u, p, ok := c.Request.BasicAuth()
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}

	if len(u) == 0 || len(p) == 0 {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
