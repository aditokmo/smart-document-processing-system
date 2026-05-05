import { createFileRoute } from '@tanstack/react-router';
import DocumentDetail from '@/modules/document/pages/Detail';

export const Route = createFileRoute('/document/$id')({
    component: DocumentDetail,
});