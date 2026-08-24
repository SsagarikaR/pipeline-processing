/** Backend API paths, kept in one place so jobService never hardcodes a URL. */
export const API_ROUTES = {
  PIPELINES: '/pipelines',
  /** URL for one job by ID. */
  pipeline: (id: string) => `/pipelines/${id}`,
  /** URL for a job's live progress. */
  pipelineProgress: (id: string) => `/pipelines/${id}/progress`,
  /** URL for a job's aggregated results. */
  pipelineResults: (id: string) => `/pipelines/${id}/results`,
  /** URL for a job's failed records. */
  pipelineErrors: (id: string) => `/pipelines/${id}/errors`,
  /** URL to cancel a running job. */
  pipelineCancel: (id: string) => `/pipelines/${id}/cancel`,
};
