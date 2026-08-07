package httpserver

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProjectEndpointReturnsJSONContract(t *testing.T) {
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	before := time.Now().UTC()
	response, err := http.Get(server.URL + "/projeto-korp")
	if err != nil {
		t.Fatalf("GET /projeto-korp: %v", err)
	}
	defer response.Body.Close()
	after := time.Now().UTC()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if got := response.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want %q", got, "application/json")
	}

	var body projectResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Name != "Projeto Korp" {
		t.Fatalf("nome = %q, want %q", body.Name, "Projeto Korp")
	}

	parsedTime, err := time.Parse(time.RFC3339, body.TimeUTC)
	if err != nil {
		t.Fatalf("parse horario %q as RFC3339: %v", body.TimeUTC, err)
	}
	if parsedTime.Location() != time.UTC {
		t.Fatalf("horario location = %v, want UTC", parsedTime.Location())
	}
	if parsedTime.Before(before.Add(-time.Second)) || parsedTime.After(after.Add(time.Second)) {
		t.Fatalf("horario %v is outside request-time window [%v, %v]", parsedTime, before, after)
	}
}

func TestProjectEndpointResolvesTimePerRequest(t *testing.T) {
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	first := requestProjectTime(t, server.URL)
	time.Sleep(1100 * time.Millisecond)
	second := requestProjectTime(t, server.URL)

	if !second.After(first) {
		t.Fatalf("second horario %v is not after first horario %v", second, first)
	}
}

func TestProjectEndpointUsesNativeHTTPMethodSemantics(t *testing.T) {
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	request, err := http.NewRequest(http.MethodHead, server.URL+"/projeto-korp", nil)
	if err != nil {
		t.Fatalf("create HEAD request: %v", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("HEAD /projeto-korp: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("HEAD status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if got := response.Header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("HEAD Content-Type = %q, want %q", got, "application/json")
	}
	if body, err := io.ReadAll(response.Body); err != nil {
		t.Fatalf("read HEAD body: %v", err)
	} else if len(body) != 0 {
		t.Fatalf("HEAD body length = %d, want 0", len(body))
	}
}

func TestUnknownRouteReturnsNotFound(t *testing.T) {
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	response, err := http.Get(server.URL + "/unknown")
	if err != nil {
		t.Fatalf("GET /unknown: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusNotFound)
	}
}

func TestUnsupportedMethodReturnsMethodNotAllowed(t *testing.T) {
	server := httptest.NewServer(NewHandler())
	defer server.Close()

	request, err := http.NewRequest(http.MethodPost, server.URL+"/projeto-korp", nil)
	if err != nil {
		t.Fatalf("create POST request: %v", err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("POST /projeto-korp: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusMethodNotAllowed)
	}
	if got := response.Header.Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("Allow = %q, want %q", got, "GET, HEAD")
	}
}

func requestProjectTime(t *testing.T, serverURL string) time.Time {
	t.Helper()

	response, err := http.Get(serverURL + "/projeto-korp")
	if err != nil {
		t.Fatalf("GET /projeto-korp: %v", err)
	}
	defer response.Body.Close()

	var body projectResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	parsedTime, err := time.Parse(time.RFC3339, body.TimeUTC)
	if err != nil {
		t.Fatalf("parse horario %q: %v", body.TimeUTC, err)
	}
	return parsedTime
}
