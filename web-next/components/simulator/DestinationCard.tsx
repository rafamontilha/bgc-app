'use client';

/**
 * US-002: Destination Card Component
 * Displays a single destination with all relevant metrics
 */

import React from 'react';
import {
  Card,
  CardContent,
  Box,
  Typography,
  Chip,
  Stack,
  LinearProgress,
  Tooltip,
} from '@mui/material';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import PublicIcon from '@mui/icons-material/Public';
import LocalShippingIcon from '@mui/icons-material/LocalShipping';
import { Destination } from '@/types/simulator';

export interface DestinationCardProps {
  destination: Destination;
}

/**
 * Get color configuration for demand level
 */
function getDemandLevelColor(
  level: Destination['demand']
): 'success' | 'warning' | 'error' {
  switch (level) {
    case 'Alto':
      return 'success';
    case 'Médio':
      return 'warning';
    case 'Baixo':
      return 'error';
    default:
      return 'warning';
  }
}

/**
 * Format currency to USD
 */
function formatCurrency(value: number): string {
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'USD',
  }).format(value);
}

/**
 * Format number with thousand separators
 */
function formatNumber(value: number): string {
  return new Intl.NumberFormat('pt-BR').format(value);
}

/**
 * Format percentage
 */
function formatPercentage(value: number): string {
  return `${value.toFixed(1)}%`;
}

/**
 * Get country flag emoji from country code
 */
function getCountryFlag(countryCode: string): string {
  const codePoints = countryCode
    .toUpperCase()
    .split('')
    .map((char) => 127397 + char.charCodeAt(0));
  return String.fromCodePoint(...codePoints);
}

