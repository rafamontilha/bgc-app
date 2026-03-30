'use client';

import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Button,
  IconButton,
  Collapse,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import AssignmentIndIcon from '@mui/icons-material/AssignmentInd';
import { useRouter } from 'next/navigation';
import { REMINDER_BANNER_DISMISSED_KEY } from '@/lib/types/onboarding';

export function OnboardingReminderBanner(): React.ReactElement | null {
  const router = useRouter();
  const [visible, setVisible] = useState(false);

  useEffect(() => {
    const dismissed = localStorage.getItem(REMINDER_BANNER_DISMISSED_KEY);
    if (!dismissed) {
      setVisible(true);
    }
  }, []);

  function handleDismiss() {
    localStorage.setItem(REMINDER_BANNER_DISMISSED_KEY, '1');
    setVisible(false);
  }

  function handleComplete() {
    handleDismiss();
    router.push('/onboarding');
  }

  return (
    <Collapse in={visible} unmountOnExit>
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          gap: 2,
          px: 3,
          py: 1.5,
          bgcolor: 'var(--md-sys-color-secondary-container)',
          borderBottom: '1px solid var(--md-sys-color-outline-variant)',
          flexWrap: 'wrap',
        }}
        role="banner"
        aria-label="Lembrete de configuração de perfil"
      >
        <AssignmentIndIcon
          fontSize="small"
          sx={{ color: 'var(--md-sys-color-on-secondary-container)', flexShrink: 0 }}
        />

        <Typography
          variant="body2"
          sx={{ color: 'var(--md-sys-color-on-secondary-container)', flex: 1, minWidth: 200 }}
        >
          Complete seu perfil para receber recomendações personalizadas de destinos.
        </Typography>

        <Button
          size="small"
          variant="contained"
          onClick={handleComplete}
          sx={{
            bgcolor: 'var(--md-sys-color-secondary)',
            color: 'var(--md-sys-color-on-secondary)',
            borderRadius: '8px',
            textTransform: 'none',
            fontWeight: 600,
            flexShrink: 0,
            '&:hover': { bgcolor: 'var(--md-sys-color-secondary)', opacity: 0.9 },
          }}
        >
          Completar agora
        </Button>

        <IconButton
          size="small"
          onClick={handleDismiss}
          aria-label="Dispensar lembrete"
          sx={{ color: 'var(--md-sys-color-on-secondary-container)', flexShrink: 0 }}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      </Box>
    </Collapse>
  );
}
