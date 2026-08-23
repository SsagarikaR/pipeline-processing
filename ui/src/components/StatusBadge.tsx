import { StatusBadgeProps } from '../types/common';
import { STATUS_BADGE_STYLES } from '../constants/common';

export default function StatusBadge({ status }: StatusBadgeProps) {
  const style = STATUS_BADGE_STYLES[status] ?? STATUS_BADGE_STYLES.pending;
  return (
    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${style}`}>
      {status}
    </span>
  );
}