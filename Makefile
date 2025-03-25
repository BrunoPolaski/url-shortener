name ?= DefaultMigration

migration-create:
	@echo "Creating migration"
	@touch internal/config/migrations/$(name).go

	@echo "package migrations" > internal/config/migrations/$(name).go
	
	@echo "" >> internal/config/migrations/$(name).go
	
	@echo "import (" >> internal/config/migrations/$(name).go
	@echo "	\"database/sql\"" >> internal/config/migrations/$(name).go
	@echo ")" >> internal/config/migrations/$(name).go

	@echo "" >> internal/config/migrations/$(name).go
	
	@echo "type Migration$(name) struct {" >> internal/config/migrations/$(name).go
	@echo "	Database *sql.DB" >> internal/config/migrations/$(name).go
	@echo "}" >> internal/config/migrations/$(name).go
	
	@echo "" >> internal/config/migrations/$(name).go
	
	@echo "func (m *Migration$(name)) Up() error {" >> internal/config/migrations/$(name).go
	@echo "	return nil // Your code here..." >> internal/config/migrations/$(name).go
	@echo "}" >> internal/config/migrations/$(name).go
	
	@echo "" >> internal/config/migrations/$(name).go
	
	@echo "func (m *Migration$(name)) Down() error {" >> internal/config/migrations/$(name).go
	@echo "	return nil // Your code here..." >> internal/config/migrations/$(name).go
	@echo "}" >> internal/config/migrations/$(name).go

	@echo "Migration file created successfully"

