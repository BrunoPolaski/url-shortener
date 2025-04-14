package cmd

import (
	"fmt"
	"net/http"
	"strings"

	http2 "github.com/AdagaDigital/url-redirect-service/internal/adapters/http"
	"github.com/AdagaDigital/url-redirect-service/internal/adapters/http/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

func Handler(request *events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("event.RawPath:", request.RawPath)
	fmt.Println("event.RequestContext.Path:", request.RequestContext.HTTP.Path)

	if request.RequestContext.HTTP.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	httpRequest, err := http.NewRequest(request.RequestContext.HTTP.Method, request.RequestContext.HTTP.Path, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	for k, v := range request.Headers {
		httpRequest.Header.Set(k, v)
	}

	q := httpRequest.URL.Query()
	for k, v := range request.QueryStringParameters {
		q.Set(k, v)
	}
	httpRequest.URL.RawQuery = q.Encode()

	rr := http2.NewResponseRecorder()

	engine := gin.Default()
	routes.InitRoutes(engine)
	engine.ServeHTTP(rr, httpRequest)

	rr.Headers["Content-Type"] = "application/json"
	rr.Headers["Access-Control-Allow-Origin"] = "*"

	return &events.APIGatewayProxyResponse{
		StatusCode: rr.StatusCode,
		Body:       rr.Body,
		Headers:    rr.Headers,
	}, nil
}
