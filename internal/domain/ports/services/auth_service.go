package services

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/response"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type AuthService interface {
	RefreshToken(uid string) (*response.LoginResponse, *rest_err.RestErr)
	FindApiKey(uuid string) (entities.ApiKey, *rest_err.RestErr)
	GenerateApiKey(*request.CreateApiKeyRequest) (*response.CreateApiKeyResponse, *rest_err.RestErr)
	Login(apiKey string) (*response.LoginResponse, *rest_err.RestErr)
}
