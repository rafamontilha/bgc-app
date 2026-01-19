'use client';

/**
 * US-005: Error State Component
 * Displays user-friendly error messages for various error scenarios
 */

import React from 'react';
import { Box, Typography, Button, Stack, Alert } from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
import SearchOffIcon from '@mui/icons-material/SearchOff';
import CloudOffIcon from '@mui/icons-material/CloudOff';
import WarningAmberIcon from '@mui/icons-material/WarningAmber';
import RefreshIcon from '@mui/icons-material/Refresh';

export type ErrorType =
  | 'not_found'
  | 'validation'
  | 'server_error'
  | 'network_error'
  | 'timeout'
  | 'unknown';

export interface ErrorStateProps {
  type?: ErrorType;
  title?: string;
  message?: string;
  onRetry?: () => void;
  showRetryButton?: boolean;
}

/**
 * Error configuration for different error types
 */
const ERROR_CONFIG: Record<
  ErrorType,
  {
    icon: React.ReactElement;
    defaultTitle: string;
    defaultMessage: string;
    severity: 'error' | 'warning' | 'info';
  }
> = {
  not_found: {
    icon: <SearchOffIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'NCM não encontrado',
    defaultMessage:
      'O código NCM informado não foi encontrado em nossa base de dados. Verifique se o código está correto e tente novamente.',
    severity: 'warning',
  },
  validation: {
    icon: <WarningAmberIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'Dados inválidos',
    defaultMessage:
      'Os dados informados não são válidos. Verifique os campos e tente novamente.',
    severity: 'warning',
  },
  server_error: {
    icon: <ErrorOutlineIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'Erro no servidor',
    defaultMessage:
      'Ocorreu um erro em nossos servidores. Nossa equipe já foi notificada. Tente novamente em alguns instantes.',
    severity: 'error',
  },
  network_error: {
    icon: <CloudOffIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'Erro de conexão',
    defaultMessage:
      'Não foi possível conectar ao servidor. Verifique sua conexão com a internet e tente novamente.',
    severity: 'error',
  },
  timeout: {
    icon: <CloudOffIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'Tempo esgotado',
    defaultMessage:
      'A requisição demorou muito tempo. Verifique sua conexão e tente novamente.',
    severity: 'error',
  },
  unknown: {
    icon: <ErrorOutlineIcon sx={{ fontSize: 64 }} />,
    defaultTitle: 'Erro inesperado',
    defaultMessage:
      'Algo deu errado. Por favor, tente novamente. Se o problema persistir, entre em contato com o suporte.',
    severity: 'error',
  },
};

export function ErrorState({
  type = 'unknown',
  title,
  message,
  onRetry,
  showRetryButton = true,
}: ErrorStateProps): React.ReactElement {
  const config = ERROR_CONFIG[type];
  const displayTitle = title || config.defaultTitle;
  const displayMessage = message || config.defaultMessage;

  return (
    <Box
      sx={{
        width: '100%',
        py: 8,
        px: 3,
      }}
    >
      <Stack spacing={3} alignItems="center">
        {/* Icon */}
        <Box
          sx={{
            color:
              config.severity === 'error'
                ? 'error.main'
                : config.severity === 'warning'
                  ? 'warning.main'
                  : 'info.main',
          }}
        >
          {config.icon}
        </Box>

        {/* Title and message */}
        <Box sx={{ textAlign: 'center', maxWidth: 500 }}>
          <Typography
            variant="h4"
            sx={{
              fontWeight: 600,
              color: 'text.primary',
              mb: 1.5,
            }}
          >
            {displayTitle}
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'text.secondary',
              lineHeight: 1.6,
            }}
          >
            {displayMessage}
          </Typography>
        </Box>

        {/* Retry button */}
        {showRetryButton && onRetry && (
          <Button
            variant="contained"
            size="large"
            onClick={onRetry}
            startIcon={<RefreshIcon />}
            sx={{
              borderRadius: 999,
              px: 4,
              fontWeight: 600,
            }}
          >
            Tentar Novamente
          </Button>
        )}

        {/* Additional info for specific errors */}
        {type === 'not_found' && (
          <Alert
            severity="info"
            sx={{
              borderRadius: 3,
              maxWidth: 500,
            }}
          >
            <Typography variant="body2">
              <strong>Dica:</strong> O código NCM deve ter exatamente 8 dígitos
              numéricos. Exemplos válidos: 17011400, 08051020, 12011000
            </Typography>
          </Alert>
        )}

        {type === 'network_error' && (
          <Alert
            severity="info"
            sx={{
              borderRadius: 3,
              maxWidth: 500,
            }}
          >
            <Typography variant="body2">
              <strong>Verificações:</strong>
            </Typography>
            <Box component="ul" sx={{ mt: 1, pl: 2, mb: 0 }}>
              <li>Sua conexão com a internet está ativa?</li>
              <li>O firewall ou antivírus está bloqueando a conexão?</li>
              <li>Você está usando uma VPN que pode estar interferindo?</li>
            </Box>
          </Alert>
        )}
      </Stack>
    </Box>
  );
}
