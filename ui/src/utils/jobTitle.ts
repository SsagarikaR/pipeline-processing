import type { Job } from '../types/job';

/**
 * A job's display title: the file name/key from its first configured
 * export - the one bit of naming the user actually typed in when they
 * created it, so it's more recognizable than the job's ID. Falls back to
 * its creation time for the rare case a job has no exports at all.
 */
export function jobTitle(job: Job): string {
  return job.spec?.exports?.[0]?.path || new Date(job.created_at).toLocaleString();
}
