'use client';

import React from 'react';
import {
  Box,
  Typography,
  Autocomplete,
  TextField,
  Chip,
} from '@mui/material';
import CategoryIcon from '@mui/icons-material/Category';
import type { NcmChapter } from '@/lib/types/onboarding';
import { NCM_CHAPTERS } from '@/lib/data/ncm-chapters';

export interface OnboardingStep1NcmProps {
  value: NcmChapter | null;
  onChange: (chapter: NcmChapter | null) => void;
}

export function OnboardingStep1Ncm({
  value,
  onChange,
}: OnboardingStep1NcmProps): React.ReactElement {
  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1 }}>
        <CategoryIcon sx={{ color: 'var(--md-sys-color-primary)', fontSize: 28 }} />
        <Typography variant="h6" sx={{ fontWeight: 600, color: 'var(--md-sys-color-on-surface)' }}>
          Qual é o seu produto principal?
        </Typography>
      </Box>

      <Typography
        variant="body2"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 3, pl: '44px' }}
      >
        Selecione o capítulo NCM que melhor representa o que você exporta. Isso
        personaliza suas recomendações de destino.
      </Typography>

      <Autocomplete<NcmChapter>
        options={NCM_CHAPTERS}
        value={value}
        onChange={(_event, newValue) => onChange(newValue)}
        getOptionLabel={(option) => `${option.code} — ${option.label}`}
        isOptionEqualToValue={(option, val) => option.code === val.code}
        renderInput={(params) => (
          <TextField
            {...params}
            label="Capítulo NCM"
            placeholder="Ex: 09 — Café, chá, mate..."
            helperText="Digite o código ou palavras-chave do produto"
            sx={{
              '& .MuiOutlinedInput-root': {
                borderRadius: '12px',
              },
            }}
          />
        )}
        renderOption={({ key, ...optionProps }, option) => (
          <Box
            key={key}
            component="li"
            {...optionProps}
            sx={{ display: 'flex', alignItems: 'baseline', gap: 1 }}
          >
            <Chip
              label={option.code}
              size="small"
              sx={{
                fontFamily: 'monospace',
                fontSize: '0.75rem',
                bgcolor: 'var(--md-sys-color-secondary-container)',
                color: 'var(--md-sys-color-on-secondary-container)',
                borderRadius: '6px',
                minWidth: 36,
              }}
            />
            <Typography variant="body2" sx={{ flexShrink: 1 }}>
              {option.label}
            </Typography>
          </Box>
        )}
        noOptionsText="Nenhum capítulo encontrado"
        fullWidth
      />

      {value && (
        <Box
          sx={{
            mt: 2,
            p: 2,
            borderRadius: '12px',
            bgcolor: 'var(--md-sys-color-secondary-container)',
            border: '1px solid var(--md-sys-color-outline-variant)',
          }}
        >
          <Typography
            variant="caption"
            sx={{ color: 'var(--md-sys-color-on-surface-variant)', display: 'block', mb: 0.5 }}
          >
            Capítulo selecionado
          </Typography>
          <Typography variant="body2" sx={{ fontWeight: 600, color: 'var(--md-sys-color-on-secondary-container)' }}>
            {value.code} — {value.label}
          </Typography>
        </Box>
      )}
    </Box>
  );
}
