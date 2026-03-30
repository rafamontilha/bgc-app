'use client';

/**
 * Custom User Menu Component
 * Replaces Clerk's UserButton with MUI components following Apple HIG design system
 */

import React, { useState } from 'react';
import { useUser, useClerk } from '@clerk/nextjs';
import { useRouter } from 'next/navigation';
import {
  Avatar,
  Box,
  Menu,
  MenuItem,
  ListItemIcon,
  ListItemText,
  Divider,
  Typography,
  IconButton,
  Skeleton,
} from '@mui/material';
import PersonOutlineIcon from '@mui/icons-material/PersonOutline';
import SettingsOutlinedIcon from '@mui/icons-material/SettingsOutlined';
import LogoutIcon from '@mui/icons-material/Logout';

export function UserMenu(): React.ReactElement | null {
  const { user, isLoaded } = useUser();
  const { signOut } = useClerk();
  const router = useRouter();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);

  const handleClick = (event: React.MouseEvent<HTMLElement>): void => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = (): void => {
    setAnchorEl(null);
  };

  const handleProfile = (): void => {
    handleClose();
    router.push('/profile');
  };

  const handleSettings = (): void => {
    handleClose();
    router.push('/profile#security');
  };

  const handleSignOut = async (): Promise<void> => {
    handleClose();
    await signOut();
    router.push('/');
  };

  if (!isLoaded) {
    return <Skeleton variant="circular" width={40} height={40} />;
  }

  if (!user) {
    return null;
  }

  const userInitials = user.firstName && user.lastName
    ? `${user.firstName[0]}${user.lastName[0]}`
    : user.firstName?.[0] || user.emailAddresses[0]?.emailAddress[0]?.toUpperCase() || '?';

  return (
    <>
      <IconButton
        onClick={handleClick}
        size="small"
        aria-controls={open ? 'user-menu' : undefined}
        aria-haspopup="true"
        aria-expanded={open ? 'true' : undefined}
        sx={{
          p: 0.5,
          transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
          '&:hover': {
            transform: 'scale(1.05)',
          },
        }}
      >
        <Avatar
          src={user.imageUrl}
          alt={user.fullName || 'User avatar'}
          sx={{
            width: 40,
            height: 40,
            bgcolor: 'primary.main',
            fontSize: '1rem',
            fontWeight: 600,
          }}
        >
          {userInitials}
        </Avatar>
      </IconButton>

      <Menu
        id="user-menu"
        anchorEl={anchorEl}
        open={open}
        onClose={handleClose}
        onClick={handleClose}
        transformOrigin={{ horizontal: 'right', vertical: 'top' }}
        anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
        slotProps={{
          paper: {
            elevation: 3,
            sx: {
              mt: 1.5,
              minWidth: 240,
              borderRadius: 3,
              overflow: 'visible',
              '&::before': {
                content: '""',
                display: 'block',
                position: 'absolute',
                top: 0,
                right: 14,
                width: 10,
                height: 10,
                bgcolor: 'background.paper',
                transform: 'translateY(-50%) rotate(45deg)',
                zIndex: 0,
                boxShadow: '-2px -2px 4px rgba(0,0,0,0.04)',
              },
            },
          },
        }}
      >
        {/* User Info Header */}
        <Box sx={{ px: 2, py: 1.5 }}>
          <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>
            {user.fullName || 'Usuário'}
          </Typography>
          <Typography variant="body2" sx={{ color: 'text.secondary' }}>
            {user.primaryEmailAddress?.emailAddress}
          </Typography>
        </Box>

        <Divider sx={{ my: 1 }} />

        {/* Menu Items */}
        <MenuItem
          onClick={handleProfile}
          sx={{
            py: 1.5,
            px: 2,
            borderRadius: 2,
            mx: 1,
            '&:hover': {
              bgcolor: 'action.hover',
            },
          }}
        >
          <ListItemIcon>
            <PersonOutlineIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText
            primary="Meu Perfil"
            primaryTypographyProps={{ fontWeight: 500 }}
          />
        </MenuItem>

        <MenuItem
          onClick={handleSettings}
          sx={{
            py: 1.5,
            px: 2,
            borderRadius: 2,
            mx: 1,
            '&:hover': {
              bgcolor: 'action.hover',
            },
          }}
        >
          <ListItemIcon>
            <SettingsOutlinedIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText
            primary="Configurações"
            primaryTypographyProps={{ fontWeight: 500 }}
          />
        </MenuItem>

        <Divider sx={{ my: 1 }} />

        <MenuItem
          onClick={handleSignOut}
          sx={{
            py: 1.5,
            px: 2,
            borderRadius: 2,
            mx: 1,
            mb: 1,
            color: 'error.main',
            '&:hover': {
              bgcolor: 'error.main',
              color: 'error.contrastText',
              '& .MuiListItemIcon-root': {
                color: 'error.contrastText',
              },
            },
          }}
        >
          <ListItemIcon sx={{ color: 'error.main' }}>
            <LogoutIcon fontSize="small" />
          </ListItemIcon>
          <ListItemText
            primary="Sair"
            primaryTypographyProps={{ fontWeight: 500 }}
          />
        </MenuItem>
      </Menu>
    </>
  );
}
