import { lazy, Suspense } from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import ErrorBoundary from './components/common/ErrorBoundary';
import { Toaster } from 'react-hot-toast';
import { ROUTES, COMMON_LABELS } from './constants/common';

// Route-level code splitting: each page's JS only downloads when the
// user actually navigates to it, instead of all being bundled together.
const JobList = lazy(() => import('./pages/JobList'));
const JobDetail = lazy(() => import('./pages/JobDetail'));

/**
 * Root component: sets up client-side routing (job list + job detail
 * pages), wraps everything in an error boundary so a crash anywhere in
 * the tree shows a friendly fallback instead of a blank page, and
 * mounts the toast notification container used app-wide.
 */
export default function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <div className="min-h-screen bg-neutral-50">
          <Suspense fallback={<div className="p-8 text-neutral-500">{COMMON_LABELS.LOADING}</div>}>
            <Routes>
              <Route path={ROUTES.HOME} element={<JobList />} />
              <Route path={ROUTES.JOB_DETAIL} element={<JobDetail />} />
            </Routes>
          </Suspense>
        </div>
        <Toaster position="top-right" />
      </BrowserRouter>
    </ErrorBoundary>
  );
}
