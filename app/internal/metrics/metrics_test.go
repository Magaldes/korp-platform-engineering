package metrics

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/projeto-korp/app/internal/httpserver"
)

func TestMetricsEndpointExposesApplicationCounter(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	get(t, server.URL+"/projeto-korp")
	response, body := request(t, http.MethodGet, server.URL+metricsPath)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if !strings.Contains(body, metricName) {
		t.Fatalf("metrics exposition does not contain %q", metricName)
	}
}

func TestProjectRequestIncrementsKnownRouteSeries(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	get(t, server.URL+projectPath)
	body := scrape(t, server.URL)

	assertCounterValue(t, body, `method="GET",route="/projeto-korp",status="200"`, 1)
}

func TestUnsupportedProjectMethodIncrements405Series(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	response, _ := request(t, http.MethodPost, server.URL+projectPath)
	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("POST /projeto-korp status = %d, want %d", response.StatusCode, http.StatusMethodNotAllowed)
	}
	if got := response.Header.Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("Allow = %q, want %q", got, "GET, HEAD")
	}

	assertCounterValue(t, scrape(t, server.URL), `method="POST",route="/projeto-korp",status="405"`, 1)
}

func TestUnknownRouteUsesBoundedUnmatchedLabel(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	response, _ := request(t, http.MethodGet, server.URL+"/arbitrary-client-path")
	if response.StatusCode != http.StatusNotFound {
		t.Fatalf("unknown route status = %d, want %d", response.StatusCode, http.StatusNotFound)
	}

	body := scrape(t, server.URL)
	assertCounterValue(t, body, `method="GET",route="unmatched",status="404"`, 1)
	if strings.Contains(body, `route="/arbitrary-client-path"`) {
		t.Fatalf("raw unknown path was emitted as route label")
	}
}

func TestMetricsScrapesDoNotIncrementApplicationCounter(t *testing.T) {
	server := newTestServer()
	defer server.Close()

	get(t, server.URL+projectPath)
	before := counterValue(t, scrape(t, server.URL), `method="GET",route="/projeto-korp",status="200"`)
	scrape(t, server.URL)
	scrape(t, server.URL)
	after := counterValue(t, scrape(t, server.URL), `method="GET",route="/projeto-korp",status="200"`)

	if after != before {
		t.Fatalf("application counter changed after /metrics scrapes: before=%v after=%v", before, after)
	}
}

func newTestServer() *httptest.Server {
	component := New()
	return httptest.NewServer(component.Handler(httpserver.NewHandler()))
}

func get(t *testing.T, url string) {
	t.Helper()
	response, _ := request(t, http.MethodGet, url)
	response.Body.Close()
}

func scrape(t *testing.T, serverURL string) string {
	t.Helper()
	response, body := request(t, http.MethodGet, serverURL+metricsPath)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("GET /metrics status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	return body
}

func request(t *testing.T, method, url string) (*http.Response, string) {
	t.Helper()
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatalf("create %s request: %v", method, err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("read %s %s response: %v", method, url, err)
	}
	return response, string(body)
}

func assertCounterValue(t *testing.T, exposition, labels string, want float64) {
	t.Helper()
	if got := counterValue(t, exposition, labels); got != want {
		t.Fatalf("counter labels {%s} = %v, want %v\nexposition:\n%s", labels, got, want, exposition)
	}
}

func counterValue(t *testing.T, exposition, labels string) float64 {
	t.Helper()
	prefix := fmt.Sprintf("%s{%s} ", metricName, labels)
	for _, line := range strings.Split(exposition, "\n") {
		if strings.HasPrefix(line, prefix) {
			value, err := strconv.ParseFloat(strings.TrimPrefix(line, prefix), 64)
			if err != nil {
				t.Fatalf("parse counter line %q: %v", line, err)
			}
			return value
		}
	}
	t.Fatalf("counter labels {%s} not found in exposition", labels)
	return 0
}
