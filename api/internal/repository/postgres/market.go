package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bgc-app/internal/business/market"
	"bgc-app/internal/observability/metrics"
	"bgc-app/internal/observability/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type MarketRepository struct {
	db *sql.DB
}

func NewMarketRepository(db *sql.DB) market.Repository {
	return &MarketRepository{db: db}
}

func (r *MarketRepository) GetMarketDataByYearRange(yearFrom, yearTo int, chapters []string, ncmChapter string) ([]market.MarketItem, error) {
	ctx, span := tracing.StartSpan(context.Background(), "db.GetMarketDataByYearRange")
	defer span.End()

	span.SetAttributes(
		attribute.Int("year.from", yearFrom),
		attribute.Int("year.to", yearTo),
		attribute.Int("chapters.count", len(chapters)),
	)
	var sb strings.Builder
	args := []any{yearFrom, yearTo}

	sb.WriteString(`SELECT ano, ncm_chapter, tam_total_usd
	                FROM v_tam_by_year_chapter
	               WHERE ano BETWEEN $1 AND $2`)
	argPos := 3

	if len(chapters) > 0 {
		ph := make([]string, 0, len(chapters))
		for range chapters {
			ph = append(ph, fmt.Sprintf("$%d", argPos))
			argPos++
		}
		sb.WriteString(" AND ncm_chapter IN (" + strings.Join(ph, ",") + ")")
		for _, ch := range chapters {
			args = append(args, ch)
		}
	}

	if ncmChapter != "" {
		sb.WriteString(fmt.Sprintf(" AND ncm_chapter = $%d", argPos))
		args = append(args, ncmChapter)
	}

	sb.WriteString(" ORDER BY ano, ncm_chapter")

	// Start timing for metrics
	start := time.Now()
	rows, err := r.db.QueryContext(ctx, sb.String(), args...)
	duration := time.Since(start)
	metrics.RecordDBQuery("SELECT", "v_tam_by_year_chapter", duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	defer rows.Close()

	items := make([]market.MarketItem, 0, 64)
	for rows.Next() {
		var mi market.MarketItem
		if err := rows.Scan(&mi.Ano, &mi.NCMChapter, &mi.ValorUSD); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		items = append(items, mi)
	}

	span.SetAttributes(attribute.Int("result.count", len(items)))
	span.SetStatus(codes.Ok, "query successful")

	return items, rows.Err()
}
