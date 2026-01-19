'use client';

/**
 * Features Section Component
 * Showcases 3 main value propositions of the platform
 * Matches original globalconnect-bridge structure
 */

import React from 'react';
import {
  Box,
  Container,
  Typography,
  Card,
  CardContent,
  Avatar,
} from '@mui/material';
import TrendingUpIcon from '@mui/icons-material/TrendingUp';
import ShieldIcon from '@mui/icons-material/Shield';
import PublicIcon from '@mui/icons-material/Public';

interface Feature {
  icon: React.ReactNode;
  title: string;
  description: string;
  color: string;
}

const features: Feature[] = [
  {
    icon: <TrendingUpIcon sx={{ fontSize: 48 }} />,
    title: 'Inteligência de Mercado',
    description:
      'Dados em tempo real sobre demanda, preços e tendências dos principais mercados globais',
    color: 'primary',
  },
  {
    icon: <ShieldIcon sx={{ fontSize: 48 }} />,
    title: 'Segurança e Conformidade',
    description:
      'KYC/KYB certificado, mTLS e conformidade com LGPD para proteger suas operações',
    color: 'secondary',
  },
  {
    icon: <PublicIcon sx={{ fontSize: 48 }} />,
    title: 'Rede Global',
    description:
      'Conecte-se com importadores, distribuidores e parceiros em mais de 50 países',
    color: 'success',
  },
];

export function FeaturesSection(): React.ReactElement {
  return (
    <Box
      id="features"
      component="section"
      sx={{
        backgroundColor: 'background.default',
        pt: { xs: 8, md: 12 },
        pb: { xs: 8, md: 12 },
      }}
    >
      <Container maxWidth="lg">
        {/* Section Header */}
        <Box sx={{ mb: 8 }}>
          <Typography
            variant="h2"
            sx={{
              fontWeight: 700,
              fontSize: { xs: '2rem', sm: '2.5rem', md: '3rem' },
              lineHeight: 1.2,
              color: 'text.primary',
              textAlign: 'left',
            }}
          >
            Por que escolher o BGC?
          </Typography>
        </Box>

        {/* Feature Cards */}
        <Box
          sx={{
            display: 'grid',
            gridTemplateColumns: {
              xs: '1fr',
              sm: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
            },
            gap: 4,
          }}
        >
          {features.map((feature, index) => (
            <Card
              key={index}
              sx={{
                display: 'flex',
                flexDirection: 'column',
                minHeight: 320,
                backgroundColor: 'background.paper',
                transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                '&:hover': {
                  transform: 'scale(1.02)',
                  boxShadow: 3,
                },
              }}
            >
              <CardContent
                sx={{
                  flex: 1,
                  display: 'flex',
                  flexDirection: 'column',
                  p: 4,
                }}
              >
                {/* Icon */}
                <Avatar
                  sx={{
                    width: 72,
                    height: 72,
                    bgcolor: `${feature.color}.main`,
                    color: `${feature.color}.contrastText`,
                    mb: 3,
                  }}
                >
                  {feature.icon}
                </Avatar>

                {/* Title */}
                <Typography
                  variant="h5"
                  sx={{
                    fontWeight: 600,
                    fontSize: '1.5rem',
                    color: 'text.primary',
                    mb: 2,
                    textAlign: 'left',
                  }}
                >
                  {feature.title}
                </Typography>

                {/* Description */}
                <Typography
                  variant="body1"
                  sx={{
                    color: 'text.secondary',
                    lineHeight: 1.7,
                    fontSize: '1rem',
                    textAlign: 'left',
                  }}
                >
                  {feature.description}
                </Typography>
              </CardContent>
            </Card>
          ))}
        </Box>
      </Container>
    </Box>
  );
}
