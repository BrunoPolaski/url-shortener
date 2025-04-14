package cmd

import (
	"github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/database"
	"github.com/BrunoPolaski/migrator/migrator"
)

const migrationsPath = "internal/config/migrations"

func Migrate() {
	db := database.NewMySQLAdapter()

	migrator.Exec(
		db.Connect(),
		migrationsPath,
	)
}
