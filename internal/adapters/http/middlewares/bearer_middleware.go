package middlewares

import (
	"time"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	internal_jwt "github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func BearerMiddleware(authRepository repositories.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtAdapter := internal_jwt.NewJWTAdapter()

		token, restErr := jwtAdapter.TrimPrefix(c.GetHeader("Authorization"))
		if restErr != nil {
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		parsedToken, restErr := jwtAdapter.ParseToken(token)
		if restErr != nil {
			c.AbortWithStatusJSON(restErr.Code, restErr.Message)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			restErr := rest_err.NewUnauthorizedError("invalid token")
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil {
			restErr := rest_err.NewUnauthorizedError("invalid token")
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		if exp.Before(time.Now()) {
			restErr := rest_err.NewUnauthorizedError("invalid token")
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		sub, err := claims.GetSubject()
		if err != nil {
			restErr := rest_err.NewUnauthorizedError("invalid token")
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		_, restErr = authRepository.FindApiKey(sub)
		if restErr != nil {
			restErr := rest_err.NewUnauthorizedError("invalid token")
			c.AbortWithStatusJSON(restErr.Code, restErr)
			return
		}

		c.Set("apiKey", sub)
		c.Next()
	}
}
