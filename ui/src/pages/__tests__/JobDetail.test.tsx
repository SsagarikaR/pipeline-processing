import { render, screen, waitFor, fireEvent } from '@testing-library/react';
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
    getJob: vi.fn(),
    getProgress: vi.fn(),
    getResults: vi.fn(),
    getErrors: vi.fn(),
  }
}));

describe('JobDetail', () => {
  it('loads and displays job progress', async () => {
    vi.mocked(jobService.getJob).mockResolvedValue({
      id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
      status: 'running',
      spec: {
        sources: [{ type: 'csv', path: 'in.csv' }],
        transforms: [],
        aggregations: [{ field: 'amount', op: 'sum' }],
        exports: [{ type: 's3', path: 'results.json' }],
        concurrency: { validateWorkers: 4, transformWorkers: 4 },
      },
      total_records: 1000,
      processed_records: 100,
      error_count: 0,
      created_at: new Date().toISOString(),
      started_at: new Date().toISOString(),
      completed_at: null,
      export_url: null,
    });
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
      expect(screen.getByText('results.json')).toBeInTheDocument();
      expect(screen.getByText('100')).toBeInTheDocument(); // Processed
      expect(screen.getByText('10.5%')).toBeInTheDocument(); // percent
      expect(screen.getByText('5.0/s')).toBeInTheDocument(); // rate
    });
  });

  it('previews the export file and downloads it', async () => {
    vi.mocked(jobService.getJob).mockResolvedValue({
      id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
      status: 'completed',
      spec: {
        sources: [{ type: 'csv', path: 'in.csv' }],
        transforms: [],
        aggregations: [{ field: 'amount', op: 'sum' }],
        exports: [{ type: 's3', path: 'exports/results.json' }],
        concurrency: { validateWorkers: 4, transformWorkers: 4 },
      },
      total_records: 1000,
      processed_records: 1000,
      error_count: 0,
      created_at: new Date().toISOString(),
      started_at: new Date().toISOString(),
      completed_at: new Date().toISOString(),
      export_url: 'https://example.com/exports/results.json',
    });
    vi.mocked(jobService.getProgress).mockResolvedValue({
      jobId: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d',
      status: 'completed',
      processed: 1000,
      errorCount: 0,
      percentComplete: 100,
      recordsPerSec: 0,
      startedAt: new Date().toISOString(),
      completedAt: new Date().toISOString(),
      exportUrl: 'https://example.com/exports/results.json',
    });
    vi.mocked(jobService.getResults).mockResolvedValue([]);
    vi.mocked(jobService.getErrors).mockResolvedValue([]);

    const fetchMock = vi.fn().mockResolvedValue({
      text: () => Promise.resolve('{"total": 42}'),
    });
    vi.stubGlobal('fetch', fetchMock);

    const createObjectURL = vi.fn().mockReturnValue('blob:mock-url');
    const revokeObjectURL = vi.fn();
    vi.stubGlobal('URL', { ...URL, createObjectURL, revokeObjectURL });
    const clickSpy = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => { });

    render(
      <BrowserRouter>
        <JobDetail />
      </BrowserRouter>
    );

    const previewTrigger = await screen.findByRole('button', { name: 'Preview exports/results.json' });
    fireEvent.click(previewTrigger);

    expect(fetchMock).toHaveBeenCalledWith('https://example.com/exports/results.json');

    await waitFor(() => {
      expect(screen.getByText(/"total": 42/)).toBeInTheDocument();
    });

    const downloadBtn = screen.getByRole('button', { name: 'Download' });
    fireEvent.click(downloadBtn);

    expect(createObjectURL).toHaveBeenCalled();
    expect(clickSpy).toHaveBeenCalled();
    expect(revokeObjectURL).toHaveBeenCalledWith('blob:mock-url');

    vi.unstubAllGlobals();
  });
});
