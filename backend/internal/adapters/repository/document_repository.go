package repository

import (
	"backend/internal/domain/model"
	"backend/internal/domain/port"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DocumentRepositoryImpl struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepositoryImpl {
	return &DocumentRepositoryImpl{db: db}
}

func (r *DocumentRepositoryImpl) Save(doc *model.Document) error {
	lineItemsJSON, _ := json.Marshal(doc.LineItems)
	issuesJSON, _ := json.Marshal(doc.Issues)

	_, err := r.db.Exec(`
        INSERT INTO documents
            (id, type, supplier_name, document_number, filename, issue_date, due_date, currency, line_items, subtotal, tax, total, status, issues, created_at, updated_at)
        VALUES
            ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
        ON CONFLICT (id) DO UPDATE SET
            type = EXCLUDED.type,
            supplier_name = EXCLUDED.supplier_name,
            document_number = EXCLUDED.document_number,
            filename = EXCLUDED.filename,
            issue_date = EXCLUDED.issue_date,
            due_date = EXCLUDED.due_date,
            currency = EXCLUDED.currency,
            line_items = EXCLUDED.line_items,
            subtotal = EXCLUDED.subtotal,
            tax = EXCLUDED.tax,
            total = EXCLUDED.total,
            status = EXCLUDED.status,
            issues = EXCLUDED.issues,
            created_at = EXCLUDED.created_at,
            updated_at = EXCLUDED.updated_at
    `,
		doc.ID, doc.Type, doc.SupplierName, doc.DocumentNumber, doc.Filename, doc.IssueDate, doc.DueDate, doc.Currency,
		string(lineItemsJSON), doc.Subtotal, doc.Tax, doc.Total, doc.Status, string(issuesJSON), doc.CreatedAt, doc.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.SQLState() == "23505" {
			return port.ErrDuplicateDocumentNumber
		}
	}
	return err
}

func (r *DocumentRepositoryImpl) FindByID(id string) (*model.Document, error) {
	row := r.db.QueryRow(`
        SELECT id, type, supplier_name, document_number, filename, issue_date, due_date, currency, line_items, subtotal, tax, total, status, issues, created_at, updated_at
        FROM documents WHERE id = $1`, id)

	doc := &model.Document{}
	var lineItemsJSON, issuesJSON []byte
	var issueDate, dueDate sql.NullTime

	err := row.Scan(&doc.ID, &doc.Type, &doc.SupplierName, &doc.DocumentNumber, &doc.Filename, &issueDate, &dueDate, &doc.Currency,
		&lineItemsJSON, &doc.Subtotal, &doc.Tax, &doc.Total, &doc.Status, &issuesJSON, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		return nil, err
	}

	if issueDate.Valid {
		doc.IssueDate = &issueDate.Time
	}
	if dueDate.Valid {
		doc.DueDate = &dueDate.Time
	}

	json.Unmarshal(lineItemsJSON, &doc.LineItems)
	json.Unmarshal(issuesJSON, &doc.Issues)

	return doc, nil
}

func (r *DocumentRepositoryImpl) FindAll() ([]*model.Document, error) {
	rows, err := r.db.Query(`
        SELECT id, type, supplier_name, document_number, filename, issue_date, due_date, currency, line_items, subtotal, tax, total, status, issues, created_at, updated_at
        FROM documents ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []*model.Document
	for rows.Next() {
		doc := &model.Document{}
		var lineItemsJSON, issuesJSON []byte
		var issueDate, dueDate sql.NullTime

		err := rows.Scan(&doc.ID, &doc.Type, &doc.SupplierName, &doc.DocumentNumber, &doc.Filename, &issueDate, &dueDate, &doc.Currency,
			&lineItemsJSON, &doc.Subtotal, &doc.Tax, &doc.Total, &doc.Status, &issuesJSON, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}

		if issueDate.Valid {
			doc.IssueDate = &issueDate.Time
		}
		if dueDate.Valid {
			doc.DueDate = &dueDate.Time
		}

		json.Unmarshal(lineItemsJSON, &doc.LineItems)
		json.Unmarshal(issuesJSON, &doc.Issues)

		docs = append(docs, doc)
	}
	return docs, nil
}

func (r *DocumentRepositoryImpl) UpdateStatus(id string, status model.DocumentStatus) error {
	_, err := r.db.Exec("UPDATE documents SET status = $1, updated_at = $2 WHERE id = $3", status, time.Now(), id)
	return err
}

func (r *DocumentRepositoryImpl) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM documents WHERE id = $1", id)
	return err
}
