/**
 * MUI v7 Theme Configuration
 * Apple Human Interface Guidelines inspired design system
 */

import { createTheme, alpha } from '@mui/material/styles';

/**
 * Apple-inspired color palette
 */
const colors = {
  primary: {
    main: '#007AFF', // Apple Science Blue
    light: '#5AC8FA',
    dark: '#0051D5',
    contrastText: '#FFFFFF',
  },
  background: {
    default: '#F5F5F7', // Light neutral gray
    paper: '#FFFFFF',
    elevated: '#FAFAFA',
  },
  text: {
    primary: '#1C1C1E', // Dark neutral
    secondary: '#3A3A3C',
    disabled: '#C7C7CC',
  },
  neutrals: {
    50: '#F2F2F7',
    100: '#E5E5EA',
    200: '#D1D1D6',
    300: '#C7C7CC',
    400: '#AEAEB2',
    500: '#8E8E93',
    600: '#636366',
    700: '#48484A',
    800: '#3A3A3C',
    900: '#2C2C2E',
    950: '#1C1C1E',
  },
  success: {
    main: '#34C759',
    light: '#30D158',
    dark: '#248A3D',
  },
  warning: {
    main: '#FF9500',
    light: '#FFCC00',
    dark: '#FF6B00',
  },
  error: {
    main: '#FF3B30',
    light: '#FF453A',
    dark: '#D70015',
  },
};

/**
 * Apple typography system
 * Using SF Pro Text font family with system fallbacks
 */
const typography = {
  fontFamily:
    '-apple-system, BlinkMacSystemFont, "SF Pro Text", "SF Pro Display", system-ui, Roboto, "Helvetica Neue", sans-serif',
  h1: {
    fontWeight: 600,
    fontSize: '2.4rem',
    lineHeight: 1.2,
    letterSpacing: '-0.02em',
  },
  h2: {
    fontWeight: 600,
    fontSize: '1.8rem',
    lineHeight: 1.3,
    letterSpacing: '-0.01em',
  },
  h3: {
    fontWeight: 600,
    fontSize: '1.5rem',
    lineHeight: 1.4,
    letterSpacing: '-0.01em',
  },
  h4: {
    fontWeight: 600,
    fontSize: '1.25rem',
    lineHeight: 1.4,
  },
  h5: {
    fontWeight: 600,
    fontSize: '1.125rem',
    lineHeight: 1.5,
  },
  h6: {
    fontWeight: 600,
    fontSize: '1rem',
    lineHeight: 1.5,
  },
  body1: {
    fontWeight: 400,
    fontSize: '1rem',
    lineHeight: 1.5,
  },
  body2: {
    fontWeight: 400,
    fontSize: '0.875rem',
    lineHeight: 1.5,
  },
  subtitle1: {
    fontWeight: 500,
    fontSize: '1rem',
    lineHeight: 1.5,
  },
  subtitle2: {
    fontWeight: 500,
    fontSize: '0.875rem',
    lineHeight: 1.5,
  },
  button: {
    fontWeight: 500,
    fontSize: '0.9375rem',
    lineHeight: 1.5,
    textTransform: 'none' as const,
    letterSpacing: '0.01em',
  },
  caption: {
    fontWeight: 400,
    fontSize: '0.75rem',
    lineHeight: 1.4,
  },
  overline: {
    fontWeight: 500,
    fontSize: '0.75rem',
    lineHeight: 1.4,
    textTransform: 'uppercase' as const,
    letterSpacing: '0.05em',
  },
};

/**
 * Standard transition timing (Apple-like smooth animations)
 */
const transitions = {
  duration: {
    shortest: 150,
    shorter: 200,
    short: 220,
    standard: 300,
    complex: 375,
    enteringScreen: 225,
    leavingScreen: 195,
  },
  easing: {
    easeInOut: 'cubic-bezier(0.32, 0.72, 0, 1)',
    easeOut: 'cubic-bezier(0.0, 0, 0.2, 1)',
    easeIn: 'cubic-bezier(0.4, 0, 1, 1)',
    sharp: 'cubic-bezier(0.4, 0, 0.6, 1)',
  },
};

/**
 * Create the Apple-inspired MUI theme
 */
