'use client';

import React, { useEffect, useState } from 'react';
import {
  Dialog,
  DialogContent,
  Box,
  Typography,
  Button,
  Stack,
  MobileStepper,
  IconButton,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';
import PublicIcon from '@mui/icons-material/Public';
import InsightsIcon from '@mui/icons-material/Insights';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import { TUTORIAL_SEEN_KEY } from '@/lib/types/onboarding';

interface TutorialSlide {
  icon: React.ReactElement;
  title: string;
  body: string;
}

const SLIDES: TutorialSlide[] = [
  {
    icon: <RocketLaunchIcon sx={{ fontSize: 56, color: 'var(--md-sys-color-primary)' }} />,
    title: 'Bem-vindo ao BGC!',
    body: 'Você acaba de configurar seu perfil de exportador. Deixa a gente te mostrar as principais funcionalidades em 3 passos rápidos.',
  },
  {
    icon: <PublicIcon sx={{ fontSize: 56, color: 'var(--md-sys-color-secondary)' }} />,
    title: 'Simulador de Destinos',
    body: 'Informe um NCM e veja instantaneamente os melhores mercados para o seu produto, com score de atratividade baseado em dados reais do Comex Stat.',
  },
  {
    icon: <InsightsIcon sx={{ fontSize: 56, color: 'var(--md-sys-color-tertiary, #6750A4)' }} />,
    title: 'Insights Personalizados',
    body: 'Com base no seu NCM e destinos alvo, o BGC vai sugerir oportunidades de mercado e alertar sobre mudanças de tarifa e demanda. Em breve!',
  },
  {
    icon: <AccountCircleIcon sx={{ fontSize: 56, color: 'var(--md-sys-color-primary)' }} />,
    title: 'Seu perfil é dinâmico',
    body: 'A qualquer momento você pode atualizar seu NCM, volume e destinos em Perfil → Configurações. Quanto mais completo, melhores as recomendações.',
  },
];

export interface OnboardingTutorialProps {
  /** Called when the user closes or finishes the tutorial */
  onClose?: () => void;
}

export function OnboardingTutorial({ onClose }: OnboardingTutorialProps): React.ReactElement | null {
  const [open, setOpen] = useState(false);
  const [step, setStep] = useState(0);

  useEffect(() => {
    const seen = localStorage.getItem(TUTORIAL_SEEN_KEY);
    if (!seen) {
      setOpen(true);
    }
  }, []);

  function handleClose() {
    localStorage.setItem(TUTORIAL_SEEN_KEY, '1');
    setOpen(false);
    onClose?.();
  }

  function handleNext() {
    if (step < SLIDES.length - 1) {
      setStep((s) => s + 1);
    } else {
      handleClose();
    }
  }

  function handleBack() {
    setStep((s) => Math.max(0, s - 1));
  }

  if (!open) return null;

  const slide = SLIDES[step];
  const isLast = step === SLIDES.length - 1;

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      maxWidth="xs"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: '20px',
          overflow: 'hidden',
        },
      }}
    >
      {/* Close button */}
      <Box sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1 }}>
        <IconButton size="small" onClick={handleClose} aria-label="Fechar tutorial">
          <CloseIcon fontSize="small" />
        </IconButton>
      </Box>

      <DialogContent sx={{ p: 4, textAlign: 'center' }}>
        {/* Icon */}
        <Box
          sx={{
            width: 96,
            height: 96,
            borderRadius: '24px',
            bgcolor: 'var(--md-sys-color-primary-container)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            mx: 'auto',
            mb: 3,
          }}
        >
          {slide.icon}
        </Box>

        {/* Text */}
        <Typography variant="h6" sx={{ fontWeight: 700, mb: 1.5 }}>
          {slide.title}
        </Typography>
        <Typography
          variant="body2"
          sx={{ color: 'var(--md-sys-color-on-surface-variant)', lineHeight: 1.6 }}
        >
          {slide.body}
        </Typography>

        {/* Stepper dots */}
        <MobileStepper
          variant="dots"
          steps={SLIDES.length}
          position="static"
          activeStep={step}
          sx={{
            justifyContent: 'center',
            bgcolor: 'transparent',
            mt: 3,
            '& .MuiMobileStepper-dot': {
              bgcolor: 'var(--md-sys-color-outline)',
            },
            '& .MuiMobileStepper-dotActive': {
              bgcolor: 'var(--md-sys-color-primary)',
            },
          }}
          nextButton={null}
          backButton={null}
        />

        {/* Navigation */}
        <Stack direction="row" spacing={1.5} justifyContent="center" sx={{ mt: 2 }}>
          {step > 0 && (
            <Button
              variant="outlined"
              onClick={handleBack}
              sx={{
                borderRadius: '12px',
                textTransform: 'none',
                borderColor: 'var(--md-sys-color-outline)',
                color: 'var(--md-sys-color-on-surface-variant)',
                px: 3,
              }}
            >
              Anterior
            </Button>
          )}
          <Button
            variant="contained"
            onClick={handleNext}
            sx={{
              bgcolor: 'var(--md-sys-color-primary)',
              color: 'var(--md-sys-color-on-primary)',
              borderRadius: '12px',
              textTransform: 'none',
              fontWeight: 600,
              px: 4,
              '&:hover': { bgcolor: 'var(--md-sys-color-primary)', opacity: 0.9 },
            }}
          >
            {isLast ? 'Começar' : 'Próximo'}
          </Button>
        </Stack>
      </DialogContent>
    </Dialog>
  );
}
