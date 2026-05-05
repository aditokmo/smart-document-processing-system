package port

import (
	"backend/internal/domain/model"
	"errors"
)

var ErrDuplicateDocumentNumber = errors.New("document number already exists")

type DocumentRepository interface {
	Save(doc *model.Document) error
	FindByID(id string) (*model.Document, error)
	FindAll() ([]*model.Document, error)
	UpdateStatus(id string, status model.DocumentStatus) error
	Delete(id string) error
}

type ExtractionService interface {
	Extract(content string, docType model.DocumentType) (*model.Document, error)
}

type FileProcessor interface {
	Process(file []byte, filename string) (*model.Document, error)
}
