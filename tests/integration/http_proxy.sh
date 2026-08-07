#!/usr/bin/env bash

set -euo pipefail

BASE_URL="${KORP_BASE_URL:-http://localhost:80}"

assert_status() {
    local expected="$1"
    shift
    local actual
    actual="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' "$@")"
    if [[ "$actual" != "$expected" ]]; then
        printf 'expected HTTP %s, got HTTP %s for %s\n' "$expected" "$actual" "$1" >&2
        return 1
    fi
}

project_response="$(mktemp)"
trap 'rm -f "$project_response"' EXIT

assert_status 200 "$BASE_URL/projeto-korp"
curl --silent --show-error "$BASE_URL/projeto-korp" >"$project_response"
python3 - "$project_response" <<'PY'
import json
import sys

with open(sys.argv[1], encoding="utf-8") as response_file:
    payload = json.load(response_file)

assert payload["nome"] == "Projeto Korp"
assert payload.get("horario")
PY

allow_header="$(curl --silent --show-error --dump-header - --output /dev/null -X POST "$BASE_URL/projeto-korp" | awk 'BEGIN{IGNORECASE=1} /^Allow:/{sub("\r$", ""); print; exit}')"
post_status="$(curl --silent --show-error --output /dev/null --write-out '%{http_code}' -X POST "$BASE_URL/projeto-korp")"
if [[ "$post_status" != 405 || "$allow_header" != 'Allow: GET, HEAD' ]]; then
    printf 'expected POST /projeto-korp to return 405 with Allow: GET, HEAD; got HTTP %s and %s\n' "$post_status" "$allow_header" >&2
    exit 1
fi

assert_status 404 "$BASE_URL/rota-desconhecida"
assert_status 404 "$BASE_URL/metrics"

printf 'HTTP proxy integration checks passed for %s\n' "$BASE_URL"
