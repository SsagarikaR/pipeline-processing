import { BrowserRouter, Routes, Route } from 'react-router-dom';
import JobList from './pages/JobList';
import JobDetail from './pages/JobDetail';
import ErrorBoundary from './components/common/ErrorBoundary';
import { Toaster } from 'react-hot-toast';
import { ROUTES } from './constants/common';

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
          <Routes>
            <Route path={ROUTES.HOME} element={<JobList />} />
            <Route path={ROUTES.JOB_DETAIL} element={<JobDetail />} />
          </Routes>
        </div>
        <Toaster position="top-right" />
      </BrowserRouter>
    </ErrorBoundary>
  );
}