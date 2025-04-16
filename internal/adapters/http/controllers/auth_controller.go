package controllers

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/services"
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

func (lc *authController) RefreshToken(c *gin.Context) {

}
