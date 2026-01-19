-- Migration 0012: Update exportacao dates to 2025-2026
-- Adds recent data (last 12 months) for testing

-- Update existing 2024 data to 2025
UPDATE stg.exportacao SET co_ano = 2025 WHERE co_ano = 2024 AND co_mes = 8;
UPDATE stg.exportacao SET co_ano = 2025 WHERE co_ano = 2024 AND co_mes = 9;
UPDATE stg.exportacao SET co_ano = 2025 WHERE co_ano = 2024 AND co_mes = 10;
UPDATE stg.exportacao SET co_ano = 2025 WHERE co_ano = 2024 AND co_mes = 11;

-- Add December 2025 data (1 month ago)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
-- NCM 17011400 (Sugar)
(2025, 12, 'CN', '17011400', 16000000, 52000000),
(2025, 12, 'IN', '17011400', 8500000, 29000000),
(2025, 12, 'AE', '17011400', 5200000, 18500000),
(2025, 12, 'BD', '17011400', 3600000, 12500000),
(2025, 12, 'US', '17011400', 2600000, 8800000),
(2025, 12, 'MX', '17011400', 1800000, 6500000),

-- NCM 26011200 (Iron Ore)
(2025, 12, 'CN', '26011200', 85000000, 4500000),
(2025, 12, 'DE', '26011200', 25000000, 1200000),
(2025, 12, 'JP', '26011200', 32000000, 1500000),
(2025, 12, 'NL', '26011200', 18000000, 850000),

-- NCM 12010090 (Soybean)
(2025, 12, 'CN', '12010090', 150000000, 6500000),
(2025, 12, 'AR', '12010090', 45000000, 1800000),
(2025, 12, 'ES', '12010090', 28000000, 1100000),
(2025, 12, 'TH', '12010090', 22000000, 950000),
(2025, 12, 'VN', '12010090', 18000000, 800000),
(2025, 12, 'IR', '12010090', 15000000, 650000),
(2025, 12, 'CL', '12010090', 12000000, 520000);

-- Add January 2026 data (current month - partial)
INSERT INTO stg.exportacao (co_ano, co_mes, co_pais, co_ncm, vl_fob, kg_liquido) VALUES
-- NCM 17011400 (Sugar) - partial month
(2026, 1, 'CN', '17011400', 5000000, 17000000),
(2026, 1, 'IN', '17011400', 2800000, 9500000),
(2026, 1, 'AE', '17011400', 1700000, 6000000),
(2026, 1, 'BD', '17011400', 1200000, 4100000),
(2026, 1, 'US', '17011400', 850000, 2900000);

COMMENT ON TABLE stg.exportacao IS 'Brazilian export data from ComexStat (updated to 2025-2026)';
