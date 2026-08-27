CREATE TABLE IF NOT EXISTS doctor_locations (
    doctor_id   BIGINT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    location_id BIGINT NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    PRIMARY KEY (doctor_id, location_id)
);
CREATE INDEX IF NOT EXISTS idx_doctor_locations_location_id ON doctor_locations (location_id);

CREATE TABLE IF NOT EXISTS doctor_specialties (
    doctor_id    BIGINT NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    specialty_id BIGINT NOT NULL REFERENCES specialties(id) ON DELETE CASCADE,
    PRIMARY KEY (doctor_id, specialty_id)
);
CREATE INDEX IF NOT EXISTS idx_doctor_specialties_specialty_id ON doctor_specialties (specialty_id);