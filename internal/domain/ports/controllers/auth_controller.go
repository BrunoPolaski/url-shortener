package controllers

import "github.com/gin-gonic/gin"

type AuthController interface {
	RefreshToken(c *gin.Context)
}
