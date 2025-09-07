-- BGC  Migration 0001 (estrutura base)
-- Cria schemas e tabelas mínimas para ingestão e leitura

CREATE SCHEMA IF NOT EXISTS stg;
CREATE SCHEMA IF NOT EXISTS dim;
CREATE SCHEMA IF NOT EXISTS rpt;

-- Tabela de staging (cargas brutas)
CREATE TABLE IF NOT EXISTS stg.exportacao (
  id     BIGSERIAL PRIMARY KEY,
  ano    INTEGER      NOT NULL,
  setor  TEXT         NOT NULL,
  pais   TEXT         NOT NULL,
  ncm    TEXT         NOT NULL,
  valor  NUMERIC(18,2) NOT NULL,
  qtde   NUMERIC(18,3) NOT NULL DEFAULT 0
);

-- Índices úteis para leitura/aggregations
CREATE INDEX IF NOT EXISTS ix_export_ano_setor ON stg.exportacao (ano, setor);
CREATE INDEX IF NOT EXISTS ix_export_ano_pais  ON stg.exportacao (ano, pais);
CREATE INDEX IF NOT EXISTS ix_export_ncm       ON stg.exportacao (ncm);

-- Dimensões (mínimas) — opcionais, para futura normalização
CREATE TABLE IF NOT EXISTS dim.ncm (
  ncm TEXT PRIMARY KEY,
  descricao TEXT
);

CREATE TABLE IF NOT EXISTS dim.setor (
  setor TEXT PRIMARY KEY,
  descricao TEXT
);
