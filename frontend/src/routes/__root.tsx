import { createRootRoute, Outlet } from '@tanstack/react-router';
import { Navbar } from '@/layout/Navbar';

export const Route = createRootRoute({
    component: () => {
        return (
            <>
                <Navbar />
                <Outlet />
            </>
        );
    },
});