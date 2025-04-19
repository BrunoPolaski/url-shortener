package services

import (
	"fmt"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/response"
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

func (as *authService) GenerateApiKey(input *request.CreateApiKeyRequest) (*response.CreateApiKeyResponse, *rest_err.RestErr) {
	logger.Info("Generating API key")

	domainApiKey := entities.NewApiKey(
		uuid.NewString(),
		uuid.NewString(),
		input.Slug,
	)

	if err := domainApiKey.EncryptSecret(); err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to encrypt secret: %s", err.Error()))
	}

	apiKey, restErr := as.authRepository.CreateApiKey(domainApiKey)
	if restErr != nil {
		return nil, restErr
	}

	if restErr = apiKey.DecryptSecret(); restErr != nil {
		return nil, restErr
	}

	response := &response.CreateApiKeyResponse{
		UUID:   apiKey.GetUUID(),
		Secret: apiKey.GetSecret(),
		Slug:   apiKey.GetSlug(),
	}

	return response, nil
}

func (as *authService) RefreshToken(refreshToken string) (*response.LoginResponse, *rest_err.RestErr) {
	logger.Info("Refreshing token")

	jwtAdapter := internalJwt.NewJWTAdapter()

	token, restErr := jwtAdapter.ParseToken(refreshToken)
	if restErr != nil {
		return nil, restErr
	}

	claims := token.Claims.(jwt.MapClaims)

	tokenId, ok := claims["tid"].(string)
	if !ok {
		return nil, rest_err.NewUnauthorizedError("invalid token")
	}

	apiKey, err := claims.GetSubject()
	if err != nil {
		return nil, rest_err.NewBadRequestError("invalid token")
	}

	err = uuid.Validate(apiKey)
	if err != nil {
		return nil, rest_err.NewBadRequestError("invalid token")
	}

	restErr = as.authRepository.DeleteToken(tokenId)
	if restErr != nil {
		return nil, restErr
	}

	newRefreshToken := entities.NewToken(
		uuid.NewString(),
		apiKey,
	)

	_, restErr = as.authRepository.CreateToken(newRefreshToken)
	if restErr != nil {
		return nil, restErr
	}

	refreshToken, restErr = jwtAdapter.GenerateToken(newRefreshToken.GetUUID(), newRefreshToken.GetApiKey())
	if restErr != nil {
		return nil, restErr
	}

	newAccessToken := entities.NewToken(
		uuid.NewString(),
		apiKey,
	)

	accessToken, restErr := jwtAdapter.GenerateToken(
		newAccessToken.GetUUID(),
		newAccessToken.GetApiKey(),
	)
	if restErr != nil {
		return nil, restErr
	}

	return &response.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (as *authService) FindApiKey(apiKey string) (entities.ApiKey, *rest_err.RestErr) {
	return as.authRepository.FindApiKey(apiKey)
}

func (as *authService) Login(apiKey string) (*response.LoginResponse, *rest_err.RestErr) {
	var restErr *rest_err.RestErr
	jwtAdapter := internalJwt.NewJWTAdapter()

	newRefreshToken := entities.NewToken(
		uuid.NewString(),
		apiKey,
	)

	_, restErr = as.authRepository.CreateToken(newRefreshToken)
	if restErr != nil {
		return nil, restErr
	}

	refreshToken, restErr := jwtAdapter.GenerateToken(newRefreshToken.GetUUID(), newRefreshToken.GetApiKey())
	if restErr != nil {
		return nil, restErr
	}

	newAccessToken := entities.NewToken(
		uuid.NewString(),
		apiKey,
	)

	accessToken, restErr := jwtAdapter.GenerateToken(
		newAccessToken.GetUUID(),
		newAccessToken.GetApiKey(),
	)
	if restErr != nil {
		return nil, restErr
	}

	return &response.LoginResponse{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil

}
