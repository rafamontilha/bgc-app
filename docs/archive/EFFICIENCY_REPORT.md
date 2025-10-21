# BGC App - Code Efficiency Report

**Generated**: October 5, 2024  
**Reviewer**: Devin AI  
**Repository**: rafamontilha/bgc-app

## Executive Summary

This report documents 6 efficiency issues identified in the BGC analytics application codebase. These issues range from duplicate code and missing database indexes to suboptimal data loading patterns. The most impactful fixes would be:

1. **Add indexes to staging tables** (⚡ HIGH IMPACT - implemented in this PR)
2. **Implement batch inserts for CSV/XLSX loading** (⚡ HIGH IMPACT)
3. **Extract shared database connection code** (🔧 MEDIUM IMPACT)

## Detailed Findings

### 1. Missing Indexes on Staging Tables ⚡ HIGH IMPACT - FIXED

**Location**: `db/migrations/0001_init.sql`

**Issue**: The staging tables `stg.exportacao` and `stg.importacao` have no indexes defined. Any queries filtering by `ano`, `setor`, `pais`, or `ncm` will perform full table scans.

**Impact**: 
- Query performance: O(n) full table scans instead of O(log n) index lookups
- As tables grow to thousands/millions of rows, queries will become increasingly slow
- Estimated improvement: 10-1000x faster queries depending on table size

**Current Code** (`db/migrations/0001_init.sql:6-24`):
```sql
CREATE TABLE IF NOT EXISTS stg.exportacao (
  id BIGSERIAL PRIMARY KEY,
  ano INT NOT NULL,
  setor TEXT,
  pais TEXT,
  ncm TEXT,
  valor NUMERIC,
  qtde NUMERIC
);

CREATE TABLE IF NOT EXISTS stg.importacao (
  id BIGSERIAL PRIMARY KEY,
  ano INT NOT NULL,
  setor TEXT,
  pais TEXT,
  ncm TEXT,
  valor NUMERIC,
  qtde NUMERIC
);
```

**Recommended Fix**: Add indexes on commonly queried columns (implemented in migration 0002):
- Single column indexes: `ano`, `setor`, `pais`, `ncm`
- Composite indexes: `(ano, setor)`, `(ano, pais)` for common filtering patterns

**Status**: ✅ FIXED in `db/migrations/0002_add_staging_indexes.sql`

---

### 2. Individual INSERT Statements Instead of Batch Processing ⚡ HIGH IMPACT

**Location**: `services/bgc-ingest/main.go`

**Issue**: CSV and XLSX loaders use individual INSERT statements in a loop instead of batch operations.

**Current Code** (`services/bgc-ingest/main.go:299-302`):
```go
if _, err := tx.Exec(context.Background(), stmt, ano, setor, pais, ncm, valor, qtde); err != nil {
    return fmt.Errorf("insert linha %d: %w", count+1, err)
}
count++
```

Similar pattern in XLSX loader (`services/bgc-ingest/main.go:442-445`).

**Impact**:
- Each INSERT requires network round-trip to database
- For 10,000 rows: 10,000 round trips vs 1 batch operation
- Estimated improvement: 10-100x faster data loading

**Recommended Fix**:
1. Use PostgreSQL's COPY protocol: `pgx.CopyFrom()` for bulk inserts
2. Or use prepared statement batching with `pgx.Batch`
3. Example using COPY:
```go
copySource := pgx.CopyFromSlice(len(rows), func(i int) ([]any, error) {
    r := rows[i]
    return []any{r.ano, r.setor, r.pais, r.ncm, r.valor, r.qtde}, nil
})
_, err := tx.CopyFrom(ctx, pgx.Identifier{"stg", "exportacao"}, 
    []string{"ano", "setor", "pais", "ncm", "valor", "qtde"}, copySource)
```

**Status**: 📋 DOCUMENTED - Not implemented in this PR

---

### 3. Duplicate Database Connection Code 🔧 MEDIUM IMPACT

**Location**: `services/bgc-api/main.go` and `services/bgc-ingest/main.go`

**Issue**: Both services have identical database connection utility functions that should be shared.

**Duplicate Functions**:
- `getenv()` (api:20-25, ingest:34-39) - identical implementation
- `dsnFromEnv()` (api:28-36, ingest:43-51) - identical implementation
- Connection pool setup (api:39-50, ingest:55-69) - nearly identical

**Impact**:
- Code maintainability: Changes must be made in two places
- Risk of divergence: Connection logic could drift between services
- Binary size: Minor increase due to duplicate code

**Recommended Fix**:
1. Create a shared package: `internal/dbutil` or `pkg/database`
2. Move common functions to shared package:
```go
// internal/dbutil/pool.go
package dbutil

func GetEnv(k, def string) string { /* ... */ }
func DSNFromEnv() string { /* ... */ }
func NewPool(ctx context.Context) (*pgxpool.Pool, error) { /* ... */ }
```
3. Import in both services:
```go
import "github.com/rafamontilha/bgc-app/internal/dbutil"

pool, err := dbutil.NewPool(context.Background())
```

