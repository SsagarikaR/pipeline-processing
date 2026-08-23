export const API_ROUTES = {
  PIPELINES: '/pipelines',
  pipeline: (id: string | number) => `/pipelines/${id}`,
  pipelineProgress: (id: string | number) => `/pipelines/${id}/progress`,
  pipelineResults: (id: string | number) => `/pipelines/${id}/results`,
  pipelineErrors: (id: string | number) => `/pipelines/${id}/errors`,
  pipelineCancel: (id: string | number) => `/pipelines/${id}/cancel`,
};
