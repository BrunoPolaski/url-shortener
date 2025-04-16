package request

import (
	"strconv"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
)

type ApiKey struct {
	UUID      string `json:"uuid" binding:"required"`
	Secret    string `json:"secret,omitempty"`
	Slug      string `json:"slug,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (a *ApiKey) ToDomain() entities.ApiKey {
	return entities.NewApiKey(
		a.UUID,
		a.Secret,
		a.Slug,
	)
}

func ApiKeyFromDomain(apiKey entities.ApiKey) *ApiKey {
	createdAt := strconv.Itoa(
		int(apiKey.GetCreatedAt().Unix()),
	)

	return &ApiKey{
		UUID:      apiKey.GetUUID(),
		Secret:    apiKey.GetSecret(),
		Slug:      apiKey.GetSlug(),
		CreatedAt: createdAt,
	}
}
