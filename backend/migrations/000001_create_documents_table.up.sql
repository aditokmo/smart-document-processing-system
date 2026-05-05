CREATE TABLE IF NOT EXISTS documents (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    supplier_name TEXT,
    document_number TEXT,
    issue_date TIMESTAMPTZ,
    due_date TIMESTAMPTZ,
    currency TEXT,
    filename TEXT,
    line_items JSONB,
    subtotal DOUBLE PRECISION,
    tax DOUBLE PRECISION,
    total DOUBLE PRECISION,
    status TEXT NOT NULL,
    issues JSONB,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_document_number 
ON documents(document_number)
WHERE document_number IS NOT NULL AND document_number != '';

CREATE INDEX IF NOT EXISTS idx_document_status 
ON documents(status);

CREATE INDEX IF NOT EXISTS idx_document_created_at 
ON documents(created_at DESC);