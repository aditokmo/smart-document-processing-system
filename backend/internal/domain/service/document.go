package service

import (
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDocumentNotFound            = errors.New("document not found")
	ErrDocumentCannotBeApproved    = errors.New("document cannot be approved in current status")
	ErrDocumentHasUnresolvedIssues = errors.New("document has unresolved validation issues")
)

type DocumentService struct {
	repository  port.DocumentRepository
	extractor   port.ExtractionService
	fileProcess port.FileProcessor
}

func NewDocumentService(repository port.DocumentRepository, extractor port.ExtractionService, fileProcess port.FileProcessor) *DocumentService {
	return &DocumentService{
		repository:  repository,
		extractor:   extractor,
		fileProcess: fileProcess,
	}
}

func (s *DocumentService) Process(file []byte, filename string) (*model.Document, error) {
	doc, err := s.fileProcess.Process(file, filename)
	if err != nil {
		return nil, err
	}

	doc.Filename = filename
	doc.Status = model.StatusUploaded
	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()

	issues := doc.Validate()
	doc.Issues = issues
	if len(issues) > 0 {
		doc.Status = model.StatusNeedsReview
	} else {
		doc.Status = model.StatusValidated
	}

	if err := s.repository.Save(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (s *DocumentService) GetDocument(id string) (*model.Document, error) {
	doc, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}
	return doc, nil
}

func (s *DocumentService) GetDocuments() ([]*model.Document, error) {
	return s.repository.FindAll()
}

func (s *DocumentService) UpdateDocument(doc *model.Document) error {
	existing, err := s.repository.FindByID(doc.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDocumentNotFound
		}
		return err
	}

	if doc.Filename == "" {
		doc.Filename = existing.Filename
	}
	if doc.CreatedAt.IsZero() {
		doc.CreatedAt = existing.CreatedAt
	}

	doc.UpdatedAt = time.Now()
	issues := doc.Validate()
	doc.Issues = issues
	if len(issues) > 0 {
		doc.Status = model.StatusNeedsReview
	} else {
		doc.Status = existing.Status
	}
	return s.repository.Save(doc)
}

func (s *DocumentService) ApproveDocument(id string, corrections *model.Document) error {
	doc, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrDocumentNotFound
		}
		return err
	}
	if doc.Status != model.StatusNeedsReview && doc.Status != model.StatusValidated {
		return ErrDocumentCannotBeApproved
	}

	if corrections == nil {
		if doc.Status == model.StatusNeedsReview {
			return ErrDocumentHasUnresolvedIssues
		}
		doc.Status = model.StatusValidated
		doc.UpdatedAt = time.Now()
		return s.repository.Save(doc)
	}

	if corrections != nil {
		if corrections.Type != "" {
			doc.Type = corrections.Type
		}
		if corrections.SupplierName != "" {
			doc.SupplierName = corrections.SupplierName
		}
		if corrections.DocumentNumber != "" {
			doc.DocumentNumber = corrections.DocumentNumber
		}
		if corrections.IssueDate != nil {
			doc.IssueDate = corrections.IssueDate
		}
		if corrections.DueDate != nil {
			doc.DueDate = corrections.DueDate
		}
		if corrections.Currency != "" {
			doc.Currency = corrections.Currency
		}
		if corrections.LineItems != nil && len(corrections.LineItems) > 0 {
			doc.LineItems = corrections.LineItems
			doc.Subtotal = 0
			for _, item := range doc.LineItems {
				doc.Subtotal += item.Total
			}
		}
		if corrections.Subtotal > 0 {
			doc.Subtotal = corrections.Subtotal
		}
		if corrections.Tax > 0 {
			doc.Tax = corrections.Tax
		}
		if corrections.Total > 0 {
			doc.Total = corrections.Total
		}
	}

	issues := doc.Validate()
	if len(issues) > 0 {
		return ErrDocumentHasUnresolvedIssues
	}
	doc.Issues = issues
	doc.Status = model.StatusValidated
	doc.UpdatedAt = time.Now()
	return s.repository.Save(doc)
}

func (s *DocumentService) RejectDocument(id string) error {
	doc, err := s.repository.FindByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("document not found")
		}
		return err
	}
	doc.Status = model.StatusRejected
	doc.UpdatedAt = time.Now()
	return s.repository.Save(doc)
}
