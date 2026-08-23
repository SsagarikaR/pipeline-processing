import type { JobStatus } from '../types/job';

export const STATUS_BADGE_STYLES: Record<JobStatus, string> = {
  pending: 'bg-neutral-100 text-neutral-700 border-neutral-300',
  running: 'bg-brand-100 text-brand-700 border-brand-300 animate-pulse',
  completed: 'bg-success-100 text-success-700 border-success-300',
  failed: 'bg-danger-100 text-danger-700 border-danger-300',
  cancelled: 'bg-amber-100 text-amber-700 border-amber-300',
};

export const ROUTES = {
  HOME: '/',
  JOB_DETAIL: '/jobs/:id',
  jobDetail: (id: string | number) => `/jobs/${id}`,
};

export const COMMON_LABELS = {
  BACK_TO_JOBS: '← Back to Jobs',
  DELETE: 'Delete',
  CANCEL: 'Cancel',
  CONFIRM: 'Confirm',
  REFRESH: 'Refresh',
  LOADING: 'Loading…',
  REMOVE: '✕',
};

export const ERROR_BOUNDARY_TEXTS = {
  TITLE: 'Something went wrong',
  MESSAGE: 'An unexpected error occurred in the application.',
  RELOAD_BUTTON: 'Reload Page',
  toastMessage: (error: string) => `A critical error occurred: ${error}`,
};
