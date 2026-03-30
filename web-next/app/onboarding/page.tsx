'use client';

export const dynamic = 'force-dynamic';

import React, { useState, useMemo } from 'react';
import { useUser } from '@clerk/nextjs';
import {
  Box,
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  Stepper,
  Step,
  StepLabel,
  CircularProgress,
  LinearProgress,
  Stack,
} from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useRouter, useSearchParams } from 'next/navigation';
import type { OnboardingWizardState, OnboardingMetadata } from '@/lib/types/onboarding';
import { getNcmChapter } from '@/lib/data/ncm-chapters';
import { TRADE_REGIONS } from '@/lib/data/trade-regions';
import { OnboardingStep1Ncm } from '@/components/onboarding/OnboardingStep1Ncm';
import { OnboardingStep2Volume } from '@/components/onboarding/OnboardingStep2Volume';
import { OnboardingStep3Regions } from '@/components/onboarding/OnboardingStep3Regions';
import { OnboardingWelcome } from '@/components/onboarding/OnboardingWelcome';

const STEPS = ['Produto NCM', 'Volume', 'Destinos'];
const TOTAL_STEPS = STEPS.length; // steps 0-2; step 3 = completion screen

/** Constrói o estado inicial do wizard a partir do publicMetadata salvo (re-onboarding). */
function buildStateFromMetadata(meta: OnboardingMetadata): OnboardingWizardState {
  const ncmChapter = meta.ncmChapter ? getNcmChapter(meta.ncmChapter) ?? null : null;

  const targetRegions = TRADE_REGIONS.filter(
    (r) =>
      (r.type === 'continent' && meta.targetContinents.includes(r.code)) ||
      (r.type === 'country' && meta.targetCountries.includes(r.code))
  );

  return {
    ncmChapter,
    volumeMonthlyKg: meta.volumeMonthlyKg,
    targetRegions,
  };
}

