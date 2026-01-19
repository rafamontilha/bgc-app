'use client';

/**
 * US-003: Rate Limit Banner Component
 * Displays remaining simulations and upgrade prompt
 */

import React from 'react';
import {
  Alert,
  Box,
  Typography,
  LinearProgress,
  AlertTitle,
  Button,
} from '@mui/material';
import InfoIcon from '@mui/icons-material/Info';
import WarningIcon from '@mui/icons-material/Warning';
import { RateLimitInfo } from '@/types/simulator';

export interface RateLimitBannerProps {
  rateLimitInfo: RateLimitInfo | null;
  onUpgradeClick?: () => void;
}

/**
 * Format Unix timestamp to readable time
 */
function formatResetTime(timestamp: number): string {
  const resetDate = new Date(timestamp * 1000);
  const now = new Date();
  const diffMs = resetDate.getTime() - now.getTime();
  const diffMinutes = Math.ceil(diffMs / (1000 * 60));

  if (diffMinutes <= 0) {
    return 'em breve';
  }

  if (diffMinutes < 60) {
    return `em ${diffMinutes} ${diffMinutes === 1 ? 'minuto' : 'minutos'}`;
  }

  const diffHours = Math.ceil(diffMinutes / 60);
  return `em ${diffHours} ${diffHours === 1 ? 'hora' : 'horas'}`;
}

export function RateLimitBanner({
  rateLimitInfo,
  onUpgradeClick,
}: RateLimitBannerProps): React.ReactElement | null {
  // Don't show banner if no rate limit info
  if (!rateLimitInfo) {
    return null;
  }

  const { limit, remaining, reset } = rateLimitInfo;
  const usedPercentage = ((limit - remaining) / limit) * 100;
  const isLowRemaining = remaining < 2;
  const isNoRemaining = remaining === 0;

  // Critical state (no simulations left)
  if (isNoRemaining) {
    return (
      <Alert
        severity="error"
        icon={<WarningIcon fontSize="medium" />}
        sx={{
          borderRadius: 3,
          boxShadow: '0px 4px 20px rgba(255, 59, 48, 0.12)',
          '& .MuiAlert-message': {
            width: '100%',
          },
        }}
        action={
          onUpgradeClick && (
            <Button
              color="inherit"
              size="small"
              onClick={onUpgradeClick}
              sx={{
                fontWeight: 600,
                borderRadius: 999,
                px: 2,
              }}
            >
              Ver Planos
            </Button>
          )
        }
      >
        <AlertTitle sx={{ fontWeight: 600, mb: 1 }}>
          Limite de Simulações Atingido
        </AlertTitle>
        <Typography variant="body2" sx={{ mb: 1.5 }}>
          Você atingiu o limite de <strong>{limit} simulações</strong> do plano
          gratuito. Faça upgrade para Pro e simule ilimitado.
        </Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          Suas simulações serão renovadas {formatResetTime(reset)}
        </Typography>
      </Alert>
    );
  }

  // Warning state (low remaining)
  if (isLowRemaining) {
    return (
      <Alert
        severity="warning"
        icon={<WarningIcon fontSize="medium" />}
        sx={{
          borderRadius: 3,
          boxShadow: '0px 4px 20px rgba(255, 149, 0, 0.12)',
          '& .MuiAlert-message': {
            width: '100%',
          },
        }}
        action={
          onUpgradeClick && (
            <Button
              color="inherit"
              size="small"
              onClick={onUpgradeClick}
              sx={{
                fontWeight: 600,
                borderRadius: 999,
                px: 2,
              }}
            >
              Upgrade
            </Button>
          )
        }
      >
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
          <Typography variant="body2" sx={{ fontWeight: 600 }}>
            {remaining} de {limit} simulações restantes
          </Typography>
          <Typography variant="caption" sx={{ color: 'text.secondary' }}>
            Renova {formatResetTime(reset)}
          </Typography>
        </Box>
        <LinearProgress
          variant="determinate"
          value={usedPercentage}
          color="warning"
          sx={{
            height: 6,
            borderRadius: 3,
            backgroundColor: 'rgba(255, 149, 0, 0.12)',
          }}
        />
      </Alert>
    );
  }

  // Info state (normal usage)
  return (
    <Alert
      severity="info"
      icon={<InfoIcon fontSize="medium" />}
      sx={{
        borderRadius: 3,
        backgroundColor: 'rgba(0, 122, 255, 0.08)',
        '& .MuiAlert-message': {
          width: '100%',
        },
      }}
    >
      <Box
        sx={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          flexWrap: 'wrap',
          gap: 1,
        }}
      >
        <Typography variant="body2">
          <strong>{remaining}</strong> de <strong>{limit}</strong> simulações
          restantes hoje
        </Typography>
        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
          Renova {formatResetTime(reset)}
        </Typography>
      </Box>
    </Alert>
  );
}
