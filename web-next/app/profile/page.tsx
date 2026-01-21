'use client';

import { UserProfile } from '@clerk/nextjs';
import { Box, Container, Typography } from '@mui/material';

export default function ProfilePage() {
  return (
    <Box
      sx={{
        minHeight: '100vh',
        bgcolor: 'var(--md-sys-color-background)',
        py: 4,
      }}
    >
      <Container maxWidth="lg">
        <Box sx={{ mb: 4 }}>
          <Typography
            variant="h4"
            sx={{
              fontWeight: 600,
              color: 'var(--md-sys-color-on-background)',
              mb: 1,
            }}
          >
            Meu Perfil
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'var(--md-sys-color-on-surface-variant)',
            }}
          >
            Gerencie suas informações pessoais e preferências
          </Typography>
        </Box>

        <Box
          sx={{
            '& .cl-rootBox': {
              width: '100%',
            },
            '& .cl-card': {
              boxShadow: 'var(--md-sys-elevation-1)',
              borderRadius: '16px',
              border: '1px solid var(--md-sys-color-outline)',
            },
          }}
        >
          <UserProfile
            appearance={{
              elements: {
                rootBox: 'cl-rootBox',
                card: 'cl-card',
                navbar: {
                  borderRadius: '16px 0 0 16px',
                },
                navbarButton: {
                  borderRadius: '8px',
                  '&[data-active]': {
                    backgroundColor: 'var(--md-sys-color-surface-variant)',
                  },
                },
                formButtonPrimary: {
                  backgroundColor: 'var(--md-sys-color-primary)',
                  borderRadius: '12px',
                  fontWeight: 500,
                  '&:hover': {
                    backgroundColor: 'var(--md-sys-color-primary)',
                    opacity: 0.9,
                  },
                },
              },
            }}
          />
        </Box>
      </Container>
    </Box>
  );
}
