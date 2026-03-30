/**
 * TDD — lib/data/trade-regions.ts
 *
 * Valida a estrutura dos dados de regiões usados no Step 3 do onboarding.
 */

import {
  TRADE_REGIONS,
  CONTINENTS,
  COUNTRIES,
  getCountriesByContinent,
} from '@/lib/data/trade-regions';

describe('TRADE_REGIONS', () => {
  it('contém 7 continentes', () => {
    const continents = TRADE_REGIONS.filter((r) => r.type === 'continent');
    expect(continents).toHaveLength(7);
  });

  it('contém pelo menos 30 países', () => {
    const countries = TRADE_REGIONS.filter((r) => r.type === 'country');
    expect(countries.length).toBeGreaterThanOrEqual(30);
  });

  it('todos os países têm iso3 de 3 letras maiúsculas', () => {
    const countries = TRADE_REGIONS.filter((r) => r.type === 'country');
    for (const c of countries) {
      expect(c.iso3).toMatch(/^[A-Z]{3}$/);
    }
  });

  it('todos os países têm continent definido', () => {
    const countries = TRADE_REGIONS.filter((r) => r.type === 'country');
    for (const c of countries) {
      expect(c.continent).toBeDefined();
      expect(c.continent!.length).toBeGreaterThan(0);
    }
  });

  it('continentes não têm iso3 nem continent', () => {
    const continents = TRADE_REGIONS.filter((r) => r.type === 'continent');
    for (const c of continents) {
      expect(c.iso3).toBeUndefined();
      expect(c.continent).toBeUndefined();
    }
  });

  it('codes são únicos em toda a lista', () => {
    const codes = TRADE_REGIONS.map((r) => r.code);
    const unique = new Set(codes);
    expect(unique.size).toBe(codes.length);
  });
});

describe('CONTINENTS', () => {
  it('é subconjunto de TRADE_REGIONS com type=continent', () => {
    expect(CONTINENTS.every((c) => c.type === 'continent')).toBe(true);
  });

  it('contém América do Norte, América do Sul, Europa, Ásia, Oriente Médio, África, Oceania', () => {
    const codes = CONTINENTS.map((c) => c.code);
    expect(codes).toContain('america_norte');
    expect(codes).toContain('america_sul');
    expect(codes).toContain('europa');
    expect(codes).toContain('asia');
    expect(codes).toContain('oriente_medio');
    expect(codes).toContain('africa');
    expect(codes).toContain('oceania');
  });
});

describe('COUNTRIES', () => {
  it('todos têm type=country', () => {
    expect(COUNTRIES.every((c) => c.type === 'country')).toBe(true);
  });

  it('inclui principais parceiros do Brasil: China, EUA, Argentina', () => {
    const isos = COUNTRIES.map((c) => c.iso3);
    expect(isos).toContain('CHN');
    expect(isos).toContain('USA');
    expect(isos).toContain('ARG');
  });
});

describe('getCountriesByContinent', () => {
  it('retorna apenas países do continente solicitado', () => {
    const result = getCountriesByContinent('europa');
    expect(result.length).toBeGreaterThan(0);
    expect(result.every((r) => r.continent === 'europa')).toBe(true);
    expect(result.every((r) => r.type === 'country')).toBe(true);
  });

  it('Europa inclui Alemanha, França, Portugal', () => {
    const europa = getCountriesByContinent('europa');
    const isos = europa.map((c) => c.iso3);
    expect(isos).toContain('DEU');
    expect(isos).toContain('FRA');
    expect(isos).toContain('PRT');
  });

  it('retorna lista vazia para continente inexistente', () => {
    expect(getCountriesByContinent('antartida')).toHaveLength(0);
  });

  it('não retorna continentes na lista de países', () => {
    const result = getCountriesByContinent('asia');
    expect(result.every((r) => r.type !== 'continent')).toBe(true);
  });
});