export function DestinationCard({
  destination,
}: DestinationCardProps): React.ReactElement {
  const demandColor = getDemandLevelColor(destination.demand);
  const scorePercentage = (destination.score / 10) * 100;

  return (
    <Card
      sx={{
        height: '100%',
        transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: '0px 12px 40px rgba(0, 0, 0, 0.16)',
        },
      }}
    >
      <CardContent sx={{ p: 3 }}>
        <Stack spacing={2.5}>
          {/* Header with rank, flag, and country */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <Box
              sx={{
                width: 48,
                height: 48,
                borderRadius: '50%',
                backgroundColor: 'primary.main',
                color: 'primary.contrastText',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontWeight: 700,
                fontSize: '1.25rem',
                flexShrink: 0,
              }}
            >
              {destination.rank}
            </Box>

            <Box sx={{ flex: 1, minWidth: 0 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5 }}>
                <Typography
                  variant="h4"
                  component="h3"
                  sx={{
                    fontWeight: 600,
                    fontSize: '1.5rem',
                  }}
                >
                  {getCountryFlag(destination.country_code)}{' '}
                  {destination.country_name}
                </Typography>
              </Box>
              <Chip
                label={`Demanda: ${destination.demand}`}
                color={demandColor}
                size="small"
                sx={{ fontWeight: 500 }}
              />
            </Box>
          </Box>

          {/* Score visual */}
          <Box>
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mb: 1,
              }}
            >
              <Typography
                variant="subtitle2"
                sx={{ color: 'text.secondary', fontWeight: 600 }}
              >
                Score de Adequação
              </Typography>
              <Typography
                variant="h6"
                sx={{ color: 'primary.main', fontWeight: 700 }}
              >
                {destination.score.toFixed(1)}/10
              </Typography>
            </Box>
            <LinearProgress
              variant="determinate"
              value={scorePercentage}
              sx={{
                height: 8,
                borderRadius: 4,
                backgroundColor: 'rgba(0, 122, 255, 0.12)',
                '& .MuiLinearProgress-bar': {
                  borderRadius: 4,
                },
              }}
            />
          </Box>

          {/* Key metrics */}
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: 'repeat(2, 1fr)',
              gap: 2,
            }}
          >
            <Tooltip title="Tamanho total do mercado em dólares" arrow>
              <Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 0.5 }}>
                  <PublicIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                  <Typography
                    variant="caption"
                    sx={{ color: 'text.secondary', fontWeight: 500 }}
                  >
                    Mercado
                  </Typography>
                </Box>
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  {formatCurrency(destination.market_size_usd)}
                </Typography>
              </Box>
            </Tooltip>

            <Tooltip title="Taxa de crescimento anual do mercado" arrow>
              <Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 0.5 }}>
                  <TrendingUpIcon sx={{ fontSize: 18, color: 'success.main' }} />
                  <Typography
                    variant="caption"
                    sx={{ color: 'text.secondary', fontWeight: 500 }}
                  >
                    Crescimento
                  </Typography>
                </Box>
                <Typography
                  variant="body1"
                  sx={{ fontWeight: 600, color: 'success.main' }}
                >
                  +{formatPercentage(destination.growth_rate_pct)}
                </Typography>
              </Box>
            </Tooltip>

            <Tooltip title="Preço médio por quilograma" arrow>
              <Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 0.5 }}>
                  <AttachMoneyIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                  <Typography
                    variant="caption"
                    sx={{ color: 'text.secondary', fontWeight: 500 }}
                  >
                    Preço/kg
                  </Typography>
                </Box>
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  {formatCurrency(destination.price_per_kg_usd)}
                </Typography>
              </Box>
            </Tooltip>

            <Tooltip title="Distância aproximada em quilômetros" arrow>
              <Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5, mb: 0.5 }}>
                  <LocalShippingIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                  <Typography
                    variant="caption"
                    sx={{ color: 'text.secondary', fontWeight: 500 }}
                  >
                    Distância
                  </Typography>
                </Box>
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  {formatNumber(destination.distance_km)} km
                </Typography>
              </Box>
            </Tooltip>
          </Box>

          {/* Recommendation reason */}
          <Box
            sx={{
              backgroundColor: 'rgba(0, 122, 255, 0.06)',
              borderRadius: 2,
              p: 2,
              borderLeft: 4,
              borderColor: 'primary.main',
            }}
          >
            <Typography
              variant="body2"
              sx={{
                color: 'text.primary',
                fontWeight: 500,
                lineHeight: 1.6,
              }}
            >
              {destination.recommendation_reason}
            </Typography>
          </Box>

          {/* Additional details */}
          <Box
            sx={{
              display: 'flex',
              flexWrap: 'wrap',
              gap: 1.5,
              pt: 1,
              borderTop: 1,
              borderColor: 'divider',
            }}
          >
            <Box sx={{ flex: '1 1 auto', minWidth: 100 }}>
              <Typography
                variant="caption"
                sx={{ color: 'text.secondary', display: 'block', mb: 0.5 }}
              >
                Margem Est.
              </Typography>
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                {formatPercentage(destination.estimated_margin_pct)}
              </Typography>
            </Box>
            <Box sx={{ flex: '1 1 auto', minWidth: 100 }}>
              <Typography
                variant="caption"
                sx={{ color: 'text.secondary', display: 'block', mb: 0.5 }}
              >
                Custo Logística
              </Typography>
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                {formatCurrency(destination.logistics_cost_usd)}
              </Typography>
            </Box>
            <Box sx={{ flex: '1 1 auto', minWidth: 100 }}>
              <Typography
                variant="caption"
                sx={{ color: 'text.secondary', display: 'block', mb: 0.5 }}
              >
                Tarifa
              </Typography>
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                {formatPercentage(destination.tariff_rate_pct)}
              </Typography>
            </Box>
            <Box sx={{ flex: '1 1 auto', minWidth: 100 }}>
              <Typography
                variant="caption"
                sx={{ color: 'text.secondary', display: 'block', mb: 0.5 }}
              >
                Lead Time
              </Typography>
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                {destination.lead_time_days} dias
              </Typography>
            </Box>
          </Box>
        </Stack>
      </CardContent>
    </Card>
  );
}
