'use client';

import { SignUp } from '@clerk/nextjs';
import { Box, Container, Typography } from '@mui/material';

export default function SignUpPage() {
  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'var(--md-sys-color-background)',
        py: 4,
      }}
    >
      <Container maxWidth="sm">
        <Box sx={{ mb: 4, textAlign: 'center' }}>
          <Typography
            variant="h3"
            sx={{
              fontWeight: 600,
              color: 'var(--md-sys-color-on-background)',
              mb: 1,
            }}
          >
            Criar Conta
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'var(--md-sys-color-on-surface-variant)',
            }}
          >
            Junte-se à BGC e descubra os melhores destinos para suas exportações
          </Typography>
        </Box>

        <Box
          sx={{
            display: 'flex',
            justifyContent: 'center',
            '& .cl-rootBox': {
              width: '100%',
            },
            '& .cl-card': {
              boxShadow: 'var(--md-sys-elevation-2)',
              borderRadius: '16px',
              border: '1px solid var(--md-sys-color-outline)',
            },
          }}
        >
          <SignUp
            appearance={{
              elements: {
                rootBox: 'cl-rootBox',
                card: 'cl-card',
                headerTitle: { display: 'none' },
                headerSubtitle: { display: 'none' },
                socialButtonsBlockButton: {
                  borderRadius: '12px',
                  border: '1px solid var(--md-sys-color-outline)',
                  fontWeight: 500,
                },
                formButtonPrimary: {
                  backgroundColor: 'var(--md-sys-color-primary)',
                  borderRadius: '12px',
                  fontWeight: 500,
                  fontSize: '16px',
                  padding: '12px',
                  '&:hover': {
                    backgroundColor: 'var(--md-sys-color-primary)',
                    opacity: 0.9,
                  },
                },
                footerActionLink: {
                  color: 'var(--md-sys-color-primary)',
                },
              },
            }}
          />
        </Box>
      </Container>
    </Box>
  );
}
