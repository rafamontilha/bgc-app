'use client';

/**
 * Header Component
 * Navigation bar with logo and primary CTAs
 */

import React from 'react';
import {
  AppBar,
  Toolbar,
  Box,
  Button,
  Container,
  Typography,
} from '@mui/material';
import LoginIcon from '@mui/icons-material/Login';

export interface HeaderProps {
  onLoginClick?: () => void;
}

export function Header({
  onLoginClick,
}: HeaderProps): React.ReactElement {
  const handleLoginClick = (): void => {
    if (onLoginClick) {
      onLoginClick();
    } else {
      // Default behavior: navigate to login page
      window.location.href = '/auth/login';
    }
  };

  return (
    <AppBar
      position="sticky"
      elevation={0}
      sx={{
        backgroundColor: 'background.paper',
        borderBottom: 1,
        borderColor: 'divider',
      }}
    >
      <Container maxWidth="lg">
        <Toolbar
          sx={{
            justifyContent: 'space-between',
            py: 1.5,
            px: { xs: 0, sm: 2 },
          }}
        >
          {/* Logo */}
          <Box
            sx={{
              display: 'flex',
              alignItems: 'center',
              gap: 1.5,
            }}
          >
            <Box
              sx={{
                width: 40,
                height: 40,
                borderRadius: 2,
                backgroundColor: 'primary.main',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <Typography
                variant="h6"
                sx={{
                  color: 'primary.contrastText',
                  fontWeight: 700,
                  fontSize: '1.25rem',
                }}
              >
                BGC
              </Typography>
            </Box>
            <Typography
              variant="h6"
              sx={{
                fontWeight: 600,
                color: 'text.primary',
                display: { xs: 'none', sm: 'block' },
              }}
            >
              Brasil Global Connect
            </Typography>
          </Box>

          {/* Navigation */}
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            {/* Desktop Navigation Links */}
            <Box sx={{ display: { xs: 'none', md: 'flex' }, alignItems: 'center', gap: 1, mr: 2 }}>
              <Button
                href="/"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  px: 2.5,
                  transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                  '&:hover': {
                    backgroundColor: 'action.hover',
                  },
                }}
              >
                Início
              </Button>

              <Button
                href="/mapa"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  px: 2.5,
                  transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                  '&:hover': {
                    backgroundColor: 'action.hover',
                  },
                }}
              >
                Mapa
              </Button>

              <Button
                href="/conteudos"
                sx={{
                  color: 'text.primary',
                  fontWeight: 500,
                  px: 2.5,
                  transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                  '&:hover': {
                    backgroundColor: 'action.hover',
                  },
                }}
              >
                Conteúdos
              </Button>
            </Box>

            {/* Auth Buttons */}
            <Button
              onClick={handleLoginClick}
              variant="outlined"
              sx={{
                fontWeight: 500,
                px: 3,
                borderColor: 'divider',
                color: 'text.primary',
                transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                '&:hover': {
                  borderColor: 'primary.main',
                  backgroundColor: 'action.hover',
                },
              }}
              startIcon={<LoginIcon />}
            >
              Login
            </Button>

            <Button
              href="/register"
              variant="contained"
              sx={{
                fontWeight: 500,
                px: 3,
                transition: 'all 220ms cubic-bezier(0.32, 0.72, 0, 1)',
                '&:hover': {
                  backgroundColor: 'primary.dark',
                },
              }}
            >
              Cadastrar
            </Button>
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
}
