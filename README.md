# Pokemon TCG Search Engine

An unofficial Pokemon Trading Card Game Recommendation System for
Agentic Workflows

## Database migrations

Goose does not load `.env` files automatically. Use the Make targets in this
repository, which load `.env` before invoking Goose:

```sh
make db-status
make db-up
make db-down
make db-validate
```

For another Goose command, pass it through the generic target:

```sh
make migration COMMAND="status"
make migration COMMAND="up-to 20260929152133"
```

Alternatively, export the variables in your current shell first:

```sh
set -a
. ./.env
set +a
goose status
```

`.env` is intentionally ignored by git because it can contain database
credentials. Keep a local copy with `GOOSE_DRIVER`, `GOOSE_DBSTRING`, and
`GOOSE_MIGRATION_DIR` defined.
