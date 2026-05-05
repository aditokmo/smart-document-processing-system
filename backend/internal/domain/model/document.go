package model

import (
	"fmt"
	"log"
	"math"
	"time"
)

type DocumentStatus string

const (
	StatusUploaded    DocumentStatus = "uploaded"
	StatusNeedsReview DocumentStatus = "needs_review"
	StatusValidated   DocumentStatus = "validated"
	StatusRejected    DocumentStatus = "rejected"
)

type DocumentType string

const (
	TypeInvoice       DocumentType = "invoice"
	TypePurchaseOrder DocumentType = "purchase_order"
)

type Currency string

const (
	CurrencyUSD  Currency = "USD"
	CurrencyBAM  Currency = "BAM"
	CurrencyEURO Currency = "EURO"
)

var ValidCurrencies = map[string]bool{
	"USD":  true,
	"BAM":  true,
	"EURO": true,
}

type Document struct {
	ID             string            `json:"id"`
	Type           DocumentType      `json:"type"`
	SupplierName   string            `json:"supplier_name"`
	DocumentNumber string            `json:"document_number"`
	Filename       string            `json:"filename"`
	IssueDate      *time.Time        `json:"issue_date,omitempty"`
	DueDate        *time.Time        `json:"due_date,omitempty"`
	Currency       string            `json:"currency"`
	LineItems      []LineItem        `json:"line_items"`
	Subtotal       float64           `json:"subtotal"`
	Tax            float64           `json:"tax"`
	Total          float64           `json:"total"`
	Status         DocumentStatus    `json:"status"`
	Issues         []ValidationIssue `json:"issues,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type LineItem struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	Total       float64 `json:"total"`
}

type ValidationIssue struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (d *Document) Validate() []ValidationIssue {
	var issues []ValidationIssue

	if d.Type == "" {
		issues = append(issues, ValidationIssue{Field: "type", Message: "Document type is required"})
	}
	if d.SupplierName == "" {
		issues = append(issues, ValidationIssue{Field: "supplier_name", Message: "Supplier name is required"})
	}
	if d.DocumentNumber == "" {
		issues = append(issues, ValidationIssue{Field: "document_number", Message: "Document number is required"})
	}
	if d.Currency == "" {
		issues = append(issues, ValidationIssue{Field: "currency", Message: "Currency is required"})
	} else if !ValidCurrencies[d.Currency] {
		issues = append(issues, ValidationIssue{Field: "currency", Message: "Currency must be one of: USD, BAM, EURO"})
	}

	if d.IssueDate != nil && d.DueDate != nil && d.IssueDate.After(*d.DueDate) {
		issues = append(issues, ValidationIssue{Field: "due_date", Message: "Due date cannot be before issue date"})
	}
	if d.IssueDate != nil && d.IssueDate.After(time.Now()) {
		issues = append(issues, ValidationIssue{Field: "issue_date", Message: "Issue date cannot be in the future"})
	}

	calculatedSubtotal := 0.0
	for _, item := range d.LineItems {
		calculatedSubtotal += item.Total
		expectedItemTotal := float64(item.Quantity) * item.Price
		if math.Abs(item.Total-expectedItemTotal) > 0.01 {
			issues = append(issues, ValidationIssue{
				Field:   "line_items",
				Message: fmt.Sprintf("Line item total mismatch for %s; expected %v", item.Description, expectedItemTotal),
			})
		}
	}
	if len(d.LineItems) > 0 && math.Abs(calculatedSubtotal-d.Subtotal) > 0.01 {
		expectedSubtotal := calculatedSubtotal
		issues = append(issues, ValidationIssue{
			Field:   "subtotal",
			Message: fmt.Sprintf("Subtotal does not match line items; expected %v", expectedSubtotal),
		})
	}
	if math.Abs((d.Subtotal+d.Tax)-d.Total) > 0.01 {
		expectedTotal := d.Subtotal + d.Tax
		log.Println(expectedTotal)
		issues = append(issues, ValidationIssue{
			Field:   "total",
			Message: fmt.Sprintf("Total does not match subtotal + tax; expected %s %v", d.Currency, expectedTotal),
		})
	}

	return issues
}
