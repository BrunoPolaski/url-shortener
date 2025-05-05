package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/AdagaDigital/url-redirect-service/internal/adapters/http/routes"
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

func Handler(request *events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
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

	rr := httptest.NewRecorder()

	engine := gin.Default()
	routes.InitRoutes(engine)
	engine.ServeHTTP(rr, httpRequest)

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}
	for k, v := range rr.Header() {
		headers[k] = strings.Join(v, ",")
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: rr.Code,
		Body:       rr.Body.String(),
		Headers:    headers,
	}, nil
}
