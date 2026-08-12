package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/projeto-korp/app/internal/audit"
)

type fakeRepository struct {
	events []audit.Event
	err    error
}

type statisticsCache struct{ err error }

func (c statisticsCache) Get(context.Context, string) ([]byte, error) { return nil, c.err }
func (c statisticsCache) Set(context.Context, string, []byte, time.Duration) error {
	return nil
}

func (f *fakeRepository) Append(_ context.Context, event audit.Event) error {
	if f.err != nil {
		return f.err
	}
	event.ID = int64(len(f.events) + 1)
	f.events = append(f.events, event)
	return nil
}
func (f *fakeRepository) List(_ context.Context, _ audit.ListFilter) ([]audit.Event, bool, error) {
	return f.events, false, f.err
}
func (f *fakeRepository) Statistics(_ context.Context, _ audit.StatisticsFilter) (audit.Statistics, error) {
	return audit.Statistics{ByStatus: []audit.StatusCount{}, ByRoute: []audit.RouteCount{}}, f.err
}
func (f *fakeRepository) PurgeBefore(context.Context, time.Time) (int64, error) { return 0, f.err }

func TestAuditFailureDoesNotChangeOriginalResponse(t *testing.T) {
	repository := &fakeRepository{err: errors.New("database unavailable")}
	server := httptest.NewServer(NewHandler(audit.NewService(repository)))
	defer server.Close()
	response, err := http.Get(server.URL + "/projeto-korp")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
	if len(repository.events) != 0 {
		t.Fatal("failed audit was persisted")
	}
}

func TestAuditMiddlewareRecordsKnownAndUnknownRoutes(t *testing.T) {
	repository := &fakeRepository{}
	server := httptest.NewServer(NewHandler(audit.NewService(repository)))
	defer server.Close()
	for _, path := range []string{"/projeto-korp", "/unknown"} {
		response, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
	if len(repository.events) != 2 || repository.events[0].CanonicalRoute != "/projeto-korp" || repository.events[1].CanonicalRoute != "unmatched" {
		t.Fatalf("events = %+v", repository.events)
	}
	if repository.events[0].DurationMS < 0 || repository.events[0].OccurredAt.Before(time.Now().Add(-time.Minute)) {
		t.Fatalf("invalid event timing: %+v", repository.events[0])
	}
}

func TestAuditMiddlewareRecordsHead405404AndExcludesMetrics(t *testing.T) {
	repository := &fakeRepository{}
	server := httptest.NewServer(NewHandler(audit.NewService(repository)))
	defer server.Close()
	requests := []struct{ method, path string }{{http.MethodHead, "/projeto-korp"}, {http.MethodPost, "/projeto-korp"}, {http.MethodGet, "/missing"}, {http.MethodGet, "/metrics"}}
	for _, item := range requests {
		request, err := http.NewRequest(item.method, server.URL+item.path, nil)
		if err != nil {
			t.Fatal(err)
		}
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
	}
	if len(repository.events) != 3 {
		t.Fatalf("event count = %d, want 3", len(repository.events))
	}
	if repository.events[0].HTTPStatus != http.StatusOK || repository.events[1].HTTPStatus != http.StatusMethodNotAllowed || repository.events[2].HTTPStatus != http.StatusNotFound {
		t.Fatalf("statuses = %+v", repository.events)
	}
}

func TestAuditAndStatisticsSuccessResponses(t *testing.T) {
	repository := &fakeRepository{events: []audit.Event{{ID: 1, OccurredAt: time.Now().UTC(), Method: "GET", CanonicalRoute: "/projeto-korp", HTTPStatus: 200, DurationMS: 3}}}
	server := httptest.NewServer(NewHandler(audit.NewService(repository)))
	defer server.Close()
	response, err := http.Get(server.URL + "/audit-events")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("audit status = %d", response.StatusCode)
	}
	response, err = http.Get(server.URL + "/request-statistics?from=2026-08-11T00:00:00Z&to=2026-08-11T01:00:00Z")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("statistics status = %d", response.StatusCode)
	}
}

func TestAuditEndpointsMapValidationAndRepositoryErrors(t *testing.T) {
	for _, tc := range []struct {
		path string
		want int
	}{
		{"/audit-events?limit=0", http.StatusBadRequest},
		{"/request-statistics", http.StatusBadRequest},
		{"/audit-events", http.StatusServiceUnavailable},
	} {
		server := httptest.NewServer(NewHandler(audit.NewService(&fakeRepository{err: errors.New("down")})))
		response, err := http.Get(server.URL + tc.path)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		server.Close()
		if response.StatusCode != tc.want {
			t.Errorf("%s status = %d, want %d", tc.path, response.StatusCode, tc.want)
		}
	}
}

func TestAuditEventsReturnsUnavailableWithoutRepository(t *testing.T) {
	server := httptest.NewServer(NewHandler(audit.NewStatisticsService(nil, nil, 0, nil)))
	defer server.Close()
	response, err := http.Get(server.URL + "/audit-events")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusServiceUnavailable)
	}
}

func TestProjectEndpointWorksWithoutRepository(t *testing.T) {
	server := httptest.NewServer(NewHandler(audit.NewStatisticsService(nil, nil, 0, nil)))
	defer server.Close()
	response, err := http.Get(server.URL + "/projeto-korp")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
}

func TestStatisticsReturnsUnavailableWithoutRepositoryAfterCacheFailure(t *testing.T) {
	for _, cacheErr := range []error{audit.ErrCacheMiss, errors.New("redis unavailable")} {
		server := httptest.NewServer(NewHandler(audit.NewStatisticsService(nil, statisticsCache{err: cacheErr}, time.Minute, nil)))
		response, err := http.Get(server.URL + "/request-statistics?from=2026-08-11T00:00:00Z&to=2026-08-11T01:00:00Z")
		if err != nil {
			server.Close()
			t.Fatal(err)
		}
		response.Body.Close()
		server.Close()
		if response.StatusCode != http.StatusServiceUnavailable {
			t.Fatalf("cache error %v status = %d, want %d", cacheErr, response.StatusCode, http.StatusServiceUnavailable)
		}
	}
}
