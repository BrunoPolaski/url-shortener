package request

import (
	"strconv"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
)

type Token struct {
	UUID      string `json:"uuid" binding:"required"`
	ApiKey    string `json:"token,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

func (a *Token) ToDomain() entities.Token {
	return entities.NewToken(
		a.UUID,
		a.ApiKey,
	)
}

func TokenFromDomain(token entities.Token) *Token {
	createdAt := strconv.Itoa(
		int(token.GetCreatedAt().Unix()),
	)

	return &Token{
		UUID:      token.GetUUID(),
		ApiKey:    token.GetApiKey(),
		CreatedAt: createdAt,
	}
}
