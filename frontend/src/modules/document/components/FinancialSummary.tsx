import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import type { Document as DocType } from '../types';

interface FinancialSummaryProps {
  doc: DocType;
  formData: Partial<DocType> & { tax_rate?: number };
  isEditable: boolean;
  onUpdate: (key: keyof DocType | 'tax_rate', value: unknown) => void;
}

export function FinancialSummary({ doc, formData, isEditable, onUpdate }: FinancialSummaryProps) {
  const subtotal = formData.subtotal ?? doc.subtotal;
  const taxRate = formData.tax_rate ?? (doc.subtotal ? (doc.tax / doc.subtotal) * 100 : 0);
  const taxAmount = subtotal * (taxRate / 100);
  const total = formData.total ?? subtotal + taxAmount;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-lg">Financial Summary</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2">Subtotal</label>
          {isEditable ? (
            <Input
              type="number"
              value={subtotal}
              onChange={(e) => onUpdate('subtotal', parseFloat(e.target.value) || 0)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{doc.currency} {subtotal}</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Tax rate (%)</label>
          {isEditable ? (
            <Input
              type="number"
              value={taxRate}
              onChange={(e) => onUpdate('tax_rate', parseFloat(e.target.value) || 0)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600">{taxRate.toFixed(2)}%</p>
          )}
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Total</label>
          {isEditable ? (
            <Input
              type="number"
              value={total}
              onChange={(e) => onUpdate('total', parseFloat(e.target.value) || 0)}
              className="text-sm"
            />
          ) : (
            <p className="text-sm text-gray-600 font-bold">{doc.currency} {total}</p>
          )}
        </div>
      </CardContent>
    </Card>
  );
}
