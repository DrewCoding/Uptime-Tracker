import { useState, useEffect, useCallback } from 'react';

export interface HealthCheck {
  id: number;
  url: string;
  status_code: number | null;
  latency_ms: number;
  checked_at: string;
}

interface UseChecksResult {
  checks: HealthCheck[];
  loading: boolean;
  error: string | null;
  refresh: () => void;
}

interface UseCheckHistoryResult {
  history: HealthCheck[];
  loading: boolean;
  error: string | null;
}

const API_BASE = '/api';

/**
 * Fetches the latest health check for every monitored URL.
 * Auto-refreshes every `intervalMs` milliseconds (default 30s).
 */
export function useChecks(intervalMs = 30000): UseChecksResult {
  const [checks, setChecks] = useState<HealthCheck[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchChecks = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/checks`);
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data: HealthCheck[] = await res.json();
      setChecks(data);
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchChecks();
    const id = setInterval(fetchChecks, intervalMs);
    return () => clearInterval(id);
  }, [fetchChecks, intervalMs]);

  return { checks, loading, error, refresh: fetchChecks };
}

/**
 * Fetches the check history for a specific URL.
 */
export function useCheckHistory(url: string | null, limit = 50): UseCheckHistoryResult {
  const [history, setHistory] = useState<HealthCheck[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!url) {
      setHistory([]);
      return;
    }

    setLoading(true);
    fetch(`${API_BASE}/checks?url=${encodeURIComponent(url)}&limit=${limit}`)
      .then((res) => {
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        return res.json();
      })
      .then((data: HealthCheck[]) => {
        setHistory(data);
        setError(null);
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : 'Failed to fetch');
      })
      .finally(() => setLoading(false));
  }, [url, limit]);

  return { history, loading, error };
}
