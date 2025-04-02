name ?= DefaultMigration
timesToDownMigration ?= 1 

migration-gen:
	@echo "Creating migration"
	$(eval timestamp := $(shell date +%s))
	@touch internal/config/migrations/$(timestamp)_$(name)_down.sql
	@touch internal/config/migrations/$(timestamp)_$(name)_up.sql
	@echo "Migration file created successfully"

.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/main ./main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

migration-run:
	go run ./internal/config/migrations/migrator.go up

migration-down:
	go run ./internal/config/migrations/migrator.go down $(timesToDownMigration)

migration-status:
	go run ./internal/config/migrations/migrator.go status