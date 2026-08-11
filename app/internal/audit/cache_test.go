package audit

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type cacheTestRepository struct {
	statistics Statistics
	err        error
	calls      int
}

func (r *cacheTestRepository) Append(context.Context, Event) error { return nil }
func (r *cacheTestRepository) List(context.Context, ListFilter) ([]Event, bool, error) {
	return nil, false, nil
}
func (r *cacheTestRepository) Statistics(context.Context, StatisticsFilter) (Statistics, error) {
	r.calls++
	return r.statistics, r.err
}
func (r *cacheTestRepository) PurgeBefore(context.Context, time.Time) (int64, error) { return 0, nil }

type cacheTestStore struct {
	payload []byte
	err     error
	ttl     time.Duration
}

func (c *cacheTestStore) Get(context.Context, string) ([]byte, error) { return c.payload, c.err }
func (c *cacheTestStore) Set(_ context.Context, _ string, payload []byte, ttl time.Duration) error {
	c.payload, c.ttl = payload, ttl
	return c.err
}

func testFilter() StatisticsFilter {
	return StatisticsFilter{From: time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC), To: time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC), Method: "get", Route: "/projeto-korp"}
}

func TestStatisticsCacheKeyNormalizesEquivalentFilters(t *testing.T) {
	left := testFilter()
	right := left
	right.Method = "GET"
	if StatisticsCacheKey(left) != StatisticsCacheKey(right) {
		t.Fatalf("equivalent methods produced different keys")
	}
	if got := StatisticsCacheKey(StatisticsFilter{From: left.From, To: left.To}); got != "request-statistics:v1:2026-08-11T00:00:00Z:2026-08-11T01:00:00Z:::" {
		t.Fatalf("key = %q", got)
	}
}

func TestStatisticsCacheHitPreservesGeneratedAt(t *testing.T) {
	generatedAt := time.Date(2026, 8, 11, 0, 30, 0, 0, time.UTC)
	payload, _ := json.Marshal(CachedStatistics{Statistics: Statistics{ByStatus: []StatusCount{}, ByRoute: []RouteCount{}, TotalRequests: 4}, GeneratedAt: generatedAt})
	store := &cacheTestStore{payload: payload}
	repository := &cacheTestRepository{err: errors.New("postgres unavailable")}
	service := NewStatisticsService(repository, store, time.Minute, nil)
	statistics, gotGeneratedAt, source, err := service.StatisticsWithSource(context.Background(), testFilter())
	if err != nil || source != "cache" || statistics.TotalRequests != 4 || !gotGeneratedAt.Equal(generatedAt) {
		t.Fatalf("result = %+v %s %v %v", statistics, source, gotGeneratedAt, err)
	}
	if repository.calls != 0 {
		t.Fatal("cache hit queried PostgreSQL")
	}
}

func TestStatisticsCacheMissFallsBackAndWritesTTL(t *testing.T) {
	store := &cacheTestStore{err: ErrCacheMiss}
	repository := &cacheTestRepository{statistics: Statistics{ByStatus: []StatusCount{}, ByRoute: []RouteCount{}, TotalRequests: 2}}
	service := NewStatisticsService(repository, store, 60*time.Second, nil)
	statistics, _, source, err := service.StatisticsWithSource(context.Background(), testFilter())
	if err != nil || source != "repository" || statistics.TotalRequests != 2 {
		t.Fatalf("result = %+v %s %v", statistics, source, err)
	}
	if store.ttl <= 0 || store.ttl > 60*time.Second || len(store.payload) == 0 {
		t.Fatalf("cache write ttl=%s payload=%d", store.ttl, len(store.payload))
	}
}

func TestStatisticsCacheInvalidPayloadFallsBack(t *testing.T) {
	store := &cacheTestStore{payload: []byte(`{"generated_at":"2026-08-11T00:00:00Z","statistics":{"total_requests":1}}`)}
	repository := &cacheTestRepository{statistics: Statistics{ByStatus: []StatusCount{}, ByRoute: []RouteCount{}, TotalRequests: 3}}
	service := NewStatisticsService(repository, store, time.Minute, nil)
	statistics, _, source, err := service.StatisticsWithSource(context.Background(), testFilter())
	if err != nil || source != "repository" || statistics.TotalRequests != 3 {
		t.Fatalf("result = %+v %s %v", statistics, source, err)
	}
}
