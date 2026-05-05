# Smart Document Processing System

System that processes real-world business documents (invoices and purchase orders), extracts structured data, validates it, and presents it through a simple interface.

## Features

- Document ingestion from PDF, CSV, TXT files
- Data extraction for key fields
- Validation engine with issue detection
- Review and approval workflow
- REST API with Swagger documentation
- Docker support
- Database migrations with go-migrate

## Architecture

### Backend
Follows **Clean Architecture** with a hexagonal architecture design, separating concerns into two layers:

- **Domain** — core business logic, models, and validation rules with no external dependencies
- **Adapters** — all external concerns including HTTP handlers, PostgreSQL persistence, file processing, and database migrations

### Frontend
Built with a custom **React CLI** ([`@aditokmo/create-react-project`](https://www.npmjs.com/package/@aditokmo/create-react-project)) that builds a feature-based, pre-configured architecture. Development was done  by a **GitHub Copilot** local agent configured with a custom ruleset to enforce consistency across the codebase.

## API Endpoints

- `POST /documents` - Upload and process a document
- `GET /documents` - List all documents
- `GET /documents/{id}` - Get document details
- `PUT /documents/{id}` - Update document
- `POST /documents/{id}/approve` - Approve document
- `POST /documents/{id}/reject` - Reject document


## Prerequisites

- Go 1.25.5 or higher
- Node.js 18.0 or higher
- pnpm 8.0 or higher
- Docker

## Quick Start

```bash
# Clone repo
git clone https://github.com/aditokmo/smart-document-processing-system.git

# Enter project folder
cd smart-document-processing-system

# Install dependencies
make install

# Start all services (database, backend, frontend)
make dev
```

Then open:
- Frontend: http://localhost:5174
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger

## Make Commands

- `make install` - Install dependencies
- `make dev` - Start all services in development mode
- `make start` - Start all services in production mode
- `make stop` - Stop all services
- `make restart` - Restart all services
- `make clean` - Stop and clean up
