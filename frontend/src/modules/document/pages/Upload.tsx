import { useForm } from 'react-hook-form';
import { useUploadDocument } from '../hooks';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Upload } from 'lucide-react';
import { toast } from 'react-hot-toast';

export default function UploadPage() {
  const { register, handleSubmit } = useForm<{ file?: FileList }>();
  const uploadMutation = useUploadDocument();

  const onSubmit = (data: { file?: FileList }) => {
    const file = data.file?.[0];
    if (!file) {
      toast.error('Please select a file');
      return;
    }
    if (file.size === 0) {
      toast.error('File is empty');
      return;
    }
    uploadMutation.mutate({ file });
  };

  return (
    <div className="p-6 max-w-md mx-auto">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Upload className="mr-2" />
            Upload Document
          </CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <input
                {...register('file')}
                type="file"
                accept=".pdf,.jpg,.jpeg,.png,.csv,.txt"
                className="h-8 w-full min-w-0 rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none file:inline-flex file:h-6 file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:cursor-not-allowed disabled:bg-input/50 disabled:opacity-50 aria-invalid:border-destructive aria-invalid:ring-3 aria-invalid:ring-destructive/20 md:text-sm dark:bg-input/30 dark:disabled:bg-input/80 dark:aria-invalid:border-destructive/50 dark:aria-invalid:ring-destructive/40"
              />
            </div>
            <Button type="submit" disabled={uploadMutation.isPending} className="w-full">
              {uploadMutation.isPending ? 'Uploading...' : 'Upload'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}