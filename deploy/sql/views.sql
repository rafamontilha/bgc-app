-- Schema de relatórios 
CREATE SCHEMA IF NOT EXISTS rpt;

-- Resumo por ano/setor
CREATE OR REPLACE VIEW rpt.vw_exportacao_resumo AS
Select
  ano,
  setor,
  SUM(valor) AS valor_total,
  SUM(qtde) AS qtde_total
FROM stg.exportacao
GROUP BY ano, setor
ORDER BY ano, setor;

-- Resumo por ano/pais
CREATE OR REPLACE VIEW rpt.vw_exportacao_por_pais AS
SELECT
  ano,
  pais,
  SUM(valor) AS valor_total,
  SUM(qtde) AS qtde_total
FROM stg.exportacao
GROUP BY ano, pais
ORDER BY ano, pais;
