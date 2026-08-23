/// <reference types="node" />
import { createServer, type Server } from 'http';
import type { AddressInfo } from 'net';
import { describe, it, expect, vi, beforeAll, afterAll } from 'vitest';
import toast from 'react-hot-toast';
import axiosInstance from '../axios';

vi.mock('react-hot-toast', () => ({
  default: { error: vi.fn() },
}));

describe('axiosInstance', () => {
  let server: Server;
  let baseURL: string;
  let receivedApiKeyHeader: string | undefined;

  beforeAll(async () => {
    server = createServer((req, res) => {
      res.setHeader('Access-Control-Allow-Origin', '*');
      res.setHeader('Access-Control-Allow-Headers', '*');
      if (req.method === 'OPTIONS') {
        res.writeHead(204);
        res.end();
        return;
      }
      receivedApiKeyHeader = req.headers['x-api-key'] as string | undefined;
      res.writeHead(401, { 'Content-Type': 'text/plain; charset=utf-8' });
      res.end('401 Unauthorized - Invalid or missing API Key\n');
    });
    await new Promise<void>((resolve) => server.listen(0, resolve));
    const { port } = server.address() as AddressInfo;
    baseURL = `http://localhost:${port}`;
    axiosInstance.defaults.baseURL = baseURL;
  });

  afterAll(() => {
    server.close();
  });

  it('sends an X-API-Key header on requests', async () => {
    await axiosInstance.get('/pipelines').catch(() => undefined);
    expect(receivedApiKeyHeader).toBeTruthy();
  });

  it('surfaces the backend plain-text error body in the toast, not the generic axios message', async () => {
    vi.mocked(toast.error).mockClear();

    await axiosInstance.get('/pipelines').catch(() => undefined);

    expect(toast.error).toHaveBeenCalledWith('401 Unauthorized - Invalid or missing API Key');
  });
});
