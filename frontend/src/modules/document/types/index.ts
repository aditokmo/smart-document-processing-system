export interface LineItem {
  description: string;
  quantity: number;
  price: number;
  total: number;
}

type IssueField = 'supplier_name' | 'document_number' | 'issue_date' | 'due_date' | 'currency' | 'line_items' | 'subtotal' | 'tax' | 'total';

export interface Issues {
  field: IssueField
  message: string;
}

export interface Document {
  id: string;
  type: 'invoice' | 'purchase_order';
  supplier_name: string;
  document_number: string;
  filename?: string;
  issue_date: string;
  due_date?: string;
  currency: string;
  line_items: LineItem[];
  subtotal: number;
  tax: number;
  tax_rate?: number;
  total: number;
  status: 'uploaded' | 'needs_review' | 'validated' | 'rejected';
  issues?: Issues[];
  createdAt: string;
  updatedAt: string;
}

export interface UploadDocumentRequest {
  file: File;
}

export interface ApproveDocumentRequest {
  corrections?: Partial<Document>;
}

export interface RejectDocumentRequest {
  reason: string;
}