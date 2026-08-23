import axiosInstance from '../config/axios';
import type { Job, JobSpec, ProgressResponse, Result, JobError } from '../types/job';
import { API_ROUTES } from '../constants/api';

/**
 * Thin wrapper around the backend's `/pipelines` API. Every method
 * returns the parsed response body directly (errors are handled by the
 * axios interceptor, so callers just need a try/catch).
 */
export const jobService = {
  /** Submits a new job spec and returns the created job. */
  createJob: (spec: JobSpec): Promise<Job> =>
    axiosInstance.post<Job>(API_ROUTES.PIPELINES, spec).then((r) => r.data),

  /** Lists every job. */
  getAllJobs: (): Promise<Job[]> =>
    axiosInstance.get<Job[]>(API_ROUTES.PIPELINES).then((r) => r.data),

  /** Fetches a single job by ID. */
  getJob: (id: string): Promise<Job> =>
    axiosInstance.get<Job>(API_ROUTES.pipeline(id)).then((r) => r.data),

  /** Fetches a job's live progress (processed/error counts, percent complete). */
  getProgress: (id: string): Promise<ProgressResponse> =>
    axiosInstance.get<ProgressResponse>(API_ROUTES.pipelineProgress(id)).then((r) => r.data),

  /** Fetches a job's aggregated results. */
  getResults: (id: string): Promise<Result[]> =>
    axiosInstance.get<Result[]>(API_ROUTES.pipelineResults(id)).then((r) => r.data),

  /** Fetches the records that failed processing for a job. */
  getErrors: (id: string): Promise<JobError[]> =>
    axiosInstance.get<JobError[]>(API_ROUTES.pipelineErrors(id)).then((r) => r.data),

  /** Cancels a running job. */
  cancelJob: (id: string): Promise<void> =>
    axiosInstance.patch(API_ROUTES.pipelineCancel(id)).then(() => undefined),

  /** Deletes a job and its results/errors. */
  deleteJob: (id: string): Promise<void> =>
    axiosInstance.delete(API_ROUTES.pipeline(id)).then(() => undefined),
};
