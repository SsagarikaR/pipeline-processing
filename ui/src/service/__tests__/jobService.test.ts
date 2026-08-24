import { describe, it, expect, vi, beforeEach } from 'vitest';
import { jobService } from '../jobService';
import axiosInstance from '../../config/axios';

// Mock axiosInstance
vi.mock('../../config/axios', () => {
  return {
    default: {
      get: vi.fn(),
      post: vi.fn(),
      patch: vi.fn(),
      delete: vi.fn(),
    }
  };
});

describe('jobService (Unit Test)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('getAllJobs fetches successfully', async () => {
    const mockJobs = [{ id: 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', status: 'completed' }];
    vi.mocked(axiosInstance.get).mockResolvedValueOnce({ data: mockJobs });

    const jobs = await jobService.getAllJobs();

    expect(axiosInstance.get).toHaveBeenCalledWith('/pipelines', {
      params: { limit: 50, offset: 0 },
    });
    expect(jobs).toEqual(mockJobs);
  });

  it('cancelJob issues a patch request', async () => {
    vi.mocked(axiosInstance.patch).mockResolvedValueOnce({});

    await jobService.cancelJob('a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d');

    expect(axiosInstance.patch).toHaveBeenCalledWith('/pipelines/a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d/cancel');
  });
});
