/**
 * Custom hook para carregar cenários disponíveis
 */

import useSWR from 'swr';
import { apiClient } from '@/lib/api-client';

export function useScenarios() {
  const { data, error, isLoading } = useSWR(
    'healthz',
    async () => {
      const response = await apiClient.healthz();
      return response.available_scenarios || ['base'];
    },
    {
      revalidateOnFocus: false,
      fallbackData: ['base'],
    }
  );

  return {
    scenarios: data || ['base'],
    error,
    isLoading,
  };
}
