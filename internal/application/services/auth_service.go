package services

import (
	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/services"
	internalJwt "github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/jwt"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authService struct {
	authRepository repositories.AuthRepository
}

func NewAuthService(ar repositories.AuthRepository) services.AuthService {
	return &authService{
		authRepository: ar,
	}
}

func (ar *authService) GenerateApiKey(input *request.CreateApiKeyRequest) (entities.ApiKey, *rest_err.RestErr) {
	logger.Info("Generating API key")

	apiKey := entities.NewApiKey(
		uuid.NewString(),
		uuid.NewString(),
		input.Slug,
	)
	_, restErr := ar.authRepository.CreateApiKey(apiKey)
	if restErr != nil {
		return nil, restErr
	}

	return apiKey, nil
}

func (ar *authService) RefreshToken(refreshToken string) (string, *rest_err.RestErr) {
	logger.Info("Refreshing token")

	jwtAdapter := internalJwt.NewJWTAdapter()
	token, err := jwtAdapter.ParseToken(refreshToken)
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)

	domainToken := entities.NewToken(
		claims["tid"].(string),
		claims["sub"].(string),
	)

	_, err = ar.authRepository.FindApiKey(domainToken.GetApiKey())
	if err != nil {
		return "", err
	}

	newRefreshToken := entities.NewToken(
		uuid.NewString(),
		domainToken.GetApiKey(),
	)

	token, err := internalJwt.GenerateToken(uid)
	if err != nil {
		return "", err
	}

	_, restErr := ar.authRepository.CreateToken(token)
	if restErr != nil {
		return "", restErr
	}

	return token.GetUUID(), nil
}
