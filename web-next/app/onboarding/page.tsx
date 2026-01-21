'use client';

import { useUser } from '@clerk/nextjs';
import { Box, Container, Typography, Button, Card, CardContent, Stepper, Step, StepLabel } from '@mui/material';
import { useRouter } from 'next/navigation';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';

const steps = ['Informações da Empresa', 'Produto Principal (NCM)', 'Preferências de Mercado'];

export default function OnboardingPage() {
  const { user, isLoaded } = useUser();
  const router = useRouter();

  if (!isLoaded) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
        <Typography>Carregando...</Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        bgcolor: 'var(--md-sys-color-background)',
        py: 4,
      }}
    >
      <Container maxWidth="md">
        <Box sx={{ mb: 4, textAlign: 'center' }}>
          <RocketLaunchIcon sx={{ fontSize: 64, color: 'var(--md-sys-color-primary)', mb: 2 }} />
          <Typography
            variant="h4"
            sx={{
              fontWeight: 600,
              color: 'var(--md-sys-color-on-background)',
              mb: 1,
            }}
          >
            Bem-vindo, {user?.firstName}!
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'var(--md-sys-color-on-surface-variant)',
            }}
          >
            Vamos configurar seu perfil para personalizar suas recomendações
          </Typography>
        </Box>

        <Card
          sx={{
            borderRadius: '16px',
            border: '1px solid var(--md-sys-color-outline)',
            boxShadow: 'var(--md-sys-elevation-2)',
            mb: 4,
          }}
        >
          <CardContent sx={{ p: 4 }}>
            <Stepper activeStep={0} sx={{ mb: 4 }}>
              {steps.map((label) => (
                <Step key={label}>
                  <StepLabel>{label}</StepLabel>
                </Step>
              ))}
            </Stepper>

            <Box sx={{ textAlign: 'center', py: 4 }}>
              <Typography variant="h6" sx={{ mb: 2, fontWeight: 600 }}>
                Wizard de Onboarding
              </Typography>
              <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 4 }}>
                Esta página será implementada na J-AC01 (próxima sprint)
              </Typography>
              <Typography variant="body2" sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 4 }}>
                Por enquanto, você já pode usar o simulador!
              </Typography>
              <Button
                variant="contained"
                onClick={() => router.push('/dashboard')}
                sx={{
                  bgcolor: 'var(--md-sys-color-primary)',
                  color: 'var(--md-sys-color-on-primary)',
                  borderRadius: '12px',
                  textTransform: 'none',
                  fontWeight: 500,
                  px: 4,
                  py: 1.5,
                  '&:hover': {
                    bgcolor: 'var(--md-sys-color-primary)',
                    opacity: 0.9,
                  },
                }}
              >
                Pular para Dashboard
              </Button>
            </Box>
          </CardContent>
        </Card>
      </Container>
    </Box>
  );
}
