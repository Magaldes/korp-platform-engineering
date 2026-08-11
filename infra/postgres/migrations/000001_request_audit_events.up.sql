CREATE TABLE request_audit_events (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    occurred_at TIMESTAMPTZ NOT NULL,
    method TEXT NOT NULL,
    canonical_route TEXT NOT NULL,
    http_status INTEGER NOT NULL,
    duration_ms BIGINT NOT NULL
);

CREATE INDEX request_audit_events_ordering_idx
    ON request_audit_events (occurred_at DESC, id DESC);

CREATE INDEX request_audit_events_filters_idx
    ON request_audit_events (method, canonical_route, http_status, occurred_at DESC);
