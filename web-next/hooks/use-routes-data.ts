/**
 * Custom hook para dados de Routes Comparison com SWR
 */

import useSWR from 'swr';
import { apiClient } from '@/lib/api-client';
import type { RouteComparisonParams, RouteComparisonResponse } from '@/types/api';

export function useRoutesData(params: RouteComparisonParams | null) {
  const key = params ? ['routes/compare', params] : null;

  const { data, error, isLoading, mutate } = useSWR(
    key,
    async ([_, p]) => {
      return apiClient.routes.compare(p as RouteComparisonParams);
    },
    {
      revalidateOnFocus: false,
      shouldRetryOnError: false,
    }
  );

  return {
    data,
    error,
    isLoading,
    refetch: mutate,
  };
}
