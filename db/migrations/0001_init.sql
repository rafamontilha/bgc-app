-- Schemas
CREATE SCHEMA IF NOT EXISTS stg;
CREATE SCHEMA IF NOT EXISTS dim;

-- Staging (mínimo viável)
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

-- Dimensões
CREATE TABLE IF NOT EXISTS dim.setor (
  id SERIAL PRIMARY KEY,
  nome TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS dim.ncm (
  id SERIAL PRIMARY KEY,
  codigo TEXT UNIQUE NOT NULL,
  descricao TEXT
);

