import type { ErrorsTableProps } from '../../types/common';
import { JOB_DETAIL_TEXTS } from '../../constants/jobDetail';

/** The "Errors" tab: a list of records that failed processing, with the stage and reason. */
export default function ErrorsTable({ errors }: ErrorsTableProps) {
  if (!errors || errors.length === 0) {
    return <div className="text-sm text-neutral-400 py-8 text-center">{JOB_DETAIL_TEXTS.ERRORS_EMPTY}</div>;
  }
  return (
    <div className="space-y-2">
      {errors.map((e) => (
        <div key={e.id} className="bg-danger-50 border border-danger-200 rounded-lg p-3 text-sm">
          <div className="flex items-center gap-2 mb-1">
            <span className="px-2 py-0.5 bg-danger-100 text-danger-700 rounded text-xs font-medium">{e.stage}</span>
            <span className="text-xs text-neutral-400">{new Date(e.created_at).toLocaleTimeString()}</span>
          </div>
          <div className="text-danger-800">{e.error_message}</div>
          {e.record_data && <pre className="mt-1 text-xs text-neutral-500 overflow-x-auto">{e.record_data}</pre>}
        </div>
      ))}
    </div>
  );
}
