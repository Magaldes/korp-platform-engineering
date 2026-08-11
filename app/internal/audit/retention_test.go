package audit

import (
	"context"
	"errors"
	"testing"
	"time"
)

type retentionRepository struct {
	cutoff time.Time
	err    error
}

func (r *retentionRepository) Append(context.Context, Event) error { return nil }
func (r *retentionRepository) List(context.Context, ListFilter) ([]Event, bool, error) {
	return nil, false, nil
}
func (r *retentionRepository) Statistics(context.Context, StatisticsFilter) (Statistics, error) {
	return Statistics{}, nil
}
func (r *retentionRepository) PurgeBefore(_ context.Context, cutoff time.Time) (int64, error) {
	r.cutoff = cutoff
	return 1, r.err
}

type retentionObserver struct {
	operation string
	result    string
}

func (o *retentionObserver) ObserveCacheLookup(string, time.Duration) {}
func (o *retentionObserver) ObserveCacheWrite(string, time.Duration)  {}
func (o *retentionObserver) ObserveRepository(operation, result string, _ time.Duration) {
	o.operation, o.result = operation, result
}

func TestRetentionUsesUTCSevenDayCutoffAndRecordsSuccess(t *testing.T) {
	repository := &retentionRepository{}
	observer := &retentionObserver{}
	service := NewStatisticsService(repository, nil, 0, observer)
	before := time.Now().UTC().Add(-7 * 24 * time.Hour)
	service.runRetention(context.Background())
	after := time.Now().UTC().Add(-7 * 24 * time.Hour)

	if repository.cutoff.Before(before) || repository.cutoff.After(after) {
		t.Fatalf("cutoff = %s, outside expected interval [%s, %s]", repository.cutoff, before, after)
	}
	if observer.operation != "retention" || observer.result != "success" {
		t.Fatalf("observation = %s/%s", observer.operation, observer.result)
	}
}

func TestRetentionFailureIsBestEffort(t *testing.T) {
	repository := &retentionRepository{err: errors.New("postgres unavailable")}
	observer := &retentionObserver{}
	service := NewStatisticsService(repository, nil, 0, observer)
	service.runRetention(context.Background())

	if observer.operation != "retention" || observer.result != "error" {
		t.Fatalf("failure observation = %s/%s", observer.operation, observer.result)
	}
	if repository.cutoff.IsZero() {
		t.Fatal("retention did not attempt a purge")
	}
}
