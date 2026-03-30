'use client';

export const dynamic = 'force-dynamic';

/**
 * Custom Profile Page
 * User profile management using Clerk hooks with MUI components
 */

import React, { useState } from 'react';
import { useUser, useClerk } from '@clerk/nextjs';
import {
  Box,
  Container,
  Typography,
  Paper,
  Avatar,
  Stack,
  TextField,
  Button,
  Divider,
  Alert,
  CircularProgress,
  Chip,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import GoogleIcon from '@mui/icons-material/Google';
import EmailIcon from '@mui/icons-material/Email';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import CategoryIcon from '@mui/icons-material/Category';
import { Header } from '@/components/home/Header';
import { useRouter } from 'next/navigation';
import type { OnboardingMetadata } from '@/lib/types/onboarding';

export default function ProfilePage(): React.ReactElement {
  const { user, isLoaded } = useUser();
  const { signOut } = useClerk();
  const router = useRouter();

  const [isEditingName, setIsEditingName] = useState(false);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  const handleEditName = (): void => {
    setFirstName(user?.firstName || '');
    setLastName(user?.lastName || '');
    setIsEditingName(true);
    setError(null);
    setSuccess(null);
  };

  const handleCancelEdit = (): void => {
    setIsEditingName(false);
    setError(null);
  };

  const handleSaveName = async (): Promise<void> => {
    if (!user) return;

    setIsSaving(true);
    setError(null);

    try {
      await user.update({
        firstName,
        lastName,
      });
      setIsEditingName(false);
      setSuccess('Nome atualizado com sucesso!');
      setTimeout(() => setSuccess(null), 3000);
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string }> };
      setError(clerkError.errors?.[0]?.message || 'Erro ao atualizar nome.');
    } finally {
      setIsSaving(false);
    }
  };

  const handleDeleteAccount = async (): Promise<void> => {
    if (!user) return;

    setIsDeleting(true);

    try {
      await user.delete();
      await signOut();
    } catch (err) {
      const clerkError = err as { errors?: Array<{ message: string }> };
      setError(clerkError.errors?.[0]?.message || 'Erro ao excluir conta.');
      setIsDeleting(false);
      setDeleteDialogOpen(false);
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

  if (!user) {
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
        <Typography>Usuário não encontrado</Typography>
      </Box>
    );
  }

  const userInitials = user.firstName && user.lastName
    ? `${user.firstName[0]}${user.lastName[0]}`
    : user.firstName?.[0] || user.emailAddresses[0]?.emailAddress[0]?.toUpperCase() || '?';

  const hasGoogleAccount = user.externalAccounts?.some(
    (account) => account.provider === 'google'
  );

  const hasPasswordAuth = user.passwordEnabled;

  return (
    <Box sx={{ minHeight: '100vh', bgcolor: 'background.default' }}>
      <Header />

      <Container maxWidth="md" sx={{ py: 4 }}>
        <Typography
          variant="h4"
          sx={{
            fontWeight: 600,
            color: 'text.primary',
            mb: 4,
          }}
        >
          Meu Perfil
        </Typography>

        <Stack spacing={3}>
          {/* Success/Error Messages */}
          {success && (
            <Alert severity="success" onClose={() => setSuccess(null)}>
              {success}
            </Alert>
          )}
          {error && (
            <Alert severity="error" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {/* Profile Card */}
          <Paper elevation={2} sx={{ p: 4, borderRadius: 4 }}>
            <Stack spacing={4}>
              {/* Avatar and Basic Info */}
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 3 }}>
                <Avatar
                  src={user.imageUrl}
                  alt={user.fullName || 'User avatar'}
                  sx={{
                    width: 80,
                    height: 80,
                    bgcolor: 'primary.main',
                    fontSize: '2rem',
                    fontWeight: 600,
                  }}
                >
                  {userInitials}
                </Avatar>
                <Box sx={{ flex: 1 }}>
                  {isEditingName ? (
                    <Stack direction="row" spacing={2} alignItems="center">
                      <TextField
                        label="Nome"
                        value={firstName}
                        onChange={(e) => setFirstName(e.target.value)}
                        size="small"
                        disabled={isSaving}
                      />
                      <TextField
                        label="Sobrenome"
                        value={lastName}
                        onChange={(e) => setLastName(e.target.value)}
                        size="small"
                        disabled={isSaving}
                      />
                      <IconButton
                        onClick={handleSaveName}
                        disabled={isSaving}
                        color="primary"
                        size="small"
                      >
                        {isSaving ? <CircularProgress size={20} /> : <CheckIcon />}
                      </IconButton>
                      <IconButton
                        onClick={handleCancelEdit}
                        disabled={isSaving}
                        size="small"
                      >
                        <CloseIcon />
                      </IconButton>
                    </Stack>
                  ) : (
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="h5" sx={{ fontWeight: 600 }}>
                        {user.fullName || 'Sem nome'}
                      </Typography>
                      <IconButton onClick={handleEditName} size="small">
                        <EditIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  )}
                  <Typography variant="body2" sx={{ color: 'text.secondary', mt: 0.5 }}>
                    Membro desde {new Date(user.createdAt!).toLocaleDateString('pt-BR', {
                      month: 'long',
                      year: 'numeric',
                    })}
                  </Typography>
                </Box>
              </Box>

              <Divider />

              {/* Email */}
              <Box>
                <Typography variant="subtitle2" sx={{ color: 'text.secondary', mb: 1 }}>
                  Email
                </Typography>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <EmailIcon sx={{ color: 'text.secondary' }} />
                  <Typography variant="body1">
                    {user.primaryEmailAddress?.emailAddress}
                  </Typography>
                  {user.primaryEmailAddress?.verification?.status === 'verified' && (
                    <Chip
                      label="Verificado"
                      size="small"
                      color="success"
                      variant="outlined"
                    />
                  )}
                </Box>
              </Box>

              <Divider />

              {/* Connected Accounts */}
              <Box>
                <Typography variant="subtitle2" sx={{ color: 'text.secondary', mb: 2 }}>
                  Métodos de Autenticação
                </Typography>
                <Stack spacing={1.5}>
                  {hasGoogleAccount && (
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        p: 2,
                        bgcolor: 'background.default',
                        borderRadius: 2,
                      }}
                    >
                      <GoogleIcon sx={{ color: '#4285F4' }} />
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          Google
                        </Typography>
                        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                          Conectado via Google OAuth
                        </Typography>
                      </Box>
                      <Chip label="Ativo" size="small" color="success" variant="outlined" />
                    </Box>
                  )}
                  {hasPasswordAuth && (
                    <Box
                      sx={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: 2,
                        p: 2,
                        bgcolor: 'background.default',
                        borderRadius: 2,
                      }}
                    >
                      <EmailIcon sx={{ color: 'text.secondary' }} />
                      <Box sx={{ flex: 1 }}>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          Email e Senha
                        </Typography>
                        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                          Autenticação por email
                        </Typography>
                      </Box>
                      <Chip label="Ativo" size="small" color="success" variant="outlined" />
                    </Box>
                  )}
                </Stack>
              </Box>
            </Stack>
          </Paper>

          {/* Export Settings */}
          {(() => {
            const meta = user.publicMetadata as unknown as OnboardingMetadata | undefined;
            const hasOnboarding = meta?.onboardingCompleted && !meta?.onboardingSkipped;
            return (
              <Paper elevation={2} sx={{ p: 4, borderRadius: 4 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 2 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                    <CategoryIcon sx={{ color: 'primary.main' }} />
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Configurações de Exportação
                    </Typography>
                  </Box>
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<EditIcon />}
                    onClick={() => router.push('/onboarding?re=1')}
                    sx={{ borderRadius: 2, textTransform: 'none' }}
                  >
                    Alterar
                  </Button>
                </Box>

                {hasOnboarding ? (
                  <Stack spacing={1}>
                    <Box>
                      <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                        Produto principal (NCM)
                      </Typography>
                      <Typography variant="body2" sx={{ fontWeight: 500 }}>
                        {meta!.ncmChapter} — {meta!.ncmLabel}
                      </Typography>
                    </Box>
                    {meta!.volumeMonthlyKg && (
                      <Box>
                        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                          Volume médio mensal
                        </Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {meta!.volumeMonthlyKg.toLocaleString('pt-BR')} kg/mês
                        </Typography>
                      </Box>
                    )}
                    {(meta!.targetContinents.length > 0 || meta!.targetCountries.length > 0) && (
                      <Box>
                        <Typography variant="caption" sx={{ color: 'text.secondary' }}>
                          Destinos alvo
                        </Typography>
                        <Typography variant="body2" sx={{ fontWeight: 500 }}>
                          {[...meta!.targetContinents, ...meta!.targetCountries].join(', ')}
                        </Typography>
                      </Box>
                    )}
                  </Stack>
                ) : (
                  <Typography variant="body2" sx={{ color: 'text.secondary' }}>
                    Perfil de exportação não configurado.{' '}
                    <Button
                      variant="text"
                      size="small"
                      onClick={() => router.push('/onboarding')}
                      sx={{ textTransform: 'none', p: 0, minWidth: 0, verticalAlign: 'baseline' }}
                    >
                      Configurar agora
                    </Button>
                  </Typography>
                )}
              </Paper>
            );
          })()}

          {/* Danger Zone */}
          <Paper
            elevation={0}
            sx={{
              p: 4,
              borderRadius: 4,
              border: '1px solid',
              borderColor: 'error.light',
              bgcolor: 'transparent',
            }}
          >
            <Typography
              variant="h6"
              sx={{ fontWeight: 600, color: 'error.main', mb: 2 }}
            >
              Zona de Perigo
            </Typography>
            <Typography variant="body2" sx={{ color: 'text.secondary', mb: 3 }}>
              Ao excluir sua conta, todos os seus dados serão permanentemente removidos.
              Esta ação não pode ser desfeita.
            </Typography>
            <Button
              variant="outlined"
              color="error"
              startIcon={<DeleteOutlineIcon />}
              onClick={() => setDeleteDialogOpen(true)}
            >
              Excluir minha conta
            </Button>
          </Paper>
        </Stack>
      </Container>

      {/* Delete Confirmation Dialog */}
      <Dialog
        open={deleteDialogOpen}
        onClose={() => !isDeleting && setDeleteDialogOpen(false)}
      >
        <DialogTitle sx={{ fontWeight: 600 }}>
          Excluir conta permanentemente?
        </DialogTitle>
        <DialogContent>
          <DialogContentText>
            Todos os seus dados serão excluídos e esta ação não pode ser desfeita.
            Tem certeza que deseja continuar?
          </DialogContentText>
        </DialogContent>
        <DialogActions sx={{ p: 2, pt: 0 }}>
          <Button
            onClick={() => setDeleteDialogOpen(false)}
            disabled={isDeleting}
          >
            Cancelar
          </Button>
          <Button
            onClick={handleDeleteAccount}
            color="error"
            variant="contained"
            disabled={isDeleting}
            startIcon={isDeleting ? <CircularProgress size={16} color="inherit" /> : null}
          >
            {isDeleting ? 'Excluindo...' : 'Sim, excluir'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
