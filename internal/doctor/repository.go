package doctor

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListParams struct {
	Page        int
	PageSize    int
	Query       string
	LocationID  *int64
	SpecialtyID *int64
}

// Repository handles persistence of Doctor/Chamber records.
type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// Upsert inserts a doctor or updates an existing doctor
// using profile_url as the unique key.
//
// Chambers are replaced on every scrape.
func (r *Repository) Upsert(ctx context.Context, d Doctor) (int64, error) {

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	var id int64

	err = tx.QueryRow(ctx, `
		INSERT INTO doctors (
			name,
			bmdc_reg_no,
			degrees,
			experience_years,
			specialties,
			designation,
			workplace,
			image_url,
			profile_url,
			updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, now()
		)

		ON CONFLICT (profile_url)
		DO UPDATE SET
			name             = EXCLUDED.name,
			bmdc_reg_no      = EXCLUDED.bmdc_reg_no,
			degrees          = EXCLUDED.degrees,
			experience_years = EXCLUDED.experience_years,
			specialties      = EXCLUDED.specialties,
			designation      = EXCLUDED.designation,
			workplace        = EXCLUDED.workplace,
			image_url        = EXCLUDED.image_url,
			updated_at       = now()

		RETURNING id
	`,
		nullIfEmpty(d.Name),
		nullIfEmpty(d.BMDCRegNo),
		nullIfEmpty(d.Degrees),
		d.ExperienceYears,
		nullIfEmpty(d.Specialties),
		nullIfEmpty(d.Designation),
		nullIfEmpty(d.Workplace),
		nullIfEmpty(d.ImageURL),
		d.ProfileURL,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf(
			"upsert doctor %s: %w",
			d.ProfileURL,
			err,
		)
	}

	// Remove old chambers.
	_, err = tx.Exec(
		ctx,
		`DELETE FROM chambers WHERE doctor_id = $1`,
		id,
	)

	if err != nil {
		return 0, fmt.Errorf(
			"clearing old chambers for doctor %d: %w",
			id,
			err,
		)
	}

	// Insert fresh chambers.
	for _, ch := range d.Chambers {

		_, err = tx.Exec(ctx, `
			INSERT INTO chambers (
				doctor_id,
				name,
				address,
				visiting_hour,
				appointment_phone
			)
			VALUES ($1, $2, $3, $4, $5)
		`,
			id,
			nullIfEmpty(ch.Name),
			nullIfEmpty(ch.Address),
			nullIfEmpty(ch.VisitingHour),
			nullIfEmpty(ch.AppointmentPhone),
		)

		if err != nil {
			return 0, fmt.Errorf(
				"inserting chamber for doctor %d: %w",
				id,
				err,
			)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf(
			"commit tx: %w",
			err,
		)
	}

	return id, nil
}

// GetByID fetches one doctor with all chambers.
func (r *Repository) GetByID(ctx context.Context, id int64) (Doctor, error) {

	var d Doctor

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			COALESCE(name, ''),
			COALESCE(bmdc_reg_no, ''),
			COALESCE(degrees, ''),
			COALESCE(experience_years, 0),
			COALESCE(specialties, ''),
			COALESCE(designation, ''),
			COALESCE(workplace, ''),
			COALESCE(image_url, ''),
			profile_url,
			created_at,
			updated_at

		FROM doctors
		WHERE id = $1
	`, id).Scan(
		&d.ID,
		&d.Name,
		&d.BMDCRegNo,
		&d.Degrees,
		&d.ExperienceYears,
		&d.Specialties,
		&d.Designation,
		&d.Workplace,
		&d.ImageURL,
		&d.ProfileURL,
		&d.CreatedAt,
		&d.UpdatedAt,
	)

	if err != nil {
		return Doctor{}, err
	}

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			doctor_id,
			COALESCE(name, ''),
			COALESCE(address, ''),
			COALESCE(visiting_hour, ''),
			COALESCE(appointment_phone, '')

		FROM chambers
		WHERE doctor_id = $1
		ORDER BY id
	`, id)

	if err != nil {
		return Doctor{}, fmt.Errorf(
			"fetching chambers for doctor %d: %w",
			id,
			err,
		)
	}

	defer rows.Close()

	for rows.Next() {

		var ch Chamber

		if err := rows.Scan(
			&ch.ID,
			&ch.DoctorID,
			&ch.Name,
			&ch.Address,
			&ch.VisitingHour,
			&ch.AppointmentPhone,
		); err != nil {
			return Doctor{}, fmt.Errorf(
				"scanning chamber: %w",
				err,
			)
		}

		d.Chambers = append(
			d.Chambers,
			ch,
		)
	}

	if err := rows.Err(); err != nil {
		return Doctor{}, err
	}

	return d, nil
}

