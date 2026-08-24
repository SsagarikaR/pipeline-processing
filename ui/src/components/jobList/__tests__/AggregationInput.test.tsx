import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import AggregationInput from '../AggregationInput';
import { jobSpecSchema, type JobSpecFormValues } from '../../../schemas/jobSpec';

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
        <AggregationInput />
        <button type="submit">Submit</button>
      </form>
    </FormProvider>
  );
}

describe('AggregationInput', () => {
  it('renders a row per aggregation and handles additions', () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, aggregations: [] }} />);

    expect(screen.getByText('Aggregations')).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('field')).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('+ Add aggregation'));
    expect(screen.getByPlaceholderText('field')).toBeInTheDocument();
  });

  it('removes a row', () => {
    render(<Harness />);

    expect(screen.getByPlaceholderText('field')).toBeInTheDocument();
    fireEvent.click(screen.getByText('✕'));
    expect(screen.queryByPlaceholderText('field')).not.toBeInTheDocument();
  });

  it('shows a validation error for an empty aggregation field on submit', async () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, aggregations: [{ field: '', op: 'sum', groupBy: '' }] }} />);

    fireEvent.click(screen.getByText('Submit'));

    expect(await screen.findByText('Field is required')).toBeInTheDocument();
  });
});
