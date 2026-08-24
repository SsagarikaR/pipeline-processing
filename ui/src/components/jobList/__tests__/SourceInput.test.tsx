import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import SourceInput from '../SourceInput';
import { jobSpecSchema, type JobSpecFormValues } from '../../../schemas/jobSpec';

const VALID_SPEC: JobSpecFormValues = {
  sources: [{ type: 'csv', path: '' }],
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
        <SourceInput />
        <button type="submit">Submit</button>
      </form>
    </FormProvider>
  );
}

describe('SourceInput', () => {
  it('renders a row per source and handles additions', () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, sources: [] }} />);

    expect(screen.getByText('Sources')).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('/path/to/file or Data URI')).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('+ Add source'));
    expect(screen.getByPlaceholderText('/path/to/file or Data URI')).toBeInTheDocument();
  });

  it('removes a row', () => {
    render(<Harness />);

    expect(screen.getByPlaceholderText('/path/to/file or Data URI')).toBeInTheDocument();
    fireEvent.click(screen.getByText('✕'));
    expect(screen.queryByPlaceholderText('/path/to/file or Data URI')).not.toBeInTheDocument();
  });

  it('shows a validation error for an empty path on submit', async () => {
    render(<Harness />);

    fireEvent.click(screen.getByText('Submit'));

    expect(await screen.findByText('Path is required')).toBeInTheDocument();
  });
});
