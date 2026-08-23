import { z } from 'zod';

export const sourceSchema = z.object({
  type: z.enum(['csv', 'json']),
  path: z.string().min(1, 'Path is required'),
});

export const transformSchema = z.object({
  name: z.enum(['uppercase', 'lowercase']),
  params: z.object({ field: z.string().min(1, 'Field is required') }).catchall(z.string()),
});

export const aggregationSchema = z.object({
  field: z.string().min(1, 'Field is required'),
  op: z.enum(['sum', 'avg', 'count', 'min', 'max']),
  groupBy: z.string().optional(),
});

export const exportSchema = z.object({
  type: z.enum(['s3', 'db']),
  path: z.string().min(1, 'Path is required'),
});

export const jobSpecSchema = z.object({
  sources: z.array(sourceSchema).min(1, 'At least one source is required'),
  transforms: z.array(transformSchema),
  aggregations: z.array(aggregationSchema).min(1, 'At least one aggregation is required'),
  exports: z.array(exportSchema).min(1, 'At least one export is required'),
  concurrency: z.object({
    validateWorkers: z.number('Must be a number').int().min(1, 'Must be at least 1'),
    transformWorkers: z.number('Must be a number').int().min(1, 'Must be at least 1'),
  }),
});

export type JobSpecFormValues = z.infer<typeof jobSpecSchema>;
