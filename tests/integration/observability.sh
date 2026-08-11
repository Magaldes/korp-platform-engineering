#!/usr/bin/env bash
set -euo pipefail

: "${PROMETHEUS_BASE_URL:?PROMETHEUS_BASE_URL is required}"
: "${GRAFANA_BASE_URL:?GRAFANA_BASE_URL is required}"
: "${KORP_BASE_URL:=http://localhost:80}"

prom_json() { curl --fail --silent --show-error "$PROMETHEUS_BASE_URL$1"; }
grafana_json() { curl --fail --silent --show-error -u "${GF_ADMIN_USER:-admin}:${GF_ADMIN_PASSWORD:-admin}" "$GRAFANA_BASE_URL$1"; }

targets=$(prom_json /api/v1/targets)
jq -e '.data.activeTargets[] | select(.labels.job == "http-server-projeto-korp" and .health == "up")' <<<"$targets" >/dev/null

up=$(prom_json '/api/v1/query?query=up%7Bjob%3D%22http-server-projeto-korp%22%7D')
jq -e '.data.result[] | select(.value[1] == "1")' <<<"$up" >/dev/null

requests=$(prom_json '/api/v1/query?query=sum%28projeto_korp_http_requests_total%29')
jq -e '.data.result | length > 0' <<<"$requests" >/dev/null

query() {
  curl --fail --silent --show-error --get "$PROMETHEUS_BASE_URL/api/v1/query" --data-urlencode "query=$1"
}

assert_query_has_result() {
  local expression="$1"
  local result
  result="$(query "$expression")"
  jq -e '.status == "success" and (.data.result | length > 0)' <<<"$result" >/dev/null
}

assert_query_value_positive() {
  local expression="$1"
  local result
  result="$(query "$expression")"
  jq -e '.status == "success" and any(.data.result[]?; ((.value[1] // "0") | tonumber) > 0)' <<<"$result" >/dev/null
}

before_total="$(query 'sum(projeto_korp_http_requests_total)' | jq -r '.data.result[0].value[1] // "0"')"
curl --fail --silent --show-error "$KORP_BASE_URL/projeto-korp" >/dev/null
post_status="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' -X POST "$KORP_BASE_URL/projeto-korp")"
[[ "$post_status" == "405" ]]
unknown_status="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' "$KORP_BASE_URL/rota-desconhecida")"
[[ "$unknown_status" == "404" ]]
sleep 6
after_total="$(query 'sum(projeto_korp_http_requests_total)' | jq -r '.data.result[0].value[1] // "0"')"
awk -v before="$before_total" -v after="$after_total" 'BEGIN { exit !(after > before) }'

assert_query_has_result 'avg_over_time(up{job="http-server-projeto-korp"}[15m]) * 100'
assert_query_value_positive 'sum(increase(projeto_korp_http_requests_total[15m]))'
assert_query_has_result 'sum(rate(projeto_korp_http_requests_total[1m]))'
assert_query_has_result 'sum by (status) (rate(projeto_korp_http_requests_total[1m]))'
assert_query_value_positive 'sum(increase(projeto_korp_http_requests_total{route="unmatched",status="404"}[15m]))'
assert_query_value_positive 'sum(increase(projeto_korp_http_requests_total{route="/projeto-korp",status="405"}[15m]))'
assert_query_has_result 'scrape_duration_seconds{job="http-server-projeto-korp"}'
assert_query_has_result 'projeto_korp_statistics_cache_operations_total'
assert_query_has_result 'projeto_korp_statistics_cache_operation_duration_seconds_count'
assert_query_has_result 'projeto_korp_data_dependency_operations_total'
assert_query_has_result 'projeto_korp_data_dependency_operation_duration_seconds_count'

series=$(prom_json '/api/v1/series?match%5B%5D=projeto_korp_http_requests_total')
jq -e '.data[] | select(.route == "unmatched" and .status == "404")' <<<"$series" >/dev/null
jq -e '.data[] | select(.route == "/projeto-korp" and .status == "405")' <<<"$series" >/dev/null

datasource=$(grafana_json /api/datasources/uid/prometheus)
jq -e '.uid == "prometheus" and .url == "http://prometheus:9090"' <<<"$datasource" >/dev/null

dashboard=$(grafana_json /api/dashboards/uid/korp-observability)
jq -e '
  def expected_panels: [
    {title: "Service Availability", type: "stat"},
    {title: "Availability Over Selected Period", type: "stat"},
    {title: "Requests During Selected Period", type: "stat"},
    {title: "Request Rate", type: "timeseries"},
    {title: "Requests by HTTP Status", type: "timeseries"},
    {title: "Expected 404 - Unmatched Route", type: "stat"},
    {title: "Expected 405 - Invalid Method", type: "stat"},
    {title: "Prometheus Scrape Duration", type: "timeseries"},
    {title: "Statistics Cache Hit Ratio", type: "timeseries"},
    {title: "Statistics Cache Results", type: "timeseries"},
    {title: "Data Dependency Operation Duration", type: "timeseries"},
    {title: "Data Dependency Errors", type: "timeseries"}
  ];

  .dashboard.title == "Projeto Korp - Observability" and
  ((.dashboard.panels | map({title, type})) as $actual |
    (expected_panels) as $expected |
    (($actual | map(.title) | sort | group_by(.) | all(length == 1)) and
      (($actual | sort_by(.title)) == ($expected | sort_by(.title))))) and
  ([.dashboard.panels[].targets[].expr] | sort) ==
  (["avg_over_time(up{job=\"http-server-projeto-korp\"}[$__range]) * 100", "histogram_quantile(0.95, sum by (le, dependency, operation) (rate(projeto_korp_data_dependency_operation_duration_seconds_bucket[$__rate_interval])))", "scrape_duration_seconds{job=\"http-server-projeto-korp\"}", "sum by (dependency, operation) (rate(projeto_korp_data_dependency_operations_total{result=\"error\"}[$__rate_interval]))", "sum by (operation, result) (rate(projeto_korp_statistics_cache_operations_total[$__rate_interval]))", "sum by (status) (rate(projeto_korp_http_requests_total[$__rate_interval]))", "sum(increase(projeto_korp_http_requests_total{route=\"/projeto-korp\",status=\"405\"}[$__range]))", "sum(increase(projeto_korp_http_requests_total{route=\"unmatched\",status=\"404\"}[$__range]))", "sum(increase(projeto_korp_http_requests_total[$__range]))", "sum(rate(projeto_korp_http_requests_total[$__rate_interval]))", "sum(rate(projeto_korp_statistics_cache_operations_total{operation=\"lookup\",result=\"hit\"}[$__rate_interval])) / clamp_min(sum(rate(projeto_korp_statistics_cache_operations_total{operation=\"lookup\"}[$__rate_interval])), 1e-9) * 100", "up{job=\"http-server-projeto-korp\"}"] | sort)
' <<<"$dashboard" >/dev/null

printf '%s\n' 'observability integration checks passed'
