package services

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type AuthService interface {
	RefreshToken(uid string) (string, *rest_err.RestErr)
	GenerateApiKey(*request.CreateApiKeyRequest) (entities.ApiKey, *rest_err.RestErr)
}
