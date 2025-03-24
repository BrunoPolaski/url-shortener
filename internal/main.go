package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/interfaces/http/routes"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %s", err)
	}

	logger.Init()
	logger.Info("Starting application")

	if os.Getenv("ENV") == "local" {
		r := routes.InitRoutes()

		log.Fatal(http.ListenAndServe(":8080", r))
	}

	lambda.Start(app.Handler)
}
