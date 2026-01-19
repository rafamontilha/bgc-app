/**
 * TypeScript Types for Export Destination Simulator
 * Backend API: POST /v1/simulator/destinations
 */

/**
 * Request payload for simulator endpoint
 * Backend expects snake_case field names
 */
export interface SimulatorRequest {
  ncm: string;
  volume_kg?: number;
  countries?: string[];
  max_results?: number;
  include_all?: boolean;
}

/**
 * Single destination result from the simulator
 * Backend returns snake_case field names
 */
export interface Destination {
  rank: number;
  country_code: string;
  country_name: string;
  score: number;
  demand: 'Alto' | 'Médio' | 'Baixo';
  market_size_usd: number;
  growth_rate_pct: number;
  price_per_kg_usd: number;
  distance_km: number;
  estimated_margin_pct: number;
  logistics_cost_usd: number;
  tariff_rate_pct: number;
  lead_time_days: number;
  recommendation_reason: string;
  region: string;
  flag_emoji: string;
}

/**
 * Metadata about the simulation analysis
 */
export interface SimulatorMetadata {
  ncm: string;
  product_name: string;
  analysis_date: string;
  total_destinations: number;
  cache_hit: boolean;
  processing_time_ms: number;
}

/**
 * Complete response from simulator endpoint
 */
export interface SimulatorResponse {
  destinations: Destination[];
  metadata: SimulatorMetadata;
}

/**
 * Rate limit information from response headers
 */
export interface RateLimitInfo {
  limit: number;
  remaining: number;
  reset: number;
}

/**
 * Error response from backend API
 */
export interface ApiError {
  error: string;
  message: string;
  status_code: number;
  details?: Record<string, unknown>;
}

/**
 * Demand level configuration for visual representation
 */
export interface DemandLevelConfig {
  label: 'Alto' | 'Médio' | 'Baixo';
  color: 'success' | 'warning' | 'error';
  bgColor: string;
}

/**
 * Form validation state
 */
export interface FormErrors {
  ncm?: string;
  volume_kg?: string;
  countries?: string;
  max_results?: string;
}

/**
 * UI state for simulator component
 */
export interface SimulatorState {
  isLoading: boolean;
  error: string | null;
  data: SimulatorResponse | null;
  rateLimitInfo: RateLimitInfo | null;
}
