'use client';

export const dynamic = 'force-dynamic';

/**
 * Custom Sign Up Page
 * Uses Clerk hooks with MUI components following Apple HIG design system
 */

import React, { useState } from 'react';
import { useSignUp } from '@clerk/nextjs';
import { useRouter } from 'next/navigation';
import {
  Box,
  TextField,
  Button,
  Typography,
  Paper,
  Stack,
  Divider,
  CircularProgress,
  Alert,
  Link as MuiLink,
  InputAdornment,
  IconButton,
} from '@mui/material';
import GoogleIcon from '@mui/icons-material/Google';
import EmailOutlinedIcon from '@mui/icons-material/EmailOutlined';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import VisibilityIcon from '@mui/icons-material/Visibility';
import VisibilityOffIcon from '@mui/icons-material/VisibilityOff';
import Link from 'next/link';
import type { OAuthStrategy } from '@clerk/types';

type SignUpStep = 'form' | 'verify';

export default function SignUpPage(): React.ReactElement {
  const { isLoaded, signUp, setActive } = useSignUp();
  const router = useRouter();

  const [step, setStep] = useState<SignUpStep>('form');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [verificationCode, setVerificationCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Handle OAuth sign up (Google)
   */
  const handleOAuthSignUp = async (strategy: OAuthStrategy): Promise<void> => {
    if (!isLoaded || !signUp) return;

    try {
      await signUp.authenticateWithRedirect({
        strategy,
        redirectUrl: '/sso-callback',
        redirectUrlComplete: '/onboarding',
      });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao conectar com o provedor';
      setError(errorMessage);
    }
  };

  /**
   * Handle email/password sign up
   */
  const handleEmailSignUp = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!isLoaded || !signUp) return;

    setIsLoading(true);
    setError(null);

    try {
      await signUp.create({
        emailAddress: email,
        password,
      });

      await signUp.prepareEmailAddressVerification({
        strategy: 'email_code',
      });

      setStep('verify');
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string; code: string }> };
      if (clerkError.errors?.[0]) {
        const errorCode = clerkError.errors[0].code;
        if (errorCode === 'form_identifier_exists') {
          setError('Este email já está cadastrado. Tente fazer login.');
        } else if (errorCode === 'form_password_pwned') {
          setError('Esta senha foi comprometida. Escolha uma senha mais segura.');
        } else if (errorCode === 'form_password_length_too_short') {
          setError('A senha deve ter pelo menos 8 caracteres.');
        } else {
          setError(clerkError.errors[0].message);
        }
      } else {
        setError('Erro ao criar conta. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle email verification
   */
  const handleVerification = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!isLoaded || !signUp) return;

    setIsLoading(true);
    setError(null);

    try {
      const result = await signUp.attemptEmailAddressVerification({
        code: verificationCode,
      });

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId });
        router.push('/onboarding');
      }
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string; code: string }> };
      if (clerkError.errors?.[0]?.code === 'form_code_incorrect') {
        setError('Código incorreto. Verifique e tente novamente.');
      } else {
        setError('Erro na verificação. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Resend verification code
   */
  const handleResendCode = async (): Promise<void> => {
    if (!isLoaded || !signUp) return;

    try {
      await signUp.prepareEmailAddressVerification({
        strategy: 'email_code',
      });
      setError(null);
    } catch {
      setError('Erro ao reenviar código.');
    }
  };

  if (!isLoaded) {
    return (
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          bgcolor: 'background.default',
        }}
      >
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box
      sx={{
        minHeight: '100vh',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
        py: 4,
        px: 2,
      }}
    >
      <Paper
        elevation={2}
        sx={{
          width: '100%',
          maxWidth: 440,
          p: 4,
          borderRadius: 4,
        }}
      >
        {step === 'form' ? (
          <Stack spacing={3}>
            {/* Header */}
            <Box sx={{ textAlign: 'center' }}>
              <Typography
                variant="h4"
                sx={{
                  fontWeight: 600,
                  color: 'text.primary',
                  mb: 1,
                }}
              >
                Criar Conta
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: 'text.secondary' }}
              >
                Descubra os melhores destinos para suas exportações
              </Typography>
            </Box>

            {/* Error Alert */}
            {error && (
              <Alert severity="error" onClose={() => setError(null)}>
                {error}
              </Alert>
            )}

            {/* Google OAuth Button */}
            <Button
              variant="outlined"
              size="large"
              fullWidth
              onClick={() => handleOAuthSignUp('oauth_google')}
              startIcon={<GoogleIcon />}
              sx={{
                py: 1.5,
                borderColor: 'divider',
                color: 'text.primary',
                '&:hover': {
                  borderColor: 'primary.main',
                  bgcolor: 'action.hover',
                },
              }}
            >
              Continuar com Google
            </Button>

            {/* Divider */}
            <Divider>
              <Typography variant="body2" sx={{ color: 'text.secondary', px: 2 }}>
                ou
              </Typography>
            </Divider>

            {/* Email/Password Form */}
            <Box component="form" onSubmit={handleEmailSignUp}>
              <Stack spacing={2.5}>
                <TextField
                  label="Email"
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  fullWidth
                  required
                  disabled={isLoading}
                  autoComplete="email"
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">
                          <EmailOutlinedIcon sx={{ color: 'text.secondary' }} />
                        </InputAdornment>
                      ),
                    },
                  }}
                />

                <TextField
                  label="Senha"
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  fullWidth
                  required
                  disabled={isLoading}
                  autoComplete="new-password"
                  helperText="Mínimo 8 caracteres"
                  slotProps={{
                    input: {
                      startAdornment: (
                        <InputAdornment position="start">
                          <LockOutlinedIcon sx={{ color: 'text.secondary' }} />
                        </InputAdornment>
                      ),
                      endAdornment: (
                        <InputAdornment position="end">
                          <IconButton
                            onClick={() => setShowPassword(!showPassword)}
                            edge="end"
                            size="small"
                          >
                            {showPassword ? <VisibilityOffIcon /> : <VisibilityIcon />}
                          </IconButton>
                        </InputAdornment>
                      ),
                    },
                  }}
                />

                {/* Ponto de montagem do Smart CAPTCHA do Clerk (bot protection em custom flows) */}
                <div id="clerk-captcha" />

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isLoading || !email || !password}
                  sx={{ py: 1.5, mt: 1 }}
                >
                  {isLoading ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    'Criar conta'
                  )}
                </Button>
              </Stack>
            </Box>

            {/* Footer */}
            <Typography
              variant="body2"
              sx={{ textAlign: 'center', color: 'text.secondary' }}
            >
              Já tem uma conta?{' '}
              <MuiLink
                component={Link}
                href="/login"
                sx={{ fontWeight: 500 }}
              >
                Fazer login
              </MuiLink>
            </Typography>
          </Stack>
        ) : (
          /* Verification Step */
          <Stack spacing={3}>
            <Box sx={{ textAlign: 'center' }}>
              <Typography
                variant="h4"
                sx={{
                  fontWeight: 600,
                  color: 'text.primary',
                  mb: 1,
                }}
              >
                Verifique seu email
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: 'text.secondary' }}
              >
                Enviamos um código de 6 dígitos para{' '}
                <strong>{email}</strong>
              </Typography>
            </Box>

            {error && (
              <Alert severity="error" onClose={() => setError(null)}>
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleVerification}>
              <Stack spacing={2.5}>
                <TextField
                  label="Código de verificação"
                  value={verificationCode}
                  onChange={(e) => setVerificationCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  fullWidth
                  required
                  disabled={isLoading}
                  autoComplete="one-time-code"
                  placeholder="000000"
                  slotProps={{
                    htmlInput: {
                      maxLength: 6,
                      style: {
                        textAlign: 'center',
                        letterSpacing: '0.5em',
                        fontSize: '1.5rem',
                        fontWeight: 600,
                      },
                    },
                  }}
                />

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isLoading || verificationCode.length !== 6}
                  sx={{ py: 1.5 }}
                >
                  {isLoading ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    'Verificar'
                  )}
                </Button>
              </Stack>
            </Box>

            <Typography
              variant="body2"
              sx={{ textAlign: 'center', color: 'text.secondary' }}
            >
              Não recebeu o código?{' '}
              <MuiLink
                component="button"
                type="button"
                onClick={handleResendCode}
                sx={{ fontWeight: 500 }}
              >
                Reenviar
              </MuiLink>
            </Typography>

            <Button
              variant="text"
              onClick={() => setStep('form')}
              sx={{ color: 'text.secondary' }}
            >
              Voltar
            </Button>
          </Stack>
        )}
      </Paper>
    </Box>
  );
}
