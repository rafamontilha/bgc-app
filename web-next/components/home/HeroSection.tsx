'use client';

/**
 * Hero Section Component
 * Main value proposition section (above the fold)
 * Matches original globalconnect-bridge structure
 */

import React from 'react';
import { Box, Container, Typography, Button, Stack } from '@mui/material';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';

export interface HeroSectionProps {
  onGetStartedClick?: () => void;
  onLearnMoreClick?: () => void;
}

export function HeroSection({
  onGetStartedClick,
  onLearnMoreClick,
}: HeroSectionProps): React.ReactElement {
  const handleGetStarted = (): void => {
    if (onGetStartedClick) {
      onGetStartedClick();
    } else {
      // Default: scroll to simulator section
      const simulatorSection = document.getElementById('simulator');
      if (simulatorSection) {
        simulatorSection.scrollIntoView({ behavior: 'smooth' });
      }
    }
  };

  const handleLearnMore = (): void => {
    if (onLearnMoreClick) {
      onLearnMoreClick();
    } else {
      // Default: scroll to features section
      const featuresSection = document.getElementById('features');
      if (featuresSection) {
        featuresSection.scrollIntoView({ behavior: 'smooth' });
      }
    }
  };

  return (
    <Box
      component="section"
      sx={{
        background: 'linear-gradient(135deg, rgba(0, 122, 255, 0.08) 0%, rgba(88, 86, 214, 0.08) 100%)',
        pt: { xs: 10, md: 16 },
        pb: { xs: 10, md: 16 },
      }}
    >
      <Container maxWidth="lg">
        <Box
          sx={{
            textAlign: 'center',
            maxWidth: 900,
            mx: 'auto',
          }}
        >
          {/* Main Heading */}
          <Typography
            variant="h1"
            sx={{
              fontWeight: 700,
              fontSize: { xs: '2.5rem', sm: '3.5rem', md: '4rem' },
              lineHeight: 1.1,
              color: 'text.primary',
              mb: 3,
              letterSpacing: '-0.02em',
            }}
          >
            Conecte Seu Negócio ao Mundo
          </Typography>

          {/* Subtitle */}
          <Typography
            variant="h5"
            sx={{
              fontWeight: 400,
              fontSize: { xs: '1.125rem', sm: '1.25rem', md: '1.5rem' },
              color: 'text.secondary',
              lineHeight: 1.6,
              maxWidth: 800,
              mx: 'auto',
              mb: 5,
            }}
          >
            Descubra oportunidades de exportação, analise mercados e conecte-se
            com compradores globais
          </Typography>

          {/* CTA Buttons */}
          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            spacing={2}
            justifyContent="center"
            alignItems="center"
          >
            <Button
              variant="contained"
              size="large"
              onClick={handleGetStarted}
              endIcon={<ArrowForwardIcon />}
              sx={{
                px: 4,
                py: 1.75,
                fontSize: '1.125rem',
                fontWeight: 600,
                transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                '&:hover': {
                  transform: 'scale(1.02)',
                  boxShadow: 2,
                },
              }}
            >
              Comece Grátis
            </Button>

            <Button
              variant="outlined"
              size="large"
              onClick={handleLearnMore}
              sx={{
                px: 4,
                py: 1.75,
                fontSize: '1.125rem',
                fontWeight: 600,
                borderColor: 'primary.main',
                color: 'primary.main',
                transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                '&:hover': {
                  borderColor: 'primary.dark',
                  backgroundColor: 'action.hover',
                  transform: 'scale(1.02)',
                },
              }}
            >
              Saiba Mais
            </Button>
          </Stack>
        </Box>
      </Container>
    </Box>
  );
}
