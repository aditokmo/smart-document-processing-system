// Package main Smart Document Processing System API
//
//	@title			Smart Document Processing System API
//	@version		1.0
//	@description	A system for processing business documents with data extraction and validation
//	@termsOfService	http://swagger.io/terms/
//
//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io
//
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@host		localhost:8080
//	@BasePath	/api/v1
//
//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
package main

import (
	"backend/internal/adapters/config"
	"backend/internal/adapters/migration"
	"backend/internal/adapters/postgres"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	server "backend/internal/adapters/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		os.Exit(1)
	}

	if err := migration.Run(cfg.Database.URL(), "migrations"); err != nil {
		slog.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	slog.Info("Migrations applied")

	db, err := postgres.Connection(cfg.Database.ConnectionURL())
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv, err := server.New(cfg, db)
	if err != nil {
		slog.Error("Failed to build server", "error", err)
		os.Exit(1)
	}

	// Create context that listens for OS signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run server in goroutine
	go func() {
		slog.Info("Starting server", "port", cfg.Port)

		if err := srv.Run(":" + cfg.Port); err != nil {
			slog.Error("Server failed", "error", err)
			stop()
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutdown signal received")

	// Create timeout context for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
	} else {
		slog.Info("Server stopped gracefully")
	}
}
