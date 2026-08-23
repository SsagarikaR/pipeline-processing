import { render, screen, fireEvent } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { FormProvider, useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import TransformInput from '../TransformInput';
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
        <TransformInput />
        <button type="submit">Submit</button>
      </form>
    </FormProvider>
  );
}

describe('TransformInput', () => {
  it('renders a row per transform and handles additions', () => {
    render(<Harness />);

    expect(screen.getByText('Transforms')).toBeInTheDocument();
    expect(screen.queryByPlaceholderText('field name')).not.toBeInTheDocument();

    fireEvent.click(screen.getByText('+ Add transform'));
    expect(screen.getByPlaceholderText('field name')).toBeInTheDocument();
  });

  it('removes a row', () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, transforms: [{ name: 'uppercase', params: { field: 'name' } }] }} />);

    expect(screen.getByPlaceholderText('field name')).toBeInTheDocument();
    fireEvent.click(screen.getByText('✕'));
    expect(screen.queryByPlaceholderText('field name')).not.toBeInTheDocument();
  });

  it('shows a validation error for an empty transform field on submit', async () => {
    render(<Harness defaultValues={{ ...VALID_SPEC, transforms: [{ name: 'uppercase', params: { field: '' } }] }} />);

    fireEvent.click(screen.getByText('Submit'));

    expect(await screen.findByText('Field is required')).toBeInTheDocument();
  });
});
