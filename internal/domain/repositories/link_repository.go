package repositories

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type LinkRepository interface {
	GetByUUID(UUID string) (string, *rest_err.RestErr)
	Create(redirect *entities.Redirect) (*entities.Redirect, *rest_err.RestErr)
}
