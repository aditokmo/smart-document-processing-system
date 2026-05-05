import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertTriangle } from 'lucide-react';
import type { Issues } from '../types';

interface ValidationIssuesProps {
  issues: Issues[];
}

export function ValidationIssues({ issues }: ValidationIssuesProps) {
  if (!issues || issues?.length === 0) return null;

  return (
    <Card className="border-yellow-200 bg-yellow-50">
      <CardHeader>
        <CardTitle className="text-lg flex items-center text-yellow-800">
          <AlertTriangle className="mr-2 h-5 w-5" />
          Validation Issues
        </CardTitle>
      </CardHeader>
      <CardContent>
        <ul className="list-disc list-inside space-y-1">
          {issues?.map((issue: Issues, index) => (
            <li key={index} className="text-yellow-700">
              {issue.message}
            </li>
          ))}
        </ul>
      </CardContent>
    </Card>
  );
}
