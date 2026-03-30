'use client';

import React, { useState } from 'react';
import {
  Box,
  Typography,
  TextField,
  InputAdornment,
  ToggleButtonGroup,
  ToggleButton,
  FormHelperText,
} from '@mui/material';
import ScaleIcon from '@mui/icons-material/Scale';

export interface OnboardingStep2VolumeProps {
  value: number | undefined;
  onChange: (volume: number | undefined) => void;
}

const QUICK_OPTIONS: { label: string; value: number }[] = [
  { label: '< 500 kg', value: 250 },
  { label: '500 kg–2 t', value: 1000 },
  { label: '2–10 t', value: 5000 },
  { label: '10–50 t', value: 25000 },
  { label: '> 50 t', value: 75000 },
];

export function OnboardingStep2Volume({
  value,
  onChange,
}: OnboardingStep2VolumeProps): React.ReactElement {
  const [rawInput, setRawInput] = useState<string>(
    value !== undefined ? String(value) : ''
  );
  const [error, setError] = useState<string>('');

  const quickValue =
    QUICK_OPTIONS.find((o) => o.value === value)?.value ?? null;

  function handleQuickSelect(
    _event: React.MouseEvent<HTMLElement>,
    selected: number | null
  ) {
    if (selected === null) return;
    setRawInput(String(selected));
    setError('');
    onChange(selected);
  }

  function handleTextChange(e: React.ChangeEvent<HTMLInputElement>) {
    const raw = e.target.value;
    setRawInput(raw);

    if (raw === '') {
      setError('');
      onChange(undefined);
      return;
    }

    const parsed = Number(raw);
    if (!Number.isFinite(parsed) || parsed <= 0) {
      setError('Informe um número positivo');
      onChange(undefined);
    } else if (parsed > 10_000_000) {
      setError('Volume máximo: 10.000.000 kg');
      onChange(undefined);
    } else {
      setError('');
      onChange(parsed);
    }
  }

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1 }}>
        <ScaleIcon sx={{ color: 'var(--md-sys-color-primary)', fontSize: 28 }} />
        <Typography variant="h6" sx={{ fontWeight: 600, color: 'var(--md-sys-color-on-surface)' }}>
          Qual é o seu volume médio mensal?
        </Typography>
      </Box>

      <Typography
        variant="body2"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 3, pl: '44px' }}
      >
        Uma estimativa ajuda a calibrar os mercados mais adequados para o seu
        porte. Você pode pular esta etapa se preferir.
      </Typography>

      {/* Quick-select chips */}
      <Typography
        variant="caption"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 1, display: 'block' }}
      >
        Selecione uma faixa
      </Typography>
      <ToggleButtonGroup
        value={quickValue}
        exclusive
        onChange={handleQuickSelect}
        size="small"
        sx={{
          flexWrap: 'wrap',
          gap: 1,
          mb: 3,
          '& .MuiToggleButtonGroup-grouped': {
            borderRadius: '20px !important',
            border: '1px solid var(--md-sys-color-outline) !important',
            px: 2,
            py: 0.75,
            textTransform: 'none',
            '&.Mui-selected': {
              bgcolor: 'var(--md-sys-color-primary)',
              color: 'var(--md-sys-color-on-primary)',
              '&:hover': {
                bgcolor: 'var(--md-sys-color-primary)',
                opacity: 0.9,
              },
            },
          },
        }}
      >
        {QUICK_OPTIONS.map((opt) => (
          <ToggleButton key={opt.value} value={opt.value}>
            {opt.label}
          </ToggleButton>
        ))}
      </ToggleButtonGroup>

      {/* Free-text input */}
      <Typography
        variant="caption"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 1, display: 'block' }}
      >
        Ou insira um valor exato
      </Typography>
      <TextField
        value={rawInput}
        onChange={handleTextChange}
        label="Volume mensal"
        type="number"
        inputProps={{ min: 1, step: 100 }}
        InputProps={{
          endAdornment: <InputAdornment position="end">kg</InputAdornment>,
        }}
        error={!!error}
        helperText={error || 'Campo opcional'}
        fullWidth
        sx={{
          '& .MuiOutlinedInput-root': {
            borderRadius: '12px',
          },
        }}
      />

      {value !== undefined && !error && (
        <FormHelperText sx={{ mt: 1, color: 'var(--md-sys-color-primary)', fontWeight: 500 }}>
          {value.toLocaleString('pt-BR')} kg/mês selecionado
        </FormHelperText>
      )}
    </Box>
  );
}
