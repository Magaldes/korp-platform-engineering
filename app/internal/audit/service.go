package audit

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidInput = errors.New("invalid audit input")

var knownRoutes = map[string]bool{
	"/projeto-korp":       true,
	"/audit-events":       true,
	"/request-statistics": true,
}

type Service struct {
	repository Repository
	cache      StatisticsCache
	ttl        time.Duration
	observer   DependencyObserver
}

const (
	retentionDays    = 7 * 24 * time.Hour
	retentionCadence = time.Hour
)

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func NewStatisticsService(repository Repository, cache StatisticsCache, ttl time.Duration, observer DependencyObserver) *Service {
	return &Service{repository: repository, cache: cache, ttl: ttl, observer: observer}
}

// StartRetention starts the best-effort retention worker. The first attempt is
// launched asynchronously and subsequent attempts run once per hour.
func (s *Service) StartRetention(ctx context.Context) {
	if s == nil || s.repository == nil {
		return
	}
	go func() {
		s.runRetention(ctx)
		ticker := time.NewTicker(retentionCadence)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.runRetention(ctx)
			}
		}
	}()
}

func (s *Service) runRetention(ctx context.Context) {
	started := time.Now()
	cutoff := time.Now().UTC().Add(-retentionDays)
	_, err := s.repository.PurgeBefore(ctx, cutoff)
	result := "success"
	if err != nil {
		result = "error"
	}
	s.observeRepository("retention", result, time.Since(started))
}

func (s *Service) Record(ctx context.Context, event Event) {
	if s == nil || s.repository == nil {
		return
	}
	_ = s.repository.Append(ctx, event)
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]Event, bool, error) {
	return s.repository.List(ctx, filter)
}

func (s *Service) Statistics(ctx context.Context, filter StatisticsFilter) (Statistics, error) {
	statistics, _, _, err := s.statisticsWithCache(ctx, filter)
	return statistics, err
}

func (s *Service) StatisticsWithSource(ctx context.Context, filter StatisticsFilter) (Statistics, time.Time, string, error) {
	return s.statisticsWithCache(ctx, filter)
}

func (s *Service) statisticsWithCache(ctx context.Context, filter StatisticsFilter) (Statistics, time.Time, string, error) {
	if s == nil || s.repository == nil {
		return Statistics{}, time.Time{}, "", errors.New("statistics repository unavailable")
	}
	key := StatisticsCacheKey(filter)
	if s.cache != nil {
		started := time.Now()
		payload, err := s.cache.Get(ctx, key)
		if err == nil {
			var cached CachedStatistics
			if json.Unmarshal(payload, &cached) == nil && validCachedStatistics(cached) {
				s.observeCacheLookup("hit", time.Since(started))
				return cached.Statistics, cached.GeneratedAt, "cache", nil
			}
			s.observeCacheLookup("error", time.Since(started))
		} else {
			s.observeCacheLookup(cacheResult(err, payload), time.Since(started))
		}
	}

	started := time.Now()
	statistics, err := s.repository.Statistics(ctx, filter)
	result := "success"
	if err != nil {
		result = "error"
	}
	s.observeRepository("statistics", result, time.Since(started))
	if err != nil {
		return Statistics{}, time.Time{}, "", err
	}
	generatedAt := time.Now().UTC()
	if statistics.ByStatus == nil {
		statistics.ByStatus = []StatusCount{}
	}
	if statistics.ByRoute == nil {
		statistics.ByRoute = []RouteCount{}
	}
	if s.cache != nil && s.ttl > 0 {
		payload, marshalErr := json.Marshal(CachedStatistics{Statistics: statistics, GeneratedAt: generatedAt})
		if marshalErr == nil {
			started = time.Now()
			setErr := s.cache.Set(ctx, key, payload, s.ttl)
			writeResult := "success"
			if setErr != nil {
				writeResult = "error"
			}
			s.observeCacheWrite(writeResult, time.Since(started))
		}
	}
	return statistics, generatedAt, "repository", nil
}

func cacheResult(err error, payload []byte) string {
	if err == nil && len(payload) > 0 {
		return "hit"
	}
	if errors.Is(err, ErrCacheMiss) {
		return "miss"
	}
	return "error"
}

func validCachedStatistics(payload CachedStatistics) bool {
	return !payload.GeneratedAt.IsZero() && payload.Statistics.TotalRequests >= 0 && payload.Statistics.ErrorRequests >= 0 && payload.Statistics.ByStatus != nil && payload.Statistics.ByRoute != nil
}

func (s *Service) observeCacheLookup(result string, duration time.Duration) {
	if s.observer != nil {
		s.observer.ObserveCacheLookup(result, duration)
	}
}
func (s *Service) observeCacheWrite(result string, duration time.Duration) {
	if s.observer != nil {
		s.observer.ObserveCacheWrite(result, duration)
	}
}
func (s *Service) observeRepository(operation, result string, duration time.Duration) {
	if s.observer != nil {
		s.observer.ObserveRepository(operation, result, duration)
	}
}

func StatisticsCacheKey(filter StatisticsFilter) string {
	from := filter.From.UTC().Format(time.RFC3339)
	to := filter.To.UTC().Format(time.RFC3339)
	status := ""
	if filter.Status != nil {
		status = strconv.Itoa(*filter.Status)
	}
	return strings.Join([]string{"request-statistics:v1", from, to, strings.ToUpper(filter.Method), filter.Route, status}, ":")
}

