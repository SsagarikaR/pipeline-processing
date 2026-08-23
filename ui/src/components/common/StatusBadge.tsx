import { StatusBadgeProps } from '../../types/common';
import { STATUS_BADGE_STYLES } from '../../constants/common';

/**
 * Small colored pill showing a job's status (pending/running/completed/
 * failed/cancelled), used on both the job list and job detail pages.
 */
export default function StatusBadge({ status }: StatusBadgeProps) {
  const style = STATUS_BADGE_STYLES[status] ?? STATUS_BADGE_STYLES.pending;
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${style}`}>
      {status}
    </span>
  );
}