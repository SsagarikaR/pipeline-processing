import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { BrowserRouter } from 'react-router-dom';
import JobDetail from '../JobDetail';
import { jobService } from '../../service/jobService';

// Mock react-router-dom to simulate URL parameter
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useParams: () => ({ id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d' }),
  };
});

vi.mock('../../service/jobService', () => ({
  jobService: {
    getProgress: vi.fn(),
    getResults: vi.fn(),
    getErrors: vi.fn(),
  }
}));

describe('JobDetail', () => {
  it('loads and displays job progress', async () => {
    vi.mocked(jobService.getProgress).mockResolvedValue({
      jobId: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
      status: 'running',
      processed: 100,
      errorCount: 0,
      percentComplete: 10.5,
      recordsPerSec: 5,
      startedAt: new Date().toISOString(),
      completedAt: null,
      exportUrl: null,
    });
    vi.mocked(jobService.getResults).mockResolvedValue([]);
    vi.mocked(jobService.getErrors).mockResolvedValue([]);

    render(
      <BrowserRouter>
        <JobDetail />
      </BrowserRouter>
    );

    await waitFor(() => {
      expect(screen.getByText('Job #a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument(); // Processed
      expect(screen.getByText('10.5%')).toBeInTheDocument(); // percent
      expect(screen.getByText('5.0/s')).toBeInTheDocument(); // rate
    });
  });
});
