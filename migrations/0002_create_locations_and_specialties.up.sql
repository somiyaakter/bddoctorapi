CREATE TABLE IF NOT EXISTS locations (
    id   BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url  TEXT NOT NULL,
    CONSTRAINT locations_url_key UNIQUE (url)
);

CREATE TABLE IF NOT EXISTS specialties (
    id          BIGSERIAL PRIMARY KEY,
    location_id BIGINT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    url         TEXT NOT NULL,
    CONSTRAINT specialties_url_key UNIQUE (url)
);

CREATE INDEX IF NOT EXISTS idx_specialties_location_id ON specialties (location_id);