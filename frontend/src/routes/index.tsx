import { createFileRoute } from '@tanstack/react-router';
import Dashboard from '@/modules/document/pages/Dashboard';

export const Route = createFileRoute('/')({
    component: Dashboard,
});