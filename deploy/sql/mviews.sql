CREATE SCHEMA IF NOT EXISTS rpt;

-- MV por ano/setor
CREATE MATERIALIZED VIEW IF NOT EXISTS rpt.mv_exportacao_resumo AS
SELECT ano, setor, SUM(valor) AS valor_total, SUM(qtde) AS qtde_total
FROM stg.exportacao
GROUP BY ano, setor
WITH NO DATA;

-- MV por ano/pais
CREATE MATERIALIZED VIEW IF NOT EXISTS rpt.mv_exportacao_por_pais AS
SELECT ano, pais, SUM(valor) AS valor_total, SUM(qtde) AS qtde_total
FROM stg.exportacao
GROUP BY ano, pais
WITH NO DATA;

-- Índices únicos (exigidos por REFRESH CONCURRENTLY)
CREATE UNIQUE INDEX IF NOT EXISTS ux_mv_resumo ON rpt.mv_exportacao_resumo (ano, setor);
CREATE UNIQUE INDEX IF NOT EXISTS ux_mv_pais   ON rpt.mv_exportacao_por_pais (ano, pais);
