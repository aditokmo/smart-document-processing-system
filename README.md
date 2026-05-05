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

## Improvements I would make

- Refining the UI and cleaning up the React codebase, as the initial version was built with a free GitHub Copilot agent
- Adding support for image formats (that I will probably do in couple days)
- Supporting multi-file upload once the UX flow is properly planned out
- Continued testing to surface and fix any remaining bugs in the system if they exist
- Adding upload validation to prevent confusion around submitting multiple files before multi-upload is fully supported

## AI Usage

- Building the frontend by combining my custom React CLI, which builds the full folder structure with pre-configured libraries, and GitHub Copilot Agents to generate code following that structure and the rules I defined through prompts
- Generating Swagger documentation and README docs
- Reviewing backend logic to catch anything I may have missed, and for occasional syntax guidance since I am still learning Go
- Setting up Docker and Makefile configuration to make it easier for others to run the project locally

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
