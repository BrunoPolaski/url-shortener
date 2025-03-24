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
	return &entities.Redirect{}, nil // TODO: implement
}
