package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// Create creates a new API key.
//
// keyHash is the SHA-256 hash of the plaintext API key.
// The plaintext key is never stored in the database.
//
// isInternal determines whether this is a first-party/internal key.
// Internal keys still have a per-minute rate limit but are exempt
// from the monthly quota.
func (r *Repository) Create(
	ctx context.Context,
	name string,
	keyHash string,
	requestsPerMinute int,
	monthlyQuota int,
	isInternal bool,
) (APIKey, error) {

	var k APIKey

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO api_keys (
			key_hash,
			name,
			requests_per_minute,
			monthly_quota,
			is_internal
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			key_hash,
			name,
			requests_per_minute,
			monthly_quota,
			is_internal,
			active,
			created_at
		`,
		keyHash,
		name,
		requestsPerMinute,
		monthlyQuota,
		isInternal,
	).Scan(
		&k.ID,
		&k.KeyHash,
		&k.Name,
		&k.RequestsPerMinute,
		&k.MonthlyQuota,
		&k.IsInternal,
		&k.Active,
		&k.CreatedAt,
	)

	if err != nil {
		return APIKey{}, fmt.Errorf(
			"creating api key: %w",
			err,
		)
	}

	return k, nil
}

// GetActiveByHash finds an active API key by its SHA-256 hash.
func (r *Repository) GetActiveByHash(
	ctx context.Context,
	keyHash string,
) (APIKey, error) {

	var k APIKey

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			key_hash,
			name,
			requests_per_minute,
			monthly_quota,
			is_internal,
			active,
			created_at
		FROM api_keys
		WHERE key_hash = $1
		  AND active = true
		`,
		keyHash,
	).Scan(
		&k.ID,
		&k.KeyHash,
		&k.Name,
		&k.RequestsPerMinute,
		&k.MonthlyQuota,
		&k.IsInternal,
		&k.Active,
		&k.CreatedAt,
	)

	if err != nil {
		return APIKey{}, err
	}

	return k, nil
}

// TryIncrementUsage atomically increments the current month's
// usage only if the API key has not reached its monthly quota.
//
// Example:
//
//	quota = 6
//
//	request 1 -> allowed, usage = 1
//	request 2 -> allowed, usage = 2
//	request 3 -> allowed, usage = 3
//	request 4 -> allowed, usage = 4
//	request 5 -> allowed, usage = 5
//	request 6 -> allowed, usage = 6
//	request 7 -> rejected, usage stays 6
//
// Rejected requests do not increase the usage count.
func (r *Repository) TryIncrementUsage(
	ctx context.Context,
	apiKeyID int64,
	monthlyQuota int,
) (
	allowed bool,
	used int,
	periodStart string,
	err error,
) {

	var count int

	err = r.db.QueryRow(
		ctx,
		`
		INSERT INTO api_key_usage (
			api_key_id,
			period_start,
			request_count,
			updated_at
		)
		VALUES (
			$1,
			date_trunc('month', now())::date,
			1,
			now()
		)
		ON CONFLICT (api_key_id, period_start)
		DO UPDATE SET
			request_count =
				api_key_usage.request_count + 1,
			updated_at = now()
		WHERE api_key_usage.request_count < $2
		RETURNING
			request_count,
			period_start::text
		`,
		apiKeyID,
		monthlyQuota,
	).Scan(
		&count,
		&periodStart,
	)

	if err == nil {
		return true, count, periodStart, nil
	}

	// No row returned means the current usage has already
	// reached the monthly quota.
	var currentCount int

	err = r.db.QueryRow(
		ctx,
		`
		SELECT
			request_count
		FROM api_key_usage
		WHERE api_key_id = $1
		  AND period_start = date_trunc('month', now())::date
		`,
		apiKeyID,
	).Scan(&currentCount)

	if err != nil {
		return false, 0, "", fmt.Errorf(
			"checking current usage for key %d: %w",
			apiKeyID,
			err,
		)
	}

	var currentPeriodStart string

	err = r.db.QueryRow(
		ctx,
		`
		SELECT
			period_start::text
		FROM api_key_usage
		WHERE api_key_id = $1
		  AND period_start = date_trunc('month', now())::date
		`,
		apiKeyID,
	).Scan(&currentPeriodStart)

	if err != nil {
		return false, currentCount, "", fmt.Errorf(
			"getting period start for key %d: %w",
			apiKeyID,
			err,
		)
	}

	return false, currentCount, currentPeriodStart, nil
}

// GetCurrentUsage returns this month's usage without incrementing it.
//
// If no usage row exists yet this month, usage is 0.
func (r *Repository) GetCurrentUsage(
	ctx context.Context,
	apiKeyID int64,
) (int, error) {

	var count int

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			request_count
		FROM api_key_usage
		WHERE api_key_id = $1
		  AND period_start = date_trunc('month', now())::date
		`,
		apiKeyID,
	).Scan(&count)

	if err != nil {
		// No row yet this month means zero usage.
		return 0, nil
	}

	return count, nil
}