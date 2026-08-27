CREATE TABLE IF NOT EXISTS api_keys (
    id                  BIGSERIAL PRIMARY KEY,
    key_hash            TEXT NOT NULL,
    name                TEXT NOT NULL,
    requests_per_minute INTEGER NOT NULL DEFAULT 60,
    monthly_quota       INTEGER NOT NULL DEFAULT 1000,
    active              BOOLEAN NOT NULL DEFAULT true,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT api_keys_key_hash_key UNIQUE (key_hash)
);

CREATE TABLE IF NOT EXISTS api_key_usage (
    id            BIGSERIAL PRIMARY KEY,
    api_key_id    BIGINT NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    period_start  DATE NOT NULL,
    request_count INTEGER NOT NULL DEFAULT 0,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT api_key_usage_key_period
        UNIQUE (api_key_id, period_start)
);

CREATE INDEX IF NOT EXISTS idx_api_key_usage_api_key_id
    ON api_key_usage (api_key_id);