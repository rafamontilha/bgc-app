/**
 * API Client for Export Destination Simulator
 * Handles HTTP requests, error handling, and rate limit parsing
 */

import {
  SimulatorRequest,
  SimulatorResponse,
  ApiError,
  RateLimitInfo,
} from '@/types/simulator';

/**
 * Base URL for API (configurable via environment)
 */
const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';

/**
 * API response wrapper including rate limit info
 */
export interface ApiResponse<T> {
  data: T;
  rateLimitInfo: RateLimitInfo | null;
}

/**
 * Custom error class for API errors
 */
export class SimulatorApiError extends Error {
  statusCode: number;
  apiError?: ApiError;

  constructor(message: string, statusCode: number, apiError?: ApiError) {
    super(message);
    this.name = 'SimulatorApiError';
    this.statusCode = statusCode;
    this.apiError = apiError;
  }
}

/**
 * Parse rate limit headers from response
 */
function parseRateLimitHeaders(headers: Headers): RateLimitInfo | null {
  const limit = headers.get('X-RateLimit-Limit');
  const remaining = headers.get('X-RateLimit-Remaining');
  const reset = headers.get('X-RateLimit-Reset');

  if (!limit || !remaining || !reset) {
    return null;
  }

  return {
    limit: parseInt(limit, 10),
    remaining: parseInt(remaining, 10),
    reset: parseInt(reset, 10),
  };
}

/**
 * Simulate export destinations
 *
 * @param request - Simulator request payload
 * @returns Promise with simulator response and rate limit info
 * @throws SimulatorApiError for all error cases
 */
export async function simulateDestinations(
  request: SimulatorRequest
): Promise<ApiResponse<SimulatorResponse>> {
  console.log('[Simulator API] Request:', {
    url: `${API_BASE_URL}/v1/simulator/destinations`,
    request,
  });

  try {
    const response = await fetch(`${API_BASE_URL}/v1/simulator/destinations`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Accept: 'application/json',
      },
      body: JSON.stringify(request),
    });

    console.log('[Simulator API] Response status:', response.status);
    console.log('[Simulator API] Response headers:', {
      contentType: response.headers.get('Content-Type'),
      rateLimit: response.headers.get('X-RateLimit-Limit'),
      remaining: response.headers.get('X-RateLimit-Remaining'),
    });

    const rateLimitInfo = parseRateLimitHeaders(response.headers);

    if (!response.ok) {
      let apiError: ApiError | undefined;

      try {
        apiError = await response.json();
      } catch {
        // If response body is not JSON, create a generic error
        apiError = {
          error: 'API Error',
          message: response.statusText || 'Unknown error occurred',
          status_code: response.status,
        };
      }

      // Handle specific error cases
      if (response.status === 429) {
        throw new SimulatorApiError(
          'Limite de simulações atingido. Faça upgrade para continuar.',
          429,
          apiError
        );
      }

      if (response.status === 404) {
        throw new SimulatorApiError(
          'NCM não encontrado. Verifique o código e tente novamente.',
          404,
          apiError
        );
      }

      if (response.status === 400) {
        throw new SimulatorApiError(
          apiError?.message || 'Dados inválidos. Verifique os campos.',
          400,
          apiError
        );
      }

      if (response.status >= 500) {
        throw new SimulatorApiError(
          'Erro no servidor. Tente novamente em alguns instantes.',
          response.status,
          apiError
        );
      }

      throw new SimulatorApiError(
        apiError?.message || 'Erro desconhecido',
        response.status,
        apiError
      );
    }

    const data: SimulatorResponse = await response.json();

    console.log('[Simulator API] Success:', {
      destinations: data.destinations.length,
      processingTime: data.metadata?.processing_time_ms,
    });

    return {
      data,
      rateLimitInfo,
    };
  } catch (error) {
    console.error('[Simulator API] Error:', error);

    // Network errors or other exceptions
    if (error instanceof SimulatorApiError) {
      throw error;
    }

    throw new SimulatorApiError(
      'Erro de conexão. Verifique sua internet e tente novamente.',
      0
    );
  }
}

/**
 * Validate NCM format (8 digits)
 */
export function validateNCM(ncm: string): boolean {
  return /^\d{8}$/.test(ncm);
}

/**
 * Validate volume (positive number)
 */
export function validateVolume(volume: number): boolean {
  return volume > 0 && volume <= 1000000;
}

/**
 * Validate country codes (2-letter ISO codes)
 */
export function validateCountryCodes(codes: string[]): boolean {
  return codes.every((code) => /^[A-Z]{2}$/.test(code));
}

/**
 * Validate max results (between 1 and 50)
 */
export function validateMaxResults(maxResults: number): boolean {
  return maxResults >= 1 && maxResults <= 50;
}
