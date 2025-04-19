package request

import (
	"strconv"
	"time"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
)

type ApiKey struct {
	UUID      string `json:"uuid" binding:"required"`
	Secret    string `json:"secret,omitempty"`
	Slug      string `json:"slug,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (a *ApiKey) ToDomain() entities.ApiKey {
	createdAt, _ := strconv.Atoi(a.CreatedAt)

	return entities.NewApiKeyWithCreatedAt(
		a.UUID,
		a.Secret,
		a.Slug,
		time.Unix(int64(createdAt), 0),
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
