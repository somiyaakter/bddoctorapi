package auth

import "time"

type APIKey struct {
	ID                int64
	KeyHash           string
	Name              string
	RequestsPerMinute int
	MonthlyQuota      int
	IsInternal        bool
	Active            bool
	CreatedAt         time.Time
}