package controllers

import (
	"net/http"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/controllers/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/services"
	"github.com/gin-gonic/gin"
)

type linkController struct {
	linkService services.LinkService
}

func NewLinkController(ls services.LinkService) controllers.LinkController {
	return &linkController{
		linkService: ls,
	}
}

func (lc *linkController) Redirect(c *gin.Context) {
	uuid, ok := c.Params.Get("uuid")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	url, err := lc.linkService.Redirect(uuid)
	if err != nil {
		c.JSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func (lc *linkController) CreateLink(c *gin.Context) {
	var body request.CreateLink
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	uuid, restErr := lc.linkService.CreateLink(body.URL)
	if restErr != nil {
		c.JSON(restErr.Code, gin.H{"error": restErr.Message})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"uuid": uuid})
}
