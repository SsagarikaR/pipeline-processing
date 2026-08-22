import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import CreateJob from '../CreateJob';
import { jobService } from '../../service/jobService';

vi.mock('../../service/jobService', () => ({
  jobService: {
    createJob: vi.fn()
  }
}));

describe('CreateJob', () => {
  it('renders form and handles submission', async () => {
    vi.mocked(jobService.createJob).mockResolvedValue({ id: 10, status: 'pending' } as any);
    
    render(
      <BrowserRouter>
        <CreateJob />
      </BrowserRouter>
    );

    expect(screen.getByText('New Pipeline Job')).toBeInTheDocument();
    expect(screen.getByText('Back to Jobs')).toBeInTheDocument();
    
    const submitBtn = screen.getByText('Create Job');
    fireEvent.click(submitBtn);

    await waitFor(() => {
      expect(jobService.createJob).toHaveBeenCalledTimes(1);
    });
  });
});
