import type { JobStatus } from '../types/job';

const STATUS_STYLES: Record<JobStatus, string> = {
  pending: 'bg-neutral-100 text-neutral-700 border-neutral-300',
  running: 'bg-brand-100 text-brand-700 border-brand-300 animate-pulse',
  completed: 'bg-success-100 text-success-700 border-success-300',
  failed: 'bg-danger-100 text-danger-700 border-danger-300',
  cancelled: 'bg-amber-100 text-amber-700 border-amber-300',
};

interface StatusBadgeProps {
  status: JobStatus;
}

export default function StatusBadge({ status }: StatusBadgeProps) {
  const style = STATUS_STYLES[status] ?? STATUS_STYLES.pending;
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${style}`}>
      {status}
    </span>
  );
}