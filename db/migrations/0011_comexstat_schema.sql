-- Migration 0011: ComexStat Schema for Exportacao
-- Converts stg.exportacao to ComexStat format and adds sample data

-- Drop old table (simple schema)
DROP TABLE IF EXISTS stg.exportacao;

-- Create exportacao table with ComexStat schema
CREATE TABLE stg.exportacao (
  id BIGSERIAL PRIMARY KEY,

  -- Time dimensions
  co_ano INTEGER NOT NULL,                 -- Year (2020-2025)
  co_mes INTEGER NOT NULL,                 -- Month (1-12)

  -- Geographic dimensions
  co_pais TEXT NOT NULL,                   -- Country code (ISO 3166-1 alpha-2)
  sg_uf_ncm TEXT,                          -- State code (SP, RJ, MG, etc)

  -- Product dimensions
  co_ncm TEXT NOT NULL,                    -- NCM code (8 digits)
  co_sh4 TEXT,                             -- SH4 code (4 digits)
  co_sh2 TEXT,                             -- SH2 code (2 digits)

  -- Trade metrics
  vl_fob NUMERIC(18,2) NOT NULL,           -- FOB value in USD
  kg_liquido NUMERIC(18,2) NOT NULL,       -- Net weight in kg

  -- Additional metrics
  qt_estat INTEGER,                        -- Statistical quantity
  vl_frete NUMERIC(18,2),                  -- Freight value
  vl_seguro NUMERIC(18,2),                 -- Insurance value

  -- Metadata
  created_at TIMESTAMPTZ DEFAULT now(),

  CONSTRAINT valid_year CHECK (co_ano >= 2020 AND co_ano <= 2030),
  CONSTRAINT valid_month CHECK (co_mes >= 1 AND co_mes <= 12),
  CONSTRAINT valid_ncm_length CHECK (LENGTH(co_ncm) = 8),
  CONSTRAINT valid_fob CHECK (vl_fob >= 0),
  CONSTRAINT valid_weight CHECK (kg_liquido >= 0)
);

-- Indexes for performance
CREATE INDEX idx_exportacao_ncm ON stg.exportacao(co_ncm);
CREATE INDEX idx_exportacao_country ON stg.exportacao(co_pais);
CREATE INDEX idx_exportacao_date ON stg.exportacao(co_ano, co_mes);
CREATE INDEX idx_exportacao_ncm_country ON stg.exportacao(co_ncm, co_pais);
CREATE INDEX idx_exportacao_ncm_date ON stg.exportacao(co_ncm, co_ano, co_mes);

-- Composite index for simulator queries
CREATE INDEX idx_exportacao_simulator
  ON stg.exportacao(co_ncm, co_pais, co_ano DESC, co_mes DESC);

COMMENT ON TABLE stg.exportacao IS 'Brazilian export data from ComexStat (staging layer)';

-- ============================================================================
-- SAMPLE DATA FOR TESTING
-- ============================================================================
-- NCM 17011400 - Açúcar de cana em bruto
-- NCM 26011200 - Minério de ferro aglomerado
-- NCM 12010090 - Soja em grão

-- Seed data for NCM 17011400 (Sugar) - Last 12 months
-- China (CN) - Major sugar importer
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'CN', '17011400', 15000000, 50000000),  -- Nov 2024
  (2024, 10, 'CN', '17011400', 14500000, 48000000),  -- Oct 2024
  (2024, 9, 'CN', '17011400', 14000000, 47000000),   -- Sep 2024
  (2024, 8, 'CN', '17011400', 13500000, 45000000);   -- Aug 2024

-- India (IN)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'IN', '17011400', 8000000, 28000000),
  (2024, 10, 'IN', '17011400', 7800000, 27500000),
  (2024, 9, 'IN', '17011400', 7500000, 27000000),
  (2024, 8, 'IN', '17011400', 7200000, 26000000);

-- United Arab Emirates (AE)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'AE', '17011400', 5000000, 18000000),
  (2024, 10, 'AE', '17011400', 4800000, 17500000),
  (2024, 9, 'AE', '17011400', 4600000, 17000000),
  (2024, 8, 'AE', '17011400', 4400000, 16500000);

