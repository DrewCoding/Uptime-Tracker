import type { HealthCheck } from '../hooks/useChecks';
import './StatusCard.css';

interface StatusCardProps {
  check: HealthCheck;
  onClick: () => void;
}

function getStatus(check: HealthCheck): 'online' | 'error' | 'offline' {
  if (check.status_code === null) return 'offline';
  if (check.status_code >= 200 && check.status_code < 400) return 'online';
  return 'error';
}

function getDisplayName(url: string): string {
  try {
    const hostname = new URL(url).hostname.replace('www.', '');
    return hostname;
  } catch {
    return url;
  }
}

function timeAgo(dateStr: string): string {
  const now = Date.now();
  const then = new Date(dateStr).getTime();
  const diffSec = Math.floor((now - then) / 1000);

  if (diffSec < 60) return `${diffSec}s ago`;
  if (diffSec < 3600) return `${Math.floor(diffSec / 60)}m ago`;
  if (diffSec < 86400) return `${Math.floor(diffSec / 3600)}h ago`;
  return `${Math.floor(diffSec / 86400)}d ago`;
}

export default function StatusCard({ check, onClick }: StatusCardProps) {
  const status = getStatus(check);

  return (
    <button className={`status-card status-card--${status}`} onClick={onClick} id={`card-${check.id}`}>
      <div className="status-card__indicator">
        <span className={`status-dot status-dot--${status}`} />
      </div>
      <div className="status-card__body">
        <h3 className="status-card__name">{getDisplayName(check.url)}</h3>
        <div className="status-card__meta">
          <span className="status-card__latency">{check.latency_ms}ms</span>
          <span className="status-card__separator">·</span>
          <span className="status-card__time">{timeAgo(check.checked_at)}</span>
        </div>
      </div>
      <div className="status-card__code">
        {check.status_code !== null ? check.status_code : '—'}
      </div>
    </button>
  );
}
