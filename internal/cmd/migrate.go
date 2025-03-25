package cmd

import "github.com/AdagaDigital/url-redirect-service/internal/config/migrations"

func Migrate() {
	migrations.RunMigrations()
}
