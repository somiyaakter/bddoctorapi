package auth

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/time/rate"

	"datalab_api/internal/httpx"
)

type contextKey string

const apiKeyContextKey contextKey = "apiKey"

type limiterStore struct {
	mu       sync.Mutex
	limiters map[int64]*rate.Limiter
}

func newLimiterStore() *limiterStore {
	return &limiterStore{
		limiters: make(map[int64]*rate.Limiter),
	}
}

func (s *limiterStore) get(keyID int64, requestsPerMinute int) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.limiters[keyID]
	if !ok {
		perSecond := rate.Limit(float64(requestsPerMinute) / 60.0)

		l = rate.NewLimiter(
			perSecond,
			requestsPerMinute,
		)

		s.limiters[keyID] = l
	}

	return l
}

type Middleware struct {
	repo     *Repository
	limiters *limiterStore
}

func NewMiddleware(repo *Repository) *Middleware {
	return &Middleware{
		repo:     repo,
		limiters: newLimiterStore(),
	}
}

// Authenticate validates the API key, applies the per-minute
// rate limit, and applies the monthly quota for external keys.
//
// Internal keys:
//   - Must still provide a valid API key.
//   - Are still subject to the per-minute rate limit.
//   - Are exempt from the monthly quota.
//
// External keys:
//   - Must provide a valid API key.
//   - Are subject to the per-minute rate limit.
//   - Are subject to the monthly quota.
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ---------------------------------------------------------
		// 1. Read API key
		// ---------------------------------------------------------

		rawKey := r.Header.Get("X-API-Key")

		if rawKey == "" {
			httpx.WriteError(
				w,
				http.StatusUnauthorized,
				"missing X-API-Key header",
			)
			return
		}

		// ---------------------------------------------------------
		// 2. Validate API key
		// ---------------------------------------------------------

		key, err := m.repo.GetActiveByHash(
			r.Context(),
			HashKey(rawKey),
		)

		if errors.Is(err, pgx.ErrNoRows) {
			httpx.WriteError(
				w,
				http.StatusUnauthorized,
				"invalid or inactive API key",
			)
			return
		}

		if err != nil {
			httpx.WriteError(
				w,
				http.StatusInternalServerError,
				"failed to validate api key",
			)
			return
		}

		// ---------------------------------------------------------
		// 3. Per-minute rate limit
		//
		// IMPORTANT:
		// This applies to BOTH internal and external keys.
		// ---------------------------------------------------------

		limiter := m.limiters.get(
			key.ID,
			key.RequestsPerMinute,
		)

		if !limiter.Allow() {
			httpx.WriteError(
				w,
				http.StatusTooManyRequests,
				"rate limit exceeded, try again shortly",
			)
			return
		}

		
		// 4. Monthly quota
		//
		// Internal keys skip monthly quota.
		// External keys are checked atomically.
	

		if !key.IsInternal {

			allowed, used, periodStart, err :=
				m.repo.TryIncrementUsage(
					r.Context(),
					key.ID,
					key.MonthlyQuota,
				)

			if err != nil {
				httpx.WriteError(
					w,
					http.StatusInternalServerError,
					"failed to check quota",
				)
				return
			}

			if !allowed {
				writeQuotaExceeded(
					w,
					key.MonthlyQuota,
					used,
					periodStart,
				)
				return
			}
		}

		// 5. Store API key in request context

		ctx := context.WithValue(
			r.Context(),
			apiKeyContextKey,
			key,
		)

		// 6. Continue to actual handler
		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

// writeQuotaExceeded returns a consistent response when
// an external API key has reached its monthly quota.
func writeQuotaExceeded(
	w http.ResponseWriter,
	quota int,
	used int,
	periodStart string,
) {
	resetsAt := nextMonthStart().Format("2006-01-02")

	httpx.WriteJSON(
		w,
		http.StatusTooManyRequests,
		map[string]interface{}{
			"error":         "monthly quota exceeded",
			"monthly_quota": quota,
			"used":          used,
			"period_start":  periodStart,
			"resets_at":     resetsAt,
		},
	)
}

// nextMonthStart returns the first day of the next month
// in UTC.
func nextMonthStart() time.Time {
	now := time.Now().UTC()

	firstOfThisMonth := time.Date(
		now.Year(),
		now.Month(),
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	return firstOfThisMonth.AddDate(0, 1, 0)
}