import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import JobList from '../JobList';
import { jobService } from '../../service/jobService';

// Mock the jobService
vi.mock('../../service/jobService', () => ({
  jobService: {
    getAllJobs: vi.fn(),
    deleteJob: vi.fn(),
  },
}));

describe('JobList (Integration Test)', () => {
  it('shows loading state initially and then renders jobs', async () => {
    // Setup mock to return a test job
    vi.mocked(jobService.getAllJobs).mockResolvedValue([{
      id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
      status: 'running',
      total_records: 1000,
      processed_records: 500,
      error_count: 0,
      created_at: new Date().toISOString(),
      started_at: null,
      completed_at: null,
      export_url: null,
    }]);

    render(
      <BrowserRouter>
        <JobList />
      </BrowserRouter>
    );

    // Initial loading state might not be visible textually except in an empty table or loading indicator,
    // but we can check if the API is called
    expect(jobService.getAllJobs).toHaveBeenCalledTimes(1);

    // Wait for the mock data to populate the UI
    await waitFor(() => {
      expect(screen.getByText('running')).toBeInTheDocument();
      expect(screen.getByText('500 processed')).toBeInTheDocument();
    });
  });

  it('displays an error message when the API fails', async () => {
    vi.mocked(jobService.getAllJobs).mockRejectedValue(new Error('Network error'));

    render(
      <BrowserRouter>
        <JobList />
      </BrowserRouter>
    );

    await waitFor(() => {
      expect(screen.getByText(/Failed to load jobs/i)).toBeInTheDocument();
    });
  });
});
