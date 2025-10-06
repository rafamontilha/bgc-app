--

CREATE INDEX IF NOT EXISTS idx_exportacao_ano ON stg.exportacao(ano);
CREATE INDEX IF NOT EXISTS idx_exportacao_setor ON stg.exportacao(setor);
CREATE INDEX IF NOT EXISTS idx_exportacao_pais ON stg.exportacao(pais);
CREATE INDEX IF NOT EXISTS idx_exportacao_ncm ON stg.exportacao(ncm);
CREATE INDEX IF NOT EXISTS idx_exportacao_ano_setor ON stg.exportacao(ano, setor);
CREATE INDEX IF NOT EXISTS idx_exportacao_ano_pais ON stg.exportacao(ano, pais);

CREATE INDEX IF NOT EXISTS idx_importacao_ano ON stg.importacao(ano);
CREATE INDEX IF NOT EXISTS idx_importacao_setor ON stg.importacao(setor);
CREATE INDEX IF NOT EXISTS idx_importacao_pais ON stg.importacao(pais);
CREATE INDEX IF NOT EXISTS idx_importacao_ncm ON stg.importacao(ncm);
CREATE INDEX IF NOT EXISTS idx_importacao_ano_setor ON stg.importacao(ano, setor);
CREATE INDEX IF NOT EXISTS idx_importacao_ano_pais ON stg.importacao(ano, pais);
