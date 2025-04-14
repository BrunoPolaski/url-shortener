package routes

import (
	"github.com/AdagaDigital/url-redirect-service/internal/adapters/http/controllers"
	"github.com/AdagaDigital/url-redirect-service/internal/adapters/http/middlewares"
	"github.com/AdagaDigital/url-redirect-service/internal/adapters/mysql"
	"github.com/AdagaDigital/url-redirect-service/internal/application/services"
	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/database"
	"github.com/gin-gonic/gin"
)

func InitRoutes(e *gin.Engine) {
	logger.Info("Setting up routes")

	mysqlAdapter := database.NewMySQLAdapter()
	db := mysqlAdapter.Connect()

	linkController := controllers.NewLinkController(
		services.NewLinkService(
			mysql.NewLinkRepositoryMySQL(
				db,
			),
		),
	)

	e.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	auth := e.Group("/auth", middlewares.BasicAuthMiddleware, middlewares.ApiKeyMiddleware)
	{
		auth.POST("/login")
	}

	link := e.Group("/link")
	{
		link.GET("/redirect/:uuid", linkController.Redirect)
		link.POST("/link", middlewares.BearerMiddleware, linkController.CreateLink)
	}

	e.NoMethod(func(c *gin.Context) {
		c.JSON(405, gin.H{"error": "Method not allowed"})
		c.Abort()
	})
}
