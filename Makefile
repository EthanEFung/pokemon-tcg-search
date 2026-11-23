.SHELL := /bin/sh

# Goose does not read .env itself. Make loads it and exports the Goose
# variables to the CLI process.
-include .env
export GOOSE_DRIVER GOOSE_DBSTRING GOOSE_MIGRATION_DIR

.PHONY: migration db-status db-up db-down db-reset db-version db-validate

# Usage: make migration COMMAND="up-to 20260929152133"
migration:
	@test -n "$(COMMAND)" || (echo 'usage: make migration COMMAND="status"' >&2; exit 1)
	goose $(COMMAND)

db-status:
	goose status

db-up:
	goose up

db-down:
	goose down

db-reset:
	goose reset

db-version:
	goose version

db-validate:
	goose validate
