package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/projeto-korp/app/internal/audit"
)

const auditWriteTimeout = 250 * time.Millisecond

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (w *responseRecorder) WriteHeader(status int) {
	if w.status == 0 {
		w.status = status
		w.ResponseWriter.WriteHeader(status)
	}
}
func (w *responseRecorder) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}
func (w *responseRecorder) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

func auditMiddleware(service *audit.Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(recorder, r)
		if service == nil || r.URL.Path == "/metrics" {
			return
		}
		duration := time.Since(started)
		ctx, cancel := context.WithTimeout(context.Background(), auditWriteTimeout)
		defer cancel()
		service.Record(ctx, audit.Event{OccurredAt: time.Now().UTC(), Method: r.Method, CanonicalRoute: audit.CanonicalRoute(r.URL.Path), HTTPStatus: recorder.Status(), DurationMS: duration.Milliseconds()})
	})
}

func (h *handler) auditEvents(w http.ResponseWriter, r *http.Request) {
	filter, err := audit.ParseListFilter(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}
	if h.auditService == nil {
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	items, hasNext, err := h.auditService.List(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	nextCursor := ""
	if hasNext && len(items) > 0 {
		nextCursor, err = audit.CursorFor(filter, items[len(items)-1])
		if err != nil {
			writeError(w, http.StatusServiceUnavailable)
			return
		}
	}
	body := struct {
		Items      []audit.Event `json:"items"`
		NextCursor *string       `json:"next_cursor"`
		Count      int           `json:"count"`
		Window     *window       `json:"window"`
	}{Items: items, Count: len(items)}
	if nextCursor != "" {
		body.NextCursor = &nextCursor
	}
	if filter.From != nil || filter.To != nil {
		body.Window = &window{From: filter.From, To: filter.To}
	}
	writeJSON(w, http.StatusOK, body)
}

func (h *handler) requestStatistics(w http.ResponseWriter, r *http.Request) {
	filter, err := audit.ParseStatisticsFilter(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest)
		return
	}
	if h.auditService == nil {
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	statistics, generatedAt, source, err := h.auditService.StatisticsWithSource(r.Context(), filter)
	if err != nil {
		writeError(w, http.StatusServiceUnavailable)
		return
	}
	body := struct {
		Window  window           `json:"window"`
		Filters statisticFilters `json:"filters"`
		audit.Statistics
		GeneratedAt time.Time `json:"generated_at"`
		Source      string    `json:"source"`
	}{Window: window{From: &filter.From, To: &filter.To}, Filters: statisticFilters{Method: filter.Method, Route: filter.Route, Status: filter.Status}, Statistics: statistics, GeneratedAt: generatedAt, Source: source}
	writeJSON(w, http.StatusOK, body)
}

type window struct {
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}
type statisticFilters struct {
	Method string `json:"method,omitempty"`
	Route  string `json:"route,omitempty"`
	Status *int   `json:"status,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int) {
	writeJSON(w, status, map[string]string{"error": http.StatusText(status)})
}
