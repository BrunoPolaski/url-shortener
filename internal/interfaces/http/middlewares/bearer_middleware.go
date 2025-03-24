package middlewares

import (
	"strings"

	internal_jwt "github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func BearerMiddleware(c *gin.Context) {
	jwtAdapter := internal_jwt.NewJWTAdapter()

	if strings.Contains(c.GetHeader("Authorization"), "Bearer ") {
		errRest := rest_err.NewBadRequestError("invalid token")
		c.JSON(errRest.Code, errRest)
		return
	}

	token := strings.Split(c.GetHeader("Authorization"), " ")[1]

	parsedToken, err := jwtAdapter.ParseToken(token)
	if err != nil {
		c.JSON(err.Code, err.Message)
		return
	}

	_, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		errRest := rest_err.NewUnauthorizedError("invalid token")
		c.JSON(errRest.Code, errRest)
		return
	}

}
