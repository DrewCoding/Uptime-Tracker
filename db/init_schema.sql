CREATE TABLE IF NOT EXISTS health_checks (
    id         SERIAL PRIMARY KEY,
    url        TEXT NOT NULL,
    status_code INT,
    latency_ms BIGINT NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_health_checks_url ON health_checks (url);

CREATE INDEX IF NOT EXISTS idx_health_checks_checked_at ON health_checks (checked_at DESC);