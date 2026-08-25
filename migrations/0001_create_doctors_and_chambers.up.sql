CREATE TABLE IF NOT EXISTS doctors (
    id                BIGSERIAL PRIMARY KEY,
    name              TEXT NOT NULL,
    bmdc_reg_no       TEXT,
    degrees           TEXT,
    experience_years  INTEGER,
    specialties       TEXT,
    designation       TEXT,
    workplace         TEXT,
    image_url         TEXT,
    profile_url       TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT doctors_profile_url_key UNIQUE (profile_url)
);

CREATE INDEX IF NOT EXISTS idx_doctors_name ON doctors (name);

-- BMDC reg no is optional (not every profile shows one, per Phase 4 findings).
-- A partial unique index enforces uniqueness only when it's actually present,
-- so multiple doctors with no BMDC number don't collide.
CREATE UNIQUE INDEX IF NOT EXISTS idx_doctors_bmdc_reg_no_unique
    ON doctors (bmdc_reg_no)
    WHERE bmdc_reg_no IS NOT NULL AND bmdc_reg_no <> '';

CREATE TABLE IF NOT EXISTS chambers (
    id                 BIGSERIAL PRIMARY KEY,
    doctor_id          BIGINT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    name               TEXT,
    address            TEXT,
    visiting_hour      TEXT,
    appointment_phone  TEXT
);

CREATE INDEX IF NOT EXISTS idx_chambers_doctor_id ON chambers (doctor_id);