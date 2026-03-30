/**
 * Tipos para o fluxo de Onboarding (J-AC01)
 *
 * Nota de arquitetura:
 * - OnboardingMetadata é persistido em user.publicMetadata do Clerk (~8KB limite).
 * - Migrar para PostgreSQL próprio em J-AC04 quando precisar de consultas/agregações.
 * - ncmChapter (2 dígitos) é o campo canônico do onboarding.
 *   Para chamar o simulador (que exige 8 dígitos), use ncmDefaultNcm8d da lista de capítulos.
 */

export interface NcmChapter {
  /** Código do capítulo NCM com 2 dígitos (ex: "09") */
  code: string;
  /** Descrição em português brasileiro */
  label: string;
  /**
   * NCM de 8 dígitos representativo do capítulo para pré-carregar o simulador.
   * Escolhido como o NCM de maior volume de exportação dentro do capítulo (Comex Stat 2024).
   */
  defaultNcm8d: string;
}

export interface TradeRegion {
  code: string;
  label: string;
  type: 'continent' | 'country';
  /** Código ISO3 do país (apenas para type === 'country') */
  iso3?: string;
  /** Código do continente pai (apenas para type === 'country') */
  continent?: string;
}

export interface OnboardingWizardState {
  /** Capítulo NCM selecionado (ex: {code: "09", label: "Café..."}) */
  ncmChapter: NcmChapter | null;
  /** Volume médio mensal em kg (opcional) */
  volumeMonthlyKg: number | undefined;
  /** Regiões/países alvo selecionados */
  targetRegions: TradeRegion[];
}

/**
 * Schema dos dados salvos em user.publicMetadata do Clerk.
 * Versão atual: 1 (incrementar em re-onboardings futuros se o schema mudar).
 */
export interface OnboardingMetadata {
  onboardingCompleted: boolean;
  /** Incrementado a cada re-onboarding */
  onboardingVersion: number;
  /** ISO 8601 */
  onboardingCompletedAt: string;
  /** Código do capítulo NCM (2 dígitos, ex: "09") */
  ncmChapter: string;
  /** Descrição do capítulo (ex: "Café, chá, mate e especiarias") */
  ncmLabel: string;
  /**
   * NCM de 8 dígitos para uso direto na API do simulador.
   * Derivado de NcmChapter.defaultNcm8d ao salvar.
   */
  ncmDefaultNcm8d: string;
  volumeMonthlyKg?: number;
  /** Códigos de continentes selecionados (ex: ["europa", "america_norte"]) */
  targetContinents: string[];
  /** Códigos ISO3 de países selecionados (ex: ["DEU", "USA", "CHN"]) */
  targetCountries: string[];
  /** true se o usuário clicou em "Pular" em vez de completar */
  onboardingSkipped: boolean;
}

/** Chave usada para flag de tutorial já visualizado no localStorage */
export const TUTORIAL_SEEN_KEY = 'bgc_tutorial_seen';

/** Chave usada para flag de banner de lembrete dispensado no localStorage */
export const REMINDER_BANNER_DISMISSED_KEY = 'bgc_profile_banner_dismissed';
