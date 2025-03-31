package services

import (
	"net/url"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
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

func (ls *linkService) Redirect(inputUUID string) (string, *rest_err.RestErr) {
	if err := uuid.Validate(inputUUID); err != nil {
		return "", rest_err.NewBadRequestError("invalid uuid")
	}

	link, err := ls.linkRepository.GetByUUID(inputUUID)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (ls *linkService) CreateLink(inputURL string) (string, *rest_err.RestErr) {
	logger.Info("Creating link")
	_, err := url.Parse(inputURL)
	if err != nil {
		return "", rest_err.NewBadRequestError("invalid url")
	}

	uuid := uuid.New().String()

	redirect := entities.NewRedirect(uuid, inputURL)
	_, restErr := ls.linkRepository.Create(redirect)
	if restErr != nil {
		return "", restErr
	}

	return uuid, nil
}
