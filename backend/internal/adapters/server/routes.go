package server

import (
	"backend/internal/adapters/handler"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Router(
	docHandler *handler.DocumentHandler,
) *httprouter.Router {
	r := httprouter.New()

	r.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Documents
	r.POST("/api/v1/documents", docHandler.UploadDocument)
	r.GET("/api/v1/documents", docHandler.GetDocuments)
	r.GET("/api/v1/documents/:id", docHandler.GetDocument)
	r.PUT("/api/v1/documents/:id", docHandler.UpdateDocument)
	r.POST("/api/v1/documents/:id/approve", docHandler.ApproveDocument)
	r.POST("/api/v1/documents/:id/reject", docHandler.RejectDocument)

	// Swagger routes
	r.GET("/swagger", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "docs/index.html")
	})

	r.ServeFiles("/swagger/*filepath", http.Dir("docs"))

	return r
}
