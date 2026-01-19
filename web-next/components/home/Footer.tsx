'use client';

/**
 * Footer Component
 * Site footer with links, contact, and social media
 */

import React from 'react';
import {
  Box,
  Container,
  Typography,
  Link,
  Stack,
  Divider,
  IconButton,
} from '@mui/material';
import GitHubIcon from '@mui/icons-material/GitHub';
import LinkedInIcon from '@mui/icons-material/LinkedIn';
import EmailIcon from '@mui/icons-material/Email';

export function Footer(): React.ReactElement {
  const currentYear = new Date().getFullYear();

  return (
    <Box
      component="footer"
      sx={{
        backgroundColor: '#1C1C1E',
        color: 'white',
        py: 8,
        mt: 'auto',
      }}
    >
      <Container maxWidth="lg">
        <Stack spacing={6}>
          {/* Main content */}
          <Box
            sx={{
              display: 'grid',
              gridTemplateColumns: {
                xs: '1fr',
                sm: 'repeat(2, 1fr)',
                md: 'repeat(4, 1fr)',
              },
              gap: 4,
            }}
          >
            {/* About */}
            <Box>
              <Typography
                variant="h6"
                sx={{
                  fontWeight: 600,
                  mb: 2,
                  color: 'white',
                }}
              >
                Brasil Global Connect
              </Typography>
              <Typography
                variant="body2"
                sx={{
                  color: 'rgba(255, 255, 255, 0.7)',
                  lineHeight: 1.7,
                }}
              >
                Plataforma inteligente para identificar os melhores destinos de
                exportação usando dados oficiais do ComexStat.
              </Typography>
            </Box>

            {/* Product */}
            <Box>
              <Typography
                variant="subtitle1"
                sx={{
                  fontWeight: 600,
                  mb: 2,
                  color: 'white',
                }}
              >
                Produto
              </Typography>
              <Stack spacing={1.5}>
                <Link
                  href="#simulator"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Simulador
                </Link>
                <Link
                  href="#pricing"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Planos e Preços
                </Link>
                <Link
                  href="/docs"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Documentação
                </Link>
              </Stack>
            </Box>

            {/* Company */}
            <Box>
              <Typography
                variant="subtitle1"
                sx={{
                  fontWeight: 600,
                  mb: 2,
                  color: 'white',
                }}
              >
                Empresa
              </Typography>
              <Stack spacing={1.5}>
                <Link
                  href="/about"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Sobre Nós
                </Link>
                <Link
                  href="/contact"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Contato
                </Link>
                <Link
                  href="/privacy"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Privacidade
                </Link>
                <Link
                  href="/terms"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    textDecoration: 'none',
                    fontSize: '0.875rem',
                    '&:hover': {
                      color: '#007AFF',
                    },
                  }}
                >
                  Termos de Uso
                </Link>
              </Stack>
            </Box>

            {/* Social */}
            <Box>
              <Typography
                variant="subtitle1"
                sx={{
                  fontWeight: 600,
                  mb: 2,
                  color: 'white',
                }}
              >
                Redes Sociais
              </Typography>
              <Stack direction="row" spacing={1}>
                <IconButton
                  href="https://github.com/brasil-global-connect"
                  target="_blank"
                  rel="noopener noreferrer"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    '&:hover': {
                      color: '#007AFF',
                      backgroundColor: 'rgba(0, 122, 255, 0.1)',
                    },
                  }}
                >
                  <GitHubIcon />
                </IconButton>
                <IconButton
                  href="https://linkedin.com/company/brasil-global-connect"
                  target="_blank"
                  rel="noopener noreferrer"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    '&:hover': {
                      color: '#007AFF',
                      backgroundColor: 'rgba(0, 122, 255, 0.1)',
                    },
                  }}
                >
                  <LinkedInIcon />
                </IconButton>
                <IconButton
                  href="mailto:contato@brasilglobalconnect.com.br"
                  sx={{
                    color: 'rgba(255, 255, 255, 0.7)',
                    '&:hover': {
                      color: '#007AFF',
                      backgroundColor: 'rgba(0, 122, 255, 0.1)',
                    },
                  }}
                >
                  <EmailIcon />
                </IconButton>
              </Stack>
            </Box>
          </Box>

          <Divider sx={{ borderColor: 'rgba(255, 255, 255, 0.12)' }} />

          {/* Copyright */}
          <Box
            sx={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              flexWrap: 'wrap',
              gap: 2,
            }}
          >
            <Typography
              variant="body2"
              sx={{
                color: 'rgba(255, 255, 255, 0.5)',
              }}
            >
              © {currentYear} Brasil Global Connect. Todos os direitos
              reservados.
            </Typography>
            <Typography
              variant="body2"
              sx={{
                color: 'rgba(255, 255, 255, 0.5)',
              }}
            >
              Dados: ComexStat (Governo Federal)
            </Typography>
          </Box>
        </Stack>
      </Container>
    </Box>
  );
}
