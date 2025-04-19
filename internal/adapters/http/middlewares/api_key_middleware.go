package middlewares

import (
	"net/http"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware(authRepository repositories.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("x-api-key")
		secretKey := c.GetHeader("x-secret-key")
		if len(apiKey) == 0 || len(secretKey) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		logger.Info(apiKey)
		domain, err := authRepository.FindApiKey(apiKey)
		if err != nil {
			c.JSON(err.Code, err)
			c.Abort()
			return
		}

		err = domain.DecryptSecret()
		if err != nil {
			c.JSON(err.Code, err)
			c.Abort()
			return
		}

		if domain.GetSecret() != secretKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}
