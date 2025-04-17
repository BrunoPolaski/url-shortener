package mysql

import (
	"database/sql"
	"fmt"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/model/request"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type authRepositoryMySQL struct {
	database *sql.DB
}

func NewAuthRepositoryMySQL(db *sql.DB) repositories.AuthRepository {
	return &authRepositoryMySQL{
		database: db,
	}
}

func (ar *authRepositoryMySQL) FindToken(uuid string) (string, *rest_err.RestErr) {
	stmt, err := ar.database.Prepare("SELECT url FROM auth WHERE uuid = ?")
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	var url string
	err = stmt.QueryRow(uuid).Scan(&url)
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to query row: %s", err.Error()))
	}

	redirect := entities.NewRedirect(uuid, url)

	if err := redirect.DecryptURL(); err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to decrypt url: %s", err.Error()))
	}

	return redirect.GetURL(), nil
}

func (ar *authRepositoryMySQL) FindApiKey(apiKey string) (entities.ApiKey, *rest_err.RestErr) {
	stmt, err := ar.database.Prepare("SELECT * FROM api_key WHERE uuid = ?")
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	var apiKeyEntity request.ApiKey
	err = stmt.QueryRow(apiKey).Scan(&apiKeyEntity)
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to query row: %s", err.Error()))
	}

	return apiKeyEntity.ToDomain(), nil
}

func (ar *authRepositoryMySQL) CreateApiKey(apiKey entities.ApiKey) (entities.ApiKey, *rest_err.RestErr) {
	stmt, err := ar.database.Prepare("INSERT INTO api_key (uuid, secret, slug) VALUES (?, ?, ?)")
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	if err := apiKey.EncryptSecret(); err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to encrypt secret: %s", err.Error()))
	}

	_, err = stmt.Exec(apiKey.GetUUID(), apiKey.GetSecret(), apiKey.GetSlug())
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to execute statement: %s", err.Error()))
	}

	return apiKey, nil
}

func (ar *authRepositoryMySQL) CreateToken(token entities.Token) (entities.Token, *rest_err.RestErr) {
	stmt, err := ar.database.Prepare("INSERT INTO refresh_token (uuid, api_key) VALUES (?, ?)")
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	_, err = stmt.Exec(token.GetUUID(), token.GetApiKey())
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to execute statement: %s", err.Error()))
	}

	return token, nil
}

func (ar *authRepositoryMySQL) FindTokenByApiKey(apiKey string) (string, *rest_err.RestErr) {
	stmt, err := ar.database.Prepare("SELECT uuid FROM refresh_token WHERE api_key = ?")
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	var tokenUUID string
	err = stmt.QueryRow(apiKey).Scan(&tokenUUID)
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to query row: %s", err.Error()))
	}

	return tokenUUID, nil
}
