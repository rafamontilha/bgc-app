'use client';

/**
 * Benefits Section Component
 * Highlights 4 key benefits of the platform
 */

import React from 'react';
import { Box, Container, Typography, Stack, Card, CardContent } from '@mui/material';
import StorageIcon from '@mui/icons-material/Storage';
import AutoAwesomeIcon from '@mui/icons-material/AutoAwesome';
import SpeedIcon from '@mui/icons-material/Speed';
import StarsIcon from '@mui/icons-material/Stars';

interface BenefitCardProps {
  icon: React.ReactElement;
  iconColor: string;
  iconBgColor: string;
  title: string;
  description: string;
}

function BenefitCard({
  icon,
  iconColor,
  iconBgColor,
  title,
  description,
}: BenefitCardProps): React.ReactElement {
  return (
    <Card
      sx={{
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
          gap: 2.5,
        }}
      >
        {/* Icon */}
        <Box
          sx={{
            width: 64,
            height: 64,
            borderRadius: 3,
            backgroundColor: iconBgColor,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
          }}
        >
          <Box sx={{ color: iconColor, fontSize: 32, display: 'flex' }}>
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

export function BenefitsSection(): React.ReactElement {
  const benefits = [
    {
      icon: <StorageIcon />,
      iconColor: '#007AFF',
      iconBgColor: 'rgba(0, 122, 255, 0.1)',
      title: 'Dados Reais',
      description:
        'Fonte oficial ComexStat (Governo Federal). Dados atualizados de exportações brasileiras para decisões confiáveis.',
    },
    {
      icon: <AutoAwesomeIcon />,
      iconColor: '#5856D6',
      iconBgColor: 'rgba(88, 86, 214, 0.1)',
      title: 'Análise Inteligente',
      description:
        'Algoritmo de scoring automatizado que pondera 4 variáveis principais: tamanho de mercado, crescimento, preço e distância.',
    },
    {
      icon: <SpeedIcon />,
      iconColor: '#34C759',
      iconBgColor: 'rgba(52, 199, 89, 0.1)',
      title: 'Decisões Rápidas',
      description:
        'Resultados em menos de 100ms. Performance 50x melhor que o esperado para análises complexas de mercado.',
    },
    {
      icon: <StarsIcon />,
      iconColor: '#FF9500',
      iconBgColor: 'rgba(255, 149, 0, 0.1)',
      title: 'Gratuito para Começar',
      description:
        '5 simulações por dia no plano gratuito. Teste a plataforma sem compromisso e veja os resultados antes de decidir.',
    },
  ];

  return (
    <Box
      id="benefits"
      component="section"
      sx={{
        backgroundColor: 'background.default',
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
              Por Que Escolher Brasil Global Connect
            </Typography>
            <Typography
              variant="body1"
              sx={{
                color: 'text.secondary',
                fontSize: '1.125rem',
                lineHeight: 1.7,
              }}
            >
              A plataforma mais completa para identificar oportunidades de
              exportação
            </Typography>
          </Box>

          {/* Benefits grid */}
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: {
                xs: '1fr',
                sm: 'repeat(2, 1fr)',
              },
              gap: 4,
            }}
          >
            {benefits.map((benefit, index) => (
              <BenefitCard
                key={index}
                icon={benefit.icon}
                iconColor={benefit.iconColor}
                iconBgColor={benefit.iconBgColor}
                title={benefit.title}
                description={benefit.description}
              />
            ))}
          </Box>
        </Stack>
      </Container>
    </Box>
  );
}
