package middlewares

import (
	"strings"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	internal_jwt "github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func BearerMiddleware(authRepository repositories.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtAdapter := internal_jwt.NewJWTAdapter()

		if strings.Contains(c.GetHeader("Authorization"), "Bearer ") {
			errRest := rest_err.NewBadRequestError("invalid token")
			c.JSON(errRest.Code, errRest)
			c.Abort()
			return
		}

		token := strings.Split(c.GetHeader("Authorization"), " ")[1]

		parsedToken, err := jwtAdapter.ParseToken(token)
		if err != nil {
			c.JSON(err.Code, err.Message)
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok || !parsedToken.Valid {
			errRest := rest_err.NewUnauthorizedError("invalid token")
			c.JSON(errRest.Code, errRest)
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok {
			errRest := rest_err.NewUnauthorizedError("invalid token")
			c.JSON(errRest.Code, errRest)
			c.Abort()
			return
		}

		_, err = authRepository.FindApiKey(sub)
		if err != nil {
			errRest := rest_err.NewUnauthorizedError("invalid token")
			c.JSON(errRest.Code, errRest)
			c.Abort()
			return
		}

		c.Next()
	}
}
