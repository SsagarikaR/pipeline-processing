import type { ResultsTableProps } from '../../types/common';
import { JOB_DETAIL_TEXTS } from '../../constants/jobDetail';

/**
 * The "Results" tab: a table of aggregated values, or an explanatory
 * placeholder if the job hasn't finished yet or produced nothing.
 */
export default function ResultsTable({ results, isTerminal }: ResultsTableProps) {
  if (!isTerminal) {
    return <div className="text-sm text-neutral-400 py-8 text-center">{JOB_DETAIL_TEXTS.RESULTS_PENDING}</div>;
  }
  if (!results || results.length === 0) {
    return <div className="text-sm text-neutral-400 py-8 text-center">{JOB_DETAIL_TEXTS.RESULTS_EMPTY}</div>;
  }
  return (
    <table className="w-full text-sm">
      <thead>
        <tr className="text-left text-neutral-500 border-b border-neutral-200">
          <th className="py-2 font-medium">Group Key</th>
          <th className="py-2 font-medium">Aggregated Value</th>
        </tr>
      </thead>
      <tbody>
        {results.map((r) => (
          <tr key={r.id} className="border-b border-neutral-100">
            <td className="py-2 font-mono text-xs text-neutral-700">{r.group_key}</td>
            <td className="py-2 text-neutral-900">{r.aggregated_value.toLocaleString()}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
