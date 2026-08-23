import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import ExportInput from '../ExportInput';
import { jobSpecSchema, type JobSpecFormValues } from '../../schemas/jobSpec';

const VALID_SPEC: JobSpecFormValues = {
  sources: [{ type: 'csv', path: 'in.csv' }],
  transforms: [],
  aggregations: [{ field: 'amount', op: 'sum', groupBy: '' }],
  exports: [{ type: 's3', path: 'out.json' }],
  concurrency: { validateWorkers: 4, transformWorkers: 4 },
};

function Harness({ defaultValues = VALID_SPEC }: { defaultValues?: JobSpecFormValues }) {
  const form = useForm<JobSpecFormValues>({
    resolver: zodResolver(jobSpecSchema),
    defaultValues,
  });
  return (
    <FormProvider {...form}>
      <form onSubmit={form.handleSubmit(() => {})}>
        <ExportInput />
        <button type="submit">Submit</button>
      </form>
    </FormProvider>
  );
}

describe('ExportInput', () => {
  it('renders a row per export and handles additions', () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, exports: [] }} />);

    expect(screen.getByText('Exports')).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('e.g., results.json')).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('+ Add export'));
    expect(screen.getByPlaceholderText('e.g., results.json')).toBeInTheDocument();
  });

  it('removes a row', () => {
    render(<Harness />);

    expect(screen.getByPlaceholderText('e.g., results.json')).toBeInTheDocument();
    fireEvent.click(screen.getByText('✕'));
    expect(screen.queryByPlaceholderText('e.g., results.json')).not.toBeInTheDocument();
  });

  it('shows a validation error for an empty export path on submit', async () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, exports: [{ type: 's3', path: '' }] }} />);

    fireEvent.click(screen.getByText('Submit'));

    expect(await screen.findByText('Path is required')).toBeInTheDocument();
  });
});
