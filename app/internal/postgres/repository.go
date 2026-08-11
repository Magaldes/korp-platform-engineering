package postgres

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projeto-korp/app/internal/audit"
)

type Repository struct{ pool *pgxpool.Pool }

func NewRepository(pool *pgxpool.Pool) *Repository { return &Repository{pool: pool} }

func (r *Repository) Close() {
	if r != nil && r.pool != nil {
		r.pool.Close()
	}
}

func (r *Repository) Append(ctx context.Context, event audit.Event) error {
	_, err := r.pool.Exec(ctx, `INSERT INTO request_audit_events (occurred_at, method, canonical_route, http_status, duration_ms) VALUES ($1, $2, $3, $4, $5)`, event.OccurredAt.UTC(), event.Method, event.CanonicalRoute, event.HTTPStatus, event.DurationMS)
	return err
}

func (r *Repository) List(ctx context.Context, filter audit.ListFilter) ([]audit.Event, bool, error) {
	args, where := listArgs(filter)
	rows, err := r.pool.Query(ctx, `SELECT id, occurred_at, method, canonical_route, http_status, duration_ms FROM request_audit_events WHERE `+where+` ORDER BY occurred_at DESC, id DESC LIMIT $`+strconv.Itoa(len(args)+1), append(args, filter.Limit+1)...)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	items := make([]audit.Event, 0, filter.Limit)
	for rows.Next() {
		var item audit.Event
		if err := rows.Scan(&item.ID, &item.OccurredAt, &item.Method, &item.CanonicalRoute, &item.HTTPStatus, &item.DurationMS); err != nil {
			return nil, false, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	hasNext := len(items) > filter.Limit
	if hasNext {
		items = items[:filter.Limit]
	}
	return items, hasNext, nil
}

func (r *Repository) Statistics(ctx context.Context, filter audit.StatisticsFilter) (audit.Statistics, error) {
	args, where := statsArgs(filter)
	var result audit.Statistics
	if err := r.pool.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE http_status >= 400), avg(duration_ms) FROM request_audit_events WHERE `+where, args...).Scan(&result.TotalRequests, &result.ErrorRequests, &result.AverageDurationMS); err != nil {
		return result, err
	}
	result.ByStatus = []audit.StatusCount{}
	rows, err := r.pool.Query(ctx, `SELECT http_status, count(*) FROM request_audit_events WHERE `+where+` GROUP BY http_status ORDER BY http_status`, args...)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item audit.StatusCount
		if err := rows.Scan(&item.Status, &item.Count); err != nil {
			rows.Close()
			return result, err
		}
		result.ByStatus = append(result.ByStatus, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return result, err
	}
	rows.Close()
	result.ByRoute = []audit.RouteCount{}
	rows, err = r.pool.Query(ctx, `SELECT canonical_route, count(*) FROM request_audit_events WHERE `+where+` GROUP BY canonical_route ORDER BY canonical_route`, args...)
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var item audit.RouteCount
		if err := rows.Scan(&item.Route, &item.Count); err != nil {
			rows.Close()
			return result, err
		}
		result.ByRoute = append(result.ByRoute, item)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return result, err
	}
	rows.Close()
	return result, nil
}

func (r *Repository) PurgeBefore(ctx context.Context, cutoff time.Time) (int64, error) {
	result, err := r.pool.Exec(ctx, `DELETE FROM request_audit_events WHERE occurred_at < $1`, cutoff.UTC())
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func listArgs(f audit.ListFilter) ([]any, string) {
	args := []any{f.From, f.To, f.Method, f.Route, f.Status}
	where := `($1::timestamptz IS NULL OR occurred_at >= $1) AND ($2::timestamptz IS NULL OR occurred_at < $2) AND ($3::text = '' OR method = $3) AND ($4::text = '' OR canonical_route = $4) AND ($5::int IS NULL OR http_status = $5)`
	if f.Cursor != "" {
		cursor, _ := audit.DecodeCursor(f.Cursor)
		args = append(args, cursor.OccurredAt, cursor.ID)
		where += ` AND (occurred_at, id) < ($6::timestamptz, $7::bigint)`
	}
	return args, where
}

func statsArgs(f audit.StatisticsFilter) ([]any, string) {
	return []any{f.From, f.To, f.Method, f.Route, f.Status}, `occurred_at >= $1 AND occurred_at < $2 AND ($3::text = '' OR method = $3) AND ($4::text = '' OR canonical_route = $4) AND ($5::int IS NULL OR http_status = $5)`
}

var _ audit.Repository = (*Repository)(nil)
