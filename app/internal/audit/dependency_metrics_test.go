package audit

import (
	"context"
	"errors"
	"testing"
	"time"
)

type dependencyMetricsRepository struct{ err error }

func (r *dependencyMetricsRepository) Append(context.Context, Event) error { return r.err }
func (r *dependencyMetricsRepository) List(context.Context, ListFilter) ([]Event, bool, error) {
	return []Event{}, false, r.err
}
func (r *dependencyMetricsRepository) Statistics(context.Context, StatisticsFilter) (Statistics, error) {
	return Statistics{ByStatus: []StatusCount{}, ByRoute: []RouteCount{}}, r.err
}
func (r *dependencyMetricsRepository) PurgeBefore(context.Context, time.Time) (int64, error) {
	return 0, r.err
}

type dependencyMetricsObserver struct{ observations []string }

func (o *dependencyMetricsObserver) ObserveCacheLookup(string, time.Duration) {}
func (o *dependencyMetricsObserver) ObserveCacheWrite(string, time.Duration)  {}
func (o *dependencyMetricsObserver) ObserveRepository(operation, result string, duration time.Duration) {
	if duration < 0 {
		panic("negative duration")
	}
	o.observations = append(o.observations, operation+":"+result)
}

func TestAppendObservesSuccessAndErrorWithoutPropagatingError(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{name: "success", want: "append:success"},
		{name: "error", err: errors.New("postgres unavailable"), want: "append:error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			observer := &dependencyMetricsObserver{}
			service := NewStatisticsService(&dependencyMetricsRepository{err: tc.err}, nil, 0, observer)
			service.Record(context.Background(), Event{})
			if len(observer.observations) != 1 || observer.observations[0] != tc.want {
				t.Fatalf("observations = %v, want %q", observer.observations, tc.want)
			}
		})
	}
}

func TestListObservesSuccessAndErrorAndPreservesError(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want string
	}{
		{name: "success", want: "list:success"},
		{name: "error", err: errors.New("postgres unavailable"), want: "list:error"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			observer := &dependencyMetricsObserver{}
			service := NewStatisticsService(&dependencyMetricsRepository{err: tc.err}, nil, 0, observer)
			_, _, err := service.List(context.Background(), ListFilter{Limit: 1})
			if !errors.Is(err, tc.err) {
				t.Fatalf("error = %v, want %v", err, tc.err)
			}
			if len(observer.observations) != 1 || observer.observations[0] != tc.want {
				t.Fatalf("observations = %v, want %q", observer.observations, tc.want)
			}
		})
	}
}
