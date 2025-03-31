name ?= DefaultMigration

migration-create:
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
