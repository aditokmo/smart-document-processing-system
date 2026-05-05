package repository

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"backend/pkg/utils"
)

type FileProcessorImpl struct {
	extractor port.ExtractionService
}

func NewFileProcessor(extractor port.ExtractionService) *FileProcessorImpl {
	return &FileProcessorImpl{extractor: extractor}
}

func (f *FileProcessorImpl) Process(file []byte, filename string) (*model.Document, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	switch ext {
	case ".pdf":
		return f.processPDF(file, filename)
	case ".csv":
		return f.processCSV(file, filename)
	case ".txt":
		return f.processTXT(file, filename)
	default:
		return nil, fmt.Errorf("Unsupported file type: %s", ext)
	}
}

func (f *FileProcessorImpl) processPDF(file []byte, filename string) (*model.Document, error) {
	text, err := utils.ExtractTextFromPDF(file)
	if err != nil {
		return nil, err
	}
	log.Printf("Extracted text from PDF: %s", text)
	docType := f.detectDocumentType(filename, text)
	return f.extractor.Extract(text, docType)
}

func (f *FileProcessorImpl) processCSV(file []byte, filename string) (*model.Document, error) {
	reader := csv.NewReader(bytes.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	fileText := string(file)
	docType := f.detectDocumentType(filename, fileText)

	doc := &model.Document{
		ID:   utils.GenerateID(),
		Type: docType,
	}

	for _, record := range records[1:] {
		if len(record) >= 4 {
			qty := utils.ParseInt(record[1])
			price := utils.ParseFloat(record[2])
			total := utils.ParseFloat(record[3])
			doc.LineItems = append(doc.LineItems, model.LineItem{
				Description: record[0],
				Quantity:    qty,
				Price:       price,
				Total:       total,
			})
		}
	}

	for _, item := range doc.LineItems {
		doc.Subtotal += item.Total
	}
	doc.Total = doc.Subtotal

	return doc, nil
}

func (f *FileProcessorImpl) processTXT(file []byte, filename string) (*model.Document, error) {
	text := string(file)
	docType := f.detectDocumentType(filename, text)
	return f.extractor.Extract(text, docType)
}

func (f *FileProcessorImpl) detectDocumentType(filename, content string) model.DocumentType {
	filename = strings.ToLower(filename)
	content = strings.ToLower(content)

	purchaseOrderPattern := regexp.MustCompile(`(?i)\b(purchase order|purchase_order|po\s*(?:number|no|#)?|p\.o\.|po #)\b`)
	invoicePattern := regexp.MustCompile(`(?i)\b(invoice|invoice number|invoice no|amount due|total due|bill to|bill from|due date)\b`)

	if purchaseOrderPattern.MatchString(filename) || purchaseOrderPattern.MatchString(content) {
		return model.TypePurchaseOrder
	}
	if invoicePattern.MatchString(filename) || invoicePattern.MatchString(content) {
		return model.TypeInvoice
	}

	return model.TypeInvoice
}
