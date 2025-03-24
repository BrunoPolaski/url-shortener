package controllers

import "github.com/gin-gonic/gin"

type LinkController interface {
	Redirect(c *gin.Context)
	CreateLink(c *gin.Context)
}
