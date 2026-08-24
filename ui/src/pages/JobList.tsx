import { lazy, Suspense, useEffect, useState, useCallback } from 'react';
import { Link } from 'react-router-dom';
import { RefreshCw, WifiOff, ChevronLeft, ChevronRight } from 'lucide-react';
import { jobService as api } from '../service/jobService';
import type { Job } from '../types/job';
import StatusBadge from '../components/common/StatusBadge';
import AppButton from '../components/common/AppButton';
import ConfirmModal from '../components/common/ConfirmModal';
import { ROUTES, COMMON_LABELS } from '../constants/common';
import { JOB_LIST_TEXTS } from '../constants/jobList';
import { jobTitle } from '../utils/jobTitle';
const CreateJobModal = lazy(() => import('../components/jobList/CreateJobModal'));

/**
 * The home page: lists every pipeline job, auto-refreshing every 15
 * minutes (or on demand), with entry points to create a new job and to
 * delete an existing one.
 */
export default function JobList() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [lastFetched, setLastFetched] = useState<Date | null>(null);

  const [deleteTarget, setDeleteTarget] = useState<Job | null>(null);
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  const [page, setPage] = useState(0);
  const pageSize = 50;
  const [hasMore, setHasMore] = useState(true);

  /** Fetches the job list from the API and updates loading/error state around it. */
  const loadJobs = useCallback(async () => {
    setLoading(true);
    try {
      const data = await api.getAllJobs(pageSize, page * pageSize);
      setJobs(data ?? []);
      setHasMore((data ?? []).length === pageSize);
      setLastFetched(new Date());
      setError(null);
    } catch {
      setError(JOB_LIST_TEXTS.FETCH_ERROR);
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    loadJobs();
    const interval = setInterval(loadJobs, 15 * 60 * 1000); // 15 minutes
    return () => clearInterval(interval);
  }, [loadJobs]);

  /**
   * Opens the delete-confirmation modal for a job. Stops the click from
   * also triggering the card's own Link navigation to the job detail page.
   */
  function promptDelete(job: Job, e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    setDeleteTarget(job);
  }

  /**
   * Deletes the job the confirmation modal is open for, then refreshes
   * the list. The modal always closes, success or failure - the axios
   * interceptor already toasts the error, so leaving it stuck open on
   * failure would just be redundant.
   */
  async function handleDelete() {
    if (deleteTarget) {
      try {
        await api.deleteJob(deleteTarget.id);
        loadJobs();
      } finally {
        setDeleteTarget(null);
      }
    }
  }

  if (loading) return <div className="p-8 text-neutral-500">{JOB_LIST_TEXTS.LOADING_JOBS}</div>;

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
          <AppButton onClick={() => setIsCreateOpen(true)}>
            {JOB_LIST_TEXTS.NEW_JOB}
          </AppButton>
        </div>
      </div>

      {error ? (
        <div className="text-center py-16 text-neutral-400 border-2 border-dashed rounded-xl">
          <WifiOff size={28} className="mx-auto mb-3 text-neutral-300" />
          <p className="text-neutral-500 font-medium">{JOB_LIST_TEXTS.NO_JOBS_FOUND}</p>
          <p className="text-xs text-neutral-400 mt-1">{error}</p>
          <AppButton variant="secondary" size="sm" onClick={loadJobs} className="mt-4">
            {COMMON_LABELS.REFRESH}
          </AppButton>
        </div>
      ) : jobs.length === 0 ? (
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
                  <span className="text-sm font-medium text-neutral-800">{jobTitle(job)}</span>
                  <StatusBadge status={job.status} />
                </div>
                <AppButton onClick={(e) => promptDelete(job, e)} className="text-xs text-neutral-400 hover:text-danger-600 transition-colors">
                  {COMMON_LABELS.DELETE}
                </AppButton>
              </div>
              <div className="mt-2 flex gap-6 text-sm text-neutral-500">
                <span>{job.processed_records} {JOB_LIST_TEXTS.PROCESSED}</span>
                <span>{job.error_count} {JOB_LIST_TEXTS.ERRORS}</span>
                <span>{new Date(job.created_at).toLocaleString()}</span>
              </div>
            </Link>
          ))}
          
          <div className="flex items-center justify-between pt-4 mt-6 border-t border-neutral-100">
            <AppButton 
              variant="secondary" 
              size="sm"
              onClick={() => setPage(p => Math.max(0, p - 1))} 
              disabled={page === 0 || loading}
            >
              <ChevronLeft size={16} className="mr-1" /> Previous
            </AppButton>
            <span className="text-sm font-medium text-neutral-500">Page {page + 1}</span>
            <AppButton 
              variant="secondary" 
              size="sm"
              onClick={() => setPage(p => p + 1)} 
              disabled={!hasMore || loading}
            >
              Next <ChevronRight size={16} className="ml-1" />
            </AppButton>
          </div>
        </div>
      )}

      <ConfirmModal
        isOpen={deleteTarget !== null}
        title={JOB_LIST_TEXTS.DELETE_MODAL_TITLE}
        message={deleteTarget ? JOB_LIST_TEXTS.deleteModalMessage(jobTitle(deleteTarget)) : ''}
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
      />
      {isCreateOpen && (
        <Suspense fallback={null}>
          <CreateJobModal isOpen={isCreateOpen} onClose={() => setIsCreateOpen(false)} />
        </Suspense>
      )}
    </div>
  );
}