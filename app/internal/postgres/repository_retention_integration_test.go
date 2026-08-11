package postgres

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projeto-korp/app/internal/audit"
)

func TestPurgeBeforePreservesCutoffAndNewerEvents(t *testing.T) {
	if os.Getenv("KORP_RETENTION_INTEGRATION") != "1" {
		t.Skip("set KORP_RETENTION_INTEGRATION=1 to run against PostgreSQL")
	}

	ctx := context.Background()
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("KORP_DB_HOST") + ":" + os.Getenv("KORP_DB_PORT")
	adminConnection := url.URL{Scheme: "postgres", User: url.UserPassword(user, password), Host: host, Path: "postgres"}
	adminConfig, err := pgxpool.ParseConfig(adminConnection.String())
	if err != nil {
		t.Fatal(err)
	}
	adminPool, err := pgxpool.NewWithConfig(ctx, adminConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer adminPool.Close()
	if err := adminPool.Ping(ctx); err != nil {
		t.Fatal(err)
	}
	testDatabase := fmt.Sprintf("korp_retention_test_%d", os.Getpid())
	if _, err := adminPool.Exec(ctx, `CREATE DATABASE `+testDatabase); err != nil {
		t.Fatal(err)
	}
	testConnection := url.URL{Scheme: "postgres", User: url.UserPassword(user, password), Host: host, Path: testDatabase}
	testConfig, err := pgxpool.ParseConfig(testConnection.String())
	if err != nil {
		t.Fatal(err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, testConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		pool.Close()
		_, _ = adminPool.Exec(ctx, `DROP DATABASE `+testDatabase)
	}()
	if _, err := pool.Exec(ctx, `CREATE TABLE request_audit_events (id BIGSERIAL PRIMARY KEY, occurred_at TIMESTAMPTZ NOT NULL, method TEXT NOT NULL, canonical_route TEXT NOT NULL, http_status INTEGER NOT NULL, duration_ms BIGINT NOT NULL)`); err != nil {
		t.Fatal(err)
	}

	repository := NewRepository(pool)
	cutoff := time.Now().UTC().Truncate(time.Microsecond)
	rows := []audit.Event{
		{OccurredAt: cutoff.Add(-time.Microsecond), Method: "GET", CanonicalRoute: "retention-test-old", HTTPStatus: 200},
		{OccurredAt: cutoff, Method: "GET", CanonicalRoute: "retention-test-cutoff", HTTPStatus: 200},
		{OccurredAt: cutoff.Add(time.Microsecond), Method: "GET", CanonicalRoute: "retention-test-new", HTTPStatus: 200},
		{OccurredAt: cutoff.Add(2 * time.Microsecond), Method: "GET", CanonicalRoute: "retention-test-unrelated", HTTPStatus: 201},
	}
	ids := make([]int64, 0, len(rows))
	for _, event := range rows {
		var id int64
		err := pool.QueryRow(ctx, `INSERT INTO request_audit_events (occurred_at, method, canonical_route, http_status, duration_ms) VALUES ($1, $2, $3, $4, $5) RETURNING id`, event.OccurredAt, event.Method, event.CanonicalRoute, event.HTTPStatus, event.DurationMS).Scan(&id)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	defer func() { _, _ = pool.Exec(ctx, `DELETE FROM request_audit_events WHERE id = ANY($1)`, ids) }()

	deleted, err := repository.PurgeBefore(ctx, cutoff)
	if err != nil {
		t.Fatal(err)
	}
	if deleted != 1 {
		t.Fatalf("deleted = %d, want 1", deleted)
	}

	var remaining int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM request_audit_events WHERE id = ANY($1)`, ids).Scan(&remaining); err != nil {
		t.Fatal(err)
	}
	if remaining != 3 {
		t.Fatalf("remaining = %d, want cutoff and newer/unrelated rows preserved", remaining)
	}
}
