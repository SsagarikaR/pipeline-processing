import axiosInstance from '../config/axios';
import type { Job, JobSpec, ProgressResponse, Result, JobError } from '../types/job';
import { API_ROUTES } from '../constants/api';

export const jobService = {
  createJob: (spec: JobSpec): Promise<Job> =>
    axiosInstance.post<Job>(API_ROUTES.PIPELINES, spec).then((r) => r.data),

  getAllJobs: (): Promise<Job[]> =>
    axiosInstance.get<Job[]>(API_ROUTES.PIPELINES).then((r) => r.data),

  getJob: (id: number | string): Promise<Job> =>
    axiosInstance.get<Job>(API_ROUTES.pipeline(id)).then((r) => r.data),

  getProgress: (id: number | string): Promise<ProgressResponse> =>
    axiosInstance.get<ProgressResponse>(API_ROUTES.pipelineProgress(id)).then((r) => r.data),

  getResults: (id: number | string): Promise<Result[]> =>
    axiosInstance.get<Result[]>(API_ROUTES.pipelineResults(id)).then((r) => r.data),

  getErrors: (id: number | string): Promise<JobError[]> =>
    axiosInstance.get<JobError[]>(API_ROUTES.pipelineErrors(id)).then((r) => r.data),

  cancelJob: (id: number | string): Promise<void> =>
    axiosInstance.patch(API_ROUTES.pipelineCancel(id)).then(() => undefined),

  deleteJob: (id: number | string): Promise<void> =>
    axiosInstance.delete(API_ROUTES.pipeline(id)).then(() => undefined),
};
