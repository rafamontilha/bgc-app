'use client';

/**
 * US-002: Destination List Component
 * Container for displaying multiple destination cards with loading and empty states
 */

import React from 'react';
import { Box, Typography, Stack, Skeleton } from '@mui/material';
import PublicOffIcon from '@mui/icons-material/PublicOff';
import { Destination } from '@/types/simulator';
import { DestinationCard } from './DestinationCard';

export interface DestinationListProps {
  destinations: Destination[];
  isLoading?: boolean;
}

/**
 * Loading skeleton for destination cards
 */
function DestinationSkeleton(): React.ReactElement {
  return (
    <Box
      sx={{
        p: 3,
        borderRadius: 4,
        backgroundColor: 'background.paper',
        boxShadow: '0px 4px 24px rgba(0, 0, 0, 0.10)',
      }}
    >
      <Stack spacing={2.5}>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <Skeleton variant="circular" width={48} height={48} />
          <Box sx={{ flex: 1 }}>
            <Skeleton variant="text" width="60%" height={32} />
            <Skeleton variant="rounded" width={120} height={24} sx={{ mt: 0.5 }} />
          </Box>
        </Box>

        <Box>
          <Skeleton variant="text" width="40%" height={24} sx={{ mb: 1 }} />
          <Skeleton variant="rounded" width="100%" height={8} />
        </Box>

        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: 'repeat(2, 1fr)',
            gap: 2,
          }}
        >
          {[1, 2, 3, 4].map((i) => (
            <Box key={i}>
              <Skeleton variant="text" width="60%" height={20} />
              <Skeleton variant="text" width="80%" height={28} />
            </Box>
          ))}
        </Box>

        <Skeleton variant="rounded" width="100%" height={60} />

        <Box sx={{ display: 'flex', gap: 1.5 }}>
          {[1, 2, 3, 4].map((i) => (
            <Box key={i} sx={{ flex: 1 }}>
              <Skeleton variant="text" width="100%" height={20} />
              <Skeleton variant="text" width="100%" height={24} />
            </Box>
          ))}
        </Box>
      </Stack>
    </Box>
  );
}

/**
 * Empty state when no destinations are found
 */
function EmptyState(): React.ReactElement {
  return (
    <Box
      sx={{
        textAlign: 'center',
        py: 8,
        px: 3,
      }}
    >
      <PublicOffIcon
        sx={{
          fontSize: 64,
          color: 'text.disabled',
          mb: 2,
        }}
      />
      <Typography
        variant="h5"
        sx={{
          fontWeight: 600,
          color: 'text.primary',
          mb: 1,
        }}
      >
        Nenhum destino encontrado
      </Typography>
      <Typography
        variant="body1"
        sx={{
          color: 'text.secondary',
          maxWidth: 400,
          mx: 'auto',
        }}
      >
        Não encontramos destinos para este NCM. Tente outro código ou ajuste os
        parâmetros de busca.
      </Typography>
    </Box>
  );
}

export function DestinationList({
  destinations,
  isLoading = false,
}: DestinationListProps): React.ReactElement {
  // Loading state
  if (isLoading) {
    return (
      <Box sx={{ width: '100%' }}>
        <Stack spacing={3}>
          {[1, 2, 3].map((i) => (
            <DestinationSkeleton key={i} />
          ))}
        </Stack>
      </Box>
    );
  }

  // Empty state
  if (!destinations || destinations.length === 0) {
    return <EmptyState />;
  }

  // Normal state with destinations
  return (
    <Box sx={{ width: '100%' }}>
      <Box sx={{ mb: 3 }}>
        <Typography
          variant="h4"
          sx={{
            fontWeight: 600,
            color: 'text.primary',
            mb: 1,
          }}
        >
          Destinos Recomendados
        </Typography>
        <Typography
          variant="body1"
          sx={{
            color: 'text.secondary',
          }}
        >
          {destinations.length}{' '}
          {destinations.length === 1 ? 'destino encontrado' : 'destinos encontrados'},{' '}
          ordenados por adequação
        </Typography>
      </Box>

      <Stack spacing={3}>
        {destinations.map((destination) => (
          <DestinationCard
            key={`${destination.country_code}-${destination.rank}`}
            destination={destination}
          />
        ))}
      </Stack>
    </Box>
  );
}
