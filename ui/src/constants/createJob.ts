import type { AggregationOp } from '../types/job';

export const AGGREGATION_FIELD_TYPES = ['Number', 'String', 'Boolean', 'Date'] as const;

export const AGGREGATION_OPS_BY_FIELD_TYPE: Record<string, AggregationOp[]> = {
  Number: ['sum', 'avg', 'min', 'max', 'count'],
  String: ['count'],
  Boolean: ['count'],
  Date: ['min', 'max', 'count'],
};

export const CREATE_JOB_TEXTS = {
  TITLE: 'New Pipeline Job',
  CLOSE_LABEL: 'Close',
  CONCURRENCY_LABEL: 'Concurrency',
  VALIDATE_WORKERS: 'Validate workers',
  TRANSFORM_WORKERS: 'Transform workers',
  BTN_CREATING: 'Creating…',
  BTN_CREATE: 'Create Job',
  CREATE_ERROR: 'Failed to create job',

  SOURCES_LABEL: 'Sources',
  SOURCE_PATH_PLACEHOLDER: '/path/to/file or Data URI',
  SOURCE_UPLOAD: 'Upload',
  ADD_SOURCE: '+ Add source',

  TRANSFORMS_LABEL: 'Transforms',
  TRANSFORM_FIELD_PLACEHOLDER: 'field name',
  ADD_TRANSFORM: '+ Add transform',

  AGGREGATIONS_LABEL: 'Aggregations',
  AGGREGATION_FIELD_PLACEHOLDER: 'field',
  AGGREGATION_GROUP_BY_PLACEHOLDER: 'group by (optional)',
  ADD_AGGREGATION: '+ Add aggregation',

  EXPORTS_LABEL: 'Exports',
  EXPORT_PATH_PLACEHOLDER: 'e.g., results.json',
  ADD_EXPORT: '+ Add export',
};
