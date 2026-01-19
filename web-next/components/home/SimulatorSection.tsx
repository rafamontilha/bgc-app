'use client';

/**
 * Simulator Section Component
 * Export destination simulator (appears after scroll from hero)
 * Matches original globalconnect-bridge structure
 */

import React, { useState, useCallback, useRef } from 'react';
import { Box, Container, Typography, Stack, Alert, AlertTitle } from '@mui/material';
import { SimulatorForm } from '@/components/simulator/SimulatorForm';
import { DestinationList } from '@/components/simulator/DestinationList';
import { RateLimitBanner } from '@/components/simulator/RateLimitBanner';
import { UpgradeModal } from '@/components/simulator/UpgradeModal';
import { ErrorState } from '@/components/simulator/ErrorState';
import { simulateDestinations, SimulatorApiError } from '@/lib/api/simulator';
import {
  SimulatorRequest,
  SimulatorResponse,
  RateLimitInfo,
} from '@/types/simulator';

export interface SimulatorSectionProps {
  onUpgradeClick?: () => void;
}

export function SimulatorSection({
  onUpgradeClick,
}: SimulatorSectionProps): React.ReactElement {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<{
    type: 'not_found' | 'validation' | 'server_error' | 'network_error' | 'unknown';
    message: string;
  } | null>(null);
  const [data, setData] = useState<SimulatorResponse | null>(null);
  const [rateLimitInfo, setRateLimitInfo] = useState<RateLimitInfo | null>(null);
  const [showUpgradeModal, setShowUpgradeModal] = useState(false);

  const resultsRef = useRef<HTMLDivElement>(null);

  const handleSubmit = useCallback(async (request: SimulatorRequest) => {
    setIsLoading(true);
    setError(null);
    setData(null);

    try {
      const response = await simulateDestinations(request);
      setData(response.data);
      setRateLimitInfo(response.rateLimitInfo);

      // Smooth scroll to results after a short delay
      setTimeout(() => {
        resultsRef.current?.scrollIntoView({
          behavior: 'smooth',
          block: 'start',
        });
      }, 100);
    } catch (err) {
      if (err instanceof SimulatorApiError) {
        // Handle rate limit exceeded
        if (err.statusCode === 429) {
          setShowUpgradeModal(true);
          setError({
            type: 'validation',
            message: err.message,
          });
        } else if (err.statusCode === 404) {
          setError({
            type: 'not_found',
            message: err.message,
          });
        } else if (err.statusCode === 400) {
          setError({
            type: 'validation',
            message: err.message,
          });
        } else if (err.statusCode >= 500) {
          setError({
            type: 'server_error',
            message: err.message,
          });
        } else if (err.statusCode === 0) {
          setError({
            type: 'network_error',
            message: err.message,
          });
        } else {
          setError({
            type: 'unknown',
            message: err.message,
          });
        }

        // Update rate limit info even on error
        if (err.apiError) {
          const headers = err.apiError as { headers?: Headers };
          if (headers.headers) {
            // Extract rate limit info if available
            const limit = headers.headers.get('X-RateLimit-Limit');
            const remaining = headers.headers.get('X-RateLimit-Remaining');
            const reset = headers.headers.get('X-RateLimit-Reset');
            if (limit && remaining && reset) {
              setRateLimitInfo({
                limit: parseInt(limit, 10),
                remaining: parseInt(remaining, 10),
                reset: parseInt(reset, 10),
              });
            }
          }
        }
      } else {
        setError({
          type: 'unknown',
          message: 'Erro desconhecido. Tente novamente.',
        });
      }
    } finally {
      setIsLoading(false);
    }
  }, []);

  const handleUpgradeClick = useCallback(() => {
    setShowUpgradeModal(false);
    if (onUpgradeClick) {
      onUpgradeClick();
    } else {
      // Default: scroll to pricing section
      const pricingSection = document.getElementById('pricing');
      if (pricingSection) {
        pricingSection.scrollIntoView({ behavior: 'smooth' });
      }
    }
  }, [onUpgradeClick]);

  const handleRetry = useCallback(() => {
    setError(null);
    setData(null);
  }, []);

  return (
    <Box
      id="simulator"
      component="section"
      sx={{
        backgroundColor: 'background.paper',
        pt: { xs: 8, md: 12 },
        pb: { xs: 8, md: 12 },
      }}
    >
      <Container maxWidth="lg">
        <Stack spacing={6}>
          {/* Section Description */}
          <Box>
            <Typography
              variant="h6"
              sx={{
                fontWeight: 400,
                fontSize: { xs: '1rem', sm: '1.125rem' },
                color: 'text.secondary',
                lineHeight: 1.6,
                textAlign: 'left',
              }}
            >
              Descubra os melhores mercados para seus produtos em segundos
            </Typography>
          </Box>

          {/* Coverage Info Banner */}
          <Alert
            severity="info"
            sx={{
              maxWidth: 800,
              mx: 'auto',
              backgroundColor: 'rgba(0, 122, 255, 0.08)',
              borderColor: 'primary.main',
              '& .MuiAlert-icon': {
                color: 'primary.main',
              },
            }}
          >
            <AlertTitle sx={{ fontWeight: 600 }}>Cobertura de Dados - v0.4.0</AlertTitle>
            Atualmente simulamos destinos para <strong>3 produtos</strong> dos setores: <strong>Sementes de Soja (12)</strong>, <strong>Açúcar (17)</strong> e <strong>Minério de Ferro (26)</strong>.
            <br />
            <Typography variant="caption" sx={{ mt: 0.5, display: 'block' }}>
              💡 Experimente: 12010090 (Soja), 17011400 (Açúcar) ou 26011200 (Minério)
            </Typography>
          </Alert>

          {/* Rate limit banner */}
          {rateLimitInfo && (
            <Box sx={{ width: '100%', maxWidth: 800, mx: 'auto' }}>
              <RateLimitBanner
                rateLimitInfo={rateLimitInfo}
                onUpgradeClick={handleUpgradeClick}
              />
            </Box>
          )}

          {/* Simulator form */}
          <Box sx={{ width: '100%' }}>
            <SimulatorForm onSubmit={handleSubmit} isLoading={isLoading} />
          </Box>

          {/* Results section */}
          <Box ref={resultsRef} sx={{ width: '100%' }}>
            {error && (
              <ErrorState
                type={error.type}
                message={error.message}
                onRetry={handleRetry}
                showRetryButton={error.type !== 'validation'}
              />
            )}

            {data && !error && (
              <DestinationList
                destinations={data.destinations}
                isLoading={isLoading}
              />
            )}
          </Box>
        </Stack>
      </Container>

      {/* Upgrade modal */}
      <UpgradeModal
        open={showUpgradeModal}
        onClose={() => setShowUpgradeModal(false)}
        onViewPlans={handleUpgradeClick}
      />
    </Box>
  );
}
