'use client';

import React from 'react';
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Chip,
  Stack,
} from '@mui/material';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import CategoryIcon from '@mui/icons-material/Category';
import PublicIcon from '@mui/icons-material/Public';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useRouter } from 'next/navigation';
import type { OnboardingMetadata } from '@/lib/types/onboarding';

export interface DashboardSimulatorPreviewProps {
  metadata: OnboardingMetadata;
}

export function DashboardSimulatorPreview({
  metadata,
}: DashboardSimulatorPreviewProps): React.ReactElement {
  const router = useRouter();

  function handleSimulate() {
    const params = new URLSearchParams();
    if (metadata.ncmDefaultNcm8d) {
      params.set('ncm', metadata.ncmDefaultNcm8d);
    }
    router.push(`/simulator?${params.toString()}`);
  }

  const hasCountries = metadata.targetCountries.length > 0;
  const hasContinents = metadata.targetContinents.length > 0;

  return (
    <Card
      sx={{
        borderRadius: '16px',
        border: '1px solid var(--md-sys-color-primary)',
        background:
          'linear-gradient(135deg, var(--md-sys-color-primary-container) 0%, var(--md-sys-color-surface) 100%)',
        boxShadow: '0 4px 20px rgba(0,0,0,0.08)',
        gridColumn: { md: '1 / -1' }, // span full width on desktop
      }}
    >
      <CardContent sx={{ p: 3 }}>
        <Stack direction={{ xs: 'column', sm: 'row' }} spacing={2} alignItems={{ sm: 'center' }}>
          {/* Icon */}
          <Box
            sx={{
              width: 56,
              height: 56,
              borderRadius: '14px',
              bgcolor: 'var(--md-sys-color-primary)',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              flexShrink: 0,
            }}
          >
            <RocketLaunchIcon sx={{ color: 'var(--md-sys-color-on-primary)', fontSize: 28 }} />
          </Box>

          {/* Info */}
          <Box sx={{ flex: 1 }}>
            <Typography
              variant="caption"
              sx={{ color: 'var(--md-sys-color-on-surface-variant)', fontWeight: 500 }}
            >
              Simulação personalizada para o seu perfil
            </Typography>

            {/* NCM */}
            {metadata.ncmChapter && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5, flexWrap: 'wrap' }}>
                <CategoryIcon fontSize="small" sx={{ color: 'var(--md-sys-color-primary)' }} />
                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                  NCM {metadata.ncmChapter} — {metadata.ncmLabel}
                </Typography>
                <Chip
                  label={metadata.ncmDefaultNcm8d}
                  size="small"
                  sx={{
                    fontFamily: 'monospace',
                    fontSize: '0.7rem',
                    bgcolor: 'var(--md-sys-color-secondary-container)',
                    color: 'var(--md-sys-color-on-secondary-container)',
                    borderRadius: '6px',
                  }}
                />
              </Box>
            )}

            {/* Target regions */}
            {(hasCountries || hasContinents) && (
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mt: 0.5, flexWrap: 'wrap' }}>
                <PublicIcon fontSize="small" sx={{ color: 'var(--md-sys-color-secondary)' }} />
                <Typography variant="caption" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
                  Destinos:
                </Typography>
                {metadata.targetContinents.slice(0, 2).map((c) => (
                  <Chip
                    key={c}
                    label={c}
                    size="small"
                    sx={{
                      fontSize: '0.7rem',
                      bgcolor: 'var(--md-sys-color-primary)',
                      color: 'var(--md-sys-color-on-primary)',
                      textTransform: 'capitalize',
                    }}
                  />
                ))}
                {metadata.targetCountries.slice(0, 3).map((iso) => (
                  <Chip
                    key={iso}
                    label={iso}
                    size="small"
                    sx={{
                      fontSize: '0.7rem',
                      fontFamily: 'monospace',
                      bgcolor: 'var(--md-sys-color-secondary-container)',
                      color: 'var(--md-sys-color-on-secondary-container)',
                    }}
                  />
                ))}
                {metadata.targetCountries.length + metadata.targetContinents.length > 5 && (
                  <Typography variant="caption" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
                    +{metadata.targetCountries.length + metadata.targetContinents.length - 5} mais
                  </Typography>
                )}
              </Box>
            )}
          </Box>

          {/* CTA */}
          <Button
            variant="contained"
            endIcon={<ArrowForwardIcon />}
            onClick={handleSimulate}
            disabled={!metadata.ncmDefaultNcm8d}
            sx={{
              bgcolor: 'var(--md-sys-color-primary)',
              color: 'var(--md-sys-color-on-primary)',
              borderRadius: '12px',
              textTransform: 'none',
              fontWeight: 600,
              px: 3,
              py: 1.25,
              whiteSpace: 'nowrap',
              flexShrink: 0,
              '&:hover': {
                bgcolor: 'var(--md-sys-color-primary)',
                opacity: 0.9,
              },
            }}
          >
            Simular meu produto
          </Button>
        </Stack>
      </CardContent>
    </Card>
  );
}
