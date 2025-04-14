package controllers

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/services"
	"github.com/gin-gonic/gin"
)

type authController struct {
	authService services.AuthService
}

func NewAuthController(ls services.AuthService) controllers.AuthController {
	return &authController{
		authService: ls,
	}
}

func (lc *authController) RefreshToken(c *gin.Context) {

}
