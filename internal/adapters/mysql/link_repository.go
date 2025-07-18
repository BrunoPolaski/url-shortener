package mysql

import (
	"database/sql"
	"fmt"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/ports/repositories"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type linkRepositoryMySQL struct {
	database *sql.DB
}

func NewLinkRepositoryMySQL(db *sql.DB) repositories.LinkRepository {
	return &linkRepositoryMySQL{
		database: db,
	}
}

func (lr *linkRepositoryMySQL) GetByUUID(uuid string) (string, *rest_err.RestErr) {
	stmt, err := lr.database.Prepare("SELECT url FROM link WHERE uuid = ?")
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

func (lr *linkRepositoryMySQL) Create(redirect entities.Redirect) (entities.Redirect, *rest_err.RestErr) {
	stmt, err := lr.database.Prepare("INSERT INTO link(uuid, url) VALUES(?, ?)")
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	if err := redirect.EncryptURL(); err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to encrypt url: %s", err.Error()))
	}

	result, err := stmt.Exec(redirect.GetUUID(), redirect.GetURL())
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to insert to database: %s", err.Error()))
	}

	rows, err := result.RowsAffected()
	if rows != 1 {
		return nil, rest_err.NewInternalServerError("Insert failed for url. Please contact the support team.")
	} else if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("Error trying to get rows affected by insert: %s", err.Error()))
	}

	return redirect, nil
}
