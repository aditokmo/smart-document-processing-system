import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { Document as DocType } from '../types';

const statusColors = {
  uploaded: 'bg-blue-500',
  needs_review: 'bg-yellow-500',
  validated: 'bg-green-500',
  rejected: 'bg-red-500',
};

interface DocumentHeaderProps {
  doc: DocType;
  documentNumber: string;
  supplierName: string;
}

export function DocumentHeader({ doc, documentNumber, supplierName }: DocumentHeaderProps) {
  const statusClass = statusColors[doc.status] ?? 'bg-gray-500';

  return (
    <Card>
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle className="text-2xl">{documentNumber === "" ? supplierName : `${documentNumber} - ${supplierName}`}</CardTitle>
          <Badge className={statusClass}>
            {doc.status.replace('_', ' ').toUpperCase()}
          </Badge>
        </div>
      </CardHeader>
    </Card>
  );
}
