/**
 * API Types - BGC App
 */

export type Metric = 'TAM' | 'SAM' | 'SOM';
export type Scenario = 'base' | 'aggressive';

export interface MarketSizeItem {
  ncm_chapter: string;
  ano: number;
  valor_usd: number;
}

export interface MarketSizeResponse {
  metric: string;
  year_from: number;
  year_to: number;
  ncm_chapter?: string;
  scenario?: string;
  items: MarketSizeItem[];
}

export interface MarketSizeParams {
  metric: Metric;
  year_from: number;
  year_to: number;
  ncm_chapter?: string;
  scenario?: Scenario;
}

export interface RouteComparisonResult {
  partner: string;
  share: number;
  factor: number;
  estimated_usd: number;
}

export interface RouteComparisonResponse {
  year: number;
  ncm_chapter: string;
  tariff_scenario: string;
  tam_total_usd: number;
  adjusted_total_usd: number;
  results: RouteComparisonResult[];
}

export interface RouteComparisonParams {
  year: number;
  ncm_chapter: string;
  from: string;
  alts: string;
  tariff_scenario: string;
}

export interface HealthzResponse {
  status: string;
  timestamp: string;
  available_scenarios?: string[];
}

export interface AggregatedYearData {
  ano: number;
  total: number;
}
