import { createFileRoute } from '@tanstack/react-router';
import Documents from '@/modules/document/pages/Documents';

export const Route = createFileRoute('/documents')({
    component: Documents,
});