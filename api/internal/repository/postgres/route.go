package postgres

import (
	"database/sql"

	"bgc-app/internal/business/route"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository(db *sql.DB) route.Repository {
	return &RouteRepository{db: db}
}

func (r *RouteRepository) GetTAMByYearAndChapter(year int, chapter string) (float64, error) {
	var tam float64
	q := `SELECT tam_total_usd FROM v_tam_by_year_chapter WHERE ano=$1 AND ncm_chapter=$2`
	err := r.db.QueryRow(q, year, chapter).Scan(&tam)
	return tam, err
}
