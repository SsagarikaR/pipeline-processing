import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import CreateJobModal from '../CreateJobModal';
import { jobService } from '../../../service/jobService';

vi.mock('../../../service/jobService', () => ({
  jobService: {
    createJob: vi.fn(),
  },
}));

describe('CreateJobModal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders nothing when closed', () => {
    render(
      <BrowserRouter>
        <CreateJobModal isOpen={false} onClose={vi.fn()} />
      </BrowserRouter>
    );

    expect(screen.queryByText('New Pipeline Job')).not.toBeInTheDocument();
  });

  it('renders form and handles submission when open', async () => {
    vi.mocked(jobService.createJob).mockResolvedValue({ id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', status: 'pending' } as any);
    const onClose = vi.fn();

    render(
      <BrowserRouter>
        <CreateJobModal isOpen={true} onClose={onClose} />
      </BrowserRouter>
    );

    expect(screen.getByText('New Pipeline Job')).toBeInTheDocument();

    fireEvent.change(screen.getByPlaceholderText('/path/to/file or Data URI'), { target: { value: 'in.csv' } });
    fireEvent.change(screen.getByPlaceholderText('field'), { target: { value: 'amount' } });
    fireEvent.change(screen.getByPlaceholderText('e.g., results.json'), { target: { value: 'out.json' } });

    fireEvent.click(screen.getByText('Create Job'));

    await waitFor(() => {
      expect(jobService.createJob).toHaveBeenCalledTimes(1);
      expect(onClose).toHaveBeenCalledTimes(1);
    });
  });

  it('blocks submission and shows field errors when required fields are empty', async () => {
    render(
      <BrowserRouter>
        <CreateJobModal isOpen={true} onClose={vi.fn()} />
      </BrowserRouter>
    );

    fireEvent.click(screen.getByText('Create Job'));

    // Both the source and export sections are empty by default, so this
    // message appears more than once.
    expect((await screen.findAllByText('Path is required')).length).toBeGreaterThan(0);
    expect(jobService.createJob).not.toHaveBeenCalled();
  });

  it('calls onClose when cancel is clicked', () => {
    const onClose = vi.fn();

    render(
      <BrowserRouter>
        <CreateJobModal isOpen={true} onClose={onClose} />
      </BrowserRouter>
    );

    fireEvent.click(screen.getByText('Cancel'));
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