export default function OnboardingPage() {
  const { user, isLoaded } = useUser();
  const router = useRouter();
  const searchParams = useSearchParams();
  const isReOnboarding = searchParams.get('re') === '1';

  const existingMeta = user?.publicMetadata as OnboardingMetadata | undefined;

  const initialState = useMemo<OnboardingWizardState>(() => {
    if (isReOnboarding && existingMeta?.onboardingCompleted) {
      return buildStateFromMetadata(existingMeta);
    }
    return { ncmChapter: null, volumeMonthlyKg: undefined, targetRegions: [] };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []); // intentionally run once on mount

  const [step, setStep] = useState(0);
  const [state, setState] = useState<OnboardingWizardState>(initialState);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);

  if (!isLoaded) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  // ── Save helpers ──────────────────────────────────────────────────────────

  async function saveMetadata(skipped: boolean) {
    setSaving(true);
    setSaveError(null);

    try {
      const continents = state.targetRegions
        .filter((r) => r.type === 'continent')
        .map((r) => r.code);
      const countries = state.targetRegions
        .filter((r) => r.type === 'country')
        .map((r) => r.code);

      const ncmChapter = state.ncmChapter;
      const chapterData = ncmChapter ? getNcmChapter(ncmChapter.code) : undefined;

      const currentVersion = existingMeta?.onboardingVersion ?? 0;
      const metadata: OnboardingMetadata = {
        onboardingCompleted: !skipped,
        onboardingVersion: isReOnboarding ? currentVersion + 1 : 1,
        onboardingCompletedAt: new Date().toISOString(),
        ncmChapter: ncmChapter?.code ?? '',
        ncmLabel: ncmChapter?.label ?? '',
        ncmDefaultNcm8d: chapterData?.defaultNcm8d ?? '',
        volumeMonthlyKg: state.volumeMonthlyKg,
        targetContinents: continents,
        targetCountries: countries,
        onboardingSkipped: skipped,
      };

      const res = await fetch('/api/onboarding', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(metadata),
      });

      if (!res.ok) {
        throw new Error('Falha ao salvar perfil');
      }

      // Reload Clerk user so publicMetadata is fresh
      await user?.reload();
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : 'Erro desconhecido');
      throw err;
    } finally {
      setSaving(false);
    }
  }

  // ── Navigation ────────────────────────────────────────────────────────────

  async function handleNext() {
    if (step < TOTAL_STEPS - 1) {
      setStep((s) => s + 1);
      return;
    }
    // Last step → save and show completion
    try {
      await saveMetadata(false);
      setStep(TOTAL_STEPS); // completion screen
    } catch {
      // error already set in saveError
    }
  }

  function handleBack() {
    setStep((s) => Math.max(0, s - 1));
  }

  async function handleSkip() {
    try {
      await saveMetadata(true);
    } catch {
      // ignore — redirect anyway
    }
    router.push('/dashboard');
  }

  function handleGoToDashboard() {
    router.push('/dashboard');
  }

  // ── Validation ────────────────────────────────────────────────────────────

  function canAdvance(): boolean {
    if (step === 0) return state.ncmChapter !== null;
    // steps 1 and 2 are optional
    return true;
  }

  // ── Render ────────────────────────────────────────────────────────────────

  const isCompletion = step === TOTAL_STEPS;
  const progress = Math.round((step / TOTAL_STEPS) * 100);

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'var(--md-sys-color-background)', py: 4 }}>
      <Container maxWidth="sm">
        {/* Header */}
        <Box sx={{ mb: 4, textAlign: 'center' }}>
          <Typography
            variant="h5"
            sx={{ fontWeight: 700, color: 'var(--md-sys-color-on-background)', mb: 0.5 }}
          >
            Olá{user?.firstName ? `, ${user.firstName}` : ''}!
          </Typography>
          <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
            {isCompletion
              ? 'Seu perfil foi configurado com sucesso'
              : 'Configure seu perfil em 3 passos rápidos'}
          </Typography>
        </Box>

        <Card
          sx={{
            borderRadius: '20px',
            border: '1px solid var(--md-sys-color-outline-variant)',
            boxShadow: '0 2px 12px rgba(0,0,0,0.08)',
          }}
        >
          {!isCompletion && (
            <>
              <LinearProgress
                variant="determinate"
                value={progress}
                sx={{
                  height: 4,
                  borderRadius: '20px 20px 0 0',
                  bgcolor: 'var(--md-sys-color-surface-variant)',
                  '& .MuiLinearProgress-bar': {
                    bgcolor: 'var(--md-sys-color-primary)',
                    borderRadius: 4,
                  },
                }}
              />
              <Box sx={{ px: 3, pt: 3, pb: 1 }}>
                <Stepper activeStep={step} alternativeLabel>
                  {STEPS.map((label) => (
                    <Step key={label}>
                      <StepLabel
                        sx={{
                          '& .MuiStepLabel-label': { fontSize: '0.75rem' },
                        }}
                      >
                        {label}
                      </StepLabel>
                    </Step>
                  ))}
                </Stepper>
              </Box>
            </>
          )}

          <CardContent sx={{ p: 3 }}>
            {/* Step content */}
            {step === 0 && (
              <OnboardingStep1Ncm
                value={state.ncmChapter}
                onChange={(chapter) => setState((s) => ({ ...s, ncmChapter: chapter }))}
              />
            )}
            {step === 1 && (
              <OnboardingStep2Volume
                value={state.volumeMonthlyKg}
                onChange={(vol) => setState((s) => ({ ...s, volumeMonthlyKg: vol }))}
              />
            )}
            {step === 2 && (
              <OnboardingStep3Regions
                value={state.targetRegions}
                onChange={(regions) => setState((s) => ({ ...s, targetRegions: regions }))}
              />
            )}
            {isCompletion && (
              <OnboardingWelcome
                state={state}
                onGoToDashboard={handleGoToDashboard}
                onSkip={handleSkip}
              />
            )}

            {/* Error */}
            {saveError && (
              <Typography
                variant="caption"
                sx={{ color: 'error.main', display: 'block', mt: 2, textAlign: 'center' }}
              >
                {saveError}
              </Typography>
            )}

            {/* Navigation */}
            {!isCompletion && (
              <Stack
                direction="row"
                justifyContent="space-between"
                alignItems="center"
                sx={{ mt: 4 }}
              >
                <Button
                  variant="text"
                  startIcon={<ArrowBackIcon />}
                  onClick={step === 0 ? handleSkip : handleBack}
                  disabled={saving}
                  sx={{
                    color: 'var(--md-sys-color-on-surface-variant)',
                    textTransform: 'none',
                  }}
                >
                  {step === 0 ? 'Pular' : 'Voltar'}
                </Button>

                <Button
                  variant="contained"
                  endIcon={saving ? <CircularProgress size={16} color="inherit" /> : <ArrowForwardIcon />}
                  onClick={handleNext}
                  disabled={!canAdvance() || saving}
                  sx={{
                    bgcolor: 'var(--md-sys-color-primary)',
                    color: 'var(--md-sys-color-on-primary)',
                    borderRadius: '12px',
                    textTransform: 'none',
                    fontWeight: 600,
                    px: 3,
                    py: 1,
                    '&:hover': {
                      bgcolor: 'var(--md-sys-color-primary)',
                      opacity: 0.9,
                    },
                    '&:disabled': {
                      opacity: 0.5,
                    },
                  }}
                >
                  {step === TOTAL_STEPS - 1 ? 'Concluir' : 'Próximo'}
                </Button>
              </Stack>
            )}
          </CardContent>
        </Card>

        {/* Bottom skip link (steps 1+) */}
        {!isCompletion && step > 0 && (
          <Box sx={{ textAlign: 'center', mt: 2 }}>
            <Button
              variant="text"
              size="small"
              onClick={handleSkip}
              disabled={saving}
              sx={{ color: 'var(--md-sys-color-on-surface-variant)', textTransform: 'none', fontSize: '0.8rem' }}
            >
              Pular configuração
            </Button>
          </Box>
        )}
      </Container>
    </Box>
  );
}
