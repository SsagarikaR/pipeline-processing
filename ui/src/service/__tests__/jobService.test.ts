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
    const mockJobs = [{ id: 1, status: 'completed' }];
    vi.mocked(axiosInstance.get).mockResolvedValueOnce({ data: mockJobs });

    const jobs = await jobService.getAllJobs();
    
    expect(axiosInstance.get).toHaveBeenCalledWith('/pipelines');
    expect(jobs).toEqual(mockJobs);
  });

  it('cancelJob issues a patch request', async () => {
    vi.mocked(axiosInstance.patch).mockResolvedValueOnce({});
    
    await jobService.cancelJob(42);
    
    expect(axiosInstance.patch).toHaveBeenCalledWith('/pipelines/42/cancel');
  });
});
