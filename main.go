package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/AdagaDigital/url-redirect-service/internal/adapters/http/routes"
	"github.com/AdagaDigital/url-redirect-service/internal/cmd"
	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		log.Fatalf("Error loading location: %v", err)
	}
	time.Local = location

	err = godotenv.Overload(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %s", err)
	}

	if len(os.Args) > 1 && strings.Contains(os.Args[1], "migration") {
		cmd.Migrate()
		return
	}

	logger.Init()

	logger.Info("Starting application")

	if os.Getenv("ENV") == "local" {
		gin.SetMode(gin.ReleaseMode)

		engine := gin.Default()

		routes.InitRoutes(engine)

		log.Fatal(http.ListenAndServe(":8080", engine))
	}

	lambda.Start(cmd.Handler)
}
