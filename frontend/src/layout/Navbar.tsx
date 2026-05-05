import { Link, useRouterState } from '@tanstack/react-router';

export function Navbar() {
  // eslint-disable-next-line react-hooks/rules-of-hooks
  const routerState = useRouterState();
  const isDetailPage = routerState.location.pathname.match(/^\/documents\/\d+$/) || routerState.location.pathname.match(/^\/document\/\d+$/);

  if (isDetailPage) {
    return null;
  }

  return (
    <nav className="sticky top-0 z-50 bg-white border-b border-slate-200 shadow-sm">
      <div className="max-w-screen-xl mx-auto px-6 py-4 flex items-center justify-between">
        <div className="flex items-center gap-8">
          <Link 
            to="/" 
            className="text-lg font-bold text-slate-900 hover:text-slate-700 transition-colors"
            activeProps={{
              className: "text-slate-900"
            }}
          >
            📄 SmartDocs
          </Link>
          <div className="hidden md:flex items-center gap-1">
            <Link 
              to="/" 
              className="px-4 py-2 rounded-md text-slate-700 hover:bg-slate-100 transition-colors text-sm font-medium"
              activeProps={{
                className: "bg-slate-100 text-slate-900 font-semibold"
              }}
            >
              Dashboard
            </Link>
            <Link 
              to="/documents" 
              className="px-4 py-2 rounded-md text-slate-700 hover:bg-slate-100 transition-colors text-sm font-medium"
              activeProps={{
                className: "bg-slate-100 text-slate-900 font-semibold"
              }}
            >
              Documents
            </Link>
          </div>
        </div>
        <div className="md:hidden flex items-center gap-2">
          <Link to="/" className="px-3 py-2 text-sm font-medium text-slate-700 hover:text-slate-900">
            Dashboard
          </Link>
          <Link to="/documents" className="px-3 py-2 text-sm font-medium text-slate-700 hover:text-slate-900">
            Documents
          </Link>
        </div>
      </div>
    </nav>
  );
}
