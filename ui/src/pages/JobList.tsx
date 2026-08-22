import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { RefreshCw } from 'lucide-react';
import { jobService as api } from '../service/jobService';
import type { Job } from '../types/job';
import StatusBadge from '../components/StatusBadge';
import AppButton from '../components/AppButton';
import ConfirmModal from '../components/ConfirmModal';
import { ROUTES, COMMON_LABELS } from '../constants/common';
import { JOB_LIST_TEXTS } from '../constants/jobList';

export default function JobList() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastFetched, setLastFetched] = useState<Date | null>(null);

  const [deleteJobId, setDeleteJobId] = useState<number | null>(null);

  async function loadJobs() {
    setLoading(true);
    try {
      const data = await api.getAllJobs();
      setJobs(data ?? []);
      setLastFetched(new Date());
      setError(null);
    } catch {
      setError(JOB_LIST_TEXTS.FETCH_ERROR);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadJobs();
    const interval = setInterval(loadJobs, 15 * 60 * 1000); // 15 minutes
    return () => clearInterval(interval);
  }, []);

  function promptDelete(id: number, e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    setDeleteJobId(id);
  }

  async function handleDelete() {
    if (deleteJobId) {
      await api.deleteJob(deleteJobId);
      setDeleteJobId(null);
      loadJobs();
    }
  }

  if (loading) return <div className="p-8 text-neutral-500">{JOB_LIST_TEXTS.LOADING_JOBS}</div>;
  if (error) return <div className="p-8 text-danger-600">{error}</div>;

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-semibold text-neutral-900">{JOB_LIST_TEXTS.TITLE}</h1>
          {lastFetched && (
            <div className="text-xs text-neutral-500 mt-1">
              {JOB_LIST_TEXTS.LAST_FETCHED}{lastFetched.toLocaleTimeString()}
            </div>
          )}
        </div>
        <div className="flex gap-3 items-center">
          <AppButton 
            variant="ghost"
            onClick={loadJobs} 
            disabled={loading}
            className="!p-2"
            title={COMMON_LABELS.REFRESH}
          >
            <RefreshCw size={18} className={loading ? "animate-spin" : ""} />
          </AppButton>
          <Link to={ROUTES.CREATE_JOB} className="px-4 py-2 bg-brand-600 text-white rounded-lg text-sm font-medium hover:bg-brand-700 transition-colors inline-flex items-center justify-center">
            {JOB_LIST_TEXTS.NEW_JOB}
          </Link>
        </div>
      </div>

      {jobs.length === 0 ? (
        <div className="text-center py-16 text-neutral-400 border-2 border-dashed rounded-xl">
          {JOB_LIST_TEXTS.NO_JOBS}
        </div>
      ) : (
        <div className="space-y-3">
          {jobs.map((job) => (
            <Link
              key={job.id}
              to={ROUTES.jobDetail(job.id)}
              className="block bg-white border border-neutral-200 rounded-xl p-4 hover:border-brand-300 hover:shadow-sm transition-all"
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="font-mono text-sm text-neutral-500">#{job.id}</span>
                  <StatusBadge status={job.status} />
                </div>
                <button onClick={(e) => promptDelete(job.id, e)} className="text-xs text-neutral-400 hover:text-danger-600 transition-colors">
                  {COMMON_LABELS.DELETE}
                </button>
              </div>
              <div className="mt-2 flex gap-6 text-sm text-neutral-500">
                <span>{job.processed_records} {JOB_LIST_TEXTS.PROCESSED}</span>
                <span>{job.error_count} {JOB_LIST_TEXTS.ERRORS}</span>
                <span>{new Date(job.created_at).toLocaleString()}</span>
              </div>
            </Link>
          ))}
        </div>
      )}

      <ConfirmModal 
        isOpen={deleteJobId !== null}
        title={JOB_LIST_TEXTS.DELETE_MODAL_TITLE}
        message={deleteJobId ? JOB_LIST_TEXTS.deleteModalMessage(deleteJobId) : ''}
        onConfirm={handleDelete}
        onCancel={() => setDeleteJobId(null)}
      />
    </div>
  );
}