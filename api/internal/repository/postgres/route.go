package postgres

import (
	"context"
	"database/sql"
	"time"

	"bgc-app/internal/business/route"
	"bgc-app/internal/observability/metrics"
	"bgc-app/internal/observability/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) route.Repository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) GetTAMByYearAndChapter(year int, chapter string) (float64, error) {
	ctx, span := tracing.StartSpan(context.Background(), "db.GetTAMByYearAndChapter")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year", year),
		attribute.String("chapter", chapter),
	)

	var tam float64
	q := `SELECT tam_total_usd FROM v_tam_by_year_chapter WHERE ano=$1 AND ncm_chapter=$2`

	// Start timing for metrics
	start := time.Now()
	err := r.db.QueryRowContext(ctx, q, year, chapter).Scan(&tam)
	duration := time.Since(start)
	metrics.RecordDBQuery("SELECT", "v_tam_by_year_chapter", duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return tam, err
	}

	span.SetAttributes(attribute.Float64("result.tam", tam))
	span.SetStatus(codes.Ok, "query successful")

	return tam, err
}
