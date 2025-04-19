package repositories

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type AuthRepository interface {
	FindToken(uuid string) (string, *rest_err.RestErr)
	FindTokenByApiKey(apiKey string) (string, *rest_err.RestErr)
	DeleteToken(uuid string) *rest_err.RestErr
	DeleteTokenByApiKey(apiKey string) *rest_err.RestErr
	FindApiKey(apiKey string) (entities.ApiKey, *rest_err.RestErr)
	CreateApiKey(apiKey entities.ApiKey) (entities.ApiKey, *rest_err.RestErr)
	CreateToken(token entities.Token) (entities.Token, *rest_err.RestErr)
}
