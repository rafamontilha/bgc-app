CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE stg.exportacao
  ADD COLUMN IF NOT EXISTS ingest_source text,
  ADD COLUMN IF NOT EXISTS ingest_at timestamptz NOT NULL DEFAULT now(),
  ADD COLUMN IF NOT EXISTS ingest_batch uuid;

-- defaults para novos inserts
ALTER TABLE stg.exportacao
  ALTER COLUMN ingest_source SET DEFAULT 'unknown',
  ALTER COLUMN ingest_batch  SET DEFAULT gen_random_uuid();

-- backfill dos registros existentes
UPDATE stg.exportacao
   SET ingest_source = COALESCE(ingest_source, 'historical'),
       ingest_batch  = COALESCE(ingest_batch, gen_random_uuid())
 WHERE ingest_source IS NULL
    OR ingest_batch  IS NULL;
