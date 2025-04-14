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
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/bootstrap ./main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

migration-run:
	go run main.go migration-up

migration-down:
	go run main.go migration-down $(timesToDownMigration)

migration-status:
	go run main.go migration-status