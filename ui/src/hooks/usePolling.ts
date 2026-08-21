import { useEffect, useRef, useState } from 'react';

interface UsePollingResult<T> {
  data: T | null;
  error: Error | null;
  refresh: () => Promise<void>;
}

/**
 * Repeatedly calls `fetchFn` every `intervalMs` while `active` is true.
 */
export default function usePolling<T>(
  fetchFn: () => Promise<T>,
  intervalMs: number,
  active: boolean
): UsePollingResult<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<Error | null>(null);
  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

  async function refresh() {
    try {
      const result = await fetchFn();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err : new Error(String(err)));
    }
  }

  useEffect(() => {
    let cancelled = false;

    async function tick() {
      try {
        const result = await fetchFn();
        if (!cancelled) setData(result);
      } catch (err) {
        if (!cancelled) setError(err instanceof Error ? err : new Error(String(err)));
      }
    }

    tick();

    if (active) {
      timerRef.current = setInterval(tick, intervalMs);
    }

    return () => {
      cancelled = true;
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, [fetchFn, intervalMs, active]);

  return { data, error, refresh };
}