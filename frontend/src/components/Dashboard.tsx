import { useState } from 'react';
import { useChecks } from '../hooks/useChecks';
import StatusCard from './StatusCard';
import SiteDetail from './SiteDetail';
import './Dashboard.css';

export default function Dashboard() {
  const { checks, loading, error, refresh } = useChecks();
  const [selectedUrl, setSelectedUrl] = useState<string | null>(null);

  const onlineCount = checks.filter(
    (c) => c.status_code !== null && c.status_code >= 200 && c.status_code < 400
  ).length;

  const avgLatency = checks.length
    ? Math.round(checks.reduce((sum, c) => sum + c.latency_ms, 0) / checks.length)
    : 0;

  return (
    <div className="dashboard">
      <header className="dashboard__header">
        <div className="dashboard__title-row">
          <div>
            <h1 className="dashboard__title">Sentinel</h1>
            <p className="dashboard__subtitle">Dashboard</p>
          </div>
          <button className="dashboard__refresh" onClick={refresh} aria-label="Refresh">
            Refresh
          </button>
        </div>

        <div className="dashboard__stats">
          <div className="stat-card" id="stat-total">
            <span className="stat-card__value">{checks.length}</span>
            <span className="stat-card__label">Sites Monitored</span>
          </div>
          <div className="stat-card stat-card--online" id="stat-online">
            <span className="stat-card__value">{onlineCount}</span>
            <span className="stat-card__label">Online</span>
          </div>
          <div className="stat-card stat-card--down" id="stat-down">
            <span className="stat-card__value">{checks.length - onlineCount}</span>
            <span className="stat-card__label">Issues</span>
          </div>
          <div className="stat-card" id="stat-latency">
            <span className="stat-card__value">{avgLatency}ms</span>
            <span className="stat-card__label">Avg Latency</span>
          </div>
        </div>
      </header>

      <main className="dashboard__grid">
        {loading && <p className="dashboard__message">Loading checks…</p>}
        {error && <p className="dashboard__message dashboard__error">Failed to load: {error}</p>}
        {!loading && !error && checks.length === 0 && (
          <p className="dashboard__message">
            No checks yet. Run <code>go run ./cmd/sentinel</code> to populate data.
          </p>
        )}
        {checks.map((check) => (
          <StatusCard
            key={check.id}
            check={check}
            onClick={() => setSelectedUrl(check.url)}
          />
        ))}
      </main>

      <footer className="dashboard__footer">
        <p>Auto-refreshes every 30 seconds</p>
      </footer>

      {selectedUrl && (
        <SiteDetail url={selectedUrl} onClose={() => setSelectedUrl(null)} />
      )}
    </div>
  );
}
