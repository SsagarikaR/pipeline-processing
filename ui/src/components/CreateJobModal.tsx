import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { X } from 'lucide-react';
import { jobService as api } from '../service/jobService';
import type { JobSpec } from '../types/job';
import type { CreateJobModalProps } from '../types/common';
import { jobSpecSchema, type JobSpecFormValues } from '../schemas/jobSpec';
import SourceInput from './SourceInput';
import TransformInput from './TransformInput';
import AggregationInput from './AggregationInput';
import ExportInput from './ExportInput';
import AppButton from './AppButton';
import { ROUTES, COMMON_LABELS } from '../constants/common';
import { CREATE_JOB_TEXTS } from '../constants/createJob';

const EMPTY_SPEC: JobSpecFormValues = {
  sources: [{ type: 'csv', path: '' }],
  transforms: [],
  aggregations: [{ field: '', op: 'sum', groupBy: '' }],
  exports: [{ type: 's3', path: '' }],
  concurrency: { validateWorkers: 4, transformWorkers: 4 },
};

export default function CreateJobModal({ isOpen, onClose }: CreateJobModalProps) {
  const [submitError, setSubmitError] = useState<string | null>(null);
  const navigate = useNavigate();

  const form = useForm<JobSpecFormValues>({
    resolver: zodResolver(jobSpecSchema),
    defaultValues: EMPTY_SPEC,
  });
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = form;

  if (!isOpen) return null;

  function handleClose() {
    reset(EMPTY_SPEC);
    setSubmitError(null);
    onClose();
  }

  const onSubmit = handleSubmit(async (spec) => {
    setSubmitError(null);
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
      reset(EMPTY_SPEC);
      onClose();
      navigate(ROUTES.jobDetail(job.id));
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : CREATE_JOB_TEXTS.CREATE_ERROR);
    }
  });

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4 z-[100]">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-2xl max-h-[90vh] overflow-y-auto p-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-semibold text-neutral-900">{CREATE_JOB_TEXTS.TITLE}</h1>
          <button
            onClick={handleClose}
            className="text-neutral-400 hover:text-neutral-600 transition-colors"
            aria-label={CREATE_JOB_TEXTS.CLOSE_LABEL}
          >
            <X size={20} />
          </button>
        </div>

        <FormProvider {...form}>
          <form onSubmit={onSubmit} className="space-y-6" noValidate>
            <SourceInput />
            <TransformInput />
            <AggregationInput />
            <ExportInput />

            <div>
              <label className="block text-sm font-medium text-neutral-700 mb-2">{CREATE_JOB_TEXTS.CONCURRENCY_LABEL}</label>
              <div className="flex gap-4">
                <div>
                  <span className="text-xs text-neutral-500">
                    {CREATE_JOB_TEXTS.VALIDATE_WORKERS} <span className="text-danger-600">*</span>
                  </span>
                  <input
                    type="number"
                    min={1}
                    {...register('concurrency.validateWorkers', { valueAsNumber: true })}
                    className="block border border-neutral-300 rounded-lg px-3 py-2 text-sm w-24"
                  />
                  {errors.concurrency?.validateWorkers && (
                    <p className="text-xs text-danger-600 mt-1">{errors.concurrency.validateWorkers.message}</p>
                  )}
                </div>
                <div>
                  <span className="text-xs text-neutral-500">
                    {CREATE_JOB_TEXTS.TRANSFORM_WORKERS} <span className="text-danger-600">*</span>
                  </span>
                  <input
                    type="number"
                    min={1}
                    {...register('concurrency.transformWorkers', { valueAsNumber: true })}
                    className="block border border-neutral-300 rounded-lg px-3 py-2 text-sm w-24"
                  />
                  {errors.concurrency?.transformWorkers && (
                    <p className="text-xs text-danger-600 mt-1">{errors.concurrency.transformWorkers.message}</p>
                  )}
                </div>
              </div>
            </div>

            {submitError && (
              <div className="text-sm text-danger-600 bg-danger-50 border border-danger-200 rounded-lg p-3">{submitError}</div>
            )}

            <div className="flex justify-end gap-3">
              <AppButton type="button" variant="secondary" onClick={handleClose}>
                {COMMON_LABELS.CANCEL}
              </AppButton>
              <AppButton type="submit" disabled={isSubmitting} size="lg">
                {isSubmitting ? CREATE_JOB_TEXTS.BTN_CREATING : CREATE_JOB_TEXTS.BTN_CREATE}
              </AppButton>
            </div>
          </form>
        </FormProvider>
      </div>
    </div>
  );
}
