package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/projeto-korp/app/internal/audit"
)

const (
	ServiceName = "http-server-projeto-korp"
	projectPath = "/projeto-korp"
)

type projectResponse struct {
	Name    string `json:"nome"`
	TimeUTC string `json:"horario"`
}

type handler struct{ auditService *audit.Service }

// NewHandler builds the application's HTTP router and project endpoint.
func NewHandler(services ...*audit.Service) http.Handler {
	var service *audit.Service
	if len(services) > 0 {
		service = services[0]
	}
	h := &handler{auditService: service}
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+projectPath, projectEndpoint)
	mux.HandleFunc("GET /audit-events", h.auditEvents)
	mux.HandleFunc("GET /request-statistics", h.requestStatistics)
	return auditMiddleware(service, mux)
}

func projectEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := projectResponse{
		Name:    "Projeto Korp",
		TimeUTC: time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(response)
}
