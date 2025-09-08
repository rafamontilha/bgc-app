-- Tabelas
CREATE TABLE IF NOT EXISTS ncm_lookup (
  co_ncm VARCHAR(8) PRIMARY KEY,
  no_ncm_por TEXT NOT NULL,
  co_sh2 VARCHAR(2) NOT NULL,
  no_sh2_por TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS trade_ncm_year (
  ncm VARCHAR(8) NOT NULL,
  ncm_desc TEXT,
  unidade TEXT,
  ano INT NOT NULL,
  valor_usd_fob NUMERIC NOT NULL DEFAULT 0,
  quantidade_estat NUMERIC NOT NULL DEFAULT 0,
  fluxo TEXT NOT NULL CHECK (fluxo IN ('exportacao','importacao')),
  ncm_chapter VARCHAR(2) NOT NULL,
  PRIMARY KEY (ncm, ano, fluxo)
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_trade_year_chapter ON trade_ncm_year (ano, ncm_chapter);
CREATE INDEX IF NOT EXISTS idx_trade_fluxo ON trade_ncm_year (fluxo);
CREATE INDEX IF NOT EXISTS idx_trade_ncm ON trade_ncm_year (ncm);

-- MView de agregados
CREATE MATERIALIZED VIEW IF NOT EXISTS v_tam_by_year_chapter AS
SELECT ano, ncm_chapter,
       SUM(CASE WHEN fluxo='exportacao' THEN valor_usd_fob ELSE 0 END) AS exp_valor_usd,
       SUM(CASE WHEN fluxo='importacao' THEN valor_usd_fob ELSE 0 END) AS imp_valor_usd,
       SUM(valor_usd_fob) AS tam_total_usd
FROM trade_ncm_year
GROUP BY ano, ncm_chapter;

-- Qualidade
CREATE OR REPLACE VIEW v_quality_orphans_ncm AS
SELECT DISTINCT t.ncm
FROM trade_ncm_year t
LEFT JOIN ncm_lookup l ON l.co_ncm = t.ncm
WHERE l.co_ncm IS NULL;
