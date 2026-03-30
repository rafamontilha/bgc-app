'use client';

export const dynamic = 'force-dynamic';

/**
 * Custom Sign In Page
 * Uses Clerk hooks with MUI components following Apple HIG design system
 */

import React, { useState } from 'react';
import { useSignIn } from '@clerk/nextjs';
import { useRouter, useSearchParams } from 'next/navigation';
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

type SignInStep = 'form' | 'forgot' | 'reset-code';

export default function LoginPage(): React.ReactElement {
  const { isLoaded, signIn, setActive } = useSignIn();
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirectUrl = searchParams.get('redirect_url') || '/dashboard';

  const [step, setStep] = useState<SignInStep>('form');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [resetCode, setResetCode] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  /**
   * Handle OAuth sign in (Google)
   */
  const handleOAuthSignIn = async (strategy: OAuthStrategy): Promise<void> => {
    if (!isLoaded || !signIn) return;

    try {
      await signIn.authenticateWithRedirect({
        strategy,
        redirectUrl: '/sso-callback',
        redirectUrlComplete: redirectUrl,
      });
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao conectar com o provedor';
      setError(errorMessage);
    }
  };

  /**
   * Handle email/password sign in
   */
  const handleEmailSignIn = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!isLoaded || !signIn) return;

    setIsLoading(true);
    setError(null);

    try {
      const result = await signIn.create({
        identifier: email,
        password,
      });

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId });
        router.push(redirectUrl);
      } else if (result.status === 'needs_second_factor') {
        setError('Autenticação de dois fatores necessária. Entre em contato com o suporte.');
      }
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string; code: string }> };
      if (clerkError.errors?.[0]) {
        const errorCode = clerkError.errors[0].code;
        if (errorCode === 'form_identifier_not_found') {
          setError('Email não encontrado. Verifique ou crie uma conta.');
        } else if (errorCode === 'form_password_incorrect') {
          setError('Senha incorreta. Tente novamente ou recupere sua senha.');
        } else if (errorCode === 'strategy_for_user_invalid') {
          setError('Esta conta usa login social. Tente "Continuar com Google".');
        } else {
          setError(clerkError.errors[0].message);
        }
      } else {
        setError('Erro ao fazer login. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle forgot password - send reset code
   */
  const handleForgotPassword = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!isLoaded || !signIn) return;

    setIsLoading(true);
    setError(null);

    try {
      await signIn.create({
        strategy: 'reset_password_email_code',
        identifier: email,
      });
      setStep('reset-code');
      setSuccessMessage('Código enviado para seu email.');
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string; code: string }> };
      if (clerkError.errors?.[0]?.code === 'form_identifier_not_found') {
        setError('Email não encontrado.');
      } else {
        setError('Erro ao enviar código. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  /**
   * Handle password reset with code
   */
  const handleResetPassword = async (e: React.FormEvent): Promise<void> => {
    e.preventDefault();
    if (!isLoaded || !signIn) return;

    setIsLoading(true);
    setError(null);

    try {
      const result = await signIn.attemptFirstFactor({
        strategy: 'reset_password_email_code',
        code: resetCode,
        password: newPassword,
      });

      if (result.status === 'complete') {
        await setActive({ session: result.createdSessionId });
        router.push(redirectUrl);
      }
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string; code: string }> };
      if (clerkError.errors?.[0]?.code === 'form_code_incorrect') {
        setError('Código incorreto. Verifique e tente novamente.');
      } else if (clerkError.errors?.[0]?.code === 'form_password_length_too_short') {
        setError('A nova senha deve ter pelo menos 8 caracteres.');
      } else {
        setError('Erro ao redefinir senha. Tente novamente.');
      }
    } finally {
      setIsLoading(false);
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
        {step === 'form' && (
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
                Bem-vindo de volta
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: 'text.secondary' }}
              >
                Acesse sua conta para continuar explorando oportunidades
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
              onClick={() => handleOAuthSignIn('oauth_google')}
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
            <Box component="form" onSubmit={handleEmailSignIn}>
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
                  autoComplete="current-password"
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

                <Box sx={{ textAlign: 'right' }}>
                  <MuiLink
                    component="button"
                    type="button"
                    variant="body2"
                    onClick={() => {
                      setStep('forgot');
                      setError(null);
                    }}
                    sx={{ fontWeight: 500 }}
                  >
                    Esqueceu a senha?
                  </MuiLink>
                </Box>

                {/* Ponto de montagem do Smart CAPTCHA do Clerk (bot protection em custom flows) */}
                <div id="clerk-captcha" />

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isLoading || !email || !password}
                  sx={{ py: 1.5 }}
                >
                  {isLoading ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    'Entrar'
                  )}
                </Button>
              </Stack>
            </Box>

            {/* Footer */}
            <Typography
              variant="body2"
              sx={{ textAlign: 'center', color: 'text.secondary' }}
            >
              Não tem uma conta?{' '}
              <MuiLink
                component={Link}
                href="/signup"
                sx={{ fontWeight: 500 }}
              >
                Criar conta
              </MuiLink>
            </Typography>
          </Stack>
        )}

        {step === 'forgot' && (
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
                Recuperar senha
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: 'text.secondary' }}
              >
                Digite seu email para receber um código de recuperação
              </Typography>
            </Box>

            {error && (
              <Alert severity="error" onClose={() => setError(null)}>
                {error}
              </Alert>
            )}

            <Box component="form" onSubmit={handleForgotPassword}>
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

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isLoading || !email}
                  sx={{ py: 1.5 }}
                >
                  {isLoading ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    'Enviar código'
                  )}
                </Button>
              </Stack>
            </Box>

            <Button
              variant="text"
              onClick={() => {
                setStep('form');
                setError(null);
              }}
              sx={{ color: 'text.secondary' }}
            >
              Voltar para login
            </Button>
          </Stack>
        )}

        {step === 'reset-code' && (
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
                Redefinir senha
              </Typography>
              <Typography
                variant="body2"
                sx={{ color: 'text.secondary' }}
              >
                Digite o código enviado para <strong>{email}</strong> e sua nova senha
              </Typography>
            </Box>

            {error && (
              <Alert severity="error" onClose={() => setError(null)}>
                {error}
              </Alert>
            )}

            {successMessage && (
              <Alert severity="success" onClose={() => setSuccessMessage(null)}>
                {successMessage}
              </Alert>
            )}

            <Box component="form" onSubmit={handleResetPassword}>
              <Stack spacing={2.5}>
                <TextField
                  label="Código de verificação"
                  value={resetCode}
                  onChange={(e) => setResetCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
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

                <TextField
                  label="Nova senha"
                  type={showPassword ? 'text' : 'password'}
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
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

                <Button
                  type="submit"
                  variant="contained"
                  size="large"
                  fullWidth
                  disabled={isLoading || resetCode.length !== 6 || !newPassword}
                  sx={{ py: 1.5 }}
                >
                  {isLoading ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    'Redefinir senha'
                  )}
                </Button>
              </Stack>
            </Box>

            <Button
              variant="text"
              onClick={() => {
                setStep('form');
                setError(null);
                setSuccessMessage(null);
                setResetCode('');
                setNewPassword('');
              }}
              sx={{ color: 'text.secondary' }}
            >
              Voltar para login
            </Button>
          </Stack>
        )}
      </Paper>
    </Box>
  );
}
