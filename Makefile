.PHONY: help install dev start stop restart clean logs db-up db-down backend frontend

help:
	@echo "Smart Document Processing System - Make Commands"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  install        Install all dependencies (Go modules + pnpm packages)"
	@echo "  dev            Start database, backend, and frontend (development mode)"
	@echo "  start          Start database, backend, and frontend"
	@echo "  stop           Stop all running services"
	@echo "  restart        Restart all services"
	@echo "  clean          Stop services and clean up"
	@echo "  db-up          Start only the database"
	@echo "  db-down        Stop only the database"
	@echo "  backend        Start only the backend server"
	@echo "  frontend       Start only the frontend dev server"
	@echo "  logs           Show logs from all running services"
	@echo ""

install:
	@echo "Installing Go dependencies..."
	cd backend && go mod tidy
	@echo "Installing pnpm dependencies..."
	cd frontend && pnpm install
	@echo "✓ Dependencies installed successfully"

dev: db-up
	@echo "Starting Smart Document Processing System..."
	@echo ""
	@echo "Starting backend on http://localhost:8080"
	@echo "Starting frontend on http://localhost:5174"
	@echo "Database: PostgreSQL on localhost:5432"
	@echo ""
	@echo "Press Ctrl+C to stop all services"
	@echo ""
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'Set-Location backend; air'"
	@cd frontend && pnpm dev
else
	@cd backend && air & cd frontend && pnpm dev
endif

start: db-up
	@echo "Starting Smart Document Processing System..."
	@echo ""
	@echo "Starting backend on http://localhost:8080"
	@echo "Starting frontend on http://localhost:5174"
	@echo "Database: PostgreSQL on localhost:5432"
	@echo ""
	@echo "Press Ctrl+C to stop all services"
	@echo ""
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Process powershell -ArgumentList '-NoExit', '-Command', 'Set-Location backend; go run cmd/server/main.go'"
	@cd frontend && pnpm dev
else
	@cd backend && go run cmd/server/main.go & cd frontend && pnpm dev
endif

db-up:
	@echo "Starting PostgreSQL database..."
	cd backend && docker-compose up -d db
	@echo "⏳ Waiting for database to be ready..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep -Seconds 5"
else
	@sleep 5
endif
	@echo "✓ Database is ready"

db-down:
	@echo "Stopping PostgreSQL database..."
	cd backend && docker-compose down
	@echo "✓ Database stopped"

backend:
	@echo "Starting backend server on http://localhost:8080"
	cd backend && go run cmd/server/main.go

frontend:
	@echo "Starting frontend dev server on http://localhost:5174"
	cd frontend && pnpm dev

stop:
	@echo "Stopping all services..."
ifeq ($(OS),Windows_NT)
	@powershell -Command "Get-Process | Where-Object {$$_.CommandLine -like '*air*' -or $$_.CommandLine -like '*pnpm*'} | Stop-Process -Force" 2>nul || true
else
	@pkill -f "air" || true
	@pkill -f "go run cmd/server/main.go" || true
	@pkill -f "pnpm dev" || true
endif
	cd backend && docker-compose down || true
	@echo "✓ All services stopped"

restart: stop
ifeq ($(OS),Windows_NT)
	@powershell -Command "Start-Sleep -Seconds 2"
else
	@sleep 2
endif
	@$(MAKE) start

clean: stop
	@echo "Cleaning up..."
	cd backend && docker-compose down -v || true
	@echo "✓ Cleanup complete"

logs:
	cd backend && docker-compose logs -f

.DEFAULT_GOAL := help