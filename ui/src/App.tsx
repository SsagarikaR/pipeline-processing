import { BrowserRouter, Routes, Route } from 'react-router-dom';
import JobList from './pages/JobList';
import JobDetail from './pages/JobDetail';
import ErrorBoundary from './components/common/ErrorBoundary';
import { Toaster } from 'react-hot-toast';
import { ROUTES } from './constants/common';

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