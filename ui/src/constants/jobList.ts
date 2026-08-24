/** Copy used on the job list page. */
export const JOB_LIST_TEXTS = {
  TITLE: 'Pipeline Jobs',
  LAST_FETCHED: 'Last fetched: ',
  NEW_JOB: '+ New Job',
  NO_JOBS: 'No jobs yet - create one to get started.',
  NO_JOBS_FOUND: 'No jobs found',
  LOADING_JOBS: 'Loading jobs…',
  PROCESSED: 'processed',
  ERRORS: 'errors',
  DELETE_MODAL_TITLE: 'Delete Job',
  /** Confirmation message naming the job by its title, so it's clear which one is about to go. */
  deleteModalMessage: (title: string) => `Delete "${title}"? This removes its results and errors too.`,
  FETCH_ERROR: 'Failed to load jobs - is the backend running ?',
};