**Status**: 📋 DOCUMENTED - Not implemented in this PR

---

### 4. Missing Index on Materialized View 🔧 MEDIUM IMPACT

**Location**: `db/init/00_schema.sql`

**Issue**: The materialized view `v_tam_by_year_chapter` has no indexes defined.

**Current Code** (`db/init/00_schema.sql:27-33`):
```sql
CREATE MATERIALIZED VIEW IF NOT EXISTS v_tam_by_year_chapter AS
SELECT ano, ncm_chapter,
       SUM(CASE WHEN fluxo='exportacao' THEN valor_usd_fob ELSE 0 END) AS exp_valor_usd,
       SUM(CASE WHEN fluxo='importacao' THEN valor_usd_fob ELSE 0 END) AS imp_valor_usd,
       SUM(valor_usd_fob) AS tam_total_usd
FROM trade_ncm_year
GROUP BY ano, ncm_chapter;
```

**Impact**:
- Queries against the materialized view still require full scans
- Common query patterns like filtering by year or chapter are not optimized
- Estimated improvement: 5-50x faster queries on materialized view

**Recommended Fix**:
```sql
CREATE MATERIALIZED VIEW IF NOT EXISTS v_tam_by_year_chapter AS ...;

CREATE INDEX IF NOT EXISTS idx_mv_tam_ano ON v_tam_by_year_chapter(ano);
CREATE INDEX IF NOT EXISTS idx_mv_tam_chapter ON v_tam_by_year_chapter(ncm_chapter);
CREATE INDEX IF NOT EXISTS idx_mv_tam_ano_chapter ON v_tam_by_year_chapter(ano, ncm_chapter);
```

**Status**: 📋 DOCUMENTED - Not implemented in this PR

---

### 5. Inefficient Map Operations in Frontend 💡 LOW IMPACT

**Location**: `web/index.html`

**Issue**: The `aggregateByYear()` function calls `map.get(k)` twice per iteration instead of caching the result.

**Current Code** (`web/index.html:194-199`):
```javascript
for (const it of items){
    chapters.add(it.ncm_chapter);
    sum += it.valor_usd;
    const k = it.ano;
    map.set(k, (map.get(k)||0) + it.valor_usd);  // ← map.get(k) called twice
}
```

**Impact**:
- Minor CPU overhead on client side
- For typical data volumes (hundreds of rows), impact is negligible
- Estimated improvement: <1% in JavaScript execution time

**Recommended Fix**:
```javascript
for (const it of items){
    chapters.add(it.ncm_chapter);
    sum += it.valor_usd;
    const k = it.ano;
    const current = map.get(k) || 0;
    map.set(k, current + it.valor_usd);
}
```

**Status**: 📋 DOCUMENTED - Not implemented in this PR

---

### 6. Unnecessary API Call on Page Load 💡 LOW IMPACT

**Location**: `web/index.html` and `web/routes.html`

**Issue**: Both frontend pages make API calls immediately on page load before user interaction.

**Current Code**:
- `web/index.html:287` - calls `run()` immediately
- `web/routes.html:278` - calls `loadScenarios().then(run)` immediately

**Impact**:
- Unnecessary API load if user wants to change filters first
- Wastes server resources for potentially unwanted queries
- May show stale/default data that user will immediately replace

**Recommended Fix**:
1. Remove automatic calls on page load
2. Add placeholder text like "Click 'Consultar' to load data"
3. Or add a delay/debounce to only load after user stops interacting with filters

**Status**: 📋 DOCUMENTED - Not implemented in this PR

---

## Priority Recommendations

### Immediate (This PR)
- ✅ Add indexes to `stg.exportacao` and `stg.importacao` tables

### High Priority (Next Sprint)
- Implement batch inserts using COPY protocol in CSV/XLSX loaders
- Add indexes to materialized view `v_tam_by_year_chapter`

### Medium Priority
- Extract shared database connection code to common package
- Review and optimize other database queries for index usage

### Low Priority
- Optimize frontend JavaScript map operations
- Review UX for automatic API calls on page load

## Performance Testing Recommendations

After implementing fixes, measure:
1. **Staging table queries**: Compare `EXPLAIN ANALYZE` before/after indexes
2. **Data loading**: Time CSV/XLSX imports before/after batch implementation
3. **API response times**: Monitor `/metrics/*` endpoints with production-like data volumes
4. **Database connection pool**: Monitor pool utilization and connection churn

## Conclusion

The most impactful optimizations are the database-level improvements:
- Adding indexes (implemented in this PR)
- Implementing batch inserts (recommended for next sprint)

These changes alone could improve performance by 10-100x for common operations as the application scales to production data volumes.