export const appleTheme = createTheme({
  palette: {
    mode: 'light',
    primary: colors.primary,
    background: colors.background,
    text: colors.text,
    success: colors.success,
    warning: colors.warning,
    error: colors.error,
  },
  typography,
  shape: {
    borderRadius: 16,
  },
  spacing: 8,
  transitions: {
    duration: transitions.duration,
    easing: transitions.easing,
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          borderRadius: 999, // Pill shape
          textTransform: 'none',
          fontWeight: 500,
          padding: '10px 24px',
          transition: `all ${transitions.duration.short}ms ${transitions.easing.easeInOut}`,
          boxShadow: 'none',
          '&:hover': {
            transform: 'scale(1.02)',
            boxShadow: '0px 4px 20px rgba(0, 0, 0, 0.12)',
          },
          '&:active': {
            transform: 'scale(0.98)',
          },
        },
        sizeLarge: {
          padding: '14px 32px',
          fontSize: '1rem',
        },
        sizeMedium: {
          padding: '10px 24px',
          fontSize: '0.9375rem',
        },
        sizeSmall: {
          padding: '6px 16px',
          fontSize: '0.875rem',
        },
      },
    },
    MuiCard: {
      styleOverrides: {
        root: {
          borderRadius: 16,
          boxShadow: '0px 4px 24px rgba(0, 0, 0, 0.10)',
          transition: `all ${transitions.duration.short}ms ${transitions.easing.easeInOut}`,
          '&:hover': {
            boxShadow: '0px 8px 32px rgba(0, 0, 0, 0.14)',
          },
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 12,
            transition: `all ${transitions.duration.short}ms ${transitions.easing.easeInOut}`,
            '& fieldset': {
              borderColor: colors.neutrals[200],
            },
            '&:hover fieldset': {
              borderColor: colors.neutrals[400],
            },
            '&.Mui-focused fieldset': {
              borderWidth: 2,
            },
          },
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          fontWeight: 500,
          fontSize: '0.875rem',
        },
        filled: {
          '&.MuiChip-colorSuccess': {
            backgroundColor: alpha(colors.success.main, 0.12),
            color: colors.success.dark,
          },
          '&.MuiChip-colorWarning': {
            backgroundColor: alpha(colors.warning.main, 0.12),
            color: colors.warning.dark,
          },
          '&.MuiChip-colorError': {
            backgroundColor: alpha(colors.error.main, 0.12),
            color: colors.error.dark,
          },
        },
      },
    },
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 16,
        },
        elevation1: {
          boxShadow: '0px 2px 8px rgba(0, 0, 0, 0.08)',
        },
        elevation2: {
          boxShadow: '0px 4px 16px rgba(0, 0, 0, 0.10)',
        },
        elevation3: {
          boxShadow: '0px 8px 24px rgba(0, 0, 0, 0.12)',
        },
      },
    },
    MuiLinearProgress: {
      styleOverrides: {
        root: {
          borderRadius: 4,
          height: 8,
        },
      },
    },
    MuiAlert: {
      styleOverrides: {
        root: {
          borderRadius: 12,
          padding: '12px 16px',
        },
        standardSuccess: {
          backgroundColor: alpha(colors.success.main, 0.12),
          color: colors.success.dark,
        },
        standardWarning: {
          backgroundColor: alpha(colors.warning.main, 0.12),
          color: colors.warning.dark,
        },
        standardError: {
          backgroundColor: alpha(colors.error.main, 0.12),
          color: colors.error.dark,
        },
        standardInfo: {
          backgroundColor: alpha(colors.primary.main, 0.12),
          color: colors.primary.dark,
        },
      },
    },
    MuiDialog: {
      styleOverrides: {
        paper: {
          borderRadius: 20,
          padding: 8,
        },
      },
    },
    MuiSkeleton: {
      styleOverrides: {
        root: {
          borderRadius: 8,
          backgroundColor: colors.neutrals[100],
        },
      },
    },
  },
});

/**
 * Custom theme extensions
 */
declare module '@mui/material/styles' {
  interface Theme {
    customColors: typeof colors;
  }
  interface ThemeOptions {
    customColors?: typeof colors;
  }
}

// Extend theme with custom colors
appleTheme.customColors = colors;

export default appleTheme;
