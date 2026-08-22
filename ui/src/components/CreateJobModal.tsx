import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { X } from 'lucide-react';
import { jobService as api } from '../service/jobService';
import type { JobSpec } from '../types/job';
import type { CreateJobModalProps } from '../types/common';
import SourceInput from './SourceInput';
import TransformInput from './TransformInput';
import AggregationInput from './AggregationInput';
import ExportInput from './ExportInput';
import AppButton from './AppButton';
import { ROUTES } from '../constants/common';
import { CREATE_JOB_TEXTS } from '../constants/createJob';

const EMPTY_SPEC: JobSpec = {
  sources: [{ type: 'csv', path: '' }],
  transforms: [],
  aggregations: [{ field: '', op: 'sum', groupBy: '' }],
  exports: [{ type: 's3', path: '' }],
  concurrency: { validateWorkers: 4, transformWorkers: 4 },
};

export default function CreateJobModal({ isOpen, onClose }: CreateJobModalProps) {
  const [spec, setSpec] = useState<JobSpec>(EMPTY_SPEC);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  if (!isOpen) return null;

  function handleClose() {
    setSpec(EMPTY_SPEC);
    setError(null);
    onClose();
  }

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
      setSpec(EMPTY_SPEC);
      onClose();
      navigate(ROUTES.jobDetail(job.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : CREATE_JOB_TEXTS.CREATE_ERROR);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-[100]">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-semibold text-neutral-900">{CREATE_JOB_TEXTS.TITLE}</h1>
          <button
            onClick={handleClose}
            className="text-neutral-400 hover:text-neutral-600 transition-colors"
            aria-label="Close"
          >
            <X size={20} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <SourceInput sources={spec.sources} onChange={(sources) => setSpec({ ...spec, sources })} />
          <TransformInput transforms={spec.transforms} onChange={(transforms) => setSpec({ ...spec, transforms })} />
          <AggregationInput
            aggregations={spec.aggregations}
            onChange={(aggregations) => setSpec({ ...spec, aggregations })}
          />
          <ExportInput exports={spec.exports} onChange={(exports) => setSpec({ ...spec, exports })} />

          <div>
            <label className="block text-sm font-medium text-neutral-700 mb-2">{CREATE_JOB_TEXTS.CONCURRENCY_LABEL}</label>
            <div className="flex gap-4">
              <div>
                <span className="text-xs text-neutral-500">{CREATE_JOB_TEXTS.VALIDATE_WORKERS}</span>
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
                <span className="text-xs text-neutral-500">{CREATE_JOB_TEXTS.TRANSFORM_WORKERS}</span>
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

          <div className="flex justify-end gap-3">
            <AppButton type="button" variant="secondary" onClick={handleClose}>
              Cancel
            </AppButton>
            <AppButton type="submit" disabled={submitting} size="lg">
              {submitting ? CREATE_JOB_TEXTS.BTN_CREATING : CREATE_JOB_TEXTS.BTN_CREATE}
            </AppButton>
          </div>
        </form>
      </div>
    </div>
  );
}
