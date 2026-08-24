import type { MetricCardProps } from '../../types/common';

/** One stat tile in the progress header (e.g. "Processed: 1,204"). */
export default function MetricCard({ label, value, tone = 'default' }: MetricCardProps) {
  const toneClass = tone === 'red' ? 'text-danger-600' : 'text-neutral-900';
  return (
    <div className="bg-white border border-neutral-200 rounded-lg p-3">
      <div className="text-xs text-neutral-500 mb-1">{label}</div>
      <div className={`text-lg font-semibold ${toneClass}`}>{value}</div>
    </div>
  );
}
