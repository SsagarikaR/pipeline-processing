/** Copy used on the job detail page. */
export const JOB_DETAIL_TEXTS = {
  JOB_TITLE_PREFIX: 'Job #',
  REFRESH_TITLE: 'Refresh Job Data',
  CANCEL_MODAL_TITLE: 'Cancel Job',
  CANCEL_MODAL_MSG: 'Cancel this job?',
  DELETE_MODAL_TITLE: 'Delete Job',
  DELETE_MODAL_MSG: 'Delete this job and its artifacts? This cannot be undone.',
  FILE_PREVIEW: 'File Preview',
  LOADING_PREVIEW: 'Loading preview...',
  PREVIEW_ERROR: 'Failed to load preview.',
  PREVIEW_BTN: 'Preview',
  STARTED: 'Started: ',
  COMPLETED: 'Completed: ',
  EXPORT: 'Export: ',
  TAB_PROGRESS: 'Live Progress',
  /** Results tab label with its current count, e.g. "Results (3)". */
  tabResults: (count: number) => `Results (${count})`,
  /** Errors tab label with its current count, e.g. "Errors (1)". */
  tabErrors: (count: number) => `Errors (${count})`,
  RUNNING_MSG: 'Job is running — this page auto-refreshes every 2 seconds.',
  /** Summary line shown once a job reaches a terminal status. */
  finishedMsg: (status: string) => `Job finished with status "${status}".`,
  RESULTS_EMPTY: 'No results produced.',
  RESULTS_PENDING: 'Results will appear once the job finishes.',
  ERRORS_EMPTY: 'No errors reported.',
  FETCH_ERROR: 'Failed to load job — check that it exists.',
  INVALID_ID: 'Invalid job ID.',
};