func ParseListFilter(values map[string][]string) (ListFilter, error) {
	filter := ListFilter{Limit: 50}
	if err := rejectUnknown(values, map[string]bool{"from": true, "to": true, "method": true, "route": true, "status": true, "limit": true, "cursor": true}); err != nil {
		return filter, err
	}
	var err error
	filter.From, err = optionalTime(values, "from")
	if err != nil {
		return filter, err
	}
	filter.To, err = optionalTime(values, "to")
	if err != nil {
		return filter, err
	}
	if filter.From != nil && filter.To != nil && !filter.From.Before(*filter.To) {
		return filter, ErrInvalidInput
	}
	filter.Method = optionalString(values, "method")
	filter.Route = optionalString(values, "route")
	if err := validateFilters(filter.Method, filter.Route); err != nil {
		return filter, err
	}
	filter.Status, err = optionalStatus(values, "status")
	if err != nil {
		return filter, err
	}
	if raw := optionalString(values, "limit"); raw != "" {
		filter.Limit, err = strconv.Atoi(raw)
		if err != nil || filter.Limit < 1 || filter.Limit > 100 {
			return filter, ErrInvalidInput
		}
	}
	filter.Cursor = optionalString(values, "cursor")
	if filter.Cursor != "" {
		cursor, err := DecodeCursor(filter.Cursor)
		if err != nil {
			return filter, err
		}
		if !cursorMatches(filter, cursor) {
			return filter, ErrInvalidInput
		}
	}
	return filter, nil
}

func ParseStatisticsFilter(values map[string][]string) (StatisticsFilter, error) {
	if err := rejectUnknown(values, map[string]bool{"from": true, "to": true, "method": true, "route": true, "status": true}); err != nil {
		return StatisticsFilter{}, err
	}
	from, err := requiredTime(values, "from")
	if err != nil {
		return StatisticsFilter{}, err
	}
	to, err := requiredTime(values, "to")
	if err != nil {
		return StatisticsFilter{}, err
	}
	if !from.Before(to) || to.Sub(from) > 24*time.Hour {
		return StatisticsFilter{}, ErrInvalidInput
	}
	filter := StatisticsFilter{From: from, To: to, Method: optionalString(values, "method"), Route: optionalString(values, "route")}
	if err := validateFilters(filter.Method, filter.Route); err != nil {
		return StatisticsFilter{}, err
	}
	filter.Status, err = optionalStatus(values, "status")
	if err != nil {
		return StatisticsFilter{}, err
	}
	return filter, nil
}

func EncodeCursor(cursor Cursor) (string, error) {
	b, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func DecodeCursor(raw string) (Cursor, error) {
	b, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		return Cursor{}, ErrInvalidInput
	}
	var cursor Cursor
	if json.Unmarshal(b, &cursor) != nil || cursor.ID <= 0 || cursor.OccurredAt.IsZero() {
		return Cursor{}, ErrInvalidInput
	}
	return cursor, nil
}

func CursorFor(filter ListFilter, event Event) (string, error) {
	return EncodeCursor(Cursor{OccurredAt: event.OccurredAt, ID: event.ID, Method: filter.Method, Route: filter.Route, Status: filter.Status, From: filter.From, To: filter.To})
}

func cursorMatches(filter ListFilter, cursor Cursor) bool {
	return filter.Method == cursor.Method && filter.Route == cursor.Route && sameStatus(filter.Status, cursor.Status) && sameTime(filter.From, cursor.From) && sameTime(filter.To, cursor.To)
}

func sameStatus(a, b *int) bool { return (a == nil && b == nil) || (a != nil && b != nil && *a == *b) }
func sameTime(a, b *time.Time) bool {
	return (a == nil && b == nil) || (a != nil && b != nil && a.Equal(*b))
}

func CanonicalRoute(path string) string {
	if knownRoutes[path] {
		return path
	}
	return "unmatched"
}

func ValidateMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions, http.MethodTrace, http.MethodConnect:
		return true
	}
	return false
}

func validateFilters(method, route string) error {
	if method != "" && (!ValidateMethod(method) || method != strings.ToUpper(method)) {
		return ErrInvalidInput
	}
	if route != "" && route != "unmatched" && !knownRoutes[route] {
		return ErrInvalidInput
	}
	return nil
}

func rejectUnknown(values map[string][]string, allowed map[string]bool) error {
	for key := range values {
		if !allowed[key] || len(values[key]) != 1 {
			return ErrInvalidInput
		}
	}
	return nil
}
func optionalString(values map[string][]string, key string) string {
	if v := values[key]; len(v) == 1 {
		return v[0]
	}
	return ""
}
func optionalTime(values map[string][]string, key string) (*time.Time, error) {
	raw := optionalString(values, key)
	if raw == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, ErrInvalidInput
	}
	u := t.UTC()
	return &u, nil
}
func requiredTime(values map[string][]string, key string) (time.Time, error) {
	t, err := optionalTime(values, key)
	if err != nil || t == nil {
		return time.Time{}, ErrInvalidInput
	}
	return *t, nil
}
func optionalStatus(values map[string][]string, key string) (*int, error) {
	raw := optionalString(values, key)
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 100 || n > 599 {
		return nil, fmt.Errorf("%w: status", ErrInvalidInput)
	}
	return &n, nil
}
