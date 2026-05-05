import { createFileRoute } from '@tanstack/react-router';
import UploadPage from '@/modules/document/pages/Upload';

export const Route = createFileRoute('/upload')({
    component: UploadPage,
});