// List returns paginated doctors with their chambers.

// List returns paginated doctors with relevance-based search and filters.
func (r *Repository) List(ctx context.Context, p ListParams) ([]Doctor, int, error) {

	if p.Page < 1 {
		p.Page = 1
	}

	if p.PageSize < 1 {
		p.PageSize = 20
	}

	if p.PageSize > 100 {
		p.PageSize = 100
	}

	offset := (p.Page - 1) * p.PageSize

	var joins []string
	var wheres []string
	var args []interface{}

	// ---------------------------------------------------------
	// SEARCH
	// ---------------------------------------------------------

	searchQuery := strings.TrimSpace(p.Query)

	var searchTokens []string

	if searchQuery != "" {

		// Normalize the query.
		searchQuery = strings.ToLower(searchQuery)

		// Split:
		// "dr md saifullah"
		// =>
		// ["dr", "md", "saifullah"]
		searchTokens = strings.Fields(searchQuery)

		// Every search token must exist somewhere in the
		// doctor's searchable information.
		//
		// This makes:
		//
		//   "dr md saifullah"
		//
		// behave much better than:
		//
		//   ILIKE "%dr md saifullah%"
		//
		// because the words don't have to appear as one
		// continuous string.

		for _, token := range searchTokens {

			args = append(args, "%"+token+"%")

			idx := len(args)

			wheres = append(wheres, fmt.Sprintf(`
				LOWER(
					CONCAT_WS(
						' ',
						COALESCE(d.name, ''),
						COALESCE(d.bmdc_reg_no, ''),
						COALESCE(d.degrees, ''),
						COALESCE(d.specialties, ''),
						COALESCE(d.designation, ''),
						COALESCE(d.workplace, '')
					)
				) LIKE $%d
			`, idx))
		}
	}

	// ---------------------------------------------------------
	// LOCATION FILTER
	// ---------------------------------------------------------

	if p.LocationID != nil {

		joins = append(
			joins,
			"JOIN doctor_locations dl ON dl.doctor_id = d.id",
		)

		args = append(args, *p.LocationID)

		wheres = append(
			wheres,
			fmt.Sprintf(
				"dl.location_id = $%d",
				len(args),
			),
		)
	}

	// ---------------------------------------------------------
	// SPECIALTY FILTER
	// ---------------------------------------------------------

	if p.SpecialtyID != nil {

		joins = append(
			joins,
			"JOIN doctor_specialties ds ON ds.doctor_id = d.id",
		)

		args = append(args, *p.SpecialtyID)

		wheres = append(
			wheres,
			fmt.Sprintf(
				"ds.specialty_id = $%d",
				len(args),
			),
		)
	}

	joinSQL := strings.Join(joins, " ")

	whereSQL := ""

	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}

	// ---------------------------------------------------------
	// COUNT
	// ---------------------------------------------------------

	var total int

	countSQL := fmt.Sprintf(
		`
		SELECT COUNT(DISTINCT d.id)
		FROM doctors d
		%s
		%s
		`,
		joinSQL,
		whereSQL,
	)

	if err := r.db.QueryRow(
		ctx,
		countSQL,
		args...,
	).Scan(&total); err != nil {

		return nil, 0, fmt.Errorf(
			"counting doctors: %w",
			err,
		)
	}

	// ---------------------------------------------------------
	// RELEVANCE ORDER
	// ---------------------------------------------------------

	// Default ordering.
	orderSQL := "d.id"

	if searchQuery != "" {

		// Add the complete normalized search query as an
		// additional SQL argument for ranking.
		args = append(args, searchQuery)

		queryArg := len(args)

		// Important:
		//
		// Exact name match gets the highest score.
		//
		// Example:
		//
		// Query:
		//   "dr md saifullah"
		//
		// "Dr. Md. Saifullah"
		// will rank above:
		// "Dr. Md. Saifullah Chowdhury"
		//
		// Then names starting with the query.
		//
		// Then names containing the query.
		//
		// Then general matches.

		orderSQL = fmt.Sprintf(`
			CASE

				-- Exact normalized name
				WHEN LOWER(
					REGEXP_REPLACE(
						COALESCE(d.name, ''),
						'[^a-zA-Z0-9]+',
						' ',
						'g'
					)
				) =
				REGEXP_REPLACE($%d, '[^a-zA-Z0-9]+', ' ', 'g')
				THEN 0

				-- Name starts with the query
				WHEN LOWER(COALESCE(d.name, '')) LIKE
					$%d || '%%'
				THEN 1

				-- Name contains the complete query
				WHEN LOWER(COALESCE(d.name, '')) LIKE
					'%%' || $%d || '%%'
				THEN 2

				-- Name contains all searched words
				ELSE 3

			END,
			LOWER(COALESCE(d.name, '')),
			d.id
		`, queryArg, queryArg, queryArg)
	}

	// ---------------------------------------------------------
	// FETCH DOCTORS
	// ---------------------------------------------------------

	listArgs := append(
		append([]interface{}{}, args...),
		p.PageSize,
		offset,
	)

	limitArg := len(args) + 1
	offsetArg := len(args) + 2

	listSQL := fmt.Sprintf(`
		SELECT
			d.id,
			COALESCE(d.name, ''),
			COALESCE(d.bmdc_reg_no, ''),
			COALESCE(d.degrees, ''),
			COALESCE(d.experience_years, 0),
			COALESCE(d.specialties, ''),
			COALESCE(d.designation, ''),
			COALESCE(d.workplace, ''),
			COALESCE(d.image_url, ''),
			d.profile_url,
			d.created_at,
			d.updated_at

		FROM doctors d

		%s

		%s

		ORDER BY %s

		LIMIT $%d
		OFFSET $%d
	`,
		joinSQL,
		whereSQL,
		orderSQL,
		limitArg,
		offsetArg,
	)

	rows, err := r.db.Query(
		ctx,
		listSQL,
		listArgs...,
	)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"listing doctors: %w",
			err,
		)
	}

	defer rows.Close()

	var doctors []Doctor
	var ids []int64

	for rows.Next() {

		var d Doctor

		err := rows.Scan(
			&d.ID,
			&d.Name,
			&d.BMDCRegNo,
			&d.Degrees,
			&d.ExperienceYears,
			&d.Specialties,
			&d.Designation,
			&d.Workplace,
			&d.ImageURL,
			&d.ProfileURL,
			&d.CreatedAt,
			&d.UpdatedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf(
				"scanning doctor: %w",
				err,
			)
		}

		doctors = append(doctors, d)
		ids = append(ids, d.ID)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if len(doctors) == 0 {
		return doctors, total, nil
	}

	// ---------------------------------------------------------
	// FETCH CHAMBERS
	// ---------------------------------------------------------

	chamberRows, err := r.db.Query(ctx, `
		SELECT
			id,
			doctor_id,
			COALESCE(name, ''),
			COALESCE(address, ''),
			COALESCE(visiting_hour, ''),
			COALESCE(appointment_phone, '')

		FROM chambers

		WHERE doctor_id = ANY($1)

		ORDER BY doctor_id, id
	`, ids)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"batch-fetching chambers: %w",
			err,
		)
	}

	defer chamberRows.Close()

	byDoctorID := make(map[int64][]Chamber)

	for chamberRows.Next() {

		var ch Chamber

		err := chamberRows.Scan(
			&ch.ID,
			&ch.DoctorID,
			&ch.Name,
			&ch.Address,
			&ch.VisitingHour,
			&ch.AppointmentPhone,
		)

		if err != nil {
			return nil, 0, fmt.Errorf(
				"scanning chamber: %w",
				err,
			)
		}

		byDoctorID[ch.DoctorID] =
			append(
				byDoctorID[ch.DoctorID],
				ch,
			)
	}

	if err := chamberRows.Err(); err != nil {
		return nil, 0, err
	}

	for i := range doctors {

		doctors[i].Chambers =
			byDoctorID[doctors[i].ID]

		if doctors[i].Chambers == nil {
			doctors[i].Chambers = []Chamber{}
		}
	}

	return doctors, total, nil
}

// nullIfEmpty converts empty strings to SQL NULL.
func nullIfEmpty(s string) interface{} {

	if s == "" {
		return nil
	}

	return s
}

// Keep pgx imported because GetByID returns pgx.ErrNoRows.
var _ = pgx.ErrNoRows
