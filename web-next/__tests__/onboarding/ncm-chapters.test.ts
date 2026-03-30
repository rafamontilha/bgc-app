/**
 * TDD — lib/data/ncm-chapters.ts
 *
 * Valida integridade dos dados NCM usados no onboarding e no simulador.
 * Erros aqui quebram o pré-carregamento automático do simulador no dashboard.
 */

import { NCM_CHAPTERS, getNcmChapter } from '@/lib/data/ncm-chapters';

describe('NCM_CHAPTERS', () => {
  it('contém exatamente 96 capítulos ativos (capítulo 77 é reservado)', () => {
    expect(NCM_CHAPTERS).toHaveLength(96);
  });

  it('não contém o capítulo 77 (reservado pela ONU)', () => {
    const ch77 = NCM_CHAPTERS.find((c) => c.code === '77');
    expect(ch77).toBeUndefined();
  });

  it('todos os codes têm exatamente 2 dígitos com zero à esquerda', () => {
    for (const ch of NCM_CHAPTERS) {
      expect(ch.code).toMatch(/^\d{2}$/);
    }
  });

  it('codes são únicos', () => {
    const codes = NCM_CHAPTERS.map((c) => c.code);
    const unique = new Set(codes);
    expect(unique.size).toBe(codes.length);
  });

  it('todos os labels são strings não-vazias', () => {
    for (const ch of NCM_CHAPTERS) {
      expect(typeof ch.label).toBe('string');
      expect(ch.label.length).toBeGreaterThan(0);
    }
  });

  it('todos os defaultNcm8d têm exatamente 8 dígitos numéricos', () => {
    for (const ch of NCM_CHAPTERS) {
      expect(ch.defaultNcm8d).toMatch(/^\d{8}$/);
    }
  });

  it('os primeiros 2 dígitos de defaultNcm8d correspondem ao code do capítulo', () => {
    for (const ch of NCM_CHAPTERS) {
      expect(ch.defaultNcm8d.substring(0, 2)).toBe(ch.code);
    }
  });

  it('capítulo 09 (café) tem defaultNcm8d correto', () => {
    const coffee = NCM_CHAPTERS.find((c) => c.code === '09');
    expect(coffee?.defaultNcm8d).toBe('09010110');
  });

  it('capítulo 12 (soja) tem defaultNcm8d correto', () => {
    const soy = NCM_CHAPTERS.find((c) => c.code === '12');
    expect(soy?.defaultNcm8d).toBe('12010090');
  });
});

describe('getNcmChapter', () => {
  it('retorna o capítulo correto pelo código', () => {
    const ch = getNcmChapter('09');
    expect(ch).toBeDefined();
    expect(ch?.code).toBe('09');
    expect(ch?.label).toContain('Café');
  });

  it('retorna undefined para código inexistente', () => {
    expect(getNcmChapter('99')).toBeUndefined();
    expect(getNcmChapter('77')).toBeUndefined();
  });

  it('retorna undefined para string vazia', () => {
    expect(getNcmChapter('')).toBeUndefined();
  });

  it('não faz lookup case-insensitive — código deve ser exato', () => {
    expect(getNcmChapter('9')).toBeUndefined(); // deve usar '09'
  });
});
