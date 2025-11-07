-- Migration 0004: Idempotency Support
-- Adds idempotency key tracking to prevent duplicate processing

-- Add idempotency_key to staging tables
ALTER TABLE stg.exportacao
  ADD COLUMN IF NOT EXISTS idempotency_key TEXT;

ALTER TABLE stg.importacao
  ADD COLUMN IF NOT EXISTS idempotency_key TEXT;

-- Create index for fast lookups
CREATE INDEX IF NOT EXISTS idx_exportacao_idempotency
  ON stg.exportacao(idempotency_key)
  WHERE idempotency_key IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_importacao_idempotency
  ON stg.importacao(idempotency_key)
  WHERE idempotency_key IS NOT NULL;

-- Create idempotency tracking table for API requests
CREATE TABLE IF NOT EXISTS public.api_idempotency (
  idempotency_key TEXT PRIMARY KEY,
  endpoint TEXT NOT NULL,
  method TEXT NOT NULL,
  request_params JSONB,
  response_status INT NOT NULL,
  response_body JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  expires_at TIMESTAMPTZ NOT NULL DEFAULT (now() + INTERVAL '24 hours')
);

-- Index for cleanup of expired keys
CREATE INDEX IF NOT EXISTS idx_api_idempotency_expires
  ON public.api_idempotency(expires_at);

-- Cleanup function for expired idempotency keys
CREATE OR REPLACE FUNCTION cleanup_expired_idempotency_keys()
RETURNS INTEGER AS $$
DECLARE
  deleted_count INTEGER;
BEGIN
  DELETE FROM public.api_idempotency
  WHERE expires_at < now();

  GET DIAGNOSTICS deleted_count = ROW_COUNT;
  RETURN deleted_count;
END;
$$ LANGUAGE plpgsql;

-- Optional: Create a scheduled job to cleanup (requires pg_cron extension)
-- SELECT cron.schedule('cleanup-idempotency', '0 * * * *', 'SELECT cleanup_expired_idempotency_keys()');

COMMENT ON TABLE public.api_idempotency IS 'Idempotency key tracking for API requests (24h retention)';
COMMENT ON COLUMN stg.exportacao.idempotency_key IS 'Client-provided idempotency key for duplicate detection';
COMMENT ON COLUMN stg.importacao.idempotency_key IS 'Client-provided idempotency key for duplicate detection';
