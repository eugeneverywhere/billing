PAYMENT_PROCESSOR_NAME 			:= billing
MIGRATION_SERVICE_NAME          := migrate
NAMESPACE	   									:= "default"
CONFIG         								:= $(wildcard local.yml)
PKG            								:= github.com/eugeneverywhere/billing
PKG_LIST       								:= $(shell go list ${PKG}/... | grep -v /vendor/)

all: setup test build

setup: ## Installing all service dependencies.
	echo "Setup..."
	GO111MODULE=on go mod vendor

config: ## Creating the local config yml.
	echo "Creating local config yml ..."
	cp config.example.yml local.yml

build: ## Build the executable file of service.
	echo "Building..."
	cd cmd/$(PAYMENT_PROCESSOR_NAME) && go build

run: ## Run service with local config.
	make build
	echo "Running..."
	cd cmd/$(PAYMENT_PROCESSOR_NAME) && ./$(PAYMENT_PROCESSOR_NAME) -config=../../local.yml


db\:migrate: ## Run migrations.
	cd cmd/$(MIGRATION_SERVICE_NAME) && go build && ./$(MIGRATION_SERVICE_NAME) -config=../../local.yml -migrate-path=../../db/migrations

db\:downgrade: ## Drop all migrations.
	cd cmd/$(MIGRATION_SERVICE_NAME) && go build && ./$(MIGRATION_SERVICE_NAME) -config=../../local.yml -drop -migrate-path=../../db/migrations

db\:recreate: ## Recreate the db
	make db:downgrade
	make db:migrate
	echo "Recreating complete"


test: ## Run tests for all packages.
	echo "Testing..."
	go test -race ${PKG_LIST}

clean: ## Cleans the temp files and etc.
	echo "Clean..."
	rm -f cmd/$(PAYMENT_PROCESSOR_NAME)/$(PAYMENT_PROCESSOR_NAME)

help: ## Display this help screen
	grep -E '^[a-zA-Z_\-\:]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ": .*?## "}; {gsub(/[\\]*/,""); printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	