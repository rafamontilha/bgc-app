/**
 * API Client - Environment-aware
 *
 * Funciona em ambos os ambientes:
 * - Localhost dev: usa localhost:8080 via Next.js rewrites
 * - Kubernetes: usa API interna via rewrites (bgc-api:8080)
 */

export const apiClient = {
  /**
   * GET request genérico
   */
  async get<T>(path: string, params?: Record<string, string | number>): Promise<T> {
    // No navegador: usa origin atual (nginx faz proxy)
    // No servidor: usa variável de ambiente ou localhost
    const baseUrl = typeof window !== 'undefined'
      ? window.location.origin
      : (process.env.API_URL || 'http://localhost:3000');

    const url = new URL(path, baseUrl);

    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        url.searchParams.set(key, String(value));
      });
    }

    const response = await fetch(url.toString(), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const text = await response.text();
      throw new Error(`HTTP ${response.status}: ${text}`);
    }

    return response.json();
  },

  /**
   * Market Size API
   */
  market: {
    async getSize(params: {
      metric: 'TAM' | 'SAM' | 'SOM';
      year_from: number;
      year_to: number;
      ncm_chapter?: string;
      scenario?: 'base' | 'aggressive';
    }) {
      const queryParams: Record<string, string | number> = {
        metric: params.metric,
        year_from: params.year_from,
        year_to: params.year_to,
      };

      if (params.ncm_chapter) {
        queryParams.ncm_chapter = params.ncm_chapter;
      }

      if (params.metric === 'SOM' && params.scenario) {
        queryParams.scenario = params.scenario;
      }

      return apiClient.get<{
        metric: string;
        year_from: number;
        year_to: number;
        ncm_chapter?: string;
        scenario?: string;
        items: Array<{
          ncm_chapter: string;
          ano: number;
          valor_usd: number;
        }>;
      }>('/market/size', queryParams);
    },
  },

  /**
   * Routes Comparison API
   */
  routes: {
    async compare(params: {
      year: number;
      ncm_chapter: string;
      from: string;
      alts: string;
      tariff_scenario: string;
    }) {
      return apiClient.get<{
        year: number;
        ncm_chapter: string;
        tariff_scenario: string;
        tam_total_usd: number;
        adjusted_total_usd: number;
        results: Array<{
          partner: string;
          share: number;
          factor: number;
          estimated_usd: number;
        }>;
      }>('/routes/compare', params);
    },
  },

  /**
   * Healthz API
   */
  async healthz() {
    return apiClient.get<{
      status: string;
      timestamp: string;
      available_scenarios?: string[];
    }>('/healthz');
  },
};
