package audit

import (
	"context"
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("statistics cache miss")

// StatisticsCache is the application boundary for the optional statistics cache.
// Implementations must translate backend-specific misses to ErrCacheMiss.
type StatisticsCache interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte, time.Duration) error
}

// DependencyObserver receives bounded application-boundary dependency events.
type DependencyObserver interface {
	ObserveCacheLookup(result string, duration time.Duration)
	ObserveCacheWrite(result string, duration time.Duration)
	ObserveRepository(operation, result string, duration time.Duration)
}

// Event is the privacy-minimized Request Audit Event persisted by the service.
type Event struct {
	ID             int64     `json:"id"`
	OccurredAt     time.Time `json:"occurred_at"`
	Method         string    `json:"method"`
	CanonicalRoute string    `json:"canonical_route"`
	HTTPStatus     int       `json:"http_status"`
	DurationMS     int64     `json:"duration_ms"`
}

type ListFilter struct {
	From   *time.Time
	To     *time.Time
	Method string
	Route  string
	Status *int
	Limit  int
	Cursor string
}

type Cursor struct {
	OccurredAt time.Time  `json:"occurred_at"`
	ID         int64      `json:"id"`
	Method     string     `json:"method,omitempty"`
	Route      string     `json:"route,omitempty"`
	Status     *int       `json:"status,omitempty"`
	From       *time.Time `json:"from,omitempty"`
	To         *time.Time `json:"to,omitempty"`
}

type StatisticsFilter struct {
	From   time.Time
	To     time.Time
	Method string
	Route  string
	Status *int
}

type StatusCount struct {
	Status int   `json:"status"`
	Count  int64 `json:"count"`
}

type RouteCount struct {
	Route string `json:"route"`
	Count int64  `json:"count"`
}

type Statistics struct {
	TotalRequests     int64         `json:"total_requests"`
	ByStatus          []StatusCount `json:"by_status"`
	ByRoute           []RouteCount  `json:"by_route"`
	ErrorRequests     int64         `json:"error_requests"`
	AverageDurationMS *float64      `json:"average_duration_ms"`
}

type CachedStatistics struct {
	Statistics  Statistics `json:"statistics"`
	GeneratedAt time.Time  `json:"generated_at"`
}

// Repository is the application boundary for audit persistence and queries.
type Repository interface {
	Append(ctx context.Context, event Event) error
	List(ctx context.Context, filter ListFilter) ([]Event, bool, error)
	Statistics(ctx context.Context, filter StatisticsFilter) (Statistics, error)
	PurgeBefore(ctx context.Context, cutoff time.Time) (int64, error)
}
