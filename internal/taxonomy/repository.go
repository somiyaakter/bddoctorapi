package taxonomy

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// UpsertLocation inserts a location, or updates its name if the URL
// (the unique key) already exists. Returns the location's id.
func (r *Repository) UpsertLocation(ctx context.Context, name, url string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO locations (name, url)
		VALUES ($1, $2)
		ON CONFLICT (url) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, name, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert location %s: %w", url, err)
	}
	return id, nil
}

// UpsertSpecialty inserts a specialty scoped to locationID, or updates
// its name if the URL already exists. Returns the specialty's id.
func (r *Repository) UpsertSpecialty(ctx context.Context, locationID int64, name, url string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
		INSERT INTO specialties (location_id, name, url)
		VALUES ($1, $2, $3)
		ON CONFLICT (url) DO UPDATE SET name = EXCLUDED.name, location_id = EXCLUDED.location_id
		RETURNING id
	`, locationID, name, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("upsert specialty %s: %w", url, err)
	}
	return id, nil
}

func (r *Repository) ListLocations(ctx context.Context) ([]Location, error) {
	rows, err := r.db.Query(ctx, `SELECT id, name, url FROM locations ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("listing locations: %w", err)
	}
	defer rows.Close()

	var locations []Location
	for rows.Next() {
		var l Location
		if err := rows.Scan(&l.ID, &l.Name, &l.URL); err != nil {
			return nil, fmt.Errorf("scanning location: %w", err)
		}
		locations = append(locations, l)
	}
	return locations, rows.Err()
}

// ListSpecialties returns all specialties, or only those for locationID
// if it's non-nil (used for GET /api/v1/specialties?location_id=).
func (r *Repository) ListSpecialties(ctx context.Context, locationID *int64) ([]Specialty, error) {
	var rows interface {
		Next() bool
		Scan(dest ...any) error
		Close()
		Err() error
	}
	var err error

	if locationID != nil {
		rows, err = r.db.Query(ctx, `
			SELECT id, location_id, name, url FROM specialties
			WHERE location_id = $1 ORDER BY name
		`, *locationID)
	} else {
		rows, err = r.db.Query(ctx, `
			SELECT id, location_id, name, url FROM specialties ORDER BY name
		`)
	}
	if err != nil {
		return nil, fmt.Errorf("listing specialties: %w", err)
	}
	defer rows.Close()

	var specialties []Specialty
	for rows.Next() {
		var s Specialty
		if err := rows.Scan(&s.ID, &s.LocationID, &s.Name, &s.URL); err != nil {
			return nil, fmt.Errorf("scanning specialty: %w", err)
		}
		specialties = append(specialties, s)
	}
	return specialties, rows.Err()
}
func (r *Repository) LinkDoctorLocation(ctx context.Context, doctorID, locationID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO doctor_locations (doctor_id, location_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, doctorID, locationID)
	if err != nil {
		return fmt.Errorf("linking doctor %d to location %d: %w", doctorID, locationID, err)
	}
	return nil
}

func (r *Repository) LinkDoctorSpecialty(ctx context.Context, doctorID, specialtyID int64) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO doctor_specialties (doctor_id, specialty_id)
		VALUES ($1, $2) ON CONFLICT DO NOTHING
	`, doctorID, specialtyID)
	if err != nil {
		return fmt.Errorf("linking doctor %d to specialty %d: %w", doctorID, specialtyID, err)
	}
	return nil
}