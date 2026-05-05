import { useMemo } from 'react';
import { useDocuments } from '../hooks';
import { Link } from '@tanstack/react-router';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Upload } from 'lucide-react';

const statusLabels = {
  needs_review: 'Needs Review',
  validated: 'Validated',
  rejected: 'Rejected',
};

const currencySymbols: Record<string, string> = {
  USD: '$',
  EUR: '€',
  BAM: 'KM',
};

export default function Dashboard() {
  const { data: documents, isLoading } = useDocuments();

  const docs = useMemo(() => documents ?? [], [documents]);

  const summary = useMemo(
    () => ({
      needs_review: docs.filter((doc) => doc.status === 'needs_review').length,
      validated: docs.filter((doc) => doc.status === 'validated').length,
      rejected: docs.filter((doc) => doc.status === 'rejected').length,
      issues: docs.filter((doc) => doc.issues && doc.issues.length > 0).length,
    }),
    [docs]
  );

  const currencyTotals = useMemo(() => {
    return docs.reduce<Record<string, { count: number; total: number }>>((acc, doc) => {
      const currency = doc.currency?.toUpperCase() || 'UNKNOWN Currency';
      if (!acc[currency]) acc[currency] = { count: 0, total: 0 };
      acc[currency].count += 1;
      acc[currency].total += doc.total || 0;
      return acc;
    }, {});
  }, [docs]);

  if (isLoading) return <div className="flex justify-center items-center h-64">Loading...</div>;

  return (
    <div className="p-6 max-w-screen-xl mx-auto space-y-6">
      <div className="flex flex-col gap-6 md:flex-row md:items-end md:justify-between">
        <div className="space-y-2">
          <p className="text-sm uppercase tracking-[0.2em] text-slate-500">Document Operations</p>
          <h1 className="text-3xl font-semibold tracking-tight text-slate-900">Document Dashboard</h1>
          <p className="max-w-2xl text-sm text-slate-600">
            Review documents, monitor validation status, and get a quick financial overview by currency.
          </p>
        </div>
        <div className="flex items-center gap-4">
          <Link to="/documents">
            <Button variant="outline" size="lg">
              View All Documents
            </Button>
          </Link>
          <Link to="/upload">
            <Button size="lg">
              <Upload className="mr-2 h-4 w-4" />
              Upload Document
            </Button>
          </Link>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <Card className="border-slate-200 bg-white shadow-sm">
          <CardContent className="space-y-2 p-6">
            <p className="text-sm font-medium text-slate-500">Total Documents</p>
            <p className="text-3xl font-semibold text-slate-900">{docs.length}</p>
          </CardContent>
        </Card>
        {(['needs_review', 'validated', 'rejected'] as const).map((status) => (
          <Card key={status} className="border-slate-200 bg-white shadow-sm">
            <CardContent className="space-y-2 p-6">
              <p className="text-sm font-medium text-slate-500">{statusLabels[status]}</p>
              <p className="text-3xl font-semibold text-slate-900">{summary[status]}</p>
            </CardContent>
          </Card>
        ))}
      </div>

      <Card className="border-slate-200 bg-white shadow-sm">
        <CardHeader>
          <CardTitle>Totals by Currency</CardTitle>
        </CardHeader>
        <CardContent className="space-y-3 p-6">
          {Object.keys(currencyTotals).length === 0 ? (
            <p className="text-sm text-slate-500">No documents available yet.</p>
          ) : (
            Object.entries(currencyTotals).map(([currency, group]) => (
              <div key={currency} className="flex items-center justify-between rounded-2xl border border-slate-200 bg-slate-50 p-4">
                <div>
                  <p className="text-sm font-medium text-slate-700">{currency}</p>
                  <p className="text-xs text-slate-500">{group.count} document{group.count === 1 ? '' : 's'}</p>
                </div>
                <p className="text-lg font-semibold text-slate-900">
                  {currencySymbols[currency] ?? ''}{group.total.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 })}
                </p>
              </div>
            ))
          )}
        </CardContent>
      </Card>
    </div>
  );
}
