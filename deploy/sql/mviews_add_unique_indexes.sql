-- Unique indexes exigidos para REFRESH MATERIALIZED VIEW CONCURRENTLY
CREATE UNIQUE INDEX IF NOT EXISTS uq_mv_resumo_key
  ON rpt.mv_exportacao_resumo(ano, setor);

CREATE UNIQUE INDEX IF NOT EXISTS uq_mv_pais_key
  ON rpt.mv_exportacao_por_pais(ano, pais);
