package audit

import (
	"net/url"
	"testing"
	"time"
)

func TestCanonicalRouteDoesNotExposeUnknownPath(t *testing.T) {
	if got := CanonicalRoute("/client/secret"); got != "unmatched" {
		t.Fatalf("route = %q", got)
	}
	if got := CanonicalRoute("/projeto-korp"); got != "/projeto-korp" {
		t.Fatalf("known route = %q", got)
	}
}

func TestListFilterAndCursorValidation(t *testing.T) {
	values := url.Values{"from": {"2026-08-11T00:00:00Z"}, "to": {"2026-08-11T01:00:00Z"}, "method": {"GET"}, "route": {"unmatched"}, "status": {"404"}, "limit": {"100"}}
	filter, err := ParseListFilter(values)
	if err != nil {
		t.Fatal(err)
	}
	event := Event{ID: 9, OccurredAt: time.Date(2026, 8, 11, 0, 30, 0, 0, time.UTC)}
	cursor, err := CursorFor(filter, event)
	if err != nil {
		t.Fatal(err)
	}
	values.Set("cursor", cursor)
	if _, err := ParseListFilter(values); err != nil {
		t.Fatalf("matching cursor rejected: %v", err)
	}
	values.Set("route", "/projeto-korp")
	if _, err := ParseListFilter(values); err == nil {
		t.Fatal("cursor bound to old filters was accepted")
	}
}

func TestStatisticsValidation(t *testing.T) {
	valid := url.Values{"from": {"2026-08-11T00:00:00Z"}, "to": {"2026-08-11T23:59:59Z"}}
	if _, err := ParseStatisticsFilter(valid); err != nil {
		t.Fatal(err)
	}
	tooWide := url.Values{"from": {"2026-08-11T00:00:00Z"}, "to": {"2026-08-12T00:00:01Z"}}
	if _, err := ParseStatisticsFilter(tooWide); err == nil {
		t.Fatal("more than 24 hours accepted")
	}
	missing := url.Values{"from": {"2026-08-11T00:00:00Z"}}
	if _, err := ParseStatisticsFilter(missing); err == nil {
		t.Fatal("missing to accepted")
	}
}
