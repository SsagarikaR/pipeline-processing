import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import CreateJobModal from '../CreateJobModal';
import { jobService } from '../../service/jobService';

vi.mock('../../service/jobService', () => ({
  jobService: {
    createJob: vi.fn(),
  },
}));

describe('CreateJobModal', () => {
  it('renders nothing when closed', () => {
    render(
      <BrowserRouter>
        <CreateJobModal isOpen={false} onClose={vi.fn()} />
      </BrowserRouter>
    );

    expect(screen.queryByText('New Pipeline Job')).not.toBeInTheDocument();
  });

  it('renders form and handles submission when open', async () => {
    vi.mocked(jobService.createJob).mockResolvedValue({ id: 10, status: 'pending' } as any);
    const onClose = vi.fn();

    render(
      <BrowserRouter>
        <CreateJobModal isOpen={true} onClose={onClose} />
      </BrowserRouter>
    );

    expect(screen.getByText('New Pipeline Job')).toBeInTheDocument();

    const submitBtn = screen.getByText('Create Job');
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(jobService.createJob).toHaveBeenCalledTimes(1);
      expect(onClose).toHaveBeenCalledTimes(1);
    });
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
