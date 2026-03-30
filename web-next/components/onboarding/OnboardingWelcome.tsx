'use client';

import React from 'react';
import { Box, Typography, Button, Stack, Chip } from '@mui/material';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CategoryIcon from '@mui/icons-material/Category';
import ScaleIcon from '@mui/icons-material/Scale';
import PublicIcon from '@mui/icons-material/Public';
import type { OnboardingWizardState } from '@/lib/types/onboarding';

export interface OnboardingWelcomeProps {
  state: OnboardingWizardState;
  onGoToDashboard: () => void;
  onSkip: () => void;
}

export function OnboardingWelcome({
  state,
  onGoToDashboard,
  onSkip,
}: OnboardingWelcomeProps): React.ReactElement {
  const hasNcm = state.ncmChapter !== null;
  const hasVolume = state.volumeMonthlyKg !== undefined;
  const hasRegions = state.targetRegions.length > 0;

  const completedCount = [hasNcm, hasVolume, hasRegions].filter(Boolean).length;

  return (
    <Box sx={{ textAlign: 'center' }}>
      <Box
        sx={{
          width: 80,
          height: 80,
          borderRadius: '50%',
          bgcolor: 'var(--md-sys-color-primary-container)',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          mx: 'auto',
          mb: 3,
        }}
      >
        <RocketLaunchIcon
          sx={{ fontSize: 40, color: 'var(--md-sys-color-primary)' }}
        />
      </Box>

      <Typography
        variant="h5"
        sx={{ fontWeight: 700, color: 'var(--md-sys-color-on-surface)', mb: 1 }}
      >
        Tudo pronto!
      </Typography>

      <Typography
        variant="body1"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 4 }}
      >
        Seu perfil foi configurado com sucesso. Veja o resumo abaixo.
      </Typography>

      {/* Summary cards */}
      <Stack spacing={1.5} sx={{ mb: 4, textAlign: 'left' }}>
        {/* NCM */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            p: 2,
            borderRadius: '12px',
            border: '1px solid var(--md-sys-color-outline-variant)',
            bgcolor: hasNcm
              ? 'var(--md-sys-color-primary-container)'
              : 'var(--md-sys-color-surface-variant)',
          }}
        >
          <CategoryIcon
            sx={{
              color: hasNcm
                ? 'var(--md-sys-color-primary)'
                : 'var(--md-sys-color-on-surface-variant)',
            }}
          />
          <Box sx={{ flex: 1 }}>
            <Typography variant="caption" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
              Produto principal (NCM)
            </Typography>
            <Typography variant="body2" sx={{ fontWeight: 600 }}>
              {hasNcm
                ? `${state.ncmChapter!.code} — ${state.ncmChapter!.label}`
                : 'Não informado'}
            </Typography>
          </Box>
          {hasNcm && (
            <CheckCircleIcon fontSize="small" sx={{ color: 'var(--md-sys-color-primary)' }} />
          )}
        </Box>

        {/* Volume */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            gap: 2,
            p: 2,
            borderRadius: '12px',
            border: '1px solid var(--md-sys-color-outline-variant)',
            bgcolor: hasVolume
              ? 'var(--md-sys-color-primary-container)'
              : 'var(--md-sys-color-surface-variant)',
          }}
        >
          <ScaleIcon
            sx={{
              color: hasVolume
                ? 'var(--md-sys-color-primary)'
                : 'var(--md-sys-color-on-surface-variant)',
            }}
          />
          <Box sx={{ flex: 1 }}>
            <Typography variant="caption" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
              Volume médio mensal
            </Typography>
            <Typography variant="body2" sx={{ fontWeight: 600 }}>
              {hasVolume
                ? `${state.volumeMonthlyKg!.toLocaleString('pt-BR')} kg/mês`
                : 'Não informado'}
            </Typography>
          </Box>
          {hasVolume && (
            <CheckCircleIcon fontSize="small" sx={{ color: 'var(--md-sys-color-primary)' }} />
          )}
        </Box>

        {/* Regions */}
        <Box
          sx={{
            display: 'flex',
            alignItems: 'flex-start',
            gap: 2,
            p: 2,
            borderRadius: '12px',
            border: '1px solid var(--md-sys-color-outline-variant)',
            bgcolor: hasRegions
              ? 'var(--md-sys-color-primary-container)'
              : 'var(--md-sys-color-surface-variant)',
          }}
        >
          <PublicIcon
            sx={{
              color: hasRegions
                ? 'var(--md-sys-color-primary)'
                : 'var(--md-sys-color-on-surface-variant)',
              mt: 0.25,
            }}
          />
          <Box sx={{ flex: 1 }}>
            <Typography variant="caption" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
              Regiões alvo
            </Typography>
            {hasRegions ? (
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 0.5, mt: 0.5 }}>
                {state.targetRegions.slice(0, 6).map((r) => (
                  <Chip
                    key={r.code}
                    label={r.label}
                    size="small"
                    sx={{
                      fontSize: '0.7rem',
                      bgcolor:
                        r.type === 'continent'
                          ? 'var(--md-sys-color-primary)'
                          : 'var(--md-sys-color-secondary-container)',
                      color:
                        r.type === 'continent'
                          ? 'var(--md-sys-color-on-primary)'
                          : 'var(--md-sys-color-on-secondary-container)',
                    }}
                  />
                ))}
                {state.targetRegions.length > 6 && (
                  <Chip
                    label={`+${state.targetRegions.length - 6}`}
                    size="small"
                    sx={{ fontSize: '0.7rem' }}
                  />
                )}
              </Box>
            ) : (
              <Typography variant="body2" sx={{ fontWeight: 600 }}>
                Não informado
              </Typography>
            )}
          </Box>
          {hasRegions && (
            <CheckCircleIcon fontSize="small" sx={{ color: 'var(--md-sys-color-primary)', mt: 0.25 }} />
          )}
        </Box>
      </Stack>

      {/* Progress indicator */}
      <Typography
        variant="caption"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 3, display: 'block' }}
      >
        {completedCount}/3 campos preenchidos
      </Typography>

      {/* CTA */}
      <Stack spacing={1.5}>
        <Button
          variant="contained"
          size="large"
          onClick={onGoToDashboard}
          fullWidth
          sx={{
            bgcolor: 'var(--md-sys-color-primary)',
            color: 'var(--md-sys-color-on-primary)',
            borderRadius: '12px',
            textTransform: 'none',
            fontWeight: 600,
            py: 1.5,
            '&:hover': {
              bgcolor: 'var(--md-sys-color-primary)',
              opacity: 0.9,
            },
          }}
        >
          Ir para o Dashboard
        </Button>

        {completedCount === 0 && (
          <Button
            variant="text"
            size="small"
            onClick={onSkip}
            sx={{
              color: 'var(--md-sys-color-on-surface-variant)',
              textTransform: 'none',
              fontSize: '0.8rem',
            }}
          >
            Pular por enquanto
          </Button>
        )}
      </Stack>
    </Box>
  );
}
