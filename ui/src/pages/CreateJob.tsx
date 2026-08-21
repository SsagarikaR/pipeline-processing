import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { jobService as api } from '../service/jobService';
import type { JobSpec } from '../types/job';
import SourceInput from '../components/SourceInput';
import TransformInput from '../components/TransformInput';
import AggregationInput from '../components/AggregationInput';
import ExportInput from '../components/ExportInput';
import AppButton from '../components/AppButton';

const EMPTY_SPEC: JobSpec = {
  sources: [{ type: 'csv', path: '' }],
  transforms: [],
  aggregations: [{ field: '', op: 'sum', groupBy: '' }],
  exports: [{ type: 's3', path: '' }],
  concurrency: { validateWorkers: 4, transformWorkers: 4 },
};

export default function CreateJob() {
  const [spec, setSpec] = useState<JobSpec>(EMPTY_SPEC);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setSubmitting(true);
    setError(null);
    try {
      const cleanSpec: JobSpec = {
        ...spec,
        aggregations: spec.aggregations.map((a) => ({
          field: a.field,
          op: a.op,
          ...(a.groupBy ? { groupBy: a.groupBy } : {}),
        })),
      };
      const job = await api.createJob(cleanSpec);
      navigate(`/jobs/${job.id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create job');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-2xl font-semibold text-neutral-900 mb-6">New Pipeline Job</h1>

      <form onSubmit={handleSubmit} className="space-y-6 bg-white border border-neutral-200 rounded-xl p-6">
        <SourceInput sources={spec.sources} onChange={(sources) => setSpec({ ...spec, sources })} />
        <TransformInput transforms={spec.transforms} onChange={(transforms) => setSpec({ ...spec, transforms })} />
        <AggregationInput
          aggregations={spec.aggregations}
          onChange={(aggregations) => setSpec({ ...spec, aggregations })}
        />
        <ExportInput exports={spec.exports} onChange={(exports) => setSpec({ ...spec, exports })} />

        <div>
          <label className="block text-sm font-medium text-neutral-700 mb-2">Concurrency</label>
          <div className="flex gap-4">
            <div>
              <span className="text-xs text-neutral-500">Validate workers</span>
              <input
                type="number"
                min={1}
                value={spec.concurrency.validateWorkers}
                onChange={(e) =>
                  setSpec({
                    ...spec,
                    concurrency: { ...spec.concurrency, validateWorkers: Number(e.target.value) },
                  })
                }
                className="block border border-neutral-300 rounded-lg px-3 py-2 text-sm w-24"
              />
            </div>
            <div>
              <span className="text-xs text-neutral-500">Transform workers</span>
              <input
                type="number"
                min={1}
                value={spec.concurrency.transformWorkers}
                onChange={(e) =>
                  setSpec({
                    ...spec,
                    concurrency: { ...spec.concurrency, transformWorkers: Number(e.target.value) },
                  })
                }
                className="block border border-neutral-300 rounded-lg px-3 py-2 text-sm w-24"
              />
            </div>
          </div>
        </div>

        {error && (
          <div className="text-sm text-danger-600 bg-danger-50 border border-danger-200 rounded-lg p-3">{error}</div>
        )}

        <AppButton
          type="submit"
          disabled={submitting}
          className="w-full"
          size="lg"
        >
          {submitting ? 'Creating…' : 'Create Job'}
        </AppButton>
      </form>
    </div>
  );
}