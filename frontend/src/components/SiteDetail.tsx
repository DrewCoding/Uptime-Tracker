import type { HealthCheck } from '../hooks/useChecks';
import { useCheckHistory } from '../hooks/useChecks';
import './SiteDetail.css';

interface SiteDetailProps {
  url: string;
  onClose: () => void;
}

function getDisplayName(url: string): string {
  try {
    return new URL(url).hostname.replace('www.', '');
  } catch {
    return url;
  }
}

function getStatusLabel(check: HealthCheck): { text: string; className: string } {
  if (check.status_code === null) return { text: 'Offline', className: 'offline' };
  if (check.status_code >= 200 && check.status_code < 400) return { text: 'Online', className: 'online' };
  return { text: 'Error', className: 'error' };
}

function formatTime(dateStr: string): string {
  return new Date(dateStr).toLocaleString();
}

export default function SiteDetail({ url, onClose }: SiteDetailProps) {
  const { history, loading, error } = useCheckHistory(url);

  return (
    <div className="site-detail-overlay" onClick={onClose}>
      <div className="site-detail" onClick={(e) => e.stopPropagation()}>
        <div className="site-detail__header">
          <div>
            <h2 className="site-detail__title">{getDisplayName(url)}</h2>
            <p className="site-detail__url">{url}</p>
          </div>
          <button className="site-detail__close" onClick={onClose} aria-label="Close">
            ✕
          </button>
        </div>

        <div className="site-detail__body">
          {loading && <p className="site-detail__message">Loading history…</p>}
          {error && <p className="site-detail__message site-detail__error">Error: {error}</p>}

          {!loading && !error && history.length === 0 && (
            <p className="site-detail__message">No history found.</p>
          )}

          {!loading && history.length > 0 && (
            <div className="site-detail__table-wrap">
              <table className="site-detail__table">
                <thead>
                  <tr>
                    <th>Status</th>
                    <th>Code</th>
                    <th>Latency</th>
                    <th>Checked At</th>
                  </tr>
                </thead>
                <tbody>
                  {history.map((check) => {
                    const status = getStatusLabel(check);
                    return (
                      <tr key={check.id}>
                        <td>
                          <span className={`site-detail__badge site-detail__badge--${status.className}`}>
                            {status.text}
                          </span>
                        </td>
                        <td className="site-detail__code">
                          {check.status_code !== null ? check.status_code : '—'}
                        </td>
                        <td>{check.latency_ms}ms</td>
                        <td className="site-detail__time">{formatTime(check.checked_at)}</td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
