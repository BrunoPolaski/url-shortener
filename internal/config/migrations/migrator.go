package migrations

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/database"
)

func RunMigrations() {
	var migrations []string
	var completed []string

	db := database.NewMySQLAdapter()
	conn := db.Connect()
	defer conn.Close()

	fmt.Println("Connected to database")

	files, err := filepath.Glob("internal/config/migrations/*.sql")
	if err != nil {
		fmt.Printf("Error reading migration files %v\n", err)
		return
	}

	sort.Strings(files)

	migrations = readMigrationsMetadata(conn)

	var migrationErr error
	for _, file := range files {
		filename := filepath.Base(file)
		if !contains(filename, migrations) {
			fmt.Printf("Running migration %s...\n", filename)
			if migrationErr = runMigration(conn, file); migrationErr != nil {
				break
			}

			completed = append(completed, file)
		}
	}

	if len(completed) == 0 {
		if migrationErr != nil {
			return
		}
		fmt.Printf("\nNo migrations to run!!! XD \n")
		return
	}

	fmt.Printf("Completed migrations:\n - %s\n", strings.Join(completed, "\n - "))
}

func readMigrationsMetadata(conn *sql.DB) []string {
	var migrations []string

	table := `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`

	_, err := conn.Exec(table)
	if err != nil {
		fmt.Printf("Error creating migrations table: %v\n", err)
		return migrations
	}

	rows, err := conn.Query("SELECT name FROM migrations")
	if err != nil {
		fmt.Printf("Error reading migrations metadata: %v\n", err)
		return migrations
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Printf("Error scanning migrations metadata: %v\n", err)
			return migrations
		}

		migrations = append(migrations, name)
	}

	return migrations
}

func runMigration(conn *sql.DB, file string) error {
	query, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Error reading migration file: %v\n", err)
		return err
	}

	content, err := os.ReadFile(query)
	if err != nil {
		fmt.Printf("Error reading migration file: %v\n", err)
		return err
	}

	transaction, err := conn.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: %v\n", err)
		return err
	}

	_, err = transaction.Exec(string(content))
	if err != nil {
		fmt.Printf("Error running migration: %v\n", err)
		transaction.Rollback()
		return err
	}

	fmt.Printf("Migration %s ran successfully\n", filepath.Base(file))

	_, err = conn.Exec("INSERT INTO migrations (name) VALUES (?)", filepath.Base(file))
	if err != nil {
		fmt.Printf("Error updating migrations metadata: %v\n", err)
		transaction.Rollback()
		return err
	}

	if err := transaction.Commit(); err != nil {
		fmt.Printf("Error committing transaction: %v\n", err)
		transaction.Rollback()
		return err
	}

	return nil
}

func contains(s string, a []string) bool {
	for _, k := range a {
		if s == k {
			return true
		}
	}
	return false
}
