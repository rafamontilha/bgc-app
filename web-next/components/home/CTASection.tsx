'use client';

/**
 * CTA Section Component
 * Final call-to-action section before footer
 * Matches original globalconnect-bridge structure
 */

import React from 'react';
import { Box, Container, Typography, Button } from '@mui/material';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';

export interface CTASectionProps {
  onCreateAccountClick?: () => void;
}

export function CTASection({
  onCreateAccountClick,
}: CTASectionProps): React.ReactElement {
  const handleCreateAccount = (): void => {
    if (onCreateAccountClick) {
      onCreateAccountClick();
    } else {
      // Default: navigate to register page
      window.location.href = '/register';
    }
  };

  return (
    <Box
      component="section"
      sx={{
        background: 'linear-gradient(135deg, #007AFF 0%, #5856D6 100%)',
        pt: { xs: 8, md: 12 },
        pb: { xs: 8, md: 12 },
      }}
    >
      <Container maxWidth="lg">
        <Box
          sx={{
            textAlign: 'center',
            color: 'white',
          }}
        >
          {/* Heading */}
          <Typography
            variant="h2"
            sx={{
              fontWeight: 700,
              fontSize: { xs: '2rem', sm: '2.5rem', md: '3rem' },
              lineHeight: 1.2,
              mb: 3,
            }}
          >
            Pronto para expandir globalmente?
          </Typography>

          {/* Subtitle */}
          <Typography
            variant="h5"
            sx={{
              fontWeight: 400,
              fontSize: { xs: '1.125rem', sm: '1.25rem' },
              lineHeight: 1.6,
              maxWidth: 700,
              mx: 'auto',
              mb: 5,
              opacity: 0.95,
            }}
          >
            Junte-se a centenas de empresas brasileiras que já exportam com o
            BGC
          </Typography>

          {/* CTA Button */}
          <Button
            variant="contained"
            size="large"
            onClick={handleCreateAccount}
            endIcon={<ArrowForwardIcon />}
            sx={{
              borderRadius: 999,
              px: 5,
              py: 2,
              fontSize: '1.125rem',
              fontWeight: 600,
              textTransform: 'none',
              backgroundColor: 'white',
              color: 'primary.main',
              boxShadow: '0 4px 20px rgba(0, 0, 0, 0.15)',
              '&:hover': {
                backgroundColor: 'rgba(255, 255, 255, 0.95)',
                boxShadow: '0 6px 24px rgba(0, 0, 0, 0.2)',
                transform: 'translateY(-2px)',
              },
              transition: 'all 220ms ease',
            }}
          >
            Criar Conta Gratuita
          </Button>
        </Box>
      </Container>
    </Box>
  );
}