-- Bangladesh (BD)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'BD', '17011400', 3500000, 12000000),
  (2024, 10, 'BD', '17011400', 3400000, 11800000),
  (2024, 9, 'BD', '17011400', 3300000, 11500000),
  (2024, 8, 'BD', '17011400', 3200000, 11200000);

-- United States (US)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'US', '17011400', 2500000, 8500000),
  (2024, 10, 'US', '17011400', 2400000, 8300000),
  (2024, 9, 'US', '17011400', 2300000, 8000000),
  (2024, 8, 'US', '17011400', 2200000, 7800000);

-- Seed data for NCM 26011200 (Iron Ore) - Last 12 months
-- China (CN) - Largest iron ore importer globally
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'CN', '26011200', 850000000, 12000000000),  -- Nov 2024
  (2024, 10, 'CN', '26011200', 840000000, 11800000000),
  (2024, 9, 'CN', '26011200', 830000000, 11600000000),
  (2024, 8, 'CN', '26011200', 820000000, 11400000000);

-- Japan (JP)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'JP', '26011200', 45000000, 650000000),
  (2024, 10, 'JP', '26011200', 44000000, 640000000),
  (2024, 9, 'JP', '26011200', 43000000, 630000000),
  (2024, 8, 'JP', '26011200', 42000000, 620000000);

-- Germany (DE)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'DE', '26011200', 38000000, 550000000),
  (2024, 10, 'DE', '26011200', 37000000, 540000000),
  (2024, 9, 'DE', '26011200', 36000000, 530000000),
  (2024, 8, 'DE', '26011200', 35000000, 520000000);

-- Netherlands (NL)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'NL', '26011200', 32000000, 480000000),
  (2024, 10, 'NL', '26011200', 31000000, 470000000),
  (2024, 9, 'NL', '26011200', 30000000, 460000000),
  (2024, 8, 'NL', '26011200', 29000000, 450000000);

-- Seed data for NCM 12010090 (Soybeans) - Last 12 months
-- China (CN) - Largest soybean importer
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'CN', '12010090', 950000000, 2800000000),
  (2024, 10, 'CN', '12010090', 940000000, 2750000000),
  (2024, 9, 'CN', '12010090', 930000000, 2700000000),
  (2024, 8, 'CN', '12010090', 920000000, 2650000000);

-- Spain (ES)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'ES', '12010090', 75000000, 220000000),
  (2024, 10, 'ES', '12010090', 74000000, 218000000),
  (2024, 9, 'ES', '12010090', 73000000, 216000000),
  (2024, 8, 'ES', '12010090', 72000000, 214000000);

-- Argentina (AR)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'AR', '12010090', 65000000, 195000000),
  (2024, 10, 'AR', '12010090', 64000000, 193000000),
  (2024, 9, 'AR', '12010090', 63000000, 191000000),
  (2024, 8, 'AR', '12010090', 62000000, 189000000);

-- Thailand (TH)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'TH', '12010090', 55000000, 165000000),
  (2024, 10, 'TH', '12010090', 54000000, 163000000),
  (2024, 9, 'TH', '12010090', 53000000, 161000000),
  (2024, 8, 'TH', '12010090', 52000000, 159000000);

-- Vietnam (VN)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'VN', '12010090', 48000000, 145000000),
  (2024, 10, 'VN', '12010090', 47000000, 143000000),
  (2024, 9, 'VN', '12010090', 46000000, 141000000),
  (2024, 8, 'VN', '12010090', 45000000, 139000000);

-- Iran (IR)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'IR', '12010090', 42000000, 128000000),
  (2024, 10, 'IR', '12010090', 41000000, 126000000),
  (2024, 9, 'IR', '12010090', 40000000, 124000000),
  (2024, 8, 'IR', '12010090', 39000000, 122000000);

-- Add some countries that exist in countries_metadata to ensure joins work
-- Mexico (MX)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'MX', '17011400', 1800000, 6200000),
  (2024, 10, 'MX', '17011400', 1750000, 6100000);

-- Chile (CL)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
  (2024, 11, 'CL', '12010090', 38000000, 115000000),
  (2024, 10, 'CL', '12010090', 37000000, 113000000);

COMMENT ON TABLE stg.exportacao IS 'Brazilian export data from ComexStat (staging) - Sample data for testing';
