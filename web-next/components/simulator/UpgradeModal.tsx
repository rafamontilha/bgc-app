'use client';

/**
 * US-004: Upgrade Modal Component
 * Displayed when rate limit is exceeded, prompting user to upgrade
 */

import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  Stack,
  Chip,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import IconButton from '@mui/material/IconButton';

export interface UpgradeModalProps {
  open: boolean;
  onClose: () => void;
  onViewPlans?: () => void;
}

/**
 * Pro plan features list
 */
const PRO_FEATURES = [
  'Simulações ilimitadas',
  'Histórico de análises',
  'Comparação lado a lado',
  'Exportação para PDF',
  'Alertas de oportunidades',
  'Suporte prioritário',
];

export function UpgradeModal({
  open,
  onClose,
  onViewPlans,
}: UpgradeModalProps): React.ReactElement {
  const handleViewPlans = (): void => {
    if (onViewPlans) {
      onViewPlans();
    } else {
      // Default behavior: navigate to pricing page
      window.location.href = '/pricing';
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      sx={{
        '& .MuiDialog-paper': {
          borderRadius: 5,
          p: 1,
        },
      }}
    >
      <DialogTitle
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          pb: 2,
        }}
      >
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
          <Box
            sx={{
              width: 48,
              height: 48,
              borderRadius: 3,
              backgroundColor: 'primary.main',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            <RocketLaunchIcon sx={{ color: 'white', fontSize: 28 }} />
          </Box>
          <Box>
            <Typography
              variant="h5"
              sx={{
                fontWeight: 600,
                color: 'text.primary',
              }}
            >
              Limite Atingido
            </Typography>
            <Chip
              label="Upgrade para Pro"
              color="primary"
              size="small"
              sx={{ mt: 0.5, fontWeight: 600 }}
            />
          </Box>
        </Box>
        <IconButton
          onClick={onClose}
          size="small"
          sx={{
            color: 'text.secondary',
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>

      <DialogContent>
        <Stack spacing={3}>
          <Box>
            <Typography
              variant="body1"
              sx={{
                color: 'text.secondary',
                mb: 3,
              }}
            >
              Você atingiu o limite de simulações do plano gratuito. Faça upgrade
              para o <strong>Brasil Global Connect Pro</strong> e desbloqueie
              todo o potencial da plataforma.
            </Typography>

            <Box
              sx={{
                backgroundColor: 'rgba(0, 122, 255, 0.06)',
                borderRadius: 3,
                p: 3,
                border: 2,
                borderColor: 'primary.main',
              }}
            >
              <Typography
                variant="overline"
                sx={{
                  color: 'primary.main',
                  fontWeight: 700,
                  letterSpacing: '0.1em',
                }}
              >
                Benefícios Pro
              </Typography>

              <Stack spacing={2} sx={{ mt: 2 }}>
                {PRO_FEATURES.map((feature, index) => (
                  <Box
                    key={index}
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      gap: 1.5,
                    }}
                  >
                    <CheckCircleIcon
                      sx={{
                        color: 'success.main',
                        fontSize: 24,
                        flexShrink: 0,
                      }}
                    />
                    <Typography
                      variant="body1"
                      sx={{
                        color: 'text.primary',
                        fontWeight: 500,
                      }}
                    >
                      {feature}
                    </Typography>
                  </Box>
                ))}
              </Stack>
            </Box>
          </Box>

          <Box
            sx={{
              textAlign: 'center',
              py: 2,
              px: 3,
              backgroundColor: 'background.default',
              borderRadius: 2,
            }}
          >
            <Typography variant="h4" sx={{ fontWeight: 700, mb: 0.5 }}>
              A partir de R$ 97/mês
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary' }}>
              Cancele a qualquer momento
            </Typography>
          </Box>
        </Stack>
      </DialogContent>

      <DialogActions
        sx={{
          px: 3,
          pb: 3,
          pt: 2,
          gap: 2,
        }}
      >
        <Button
          onClick={onClose}
          variant="outlined"
          size="large"
          fullWidth
          sx={{
            borderRadius: 999,
            fontWeight: 600,
          }}
        >
          Agora Não
        </Button>
        <Button
          onClick={handleViewPlans}
          variant="contained"
          size="large"
          fullWidth
          startIcon={<RocketLaunchIcon />}
          sx={{
            borderRadius: 999,
            fontWeight: 600,
            boxShadow: 3,
            '&:hover': {
              boxShadow: 6,
            },
          }}
        >
          Ver Planos
        </Button>
      </DialogActions>
    </Dialog>
  );
}
