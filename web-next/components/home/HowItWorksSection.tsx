'use client';

/**
 * How It Works Section Component
 * Explains the 3-step process to use the simulator
 */

import React from 'react';
import { Box, Container, Typography, Stack, Card, CardContent } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';

interface StepCardProps {
  stepNumber: number;
  icon: React.ReactElement;
  title: string;
  description: string;
}

function StepCard({
  stepNumber,
  icon,
  title,
  description,
}: StepCardProps): React.ReactElement {
  return (
    <Card
      sx={{
        position: 'relative',
        borderRadius: 4,
        boxShadow: '0px 4px 24px rgba(0, 0, 0, 0.10)',
        transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
        height: '100%',
        '&:hover': {
          transform: 'scale(1.02)',
          boxShadow: '0px 8px 32px rgba(0, 0, 0, 0.15)',
        },
      }}
    >
      <CardContent
        sx={{
          p: 4,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          textAlign: 'center',
          gap: 2.5,
        }}
      >
        {/* Step number badge */}
        <Box
          sx={{
            position: 'absolute',
            top: 16,
            right: 16,
            width: 32,
            height: 32,
            borderRadius: 999,
            backgroundColor: 'rgba(0, 122, 255, 0.1)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Typography
            variant="caption"
            sx={{
              fontWeight: 700,
              color: 'primary.main',
              fontSize: '0.875rem',
            }}
          >
            {stepNumber}
          </Typography>
        </Box>

        {/* Icon */}
        <Box
          sx={{
            width: 72,
            height: 72,
            borderRadius: 4,
            backgroundColor: 'primary.main',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            boxShadow: '0px 4px 20px rgba(0, 122, 255, 0.3)',
          }}
        >
          <Box sx={{ color: 'white', fontSize: 36, display: 'flex' }}>
            {icon}
          </Box>
        </Box>

        {/* Content */}
        <Box>
          <Typography
            variant="h5"
            sx={{
              fontWeight: 600,
              color: 'text.primary',
              mb: 1.5,
            }}
          >
            {title}
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'text.secondary',
              lineHeight: 1.7,
            }}
          >
            {description}
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
}

export function HowItWorksSection(): React.ReactElement {
  const steps = [
    {
      icon: <SearchIcon />,
      title: 'Informe seu NCM',
      description:
        'Digite o código NCM de 8 dígitos do seu produto e o volume que deseja exportar',
    },
    {
      icon: <TrendingUpIcon />,
      title: 'Veja destinos ranqueados',
      description:
        'Algoritmo analisa tamanho de mercado, crescimento, preço e distância logística',
    },
    {
      icon: <CheckCircleIcon />,
      title: 'Tome decisões baseadas em dados',
      description:
        'Escolha o melhor destino com confiança usando dados oficiais do ComexStat',
    },
  ];

  return (
    <Box
      id="how-it-works"
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
              Como Funciona
            </Typography>
            <Typography
              variant="body1"
              sx={{
                color: 'text.secondary',
                fontSize: '1.125rem',
                lineHeight: 1.7,
              }}
            >
              Três passos simples para descobrir os melhores mercados para seu
              produto
            </Typography>
          </Box>

          {/* Steps grid */}
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: {
                xs: '1fr',
                md: 'repeat(3, 1fr)',
              },
              gap: 4,
            }}
          >
            {steps.map((step, index) => (
              <StepCard
                key={index}
                stepNumber={index + 1}
                icon={step.icon}
                title={step.title}
                description={step.description}
              />
            ))}
          </Box>
        </Stack>
      </Container>
    </Box>
  );
}
