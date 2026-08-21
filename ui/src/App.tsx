import { BrowserRouter, Routes, Route } from 'react-router-dom';
import JobList from './pages/JobList';
import CreateJob from './pages/CreateJob';
import JobDetail from './pages/JobDetail';
import ErrorBoundary from './components/ErrorBoundary';
import { Toaster } from 'react-hot-toast';

export default function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <div className="min-h-screen bg-neutral-50">
          <Routes>
            <Route path="/" element={<JobList />} />
            <Route path="/create" element={<CreateJob />} />
            <Route path="/jobs/:id" element={<JobDetail />} />
          </Routes>
        </div>
        <Toaster position="top-right" />
      </BrowserRouter>
    </ErrorBoundary>
  );
}