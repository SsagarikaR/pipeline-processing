export const ROUTES = {
  HOME: '/',
  CREATE_JOB: '/create',
  JOB_DETAIL: '/jobs/:id',
  jobDetail: (id: string | number) => `/jobs/${id}`,
};

export const COMMON_LABELS = {
  BACK_TO_JOBS: '← Back to Jobs',
  DELETE: 'Delete',
  CANCEL: 'Cancel',
  REFRESH: 'Refresh',
  LOADING: 'Loading…',
};
