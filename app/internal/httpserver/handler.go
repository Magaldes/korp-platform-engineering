package httpserver

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	ServiceName = "http-server-projeto-korp"
	projectPath = "/projeto-korp"
)

type projectResponse struct {
	Name    string `json:"nome"`
	TimeUTC string `json:"horario"`
}

// NewHandler builds the application's HTTP router and project endpoint.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+projectPath, projectEndpoint)
	return mux
}

func projectEndpoint(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := projectResponse{
		Name:    "Projeto Korp",
		TimeUTC: time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(response)
}
