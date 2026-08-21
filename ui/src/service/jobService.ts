import axiosInstance from '../config/axios';
import type { Job, JobSpec, ProgressResponse, Result, JobError } from '../types/job';

export const jobService = {
  createJob: (spec: JobSpec): Promise<Job> =>
    axiosInstance.post<Job>('/pipelines', spec).then((r) => r.data),

  getAllJobs: (): Promise<Job[]> =>
    axiosInstance.get<Job[]>('/pipelines').then((r) => r.data),

  getJob: (id: number | string): Promise<Job> =>
    axiosInstance.get<Job>(`/pipelines/${id}`).then((r) => r.data),

  getProgress: (id: number | string): Promise<ProgressResponse> =>
    axiosInstance.get<ProgressResponse>(`/pipelines/${id}/progress`).then((r) => r.data),

  getResults: (id: number | string): Promise<Result[]> =>
    axiosInstance.get<Result[]>(`/pipelines/${id}/results`).then((r) => r.data),

  getErrors: (id: number | string): Promise<JobError[]> =>
    axiosInstance.get<JobError[]>(`/pipelines/${id}/errors`).then((r) => r.data),

  cancelJob: (id: number | string): Promise<void> =>
    axiosInstance.patch(`/pipelines/${id}/cancel`).then(() => undefined),

  deleteJob: (id: number | string): Promise<void> =>
    axiosInstance.delete(`/pipelines/${id}`).then(() => undefined),
};
