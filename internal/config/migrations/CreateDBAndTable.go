package migrations

import (
	"database/sql"
)

type MigrationCreateDBAndTable struct {
	Database *sql.DB
}

func (m *MigrationCreateDBAndTable) Up() error {
	_, err := m.Database.Exec(`
		CREATE DATABASE IF NOT EXISTS redirect_service;
		USE redirect_service;

		
		CREATE TABLE IF NOT EXISTS redirects (
			uuid VARCHAR(36) PRIMARY KEY,
			url TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

func (m *MigrationCreateDBAndTable) Down() error {
	return nil // Your code here...
}
