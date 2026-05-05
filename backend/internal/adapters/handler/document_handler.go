package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"backend/internal/adapters/render"
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"backend/internal/domain/service"

	"github.com/julienschmidt/httprouter"
)

type DocumentHandler struct {
	service *service.DocumentService
}

func NewDocumentHandler(service *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

// UploadDocument handles file upload
// @Summary Upload a document
// @Description Upload and process a document file
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Success 200 {object} model.Document
// @Router /documents [post]
func (h *DocumentHandler) UploadDocument(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, header, err := r.FormFile("file")
	if err != nil {
		render.Error(w, http.StatusBadRequest, "Failed to read file")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "Failed to read file data")
		return
	}

	doc, err := h.service.Process(data, header.Filename)
	if err != nil {
		if errors.Is(err, port.ErrDuplicateDocumentNumber) {
			render.Error(w, http.StatusBadRequest, "A document with this document number already exists")
			return
		}
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.JSON(w, http.StatusOK, doc)
}

// GetDocument handles getting a document by ID
// @Summary Get a document
// @Description Get document details by ID
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} model.Document
// @Router /documents/{id} [get]
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	doc, err := h.service.GetDocument(id)
	if err != nil {
		render.Error(w, http.StatusNotFound, "Document not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

// GetDocuments handles listing all documents
// @Summary Get documents
// @Description Get all documents
// @Produce json
// @Success 200 {array} model.Document
// @Router /documents [get]
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	docs, err := h.service.GetDocuments()
	if err != nil {
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if docs == nil {
		docs = []*model.Document{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

type UpdateDocumentRequest struct {
	Type           model.DocumentType `json:"type"`
	SupplierName   string             `json:"supplier_name"`
	DocumentNumber string             `json:"document_number"`
	IssueDate      *string            `json:"issue_date,omitempty"`
	DueDate        *string            `json:"due_date,omitempty"`
	Currency       string             `json:"currency"`
	LineItems      []model.LineItem   `json:"line_items"`
	Subtotal       float64            `json:"subtotal"`
	Tax            float64            `json:"tax"`
	Total          float64            `json:"total"`
}

// @Summary Update a document
// @Description Update document details
// @Accept json
// @Produce json
// @Param id path string true "Document ID"
// @Param document body UpdateDocumentRequest true "Document data"
// @Success 200 {object} model.Document
// @Router /documents/{id} [put]
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var req UpdateDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if !model.ValidCurrencies[req.Currency] {
		render.Error(w, http.StatusBadRequest, "Currency must be one of: USD, BAM, EURO")
		return
	}

	var issueDate, dueDate *time.Time
	if req.IssueDate != nil {
		if parsed, err := time.Parse("2006-01-02", *req.IssueDate); err == nil {
			issueDate = &parsed
		}
	}
	if req.DueDate != nil {
		if parsed, err := time.Parse("2006-01-02", *req.DueDate); err == nil {
			dueDate = &parsed
		}
	}

	doc := &model.Document{
		ID:             id,
		Type:           req.Type,
		SupplierName:   req.SupplierName,
		DocumentNumber: req.DocumentNumber,
		IssueDate:      issueDate,
		DueDate:        dueDate,
		Currency:       req.Currency,
		LineItems:      req.LineItems,
		Subtotal:       req.Subtotal,
		Tax:            req.Tax,
		Total:          req.Total,
	}

	err := h.service.UpdateDocument(doc)
	if err != nil {
		if errors.Is(err, port.ErrDuplicateDocumentNumber) {
			render.Error(w, http.StatusBadRequest, "A document with this document number already exists")
			return
		}
		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(doc)
}

type ApproveDocumentRequest struct {
	Corrections *UpdateDocumentRequest `json:"corrections,omitempty"`
}

// ApproveDocument handles approving a document
// @Summary Approve a document
// @Description Approve a document
// @Accept json
// @Param id path string true "Document ID"
// @Param corrections body ApproveDocumentRequest false "Corrections"
// @Success 200
// @Router /documents/{id}/approve [post]
func (h *DocumentHandler) ApproveDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var req ApproveDocumentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		render.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var corrections *model.Document
	if req.Corrections != nil {
		// Validate currency
		if !model.ValidCurrencies[req.Corrections.Currency] {
			render.Error(w, http.StatusBadRequest, "Currency must be one of: USD, BAM, EURO")
			return
		}

		// Convert string dates to time.Time
		var issueDate, dueDate *time.Time
		if req.Corrections.IssueDate != nil {
			if parsed, err := time.Parse("2006-01-02", *req.Corrections.IssueDate); err == nil {
				issueDate = &parsed
			}
		}
		if req.Corrections.DueDate != nil {
			if parsed, err := time.Parse("2006-01-02", *req.Corrections.DueDate); err == nil {
				dueDate = &parsed
			}
		}

		corrections = &model.Document{
			Type:           req.Corrections.Type,
			SupplierName:   req.Corrections.SupplierName,
			DocumentNumber: req.Corrections.DocumentNumber,
			IssueDate:      issueDate,
			DueDate:        dueDate,
			Currency:       req.Corrections.Currency,
			LineItems:      req.Corrections.LineItems,
			Subtotal:       req.Corrections.Subtotal,
			Tax:            req.Corrections.Tax,
			Total:          req.Corrections.Total,
		}
	}

	err := h.service.ApproveDocument(id, corrections)
	if err != nil {
		if errors.Is(err, service.ErrDocumentNotFound) {
			render.Error(w, http.StatusNotFound, "document not found")
			return
		}
		if errors.Is(err, port.ErrDuplicateDocumentNumber) {
			render.Error(w, http.StatusBadRequest, "A document with this document number already exists")
			return
		}
		if errors.Is(err, service.ErrDocumentHasUnresolvedIssues) {
			render.Error(w, http.StatusBadRequest, err.Error())
			return
		}
		if errors.Is(err, service.ErrDocumentCannotBeApproved) {
			render.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.Success(w, http.StatusOK, "document approved", nil)
}

// RejectDocument handles rejecting a document
// @Summary Reject a document
// @Description Reject a document
// @Param id path string true "Document ID"
// @Success 200
// @Router /documents/{id}/reject [post]
func (h *DocumentHandler) RejectDocument(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	err := h.service.RejectDocument(id)
	if err != nil {
		if err.Error() == "document not found" {
			render.Error(w, http.StatusNotFound, "document not found")
			return
		}

		render.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	render.Success(w, http.StatusOK, "document rejected", nil)
}
