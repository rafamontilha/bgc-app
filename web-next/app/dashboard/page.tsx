'use client';

export const dynamic = 'force-dynamic';

import { useUser } from '@clerk/nextjs';
import {
  Box,
  Container,
  Typography,
  Button,
  Card,
  CardContent,
  CircularProgress,
} from '@mui/material';
import { useRouter } from 'next/navigation';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import PublicIcon from '@mui/icons-material/Public';
import HistoryIcon from '@mui/icons-material/History';
import InsightsIcon from '@mui/icons-material/Insights';
import { Header } from '@/components/home/Header';
import { DashboardSimulatorPreview } from '@/components/dashboard/DashboardSimulatorPreview';
import { OnboardingTutorial } from '@/components/onboarding/OnboardingTutorial';
import { OnboardingReminderBanner } from '@/components/onboarding/OnboardingReminderBanner';
import type { OnboardingMetadata } from '@/lib/types/onboarding';

export default function DashboardPage() {
  const { user, isLoaded } = useUser();
  const router = useRouter();

  if (!isLoaded) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <CircularProgress />
      </Box>
    );
  }

  const meta = user?.publicMetadata as OnboardingMetadata | undefined;
  const onboardingDone = meta?.onboardingCompleted === true && !meta?.onboardingSkipped;

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'var(--md-sys-color-background)' }}>
      <Header />

      {/* Tutorial modal — shown once after first completed onboarding */}
      {onboardingDone && <OnboardingTutorial />}

      {/* Sticky reminder banner — shown to users who skipped onboarding */}
      {!onboardingDone && <OnboardingReminderBanner />}
      <Box sx={{ py: 4 }}>
        <Container maxWidth="lg">
          <Box sx={{ mb: 4 }}>
            <Typography
              variant="h4"
              sx={{ fontWeight: 600, color: 'var(--md-sys-color-on-background)', mb: 1 }}
            >
              Olá, {user?.firstName || 'Exportador'}! 👋
            </Typography>
            <Typography variant="body1" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}>
              Bem-vindo ao seu painel de inteligência de exportação
            </Typography>
          </Box>

          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: { xs: '1fr', md: '1fr 1fr' },
              gap: 3,
            }}
          >
            {/* Personalized simulator preview — shown when onboarding completed with NCM */}
            {onboardingDone && meta?.ncmDefaultNcm8d && (
              <DashboardSimulatorPreview metadata={meta} />
            )}

            {/* Onboarding prompt — shown when not yet done */}
            {!onboardingDone && (
              <Box sx={{ gridColumn: { md: '1 / -1' } }}>
                <Card
                  sx={{
                    borderRadius: '16px',
                    border: '1px dashed var(--md-sys-color-primary)',
                    bgcolor: 'var(--md-sys-color-primary-container)',
                    boxShadow: 'none',
                  }}
                >
                  <CardContent
                    sx={{
                      p: 3,
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                      flexWrap: 'wrap',
                      gap: 2,
                    }}
                  >
                    <Box>
                      <Typography variant="body1" sx={{ fontWeight: 600 }}>
                        Configure seu perfil para recomendações personalizadas
                      </Typography>
                      <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mt: 0.5 }}>
                        Leva menos de 2 minutos — informe seu produto NCM e destinos alvo.
                      </Typography>
                    </Box>
                    <Button
                      variant="contained"
                      onClick={() => router.push('/onboarding')}
                      sx={{
                        bgcolor: 'var(--md-sys-color-primary)',
                        color: 'var(--md-sys-color-on-primary)',
                        borderRadius: '12px',
                        textTransform: 'none',
                        fontWeight: 600,
                        px: 3,
                        '&:hover': { bgcolor: 'var(--md-sys-color-primary)', opacity: 0.9 },
                      }}
                    >
                      Completar perfil
                    </Button>
                  </CardContent>
                </Card>
              </Box>
            )}

            {/* Simulator card */}
            <Box>
              <Card
                sx={{
                  borderRadius: '16px',
                  border: '1px solid var(--md-sys-color-outline)',
                  boxShadow: 'var(--md-sys-elevation-1)',
                  height: '100%',
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <PublicIcon sx={{ fontSize: 32, color: 'var(--md-sys-color-primary)', mr: 1 }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Simulador de Destinos
                    </Typography>
                  </Box>
                  <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 2 }}>
                    Descubra os melhores destinos para suas exportações com análise de dados em tempo real
                  </Typography>
                  <Button
                    variant="contained"
                    onClick={() => router.push('/simulator')}
                    sx={{
                      bgcolor: 'var(--md-sys-color-primary)',
                      color: 'var(--md-sys-color-on-primary)',
                      borderRadius: '12px',
                      textTransform: 'none',
                      fontWeight: 500,
                      px: 3,
                      '&:hover': { bgcolor: 'var(--md-sys-color-primary)', opacity: 0.9 },
                    }}
                  >
                    Iniciar Simulação
                  </Button>
                </CardContent>
              </Card>
            </Box>

            {/* Histórico */}
            <Box>
              <Card
                sx={{
                  borderRadius: '16px',
                  border: '1px solid var(--md-sys-color-outline)',
                  boxShadow: 'var(--md-sys-elevation-1)',
                  height: '100%',
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <HistoryIcon sx={{ fontSize: 32, color: 'var(--md-sys-color-secondary)', mr: 1 }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Histórico
                    </Typography>
                  </Box>
                  <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 2 }}>
                    Acesse suas simulações anteriores e compare resultados ao longo do tempo
                  </Typography>
                  <Button
                    variant="outlined"
                    disabled
                    sx={{
                      borderColor: 'var(--md-sys-color-outline)',
                      color: 'var(--md-sys-color-on-surface-variant)',
                      borderRadius: '12px',
                      textTransform: 'none',
                      fontWeight: 500,
                      px: 3,
                    }}
                  >
                    Em breve
                  </Button>
                </CardContent>
              </Card>
            </Box>

            {/* Tendências */}
            <Box>
              <Card
                sx={{
                  borderRadius: '16px',
                  border: '1px solid var(--md-sys-color-outline)',
                  boxShadow: 'var(--md-sys-elevation-1)',
                  height: '100%',
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <TrendingUpIcon sx={{ fontSize: 32, color: 'var(--md-sys-color-success)', mr: 1 }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Tendências de Mercado
                    </Typography>
                  </Box>
                  <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 2 }}>
                    Análise de preços e demanda para seus produtos nos principais mercados
                  </Typography>
                  <Button
                    variant="outlined"
                    disabled
                    sx={{
                      borderColor: 'var(--md-sys-color-outline)',
                      color: 'var(--md-sys-color-on-surface-variant)',
                      borderRadius: '12px',
                      textTransform: 'none',
                      fontWeight: 500,
                      px: 3,
                    }}
                  >
                    Em breve
                  </Button>
                </CardContent>
              </Card>
            </Box>

            {/* Insights */}
            <Box>
              <Card
                sx={{
                  borderRadius: '16px',
                  border: '1px solid var(--md-sys-color-outline)',
                  boxShadow: 'var(--md-sys-elevation-1)',
                  height: '100%',
                }}
              >
                <CardContent sx={{ p: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                    <InsightsIcon sx={{ fontSize: 32, color: 'var(--md-sys-color-secondary)', mr: 1 }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Insights Personalizados
                    </Typography>
                  </Box>
                  <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 2 }}>
                    Recomendações baseadas no seu perfil e histórico de simulações
                  </Typography>
                  <Button
                    variant="outlined"
                    disabled
                    sx={{
                      borderColor: 'var(--md-sys-color-outline)',
                      color: 'var(--md-sys-color-on-surface-variant)',
                      borderRadius: '12px',
                      textTransform: 'none',
                      fontWeight: 500,
                      px: 3,
                    }}
                  >
                    Em breve
                  </Button>
                </CardContent>
              </Card>
            </Box>
          </Box>
        </Container>
      </Box>
    </Box>
  );
}
