import { useDocuments } from '../hooks';
import { Link } from '@tanstack/react-router';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Eye, Upload } from 'lucide-react';

const statusColors = {
  uploaded: 'bg-blue-500',
  needs_review: 'bg-yellow-500',
  validated: 'bg-green-500',
  rejected: 'bg-red-500',
};

const typeColors = {
  invoice: 'bg-blue-500',
  purchase_order: 'bg-green-500',
};

const currencySymbols: Record<string, string> = {
  USD: '$',
  EUR: '€',
  BAM: 'KM',
};

export default function Documents() {
  const { data: documents, isLoading } = useDocuments();
  const docs = documents ?? [];

  if (isLoading) return <div className="flex justify-center items-center h-64">Loading...</div>;

  return (
    <div className="p-6 max-w-screen-xl mx-auto space-y-6">
      <div className="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
        <div className="space-y-2">
          <p className="text-sm uppercase tracking-[0.2em] text-slate-500">Document Operations</p>
          <h1 className="text-3xl font-semibold tracking-tight text-slate-900">All Documents</h1>
          <p className="max-w-2xl text-sm text-slate-600">
            View and manage all uploaded documents.
          </p>
        </div>
        <Link to="/upload">
          <Button size="lg">
            <Upload className="mr-2 h-4 w-4" />
            Upload Document
          </Button>
        </Link>
      </div>

      <div className="space-y-4">
        {docs.length === 0 ? (
          <Card className="border-slate-200 bg-white shadow-sm">
            <CardContent className="text-center text-slate-600">No documents available yet.</CardContent>
          </Card>
        ) : (
          <div className="grid gap-4">
            {docs.map((doc) => {
              const filename = doc.filename || doc.document_number || 'Untitled document';
              return (
                <Card key={doc.id} className="border-slate-200 bg-white shadow-sm">
                  <CardHeader>
                    <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                      <div>
                        <p className="text-sm text-slate-500">{doc.type.replace('_', ' ').toUpperCase()}</p>
                        <h2 className="text-xl font-semibold text-slate-900">{filename}</h2>
                      </div>
                      <div className="flex flex-wrap items-center gap-2">
                        <Badge className={statusColors[doc.status]}>{doc.status.replace('_', ' ').toUpperCase()}</Badge>
                        <Badge className={typeColors[doc.type]}>{doc.type.replace('_', ' ').toUpperCase()}</Badge>
                      </div>
                    </div>
                  </CardHeader>
                  <CardContent className="grid gap-4 md:grid-cols-2">
                    <div className="space-y-2">
                      <p className="text-sm text-slate-500">Supplier</p>
                      <p className="text-base font-medium text-slate-900">{doc.supplier_name || 'N/A'}</p>
                    </div>
                    <div className="space-y-2">
                      <p className="text-sm text-slate-500">Currency / Total</p>
                      <p className="text-base font-medium text-slate-900">
                        {currencySymbols[doc.currency?.toUpperCase() ?? ''] ?? ''}{doc.total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })} {doc.currency}
                      </p>
                    </div>
                  </CardContent>
                  <CardContent className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                    <div className="text-sm text-slate-600">
                      {doc.issues && doc.issues.length > 0 ? `${doc.issues.length} validation issue${doc.issues.length === 1 ? '' : 's'}` : 'No issues detected'}
                    </div>
                    <Link to="/document/$id" params={{ id: doc.id }}>
                      <Button variant="outline">
                        <Eye className="mr-2 h-4 w-4" />
                        View Details
                      </Button>
                    </Link>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}