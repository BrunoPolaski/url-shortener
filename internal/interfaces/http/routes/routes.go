package routes

import (
	"github.com/AdagaDigital/url-redirect-service/internal/application/services"
	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/database"
	"github.com/AdagaDigital/url-redirect-service/internal/interfaces/http/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/interfaces/repositories"
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.RouterGroup) {
	logger.Info("Setting up routes")

	mysqlAdapter := database.NewMySQLAdapter()
	db := mysqlAdapter.Connect()

	controller := controllers.NewLinkController(
		services.NewLinkService(
			repositories.NewLinkRepository(
				db,
			),
		),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	r.GET("/redirect/:uuid", controller.Redirect)
	r.POST("/link", controller.CreateLink)
}
