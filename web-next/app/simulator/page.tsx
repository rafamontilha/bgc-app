'use client';

/**
 * Export Destination Simulator Page
 * Main page that orchestrates all simulator components
 */

import React, { useState } from 'react';
import { Box, Container, Stack, Divider } from '@mui/material';
import {
  SimulatorRequest,
  SimulatorResponse,
  RateLimitInfo,
} from '@/types/simulator';
import {
  simulateDestinations,
  SimulatorApiError,
} from '@/lib/api/simulator';
import { SimulatorForm } from '@/components/simulator/SimulatorForm';
import { DestinationList } from '@/components/simulator/DestinationList';
import { RateLimitBanner } from '@/components/simulator/RateLimitBanner';
import { UpgradeModal } from '@/components/simulator/UpgradeModal';
import { ErrorState, ErrorType } from '@/components/simulator/ErrorState';

interface SimulatorState {
  isLoading: boolean;
  data: SimulatorResponse | null;
  error: { type: ErrorType; message: string } | null;
  rateLimitInfo: RateLimitInfo | null;
  showUpgradeModal: boolean;
}

export default function SimulatorPage(): React.ReactElement {
  const [state, setState] = useState<SimulatorState>({
    isLoading: false,
    data: null,
    error: null,
    rateLimitInfo: null,
    showUpgradeModal: false,
  });

  /**
   * Handle form submission and API call
   */
  const handleSimulate = async (request: SimulatorRequest): Promise<void> => {
    setState((prev) => ({
      ...prev,
      isLoading: true,
      error: null,
    }));

    try {
      const response = await simulateDestinations(request);

      setState((prev) => ({
        ...prev,
        isLoading: false,
        data: response.data,
        rateLimitInfo: response.rateLimitInfo,
        error: null,
      }));
    } catch (error) {
      if (error instanceof SimulatorApiError) {
        // Handle rate limit exceeded (HTTP 429)
        if (error.statusCode === 429) {
          setState((prev) => ({
            ...prev,
            isLoading: false,
            showUpgradeModal: true,
            error: null,
          }));
          return;
        }

        // Determine error type based on status code
        let errorType: ErrorType = 'unknown';
        if (error.statusCode === 404) {
          errorType = 'not_found';
        } else if (error.statusCode === 400) {
          errorType = 'validation';
        } else if (error.statusCode >= 500) {
          errorType = 'server_error';
        } else if (error.statusCode === 0) {
          errorType = 'network_error';
        }

        setState((prev) => ({
          ...prev,
          isLoading: false,
          error: {
            type: errorType,
            message: error.message,
          },
        }));
      } else {
        // Unknown error
        setState((prev) => ({
          ...prev,
          isLoading: false,
          error: {
            type: 'unknown',
            message: 'Erro desconhecido. Tente novamente.',
          },
        }));
      }
    }
  };

  /**
   * Handle retry after error
   */
  const handleRetry = (): void => {
    setState((prev) => ({
      ...prev,
      error: null,
    }));
  };

  /**
   * Handle upgrade modal close
   */
  const handleCloseUpgradeModal = (): void => {
    setState((prev) => ({
      ...prev,
      showUpgradeModal: false,
    }));
  };

  /**
   * Handle upgrade banner click
   */
  const handleUpgradeClick = (): void => {
    setState((prev) => ({
      ...prev,
      showUpgradeModal: true,
    }));
  };

  return (
    <Box
      sx={{
        minHeight: '100vh',
        backgroundColor: 'background.default',
        py: { xs: 4, sm: 6, md: 8 },
      }}
    >
      <Container maxWidth="lg">
        <Stack spacing={4}>
          {/* Rate Limit Banner */}
          {state.rateLimitInfo && (
            <RateLimitBanner
              rateLimitInfo={state.rateLimitInfo}
              onUpgradeClick={handleUpgradeClick}
            />
          )}

          {/* Simulator Form */}
          <SimulatorForm
            onSubmit={handleSimulate}
            isLoading={state.isLoading}
          />

          {/* Divider between form and results */}
          {(state.data || state.error || state.isLoading) && (
            <Divider sx={{ my: 4 }} />
          )}

          {/* Error State */}
          {state.error && !state.isLoading && (
            <ErrorState
              type={state.error.type}
              message={state.error.message}
              onRetry={handleRetry}
            />
          )}

          {/* Results List */}
          {!state.error && (state.data || state.isLoading) && (
            <DestinationList
              destinations={state.data?.destinations || []}
              isLoading={state.isLoading}
            />
          )}

          {/* Metadata (optional debug info) */}
          {state.data?.metadata && !state.isLoading && (
            <Box
              sx={{
                textAlign: 'center',
                color: 'text.disabled',
                fontSize: '0.75rem',
                pt: 2,
              }}
            >
              Análise realizada em{' '}
              {new Date(state.data.metadata.analysis_date).toLocaleString('pt-BR')} •
              Tempo de processamento: {state.data.metadata.processing_time_ms}ms
              {state.data.metadata.cache_hit && ' • Cache'}
            </Box>
          )}
        </Stack>
      </Container>

      {/* Upgrade Modal */}
      <UpgradeModal
        open={state.showUpgradeModal}
        onClose={handleCloseUpgradeModal}
      />
    </Box>
  );
}
