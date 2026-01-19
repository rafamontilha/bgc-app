'use client';

/**
 * MUI Theme Provider Component
 * Wraps the app with MUI's ThemeProvider and CssBaseline
 */

import React from 'react';
import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import { appleTheme } from '@/lib/theme';

export interface ThemeProviderProps {
  children: React.ReactNode;
}

export function ThemeProvider({ children }: ThemeProviderProps): React.ReactElement {
  return (
    <MuiThemeProvider theme={appleTheme}>
      <CssBaseline />
      {children}
    </MuiThemeProvider>
  );
}
