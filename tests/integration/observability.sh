#!/usr/bin/env bash
set -euo pipefail

: "${PROMETHEUS_BASE_URL:?PROMETHEUS_BASE_URL is required}"
: "${GRAFANA_BASE_URL:?GRAFANA_BASE_URL is required}"

prom_json() { curl --fail --silent --show-error "$PROMETHEUS_BASE_URL$1"; }
grafana_json() { curl --fail --silent --show-error -u "${GF_ADMIN_USER:-admin}:${GF_ADMIN_PASSWORD:-admin}" "$GRAFANA_BASE_URL$1"; }

targets=$(prom_json /api/v1/targets)
jq -e '.data.activeTargets[] | select(.labels.job == "http-server-projeto-korp" and .health == "up")' <<<"$targets" >/dev/null

up=$(prom_json '/api/v1/query?query=up%7Bjob%3D%22http-server-projeto-korp%22%7D')
jq -e '.data.result[] | select(.value[1] == "1")' <<<"$up" >/dev/null

requests=$(prom_json '/api/v1/query?query=sum%28projeto_korp_http_requests_total%29')
jq -e '.data.result | length > 0' <<<"$requests" >/dev/null

datasource=$(grafana_json /api/datasources/uid/prometheus)
jq -e '.uid == "prometheus" and .url == "http://prometheus:9090"' <<<"$datasource" >/dev/null

dashboard=$(grafana_json /api/dashboards/uid/korp-observability)
jq -e '
  .dashboard.title == "Projeto Korp - Observability" and
  ([.dashboard.panels[] | {title, type}] | sort) ==
  ([{title: "Availability", type: "stat"}, {title: "Request Rate", type: "timeseries"}, {title: "Request Total", type: "stat"}] | sort) and
  ([.dashboard.panels[].targets[].expr] | sort) ==
  (["sum(projeto_korp_http_requests_total)", "sum(rate(projeto_korp_http_requests_total[1m]))", "up{job=\"http-server-projeto-korp\"}"] | sort)
' <<<"$dashboard" >/dev/null

printf '%s\n' 'observability integration checks passed'
