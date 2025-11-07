# Data Dictionary - BGC App

**Version:** 1.0
**Database:** PostgreSQL 16
**Last Updated:** 2025-10-28

## Table of Contents

1. [Overview](#overview)
2. [Schema Structure](#schema-structure)
3. [Core Tables](#core-tables)
4. [Staging Tables](#staging-tables)
5. [Dimension Tables](#dimension-tables)
6. [Materialized Views](#materialized-views)
7. [Views](#views)
8. [Indexes](#indexes)
9. [Data Provenance](#data-provenance)
10. [Constraints & Validation](#constraints--validation)

---

## Overview

The BGC App database supports Brazilian export/import trade analytics with focus on NCM (Nomenclatura Comum do Mercosul) classification system.

### Key Entities
- **NCM Codes**: Product classification system (8-digit)
- **Trade Data**: Export/import values by NCM and year
- **Market Metrics**: TAM/SAM/SOM calculations
- **Routes**: Export route comparisons

---

## Schema Structure

### Schemas

| Schema | Purpose | Access |
|--------|---------|--------|
| `public` | Core operational tables and views | Read/Write |
| `stg` | Staging tables for data ingestion | Write (ingest), Read (transform) |
| `dim` | Dimension tables for reference data | Read (mostly) |

---

## Core Tables

### `public.ncm_lookup`

NCM (Nomenclatura Comum do Mercosul) product classification lookup table.

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `co_ncm` | VARCHAR(8) | NO | - | **PK** NCM code (8 digits) |
| `no_ncm_por` | TEXT | NO | - | NCM description in Portuguese |
| `co_sh2` | VARCHAR(2) | NO | - | HS chapter code (first 2 digits) |
| `no_sh2_por` | TEXT | NO | - | HS chapter description in Portuguese |

**Primary Key:** `co_ncm`

**Purpose:** Master reference table for NCM product codes and their chapter classifications.

**Example Row:**
```sql
co_ncm: '12345678'
no_ncm_por: 'Máquinas e aparelhos elétricos'
co_sh2: '84'
no_sh2_por: 'Máquinas e aparelhos'
```

---

### `public.trade_ncm_year`

Core trade data table containing export/import values by NCM code and year.

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `ncm` | VARCHAR(8) | NO | - | **PK** NCM product code |
| `ncm_desc` | TEXT | YES | NULL | NCM description (denormalized) |
| `unidade` | TEXT | YES | NULL | Statistical unit (kg, ton, etc) |
| `ano` | INT | NO | - | **PK** Year of trade data |
| `valor_usd_fob` | NUMERIC | NO | 0 | Trade value in USD (FOB) |
| `quantidade_estat` | NUMERIC | NO | 0 | Statistical quantity |
| `fluxo` | TEXT | NO | - | **PK** Flow type: 'exportacao' or 'importacao' |
| `ncm_chapter` | VARCHAR(2) | NO | - | NCM chapter (first 2 digits) |

**Primary Key:** `(ncm, ano, fluxo)`

**Constraints:**
- `CHECK (fluxo IN ('exportacao', 'importacao'))`

**Indexes:**
- `idx_trade_year_chapter` ON `(ano, ncm_chapter)`
- `idx_trade_fluxo` ON `(fluxo)`
- `idx_trade_ncm` ON `(ncm)`

**Purpose:** Central fact table for trade analytics. Stores aggregated export/import data.

**Data Volume:** ~1M rows per year (estimated)

**Example Row:**
```sql
ncm: '84714100'
ncm_desc: 'Máquinas automáticas para processamento de dados'
unidade: 'NUMERO (UNIDADES)'
ano: 2023
valor_usd_fob: 1500000.50
quantidade_estat: 1200
fluxo: 'exportacao'
ncm_chapter: '84'
```

---

## Staging Tables

### `stg.exportacao`

Staging table for raw export data ingestion.

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `id` | BIGSERIAL | NO | AUTO | **PK** Auto-increment identifier |
| `ano` | INT | NO | - | Year of export |
| `setor` | TEXT | YES | NULL | Sector/industry |
| `pais` | TEXT | YES | NULL | Destination country |
| `ncm` | TEXT | YES | NULL | NCM code (raw format) |
| `valor` | NUMERIC | YES | NULL | Export value |
| `qtde` | NUMERIC | YES | NULL | Quantity |
| `ingest_source` | TEXT | YES | 'unknown' | Data source identifier |
| `ingest_at` | TIMESTAMPTZ | NO | now() | Ingestion timestamp |
| `ingest_batch` | UUID | YES | gen_random_uuid() | Batch identifier for provenance |

**Primary Key:** `id`

**Purpose:** Receives raw export data from CSV/XLSX files before transformation.

**Retention:** Typically cleared after successful transformation to `trade_ncm_year`.

---

### `stg.importacao`

Staging table for raw import data ingestion.

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `id` | BIGSERIAL | NO | AUTO | **PK** Auto-increment identifier |
| `ano` | INT | NO | - | Year of import |
| `setor` | TEXT | YES | NULL | Sector/industry |
| `pais` | TEXT | YES | NULL | Origin country |
| `ncm` | TEXT | YES | NULL | NCM code (raw format) |
| `valor` | NUMERIC | YES | NULL | Import value |
| `qtde` | NUMERIC | YES | NULL | Quantity |
| `ingest_source` | TEXT | YES | 'unknown' | Data source identifier |
| `ingest_at` | TIMESTAMPTZ | NO | now() | Ingestion timestamp |
| `ingest_batch` | UUID | YES | gen_random_uuid() | Batch identifier for provenance |

**Primary Key:** `id`

**Purpose:** Receives raw import data from CSV/XLSX files before transformation.

**Retention:** Typically cleared after successful transformation to `trade_ncm_year`.

---

## Dimension Tables

### `dim.setor`

Sector/industry dimension table.

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `id` | SERIAL | NO | AUTO | **PK** Auto-increment identifier |
| `nome` | TEXT | NO | - | Sector name (unique) |

**Primary Key:** `id`

**Unique Constraint:** `nome`

**Purpose:** Reference table for industry sectors.

---

### `dim.ncm`

NCM dimension table (alternative to `ncm_lookup`).

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `id` | SERIAL | NO | AUTO | **PK** Auto-increment identifier |
| `codigo` | TEXT | NO | - | NCM code (unique) |
| `descricao` | TEXT | YES | NULL | NCM description |

**Primary Key:** `id`

**Unique Constraint:** `codigo`

**Purpose:** Alternative dimension table for NCM codes.

---

## Materialized Views

### `v_tam_by_year_chapter`

Aggregated market data by year and NCM chapter.

| Column | Type | Description |
|--------|------|-------------|
| `ano` | INT | Year |
| `ncm_chapter` | VARCHAR(2) | NCM chapter (2 digits) |
| `exp_valor_usd` | NUMERIC | Total export value in USD |
| `imp_valor_usd` | NUMERIC | Total import value in USD |
| `tam_total_usd` | NUMERIC | Total Addressable Market (exp + imp) |

**Definition:**
```sql
SELECT ano, ncm_chapter,
       SUM(CASE WHEN fluxo='exportacao' THEN valor_usd_fob ELSE 0 END) AS exp_valor_usd,
       SUM(CASE WHEN fluxo='importacao' THEN valor_usd_fob ELSE 0 END) AS imp_valor_usd,
       SUM(valor_usd_fob) AS tam_total_usd
FROM trade_ncm_year
GROUP BY ano, ncm_chapter
```

**Refresh Strategy:**
- **Manual:** `REFRESH MATERIALIZED VIEW v_tam_by_year_chapter;`
- **Automated:** Kubernetes CronJob daily at 03:00 UTC

**Purpose:** Pre-aggregated data for fast TAM/SAM/SOM calculations.

**Index Recommendations:**
```sql
CREATE INDEX idx_tam_year_chapter ON v_tam_by_year_chapter (ano, ncm_chapter);
```

---

## Views

### `v_quality_orphans_ncm`

Quality check view for orphaned NCM codes.

| Column | Type | Description |
|--------|------|-------------|
| `ncm` | VARCHAR(8) | NCM code without lookup entry |

**Definition:**
```sql
SELECT DISTINCT t.ncm
FROM trade_ncm_year t
LEFT JOIN ncm_lookup l ON l.co_ncm = t.ncm
WHERE l.co_ncm IS NULL
```

**Purpose:** Data quality monitoring - identifies NCM codes in trade data without corresponding lookup entries.

**Usage:**
```sql
-- Check for orphaned NCMs
SELECT COUNT(*) FROM v_quality_orphans_ncm;

-- List orphaned NCMs
SELECT * FROM v_quality_orphans_ncm LIMIT 10;
```

---

## Indexes

### Performance Indexes

| Index Name | Table | Columns | Purpose |
|------------|-------|---------|---------|
| `idx_trade_year_chapter` | `trade_ncm_year` | `(ano, ncm_chapter)` | Fast year/chapter filtering for TAM calculations |
| `idx_trade_fluxo` | `trade_ncm_year` | `(fluxo)` | Filter by export/import flow |
| `idx_trade_ncm` | `trade_ncm_year` | `(ncm)` | NCM code lookups |

### Future Index Recommendations

```sql
-- For route comparison queries
CREATE INDEX idx_trade_country ON trade_ncm_year (pais, ncm_chapter, ano)
  WHERE pais IS NOT NULL;

-- For time-series analysis
CREATE INDEX idx_trade_time_series ON trade_ncm_year (ano DESC, ncm);

-- For value-based filtering
CREATE INDEX idx_trade_high_value ON trade_ncm_year (valor_usd_fob DESC)
  WHERE valor_usd_fob > 1000000;
```

---

## Data Provenance

### Provenance Tracking (Migration 0003)

Data provenance is tracked through three columns added to staging tables:

| Column | Type | Purpose |
|--------|------|---------|
| `ingest_source` | TEXT | Source system identifier (e.g., 'comex-stat', 'siscomex') |
| `ingest_at` | TIMESTAMPTZ | Timestamp of data ingestion |
| `ingest_batch` | UUID | Batch UUID for grouping related records |

**Use Cases:**
- **Data lineage:** Track where data originated
- **Batch processing:** Identify and reprocess specific batches
- **Audit trail:** Investigate data quality issues
- **Idempotency:** Prevent duplicate ingestion of same batch

**Example Query:**
```sql
-- Find all records from a specific batch
SELECT * FROM stg.exportacao
WHERE ingest_batch = 'uuid-here';

-- Count records by source
SELECT ingest_source, COUNT(*)
FROM stg.exportacao
GROUP BY ingest_source;

-- Recent ingestions
SELECT ingest_batch, ingest_source, MIN(ingest_at), COUNT(*)
FROM stg.exportacao
WHERE ingest_at > NOW() - INTERVAL '7 days'
GROUP BY ingest_batch, ingest_source;
```

---

## Constraints & Validation

### Check Constraints

#### `trade_ncm_year.fluxo`
```sql
CHECK (fluxo IN ('exportacao', 'importacao'))
```
Ensures flow is either export or import.

### Data Type Constraints

| Column | Type | Rationale |
|--------|------|-----------|
| `ncm` | VARCHAR(8) | Fixed 8-digit NCM code format |
| `ncm_chapter` | VARCHAR(2) | Fixed 2-digit HS chapter format |
| `ano` | INT | Year as integer (2000-2100) |
| `valor_usd_fob` | NUMERIC | High precision for financial values |
| `ingest_batch` | UUID | Standard UUID v4 format |

### Default Values

| Table | Column | Default | Purpose |
|-------|--------|---------|---------|
| `trade_ncm_year` | `valor_usd_fob` | 0 | Avoid NULL in calculations |
| `trade_ncm_year` | `quantidade_estat` | 0 | Avoid NULL in calculations |
| `stg.exportacao` | `ingest_at` | now() | Auto-timestamp |
| `stg.exportacao` | `ingest_batch` | gen_random_uuid() | Auto-generate batch ID |
| `stg.exportacao` | `ingest_source` | 'unknown' | Default source |

---

## Data Lifecycle

### 1. Ingestion Phase
```
CSV/XLSX → stg.exportacao / stg.importacao
(via bgc-ingest service)
```

### 2. Transformation Phase
```
stg.* → trade_ncm_year
(via SQL migrations or ETL job)
```

### 3. Aggregation Phase
```
trade_ncm_year → v_tam_by_year_chapter (MVIEW)
(REFRESH MATERIALIZED VIEW - daily CronJob)
```

### 4. Consumption Phase
```
v_tam_by_year_chapter → API → Frontend
(via market/service.go)
```

---

## Query Patterns

### Common Queries

#### TAM Calculation
```sql
SELECT ano, ncm_chapter, tam_total_usd
FROM v_tam_by_year_chapter
WHERE ano BETWEEN 2020 AND 2023
  AND ncm_chapter IN ('84', '85')
ORDER BY ano, tam_total_usd DESC;
```

#### SAM Calculation (Filtered Chapters)
```sql
SELECT ano, SUM(tam_total_usd) AS sam_total_usd
FROM v_tam_by_year_chapter
WHERE ano = 2023
  AND ncm_chapter IN ('02', '08', '84', '85')  -- Scope chapters
GROUP BY ano;
```

#### Export vs Import Split
```sql
SELECT ncm_chapter,
       SUM(exp_valor_usd) AS total_export,
       SUM(imp_valor_usd) AS total_import,
       SUM(tam_total_usd) AS total_market
FROM v_tam_by_year_chapter
WHERE ano = 2023
GROUP BY ncm_chapter
ORDER BY total_market DESC;
```

---

## Schema Evolution

### Migration History

| Version | File | Description |
|---------|------|-------------|
| 0001 | `0001_init.sql` | Initial schema with staging and dimension tables |
| 0002 | `0002_add_staging_indexes.sql` | Performance indexes for staging tables |
| 0003 | `0003_proveniencia.sql` | Data provenance tracking columns |

### Future Migrations

Planned schema changes:

1. **Idempotency Keys** (Épico 3)
   - Add `idempotency_key` columns to staging tables
   - Add unique constraints for deduplication

2. **Route Data** (Future)
   - New table: `route_analysis` for route comparison results
   - Add country dimension table

3. **Tariff Scenarios** (Future)
   - Move from YAML to database tables
   - Enable dynamic tariff management

---

## Backup & Recovery

### Backup Strategy

**Kubernetes CronJob:** `postgres-backup-cronjob.yaml`
- **Frequency:** Daily at 02:00 UTC
- **Retention:** 7 days
- **Format:** pg_dump SQL format
- **Location:** `/backup/` volume

**Manual Backup:**
```bash
pg_dump -h localhost -U bgc -d bgc > backup_$(date +%Y%m%d).sql
```

**Restore:**
```bash
psql -h localhost -U bgc -d bgc < backup_20231028.sql
```

---

## Performance Tuning

### Recommended Settings

```sql
-- PostgreSQL config recommendations for analytics workload
shared_buffers = '4GB'
effective_cache_size = '12GB'
work_mem = '256MB'
maintenance_work_mem = '1GB'
random_page_cost = 1.1  -- SSD storage
```

### MVIEW Refresh Performance

```sql
-- Concurrent refresh (allows queries during refresh)
REFRESH MATERIALIZED VIEW CONCURRENTLY v_tam_by_year_chapter;

-- Requires unique index
CREATE UNIQUE INDEX idx_tam_unique ON v_tam_by_year_chapter (ano, ncm_chapter);
```

---

## Security & Access Control

### Role Recommendations

```sql
-- Read-only API role
CREATE ROLE bgc_api_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO bgc_api_readonly;
GRANT SELECT ON v_tam_by_year_chapter TO bgc_api_readonly;

-- Ingest role (write to staging only)
CREATE ROLE bgc_ingest;
GRANT INSERT, SELECT ON stg.exportacao, stg.importacao TO bgc_ingest;

-- Admin role (full access)
CREATE ROLE bgc_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public, stg, dim TO bgc_admin;
```

---

## Glossary

| Term | Definition |
|------|------------|
| **NCM** | Nomenclatura Comum do Mercosul - 8-digit product classification |
| **HS** | Harmonized System - International trade classification (first 6 digits) |
| **SH2** | HS Chapter - First 2 digits of HS code |
| **TAM** | Total Addressable Market - Total market size |
| **SAM** | Serviceable Addressable Market - Filtered by scope chapters |
| **SOM** | Serviceable Obtainable Market - SAM × penetration rate |
| **FOB** | Free On Board - Trade value excluding freight/insurance |
| **MVIEW** | Materialized View - Pre-computed aggregated data |
| **Provenance** | Data lineage tracking (source, batch, timestamp) |

---

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/16/)
- [NCM/HS Classification](http://www.mdic.gov.br/comercio-exterior/estatisticas-de-comercio-exterior)
- [BGC API Documentation](../README.md)
- [Architecture Document](./architecture_doc.md)

---

**Maintained by:** BGC Development Team
**Contact:** [GitHub Issues](https://github.com/rafamontilha/bgc-app/issues)
