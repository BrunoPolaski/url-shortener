package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AdagaDigital/url-redirect-service/internal/cmd"
	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/interfaces/http/routes"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %s", err)
	}

	logger.Init()

	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		cmd.Migrate()
		return
	}

	logger.Info("Starting application")

	if os.Getenv("ENV") == "local" {
		gin.SetMode(gin.ReleaseMode)

		engine := gin.Default()

		routes.InitRoutes(engine)

		log.Fatal(http.ListenAndServe(":8080", engine))
	}

	lambda.Start(cmd.Handler)
}
