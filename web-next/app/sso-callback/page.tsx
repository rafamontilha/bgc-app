'use client';

export const dynamic = 'force-dynamic';

/**
 * SSO Callback Page
 * Handles OAuth redirect callback from providers (Google, etc.)
 */

import { AuthenticateWithRedirectCallback } from '@clerk/nextjs';
import { Box, CircularProgress, Typography } from '@mui/material';

export default function SSOCallbackPage(): React.ReactElement {
  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
        gap: 2,
      }}
    >
      <CircularProgress size={48} />
      <Typography variant="body1" sx={{ color: 'text.secondary' }}>
        Autenticando...
      </Typography>

      <AuthenticateWithRedirectCallback
        signInFallbackRedirectUrl="/dashboard"
        signUpFallbackRedirectUrl="/onboarding"
      />
    </Box>
  );
}
