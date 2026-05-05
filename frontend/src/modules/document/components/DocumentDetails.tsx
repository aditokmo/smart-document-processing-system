import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';
import { VALID_CURRENCIES } from '@/lib/constants';
import type { Document as DocType } from '../types';

interface DocumentDetailsProps {
  doc: DocType;
  formData: Partial<DocType>;
  isEditable: boolean;
  onUpdate: (key: keyof DocType, value: unknown) => void;
}

export function DocumentDetails({ doc, formData, isEditable, onUpdate }: DocumentDetailsProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Document Details</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Type</label>
          {isEditable ? (
            <Input
              value={formData.type || ''}
              onChange={(e) => onUpdate('type', e.target.value)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.type}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Supplier</label>
          {isEditable ? (
            <Input
              value={formData.supplier_name || ''}
              onChange={(e) => onUpdate('supplier_name', e.target.value)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.supplier_name}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Document Number</label>
          {isEditable ? (
            <Input
              value={formData.document_number || ''}
              onChange={(e) => onUpdate('document_number', e.target.value)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.document_number || 'N/A'}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Issue Date</label>
          {isEditable ? (
            <Input
              type="date"
              value={formData.issue_date || ''}
              onChange={(e) => onUpdate('issue_date', e.target.value)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.issue_date ? new Date(doc.issue_date).toLocaleDateString('de-DE') : 'N/A'}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Due Date</label>
          {isEditable ? (
            <Input
              type="date"
              value={formData.due_date || ''}
              onChange={(e) => onUpdate('due_date', e.target.value)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.due_date ? new Date(doc.due_date).toLocaleDateString('de-DE') : 'N/A'}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Currency</label>
          {isEditable ? (
            <Select
              value={formData.currency || ''}
              onChange={(e) => onUpdate('currency', e.target.value)}
              className="text-sm"
            >
              <option value="">Select currency</option>
              {VALID_CURRENCIES.map((currency) => (
                <option key={currency} value={currency}>
                  {currency}
                </option>
              ))}
            </Select>
          ) : (
            <p className="text-sm text-gray-600">{doc.currency}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
