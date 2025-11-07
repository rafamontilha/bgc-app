# Idempotency and Reprocessing Policy

**Version:** 1.0
**Last Updated:** 2025-10-28

## Overview

This document describes the idempotency and reprocessing policies implemented in BGC App to ensure reliable data processing and prevent duplicate operations.

---

## Table of Contents

1. [What is Idempotency?](#what-is-idempotency)
2. [Why Idempotency Matters](#why-idempotency-matters)
3. [API Idempotency](#api-idempotency)
4. [Data Ingestion Idempotency](#data-ingestion-idempotency)
5. [Reprocessing Policy](#reprocessing-policy)
6. [Implementation Details](#implementation-details)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

---

## What is Idempotency?

**Idempotency** ensures that performing the same operation multiple times produces the same result as performing it once. This is critical for:

- Network retries (timeouts, failures)
- Duplicate requests from clients
- Data reprocessing scenarios
- Recovery from system failures

---

## Why Idempotency Matters

### Without Idempotency
```
POST /ingest → Process 1000 records
POST /ingest (retry due to timeout) → Process 1000 records AGAIN
Result: 2000 duplicate records ❌
```

### With Idempotency
```
POST /ingest (Idempotency-Key: abc123) → Process 1000 records
POST /ingest (Idempotency-Key: abc123) → Return cached result
Result: 1000 unique records ✅
```

---

## API Idempotency

### HTTP Methods

| Method | Naturally Idempotent | Requires Idempotency-Key |
|--------|---------------------|-------------------------|
| GET | ✅ Yes | ❌ No |
| PUT | ✅ Yes | ❌ No |
| DELETE | ✅ Yes | ❌ No |
| POST | ❌ No | ✅ Yes (recommended) |
| PATCH | ❌ No | ✅ Yes (recommended) |

### Idempotency-Key Header

Clients can provide an `Idempotency-Key` header with POST/PATCH requests:

```http
POST /v1/ingest/export HTTP/1.1
Host: api.bgc.dev
Content-Type: application/json
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000

{
  "data": [...]
}
```

**Requirements:**
- **Format:** UUID v4 recommended (or any unique string)
- **Length:** 16-128 characters
- **Uniqueness:** Must be unique per logical operation
- **Reusability:** Same key = same operation

### Cache Behavior

| Scenario | Behavior |
|----------|----------|
| First request with key | Process normally, cache response for 24h |
| Duplicate request (within 24h) | Return cached response immediately |
| Duplicate request (after 24h) | Process as new request |
| No idempotency key | Process normally, no caching |

**Headers in Response:**
```http
HTTP/1.1 200 OK
X-Idempotency-Cached: true
X-Idempotency-Cached-At: 2025-10-28T10:30:00Z
```

---

## Data Ingestion Idempotency

### Batch Ingestion

Every batch ingestion receives a unique `ingest_batch` UUID:

```sql
SELECT id, ingest_batch, ingest_source, ingest_at
FROM stg.exportacao
WHERE ingest_batch = '550e8400-e29b-41d4-a716-446655440000';
```

### Duplicate Detection

#### Level 1: Batch ID
Prevent re-ingesting the same file:

```sql
-- Check if batch already exists
SELECT COUNT(*) FROM stg.exportacao
WHERE ingest_batch = 'batch-uuid-here';

-- If count > 0: Reject with error
```

#### Level 2: Content Hash
Prevent duplicate records within a batch:

```sql
-- Add content hash to staging (future enhancement)
ALTER TABLE stg.exportacao
ADD COLUMN content_hash TEXT;

-- Generate hash from key fields
content_hash = SHA256(ano || ncm || pais || valor)
```

#### Level 3: Business Key
Prevent logical duplicates across batches:

```sql
-- Unique constraint on business key
ALTER TABLE trade_ncm_year
ADD CONSTRAINT uq_trade_ncm_year_business_key
UNIQUE (ncm, ano, fluxo);

-- Insert with ON CONFLICT
INSERT INTO trade_ncm_year (...)
VALUES (...)
ON CONFLICT (ncm, ano, fluxo)
DO UPDATE SET valor_usd_fob = EXCLUDED.valor_usd_fob;
```

---

## Reprocessing Policy

### When to Reprocess

| Scenario | Action | Idempotency Strategy |
|----------|--------|---------------------|
| Network failure during ingestion | Retry with same batch ID | ✅ Use original batch UUID |
| Data quality issues found | Reprocess with new batch ID | ❌ New UUID (different data) |
| Schema migration | Reprocess all historical data | ✅ Track migration version |
| Correction/backfill | Reprocess specific time range | ⚠️ Depends on use case |

### Reprocessing Workflow

#### 1. Identify Batch to Reprocess
```sql
-- Find batches with errors
SELECT ingest_batch, ingest_source, COUNT(*), MIN(ingest_at)
FROM stg.exportacao
WHERE ingest_source = 'siscomex'
  AND ingest_at BETWEEN '2025-10-01' AND '2025-10-31'
GROUP BY ingest_batch, ingest_source;
```

#### 2. Delete Bad Data
```sql
-- Delete from staging
DELETE FROM stg.exportacao
WHERE ingest_batch = 'bad-batch-uuid';

-- Delete from target (if already propagated)
DELETE FROM trade_ncm_year
WHERE ncm IN (
  SELECT DISTINCT ncm FROM stg.exportacao
  WHERE ingest_batch = 'bad-batch-uuid'
);
```

#### 3. Re-ingest with New Batch ID
```bash
# Re-run ingest with new batch ID
curl -X POST https://api.bgc.dev/v1/ingest/export \
  -H "Idempotency-Key: new-uuid-here" \
  -H "X-Batch-Source: siscomex-reprocess-2025-10" \
  -d @corrected_data.json
```

---

## Implementation Details

### Database Tables

#### `public.api_idempotency`
Tracks idempotency keys for API requests:

```sql
CREATE TABLE public.api_idempotency (
  idempotency_key TEXT PRIMARY KEY,
  endpoint TEXT NOT NULL,
  method TEXT NOT NULL,
  request_params JSONB,
  response_status INT NOT NULL,
  response_body JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT (now() + INTERVAL '24 hours')
);
```

**Retention:** 24 hours (auto-cleanup)

#### Staging Tables
Columns for idempotency tracking:

```sql
ALTER TABLE stg.exportacao
ADD COLUMN idempotency_key TEXT;

CREATE INDEX idx_exportacao_idempotency
ON stg.exportacao(idempotency_key)
WHERE idempotency_key IS NOT NULL;
```

### In-Memory Cache (Current Implementation)

```go
// Cache with 24h TTL
cache := NewIdempotencyCache()
cache.Set(key, response)

// Retrieve cached response
if cached, exists := cache.Get(key); exists {
  return cached.Response
}
```

**Limitations:**
- Lost on server restart
- Not shared across multiple API instances

### Redis Cache (Recommended for Production)

```go
// Set with expiration
rdb.Set(ctx, key, responseJSON, 24*time.Hour)

// Get
val, err := rdb.Get(ctx, key).Result()
```

**Benefits:**
- Persistent across restarts
- Shared cache for horizontal scaling
- Built-in TTL and eviction

---

## Best Practices

### For API Clients

1. **Always use Idempotency-Key for POST/PATCH**
   ```javascript
   const idempotencyKey = uuidv4();

   await fetch('/api/ingest', {
     method: 'POST',
     headers: {
       'Idempotency-Key': idempotencyKey,
       'Content-Type': 'application/json'
     },
     body: JSON.stringify(data)
   });
   ```

2. **Retry with the SAME key**
   ```javascript
   const response = await retryWithBackoff(async () => {
     return fetch('/api/ingest', {
       headers: { 'Idempotency-Key': idempotencyKey }
     });
   });
   ```

3. **Use deterministic keys for automation**
   ```bash
   # Good: Include date/source in key
   KEY="siscomex-export-2025-10-28"

   # Bad: Random key every time
   KEY=$(uuidv4)  # Don't do this for retries!
   ```

### For Data Engineers

1. **Track batch provenance**
   ```sql
   INSERT INTO stg.exportacao (
     ...,
     ingest_source,
     ingest_batch,
     idempotency_key
   ) VALUES (
     ...,
     'siscomex',
     'batch-uuid-here',
     'client-provided-key'  -- Link to API request
   );
   ```

2. **Implement upsert patterns**
   ```sql
   INSERT INTO trade_ncm_year (...)
   VALUES (...)
   ON CONFLICT (ncm, ano, fluxo)
   DO UPDATE SET
     valor_usd_fob = EXCLUDED.valor_usd_fob,
     updated_at = now();
   ```

3. **Log reprocessing events**
   ```sql
   CREATE TABLE audit.reprocessing_log (
     id SERIAL PRIMARY KEY,
     original_batch UUID,
     new_batch UUID,
     reason TEXT,
     reprocessed_by TEXT,
     reprocessed_at TIMESTAMPTZ DEFAULT now()
   );
   ```

---

## Examples

### Example 1: Client Retry

```typescript
// Client-side retry with idempotency
async function ingestDataSafely(data: ExportData[]) {
  const idempotencyKey = `ingest-${Date.now()}-${generateHash(data)}`;

  for (let attempt = 0; attempt < 3; attempt++) {
    try {
      const response = await fetch('/v1/ingest/export', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Idempotency-Key': idempotencyKey  // SAME KEY on retries
        },
        body: JSON.stringify(data)
      });

      if (response.ok) {
        const cached = response.headers.get('X-Idempotency-Cached');
        if (cached === 'true') {
          console.log('Returned cached result (duplicate request)');
        }
        return response.json();
      }
    } catch (error) {
      if (attempt === 2) throw error;
      await sleep(Math.pow(2, attempt) * 1000);  // Exponential backoff
    }
  }
}
```

### Example 2: Batch Reprocessing

```bash
#!/bin/bash
# Reprocess failed batches

BATCH_DATE="2025-10-28"
SOURCE="siscomex"

# 1. Find failed batches
psql -d bgc -c "
  SELECT ingest_batch
  FROM stg.exportacao
  WHERE ingest_source = '$SOURCE'
    AND DATE(ingest_at) = '$BATCH_DATE'
  GROUP BY ingest_batch
  HAVING COUNT(*) < 1000  -- Expected count
" -t -A > failed_batches.txt

# 2. Delete failed batches
while read BATCH_UUID; do
  echo "Deleting batch: $BATCH_UUID"
  psql -d bgc -c "
    DELETE FROM stg.exportacao
    WHERE ingest_batch = '$BATCH_UUID'
  "
done < failed_batches.txt

# 3. Re-ingest with new batch IDs
for FILE in data/${BATCH_DATE}/*.csv; do
  NEW_BATCH=$(uuidgen)

  curl -X POST http://localhost:8080/v1/ingest/export \
    -H "Idempotency-Key: ${SOURCE}-reprocess-${BATCH_DATE}-${NEW_BATCH}" \
    -H "Content-Type: application/json" \
    -d @"$FILE"

  echo "Re-ingested: $FILE with batch: $NEW_BATCH"
done
```

### Example 3: Upsert Pattern

```sql
-- Upsert to handle reprocessing gracefully
INSERT INTO trade_ncm_year (
  ncm,
  ano,
  fluxo,
  valor_usd_fob,
  quantidade_estat,
  ncm_chapter,
  updated_at
)
SELECT
  ncm,
  ano,
  'exportacao' AS fluxo,
  SUM(valor) AS valor_usd_fob,
  SUM(qtde) AS quantidade_estat,
  LEFT(ncm, 2) AS ncm_chapter,
  now() AS updated_at
FROM stg.exportacao
WHERE ingest_batch = 'batch-uuid-here'
GROUP BY ncm, ano

ON CONFLICT (ncm, ano, fluxo)
DO UPDATE SET
  valor_usd_fob = EXCLUDED.valor_usd_fob,
  quantidade_estat = EXCLUDED.quantidade_estat,
  updated_at = EXCLUDED.updated_at;
```

---

## Monitoring & Alerts

### Key Metrics

```sql
-- Duplicate request rate (last 24h)
SELECT
  DATE_TRUNC('hour', created_at) AS hour,
  COUNT(*) AS total_requests,
  COUNT(DISTINCT idempotency_key) AS unique_requests,
  COUNT(*) - COUNT(DISTINCT idempotency_key) AS duplicate_requests
FROM public.api_idempotency
WHERE created_at > now() - INTERVAL '24 hours'
GROUP BY hour
ORDER BY hour DESC;

-- Batch reprocessing frequency
SELECT
  ingest_source,
  COUNT(DISTINCT ingest_batch) AS total_batches,
  COUNT(DISTINCT idempotency_key) AS unique_operations,
  AVG(COUNT(*)) AS avg_records_per_batch
FROM stg.exportacao
WHERE ingest_at > now() - INTERVAL '7 days'
GROUP BY ingest_source;
```

### Recommended Alerts

- **High duplicate rate** (>10%): May indicate client retry issues
- **Cache size growing** (>10K entries): Check TTL and cleanup
- **Zero idempotency keys**: Clients not using feature

---

## Migration Guide

### Upgrading Existing Data

```sql
-- Step 1: Apply migration
\i db/migrations/0004_idempotency.sql

-- Step 2: Backfill batch UUIDs for historical data
UPDATE stg.exportacao
SET ingest_batch = gen_random_uuid()
WHERE ingest_batch IS NULL;

-- Step 3: Enable API idempotency middleware
-- (See api/internal/app/server.go)
```

---

## References

- [RFC 7231 - HTTP Idempotency](https://tools.ietf.org/html/rfc7231#section-4.2.2)
- [Stripe API Idempotency](https://stripe.com/docs/api/idempotent_requests)
- [PostgreSQL INSERT ... ON CONFLICT](https://www.postgresql.org/docs/16/sql-insert.html)
- [UUID Best Practices](https://www.ietf.org/rfc/rfc4122.txt)

---

**Maintained by:** BGC Development Team
**Questions:** [GitHub Issues](https://github.com/rafamontilha/bgc-app/issues)
