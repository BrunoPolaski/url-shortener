package repositories

import (
	"database/sql"
	"fmt"

	"github.com/AdagaDigital/url-redirect-service/internal/domain/entities"
	"github.com/AdagaDigital/url-redirect-service/internal/domain/repositories"
	"github.com/BrunoPolaski/go-crud/src/configuration/rest_err"
)

type linkRepository struct {
	database *sql.DB
}

func NewLinkRepository(db *sql.DB) repositories.LinkRepository {
	return &linkRepository{
		database: db,
	}
}

func (lr *linkRepository) GetByUUID(uuid string) (string, *rest_err.RestErr) {
	stmt, err := lr.database.Prepare("SELECT url FROM links WHERE uuid = ?")
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	var url string
	err = stmt.QueryRow(uuid).Scan(&url)
	if err != nil {
		return "", rest_err.NewInternalServerError(fmt.Sprintf("error when trying to query row: %s", err.Error()))
	}

	return url, nil
}

func (lr *linkRepository) Create(redirect *entities.Redirect) (*entities.Redirect, *rest_err.RestErr) {
	stmt, err := lr.database.Prepare("INSERT INTO links(uuid, url) VALUES(?, ?)")
	if err != nil {
		return nil, rest_err.NewInternalServerError(fmt.Sprintf("error when trying to prepare statement: %s", err.Error()))
	}

	defer stmt.Close()

	result, err := stmt.Exec(redirect.UUID, redirect.URL)
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
