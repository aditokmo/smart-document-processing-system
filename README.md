# Smart Document Processing System

## Prerequisites

- Go 1.25.5 or higher
- Node.js 18.0 or higher
- pnpm 8.0 or higher
- Docker

## Quick Start

```bash
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