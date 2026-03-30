import type { TradeRegion } from '@/lib/types/onboarding';

/**
 * Regiões e países para o passo 3 do onboarding.
 *
 * Países: top parceiros comerciais do Brasil por volume de exportação (Comex Stat 2024).
 * Limitado a ~5 países por continente para manter o autocomplete usável.
 * Códigos ISO3 conforme padrão da ONU.
 */
export const TRADE_REGIONS: TradeRegion[] = [
  // ─────────────────────────────────────────────
  // Continentes
  // ─────────────────────────────────────────────
  { code: 'america_norte', label: 'América do Norte', type: 'continent' },
  { code: 'america_sul',   label: 'América do Sul',   type: 'continent' },
  { code: 'europa',        label: 'Europa',           type: 'continent' },
  { code: 'asia',          label: 'Ásia',             type: 'continent' },
  { code: 'oriente_medio', label: 'Oriente Médio',    type: 'continent' },
  { code: 'africa',        label: 'África',           type: 'continent' },
  { code: 'oceania',       label: 'Oceania',          type: 'continent' },

  // ─────────────────────────────────────────────
  // América do Norte
  // ─────────────────────────────────────────────
  { code: 'USA', label: 'Estados Unidos', type: 'country', iso3: 'USA', continent: 'america_norte' },
  { code: 'CAN', label: 'Canadá',         type: 'country', iso3: 'CAN', continent: 'america_norte' },
  { code: 'MEX', label: 'México',         type: 'country', iso3: 'MEX', continent: 'america_norte' },

  // ─────────────────────────────────────────────
  // América do Sul
  // ─────────────────────────────────────────────
  { code: 'ARG', label: 'Argentina', type: 'country', iso3: 'ARG', continent: 'america_sul' },
  { code: 'CHL', label: 'Chile',     type: 'country', iso3: 'CHL', continent: 'america_sul' },
  { code: 'COL', label: 'Colômbia',  type: 'country', iso3: 'COL', continent: 'america_sul' },
  { code: 'PER', label: 'Peru',      type: 'country', iso3: 'PER', continent: 'america_sul' },
  { code: 'URY', label: 'Uruguai',   type: 'country', iso3: 'URY', continent: 'america_sul' },
  { code: 'PRY', label: 'Paraguai',  type: 'country', iso3: 'PRY', continent: 'america_sul' },

  // ─────────────────────────────────────────────
  // Europa
  // ─────────────────────────────────────────────
  { code: 'DEU', label: 'Alemanha',    type: 'country', iso3: 'DEU', continent: 'europa' },
  { code: 'NLD', label: 'Holanda',     type: 'country', iso3: 'NLD', continent: 'europa' },
  { code: 'ESP', label: 'Espanha',     type: 'country', iso3: 'ESP', continent: 'europa' },
  { code: 'ITA', label: 'Itália',      type: 'country', iso3: 'ITA', continent: 'europa' },
  { code: 'FRA', label: 'França',      type: 'country', iso3: 'FRA', continent: 'europa' },
  { code: 'BEL', label: 'Bélgica',     type: 'country', iso3: 'BEL', continent: 'europa' },
  { code: 'GBR', label: 'Reino Unido', type: 'country', iso3: 'GBR', continent: 'europa' },
  { code: 'PRT', label: 'Portugal',    type: 'country', iso3: 'PRT', continent: 'europa' },

  // ─────────────────────────────────────────────
  // Ásia
  // ─────────────────────────────────────────────
  { code: 'CHN', label: 'China',        type: 'country', iso3: 'CHN', continent: 'asia' },
  { code: 'JPN', label: 'Japão',        type: 'country', iso3: 'JPN', continent: 'asia' },
  { code: 'KOR', label: 'Coreia do Sul', type: 'country', iso3: 'KOR', continent: 'asia' },
  { code: 'IND', label: 'Índia',        type: 'country', iso3: 'IND', continent: 'asia' },
  { code: 'IDN', label: 'Indonésia',    type: 'country', iso3: 'IDN', continent: 'asia' },
  { code: 'VNM', label: 'Vietnã',       type: 'country', iso3: 'VNM', continent: 'asia' },
  { code: 'SGP', label: 'Singapura',    type: 'country', iso3: 'SGP', continent: 'asia' },

  // ─────────────────────────────────────────────
  // Oriente Médio
  // ─────────────────────────────────────────────
  { code: 'SAU', label: 'Arábia Saudita',    type: 'country', iso3: 'SAU', continent: 'oriente_medio' },
  { code: 'ARE', label: 'Emirados Árabes',   type: 'country', iso3: 'ARE', continent: 'oriente_medio' },
  { code: 'TUR', label: 'Turquia',           type: 'country', iso3: 'TUR', continent: 'oriente_medio' },
  { code: 'EGY', label: 'Egito',             type: 'country', iso3: 'EGY', continent: 'oriente_medio' },
  { code: 'IRN', label: 'Irã',               type: 'country', iso3: 'IRN', continent: 'oriente_medio' },

  // ─────────────────────────────────────────────
  // África
  // ─────────────────────────────────────────────
  { code: 'ZAF', label: 'África do Sul', type: 'country', iso3: 'ZAF', continent: 'africa' },
  { code: 'NGA', label: 'Nigéria',       type: 'country', iso3: 'NGA', continent: 'africa' },
  { code: 'AGO', label: 'Angola',        type: 'country', iso3: 'AGO', continent: 'africa' },
  { code: 'MOZ', label: 'Moçambique',    type: 'country', iso3: 'MOZ', continent: 'africa' },
  { code: 'DZA', label: 'Argélia',       type: 'country', iso3: 'DZA', continent: 'africa' },

  // ─────────────────────────────────────────────
  // Oceania
  // ─────────────────────────────────────────────
  { code: 'AUS', label: 'Austrália',   type: 'country', iso3: 'AUS', continent: 'oceania' },
  { code: 'NZL', label: 'Nova Zelândia', type: 'country', iso3: 'NZL', continent: 'oceania' },
];

/** Retorna apenas continentes */
export const CONTINENTS = TRADE_REGIONS.filter((r) => r.type === 'continent');

/** Retorna países de um continente específico */
export function getCountriesByContinent(continentCode: string): TradeRegion[] {
  return TRADE_REGIONS.filter(
    (r) => r.type === 'country' && r.continent === continentCode
  );
}

/** Retorna apenas países */
export const COUNTRIES = TRADE_REGIONS.filter((r) => r.type === 'country');
