package middlewares

import "github.com/gin-gonic/gin"

func ApiKeyMiddleware(c *gin.Context) {
	apiKey := c.Request.Header.Get("x-api-key")
	if len(apiKey) == 0 {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	if apiKey != "your_api_key" {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	c.Next()
}
