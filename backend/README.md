# Smart Document Processing System

A Go-based backend system for processing business documents (invoices and purchase orders) with data extraction, validation, and persistence.

## Features

- Document ingestion from PDF, CSV, TXT files
- Data extraction for key fields
- Validation engine with issue detection
- Review and approval workflow
- REST API with Swagger documentation
- Docker support
- Database migrations with go-migrate

## Architecture

This project follows a Clean Architecture with hexagonal-inspired core:

- **Domain**: Business entities and core logic
- **Application**: Use cases and application services
- **Adapters**: External concerns (DB, HTTP, file processing, migrations)

## API Endpoints

- `POST /documents/upload` - Upload and process a document
- `GET /documents` - List all documents
- `GET /documents/{id}` - Get document details
- `PUT /documents/{id}` - Update document
- `POST /documents/{id}/approve` - Approve document
- `POST /documents/{id}/reject` - Reject document

## Running the Application

### With Docker

```bash
docker-compose up --build
```

### Locally

```bash
go mod tidy
go run cmd/server/main.go
```

## Dependencies

- github.com/julienschmidt/httprouter
- github.com/jackc/pgx/v5/stdlib
- github.com/golang-migrate/migrate/v4 (for database migrations)
- github.com/swaggo/http-swagger (for Swagger docs)

## Database Migrations

This project uses `go-migrate` for database schema management. Migration files are stored in the `migrations/` directory with the naming convention `{version}_{description}.{up|down}.sql`.

### Running Migrations

Migrations are automatically executed when the application starts. The application will:
1. Connect to PostgreSQL
2. Run all pending migrations
3. Initialize the schema

### Creating New Migrations

To create a new migration:

```bash
# Create migration files (requires go-migrate CLI)
migrate create -ext sql -dir migrations -seq create_new_table
```

Then edit the generated `.up.sql` and `.down.sql` files in the `migrations/` directory.

## Running with PostgreSQL

The application now uses PostgreSQL by default. Configure environment variables as needed:

- `DB_DRIVER=postgres`
- `DB_HOST=localhost`
- `DB_PORT=5432`
- `DB_USER=postgres`
- `DB_PASSWORD=postgres`
- `DB_NAME=documentsdb`
- `DB_SSLMODE=disable`

## Development

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Ensure PostgreSQL is running locally or use Docker: `docker-compose up -d db`
4. Run the server: `go run cmd/server/main.go`
   - Migrations will run automatically on startup
5. Access API at http://localhost:8080
6. Swagger UI at http://localhost:8080/swagger/
