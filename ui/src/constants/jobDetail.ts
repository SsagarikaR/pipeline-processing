/** Copy used on the job detail page. */
export const JOB_DETAIL_TEXTS = {
  JOB_TITLE: 'Pipeline Job',
  REFRESH_TITLE: 'Refresh Job Data',
  CANCEL_MODAL_TITLE: 'Cancel Job',
  /** Confirmation message naming the job by its title, so it's clear which one is about to be cancelled. */
  cancelModalMessage: (title: string) => `Cancel "${title}"?`,
  DELETE_MODAL_TITLE: 'Delete Job',
  /** Confirmation message naming the job by its title, so it's clear which one is about to go. */
  deleteModalMessage: (title: string) => `Delete "${title}" and its artifacts? This cannot be undone.`,
  FILE_PREVIEW: 'File Preview',
  LOADING_PREVIEW: 'Loading preview...',
  PREVIEW_ERROR: 'Failed to load preview.',
  DOWNLOAD_BTN: 'Download',
  CLOSE_PREVIEW_ARIA: 'Close preview',
  /** Accessible name for the clickable file row that opens the preview. */
  previewFileAria: (path: string) => `Preview ${path}`,
  STARTED: 'Started: ',
  COMPLETED: 'Completed: ',
  EXPORT: 'Export: ',
  TAB_PROGRESS: 'Live Progress',
  /** Results tab label with its current count, e.g. "Results (3)". */
  tabResults: (count: number) => `Results (${count})`,
  /** Errors tab label with its current count, e.g. "Errors (1)". */
  tabErrors: (count: number) => `Errors (${count})`,
  RUNNING_MSG: 'Job is running - this page auto-refreshes every 2 seconds.',
  /** Summary line shown once a job reaches a terminal status. */
  finishedMsg: (status: string) => `Job finished with status "${status}".`,
  RESULTS_EMPTY: 'No results produced.',
  RESULTS_PENDING: 'Results will appear once the job finishes.',
  ERRORS_EMPTY: 'No errors reported.',
  FETCH_ERROR: 'Failed to load job - check that it exists.',
  INVALID_ID: 'Invalid job ID.',
};
