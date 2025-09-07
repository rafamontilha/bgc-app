-- Criação das Materialized Views (sem dados iniciais)
CREATE SCHEMA IF NOT EXISTS rpt;

CREATE MATERIALIZED VIEW IF NOT EXISTS rpt.mv_exportacao_resumo AS
SELECT
  ano,
  setor,
  SUM(valor) AS valor_total,
  SUM(qtde)  AS qtde_total
FROM stg.exportacao
GROUP BY 1,2
WITH NO DATA;

CREATE MATERIALIZED VIEW IF NOT EXISTS rpt.mv_exportacao_por_pais AS
SELECT
  ano,
  pais,
  SUM(valor) AS valor_total,
  SUM(qtde)  AS qtde_total
FROM stg.exportacao
GROUP BY 1,2
WITH NO DATA;

-- Índices exigidos para REFRESH ... CONCURRENTLY
CREATE UNIQUE INDEX IF NOT EXISTS uq_mv_resumo_key
  ON rpt.mv_exportacao_resumo (ano, setor);

CREATE UNIQUE INDEX IF NOT EXISTS uq_mv_pais_key
  ON rpt.mv_exportacao_por_pais (ano, pais);
