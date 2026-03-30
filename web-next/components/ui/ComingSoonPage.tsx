'use client';

import React from 'react';
import { Box, Typography, Button, Chip } from '@mui/material';
import ConstructionIcon from '@mui/icons-material/Construction';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import { useRouter } from 'next/navigation';

interface ComingSoonPageProps {
  title: string;
  description?: string;
}

export function ComingSoonPage({ title, description }: ComingSoonPageProps): React.ReactElement {
  const router = useRouter();

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
        px: 3,
        gap: 3,
        textAlign: 'center',
      }}
    >
      <Box
        sx={{
          width: 80,
          height: 80,
          borderRadius: 4,
          bgcolor: 'primary.main',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
        }}
      >
        <ConstructionIcon sx={{ fontSize: 40, color: 'primary.contrastText' }} />
      </Box>

      <Chip label="Em breve" color="primary" variant="outlined" />

      <Box>
        <Typography variant="h4" sx={{ fontWeight: 600, color: 'text.primary', mb: 1 }}>
          {title}
        </Typography>
        <Typography variant="body1" sx={{ color: 'text.secondary', maxWidth: 480 }}>
          {description ?? 'Estamos trabalhando nesta funcionalidade. Em breve estará disponível para você.'}
        </Typography>
      </Box>

      <Button
        variant="outlined"
        startIcon={<ArrowBackIcon />}
        onClick={() => router.back()}
        sx={{ mt: 1 }}
      >
        Voltar
      </Button>
    </Box>
  );
}
