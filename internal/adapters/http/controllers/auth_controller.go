package controllers

import (
	"net/http"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/services"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/gin-gonic/gin"
)

type authController struct {
	authService services.AuthService
}

func NewAuthController(as services.AuthService) controllers.AuthController {
	return &authController{
		authService: as,
	}
}

func (ac *authController) RefreshToken(c *gin.Context) {

}

func (ac *authController) CreateApiKey(c *gin.Context) {
	var body request.CreateApiKeyRequest
	var restErr *rest_err.RestErr

	err := c.ShouldBindJSON(&body)
	if err != nil {
		restErr = rest_err.NewBadRequestError("invalid body")
		c.AbortWithStatusJSON(restErr.Code, restErr)
		return
	}

	response, restErr := ac.authService.GenerateApiKey(&body)
	if restErr != nil {
		c.AbortWithStatusJSON(restErr.Code, restErr)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ac *authController) Login(c *gin.Context) {
	apiKey := c.GetHeader("x-api-key")
	logger.Info(apiKey)
	if apiKey == "" {
		restErr := rest_err.NewBadRequestError("invalid api key")
		c.JSON(restErr.Code, restErr)
		return
	}

	obj, err := ac.authService.Login(apiKey)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, obj)
}
