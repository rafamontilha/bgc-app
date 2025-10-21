/**
 * Custom hook para dados de Market Size com SWR
 */

import useSWR from 'swr';
import { apiClient } from '@/lib/api-client';
import type { MarketSizeParams, MarketSizeResponse } from '@/types/api';

export function useMarketData(params: MarketSizeParams | null) {
  const key = params ? ['market/size', params] : null;

  const { data, error, isLoading, mutate } = useSWR(
    key,
    async ([_, p]) => {
      return apiClient.market.getSize(p as MarketSizeParams);
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
