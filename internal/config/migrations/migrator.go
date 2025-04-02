package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/AdagaDigital/url-redirect-service/internal/config/logger"
	"github.com/AdagaDigital/url-redirect-service/internal/infra/thirdparty/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file. Error: %s", err)
	}

	logger.Init()

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run migrator.go <up|down|status>")
		return
	}

	argToFunction := map[string]func(*sql.DB){
		"up":     migrationUp,
		"down":   migrationDown,
		"status": migrationStatus,
	}

	if migrationFunc, exists := argToFunction[os.Args[1]]; exists {
		if os.Getenv("ENV") == "local" {
			os.Setenv("DB_HOST", "localhost")
		}

		db := database.NewMySQLAdapter()
		conn := db.Connect()
		defer conn.Close()

		fmt.Printf("\n---------\nDATABASE CONNECTION SUCCESSFUL\n---------\n")
		migrationFunc(conn)
	} else {
		fmt.Println("Invalid argument. Use up|down|status.")
	}
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

func runMigration(tx *sql.Tx, file string) error {
	query, err := filepath.Abs(file)
	if err != nil {
		return fmt.Errorf("error getting absolute path of file %s: \n\n%v", file, err)
	}

	content, err := os.ReadFile(query)
	if err != nil {
		return fmt.Errorf("error reading file %s: \n\n%v", file, err)
	}

	_, err = tx.Exec(string(content))
	if err != nil {
		if os.Args[1] == "up" {
			return fmt.Errorf("error running migration %s: \n\n%v", file, err)
		} else {
			return fmt.Errorf("error rolling back migration %s: \n\n%v", file, err)
		}
	}

	if os.Args[1] == "up" {
		fmt.Printf("Migration %s ran successfully\n", filepath.Base(file))

		_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", filepath.Base(file))
		if err != nil {
			return fmt.Errorf("error inserting migration %s into database: \n\n%v", file, err)
		}
	} else {
		_, err = tx.Exec("DELETE FROM migrations WHERE name = ?", filepath.Base(strings.Replace(file, "_down.sql", "_up.sql", 1)))
		if err != nil {
			return fmt.Errorf("error deleting migration %s from database: \n\n%v", file, err)
		}

		fmt.Printf("Migration %s rolled back successfully\n", filepath.Base(file))
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

func migrationUp(conn *sql.DB) {
	files, err := filepath.Glob("internal/config/migrations/*_up.sql")
	if err != nil {
		fmt.Printf("Error reading migration files %v\n", err)
		return
	}

	slices.Sort(files)

	migrations := readMigrationsMetadata(conn)
	var completed []string

	var migrationErr error

	tx, err := conn.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: \n\n%v\n", err)
		return
	}

	for _, file := range files {
		filename := filepath.Base(file)
		if contains(filename, migrations) {
			continue
		}

		fmt.Printf("Running migration %s...\n", filename)
		if migrationErr = runMigration(tx, file); migrationErr != nil {
			redBackground(migrationErr.Error())
			if err := tx.Rollback(); err != nil {
				fmt.Printf("Error rolling back transaction: \n\n%v\n", err)
			}
			return
		}

		completed = append(completed, file)
	}

	if len(completed) == 0 {
		greenBackground("No migrations to run!!! XD")
		return
	}

	greenBackground(fmt.Sprintf("Completed migrations:\n - %s\n", strings.Join(completed, "\n - ")))
}

func migrationDown(conn *sql.DB) {
	files, err := filepath.Glob("internal/config/migrations/*_down.sql")
	if err != nil {
		fmt.Printf("Error reading migration files: \n\n%v\n", err)
		return
	}

	slices.Sort(files)

	migrations := readMigrationsMetadata(conn)
	var completed []string

	var migrationErr error

	tx, err := conn.Begin()
	if err != nil {
		fmt.Printf("Error starting transaction: \n\n%v\n", err)
		return
	}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		filename := strings.Replace(file, "_down.sql", "_up.sql", 1)

		if !contains(filepath.Base(filename), migrations) {
			continue
		}

		fmt.Printf("Running migration %s...\n", file)
		if migrationErr = runMigration(tx, file); migrationErr != nil {
			redBackground(migrationErr.Error())
			if err := tx.Rollback(); err != nil {
				fmt.Printf("Error rolling back transaction: \n\n%v\n", err)
			}
			return
		}

		completed = append(completed, file)
	}

	if len(completed) == 0 {
		if migrationErr != nil {
			return
		}
		greenBackground("No migrations to roll back!!! XD")
		return
	}

	greenBackground(fmt.Sprintf("Rolled back migrations:\n - %s\n", strings.Join(completed, "\n - ")))
}

func migrationStatus(conn *sql.DB) {
	files, err := filepath.Glob("internal/config/migrations/*_up.sql")
	if err != nil {
		fmt.Printf("Error reading migration files \n\n%v\n", err)
		return
	}

	sort.Strings(files)

	migrations := readMigrationsMetadata(conn)

	var completedCounter int
	var notCompletedCounter int
	fmt.Printf("Migrations status:\n")
	for _, file := range files {
		filename := filepath.Base(file)
		if contains(filename, migrations) {
			completedCounter++
		} else {
			notCompletedCounter++
		}
	}

	fmt.Printf("Completed migrations: %d\n", completedCounter)
	fmt.Printf("Not completed migrations: %d\n", notCompletedCounter)
	if len(migrations) == 0 {
		fmt.Printf("You didn't run any migration yet.\n")
		return
	}
	fmt.Printf("Latest version: %s\n", migrations[len(migrations)-1])
}

func greenBackground(slice ...string) {
	fmt.Printf("\033[0;102m\n%s\n\033[0m", strings.Join(slice, " "))
}

func redBackground(slice ...string) {
	fmt.Printf("\033[0;101m\n%s\n\033[0m", strings.Join(slice, " "))
}
