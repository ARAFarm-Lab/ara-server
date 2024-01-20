CREATE TABLE IF NOT EXISTS schedules(
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    actions JSONB,
    schedule VARCHAR(20) NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    next_run_at TIMESTAMPTZ NOT NULL,
    last_lock_at TIMESTAMPTZ,
    last_run_at TIMESTAMPTZ,
    last_run_status SMALLINT,
    last_error TEXT,
    attempts SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);