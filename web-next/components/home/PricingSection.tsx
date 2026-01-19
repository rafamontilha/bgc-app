'use client';

/**
 * Pricing Section Component
 * Displays Free and Pro pricing plans
 */

import React from 'react';
import {
  Box,
  Container,
  Typography,
  Stack,
  Card,
  CardContent,
  Button,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Chip,
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CloseIcon from '@mui/icons-material/Close';
import RocketLaunchIcon from '@mui/icons-material/RocketLaunch';

export interface PricingSectionProps {
  onUpgradeClick?: () => void;
  onFreeClick?: () => void;
}

interface PlanFeature {
  label: string;
  included: boolean;
}

interface PricingCardProps {
  title: string;
  price: string;
  priceSubtext: string;
  features: PlanFeature[];
  ctaLabel: string;
  ctaVariant: 'contained' | 'outlined';
  highlighted?: boolean;
  onCtaClick?: () => void;
}

function PricingCard({
  title,
  price,
  priceSubtext,
  features,
  ctaLabel,
  ctaVariant,
  highlighted = false,
  onCtaClick,
}: PricingCardProps): React.ReactElement {
  return (
    <Card
      sx={{
        position: 'relative',
        borderRadius: 4,
        boxShadow: highlighted
          ? '0px 12px 48px rgba(0, 122, 255, 0.25)'
          : '0px 4px 24px rgba(0, 0, 0, 0.10)',
        border: highlighted ? 2 : 0,
        borderColor: highlighted ? 'primary.main' : 'transparent',
        transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
        height: '100%',
        '&:hover': {
          transform: 'scale(1.02)',
          boxShadow: highlighted
            ? '0px 16px 56px rgba(0, 122, 255, 0.3)'
            : '0px 8px 32px rgba(0, 0, 0, 0.15)',
        },
      }}
    >
      {highlighted && (
        <Box
          sx={{
            position: 'absolute',
            top: -12,
            left: '50%',
            transform: 'translateX(-50%)',
          }}
        >
          <Chip
            label="Mais Popular"
            color="primary"
            sx={{
              fontWeight: 600,
              fontSize: '0.875rem',
              px: 1.5,
            }}
          />
        </Box>
      )}

      <CardContent
        sx={{
          p: 4,
          display: 'flex',
          flexDirection: 'column',
          gap: 3,
        }}
      >
        {/* Plan title */}
        <Box>
          <Typography
            variant="h4"
            sx={{
              fontWeight: 600,
              color: 'text.primary',
              mb: 0.5,
            }}
          >
            {title}
          </Typography>
        </Box>

        {/* Price */}
        <Box>
          <Typography
            variant="h2"
            sx={{
              fontWeight: 700,
              color: 'text.primary',
              fontSize: '2.5rem',
              mb: 0.5,
            }}
          >
            {price}
          </Typography>
          <Typography
            variant="body2"
            sx={{
              color: 'text.secondary',
            }}
          >
            {priceSubtext}
          </Typography>
        </Box>

        {/* Features list */}
        <List sx={{ py: 0 }}>
          {features.map((feature, index) => (
            <ListItem
              key={index}
              sx={{
                px: 0,
                py: 1,
              }}
            >
              <ListItemIcon sx={{ minWidth: 36 }}>
                {feature.included ? (
                  <CheckCircleIcon
                    sx={{
                      color: 'success.main',
                      fontSize: 24,
                    }}
                  />
                ) : (
                  <CloseIcon
                    sx={{
                      color: 'text.disabled',
                      fontSize: 24,
                    }}
                  />
                )}
              </ListItemIcon>
              <ListItemText
                primary={feature.label}
                primaryTypographyProps={{
                  fontSize: '1rem',
                  color: feature.included ? 'text.primary' : 'text.disabled',
                  fontWeight: feature.included ? 500 : 400,
                }}
              />
            </ListItem>
          ))}
        </List>

        {/* CTA button */}
        <Button
          variant={ctaVariant}
          size="large"
          fullWidth
          onClick={onCtaClick}
          startIcon={highlighted ? <RocketLaunchIcon /> : undefined}
          sx={{
            mt: 'auto',
            borderRadius: 999,
            py: 1.75,
            fontWeight: 600,
            fontSize: '1rem',
            boxShadow: highlighted ? 3 : 0,
            '&:hover': {
              boxShadow: highlighted ? 6 : 0,
            },
          }}
        >
          {ctaLabel}
        </Button>
      </CardContent>
    </Card>
  );
}

export function PricingSection({
  onUpgradeClick,
  onFreeClick,
}: PricingSectionProps): React.ReactElement {
  const handleUpgradeClick = (): void => {
    if (onUpgradeClick) {
      onUpgradeClick();
    } else {
      // Default: navigate to checkout or contact
      window.location.href = '/checkout?plan=pro';
    }
  };

  const handleFreeClick = (): void => {
    if (onFreeClick) {
      onFreeClick();
    } else {
      // Default: scroll to simulator
      const simulatorSection = document.getElementById('simulator');
      if (simulatorSection) {
        simulatorSection.scrollIntoView({ behavior: 'smooth' });
      }
    }
  };

  const freePlanFeatures: PlanFeature[] = [
    { label: '5 simulações por dia', included: true },
    { label: 'Dados ComexStat oficiais', included: true },
    { label: 'Scoring automatizado', included: true },
    { label: 'Exportação para PDF', included: false },
    { label: 'Dashboard de analytics', included: false },
    { label: 'Histórico de análises', included: false },
    { label: 'Suporte prioritário', included: false },
  ];

  const proPlanFeatures: PlanFeature[] = [
    { label: 'Simulações ilimitadas', included: true },
    { label: 'Dados ComexStat oficiais', included: true },
    { label: 'Scoring automatizado', included: true },
    { label: 'Exportação para PDF', included: true },
    { label: 'Dashboard de analytics', included: true },
    { label: 'Histórico de análises', included: true },
    { label: 'Suporte prioritário', included: true },
  ];

  return (
    <Box
      id="pricing"
      component="section"
      sx={{
        backgroundColor: 'white',
        py: { xs: 8, md: 10 },
      }}
    >
      <Container maxWidth="lg">
        <Stack spacing={6}>
          {/* Section header */}
          <Box sx={{ textAlign: 'center', maxWidth: 700, mx: 'auto' }}>
            <Typography
              variant="h2"
              sx={{
                fontWeight: 600,
                fontSize: { xs: '1.75rem', sm: '2rem', md: '2.25rem' },
                color: 'text.primary',
                mb: 2,
              }}
            >
              Planos Simples e Transparentes
            </Typography>
            <Typography
              variant="body1"
              sx={{
                color: 'text.secondary',
                fontSize: '1.125rem',
                lineHeight: 1.7,
              }}
            >
              Comece grátis ou desbloqueie todo o potencial com o plano Pro
            </Typography>
          </Box>

          {/* Pricing cards */}
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: {
                xs: '1fr',
                md: 'repeat(2, 1fr)',
              },
              gap: 4,
              maxWidth: 1000,
              mx: 'auto',
            }}
          >
            {/* Free plan */}
            <PricingCard
              title="Free"
              price="R$ 0"
              priceSubtext="Grátis para sempre"
              features={freePlanFeatures}
              ctaLabel="Começar Grátis"
              ctaVariant="outlined"
              onCtaClick={handleFreeClick}
            />

            {/* Pro plan */}
            <PricingCard
              title="Pro"
              price="R$ 199"
              priceSubtext="por mês, cancele quando quiser"
              features={proPlanFeatures}
              ctaLabel="Upgrade para Pro"
              ctaVariant="contained"
              highlighted
              onCtaClick={handleUpgradeClick}
            />
          </Box>

          {/* Trust badges */}
          <Box
            sx={{
              textAlign: 'center',
              pt: 4,
            }}
          >
            <Typography
              variant="body2"
              sx={{
                color: 'text.secondary',
                mb: 2,
              }}
            >
              Todos os planos incluem:
            </Typography>
            <Stack
              direction={{ xs: 'column', sm: 'row' }}
              spacing={3}
              sx={{
                justifyContent: 'center',
                alignItems: 'center',
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <CheckCircleIcon sx={{ color: 'success.main', fontSize: 20 }} />
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Dados oficiais ComexStat
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <CheckCircleIcon sx={{ color: 'success.main', fontSize: 20 }} />
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Algoritmo de scoring avançado
                </Typography>
              </Box>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                <CheckCircleIcon sx={{ color: 'success.main', fontSize: 20 }} />
                <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                  Resultados em menos de 100ms
                </Typography>
              </Box>
            </Stack>
          </Box>
        </Stack>
      </Container>
    </Box>
  );
}
