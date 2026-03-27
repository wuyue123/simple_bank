i# AGENTS.md - Simple Bank API

## Architecture Overview
- **Framework**: Gin-based REST API server for banking operations
- **Database**: PostgreSQL with sqlc-generated type-safe Go queries
- **Config**: Viper loads from `app.env` file (DB_DRIVER, DB_SOURCE, SERVER_ADDRESS)
- **Modules**:
  - `api/`: HTTP handlers and routing (server.go, account.go, user.go, transfer.go)
  - `db/sqlc/`: Auto-generated database queries and transactions (store.go extends with TransferTx)
  - `db/migration/`: Schema migrations using golang-migrate
  - `util/`: Config loading, password hashing, currency validation

## Key Patterns
- **Database Transactions**: Use `store.TransferTx()` for money transfers (creates transfer + entries + balance updates atomically)
- **Custom Validators**: Register currency validation in `api/server.go` using `validCurrency` from `util/currency.go`
- **Query Generation**: Write SQL in `db/queries/*.sql` with named parameters (:one, :many), run `make sqlc` to regenerate
- **Mocking**: Use `mockgen` to generate mocks in `db/mock/store.go`, test with gomock and testify
- **Error Handling**: Return JSON errors via `errorResponse()` in handlers

## Developer Workflows
- **Setup DB**: `make network postgres createdb migrateup`
- **Generate Code**: `make sqlc` after editing queries, `make mock` for test mocks
- **Run**: `make run` or `go run main.go`
- **Test**: `make test` (uses `-v -cover -short`)
- **Migrate**: `make migrateup1` / `migratedown1` for incremental changes

## Conventions
- **Imports**: Use `db "pxsemic.com/simplebank/db/sqlc"` for database package
- **Struct Tags**: sqlc emits `json` and `db` tags on generated structs
- **Amounts**: Stored as int64 (e.g., 100 = $1.00), use util functions for formatting
- **Passwords**: Hash with `util.HashPassword()`, verify with `util.CheckPassword()`
- **Foreign Keys**: Deferrable constraints on transfers/accounts for transaction safety

## Examples
- **Add API Route**: Register in `api/server.go` NewServer(), implement handler in separate file
- **New Query**: Add to `db/queries/account.sql`, run `make sqlc`, use in store methods
- **Transaction**: Extend `store.execTx()` for multi-step DB operations like TransferTx
- **Test Handler**: Use httptest.ResponseRecorder, mock store with gomock, assert with require</content>
<parameter name="filePath">D:\H\goproject\src\simpleblankv2\AGENTS.md
