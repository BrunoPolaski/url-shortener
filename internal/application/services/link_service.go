package services

import (
	"github.com/AdagaDigital/url-redirect-service/internal/domain/repositories"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/services"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/google/uuid"
)

type linkService struct {
	linkRepository repositories.LinkRepository
}

func NewLinkService(lr repositories.LinkRepository) services.LinkService {
	return &linkService{
		linkRepository: lr,
	}
}

func (ls *linkService) Redirect(inputUuid string) (string, *rest_err.RestErr) {
	uuid.Validate(inputUuid)
}

func (ls *linkService) CreateLink(url string) (string, *rest_err.RestErr) {
	return "", nil
}
