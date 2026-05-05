package server

import (
	"backend/internal/adapters/config"
	"backend/internal/adapters/handler"
	"backend/internal/adapters/middleware"
	"backend/internal/adapters/repository"
	"backend/internal/domain/service"
	"context"
	"database/sql"
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
}

func New(cfg *config.Config, db *sql.DB) (*Server, error) {
	// Repositories
	documentRepo := repository.NewDocumentRepository(db)

	// Services
	documentExtractor := repository.NewExtractionService()
	documentFileProc := repository.NewFileProcessor(documentExtractor)

	// Services
	documentService := service.NewDocumentService(documentRepo, documentExtractor, documentFileProc)

	// Handlers
	documentHandler := handler.NewDocumentHandler(documentService)

	// Router
	r := Router(documentHandler)

	// Middlewares
	var h http.Handler = r
	h = middleware.CORS(h)

	return &Server{
		http: &http.Server{
			Handler:           h,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       120 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}, nil
}

func (s *Server) Run(addr string) error {
	s.http.Addr = addr
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
