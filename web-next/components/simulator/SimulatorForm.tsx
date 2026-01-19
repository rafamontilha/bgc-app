'use client';

/**
 * US-001: Simulator Form Component
 * Input form for NCM and simulation parameters with client-side validation
 */

import React, { useState } from 'react';
import {
  Box,
  TextField,
  Button,
  Typography,
  CircularProgress,
  Stack,
  InputAdornment,
  Autocomplete,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import { z } from 'zod';
import { SimulatorRequest } from '@/types/simulator';
import { AVAILABLE_NCMS, formatNCM, type NCMItem } from '@/lib/data/ncm-list';

/**
 * Zod validation schema for simulator form
 */
const simulatorFormSchema = z.object({
  ncm: z
    .string()
    .length(8, 'NCM deve ter exatamente 8 dígitos')
    .regex(/^\d{8}$/, 'NCM deve conter apenas números'),
  volume_kg: z
    .number({ invalid_type_error: 'Volume deve ser um número' })
    .positive('Volume deve ser maior que zero')
    .max(1000000, 'Volume máximo é 1.000.000 kg'),
  countries: z.array(z.string()).optional(),
  max_results: z
    .number()
    .min(1, 'Mínimo 1 resultado')
    .max(50, 'Máximo 50 resultados')
    .optional(),
});

export interface SimulatorFormProps {
  onSubmit: (request: SimulatorRequest) => void;
  isLoading: boolean;
}

interface FormData {
  ncm: string;
  volume_kg: string;
}

interface FormErrors {
  ncm?: string;
  volume_kg?: string;
}

export function SimulatorForm({
  onSubmit,
  isLoading,
}: SimulatorFormProps): React.ReactElement {
  const [formData, setFormData] = useState<FormData>({
    ncm: '',
    volume_kg: '1000',
  });
  const [selectedNCM, setSelectedNCM] = useState<NCMItem | null>(null);
  const [inputDisplayValue, setInputDisplayValue] = useState<string>('');
  const [errors, setErrors] = useState<FormErrors>({});

  const handleInputChange = (field: keyof FormData, value: string): void => {
    setFormData((prev) => ({ ...prev, [field]: value }));
    // Clear error for this field when user starts typing
    if (errors[field]) {
      setErrors((prev) => ({ ...prev, [field]: undefined }));
    }
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>): void => {
    event.preventDefault();
    setErrors({});

    // Prepare data for validation
    const dataToValidate = {
      ncm: formData.ncm.trim(),
      volume_kg: parseFloat(formData.volume_kg),
    };

    // Validate using Zod
    const result = simulatorFormSchema.safeParse(dataToValidate);

    if (!result.success) {
      // Extract errors from Zod validation
      const newErrors: FormErrors = {};
      result.error.errors.forEach((error) => {
        const field = error.path[0] as keyof FormErrors;
        newErrors[field] = error.message;
      });
      setErrors(newErrors);
      return;
    }

    // Submit validated data (snake_case for backend API)
    onSubmit({
      ncm: result.data.ncm,
      volume_kg: result.data.volume_kg,
      max_results: 10, // Default value
    });
  };

  return (
    <Box
      component="form"
      onSubmit={handleSubmit}
      sx={{
        width: '100%',
        maxWidth: 600,
        mx: 'auto',
      }}
    >
      <Stack spacing={3}>
        <Box>
          <Typography
            variant="h3"
            component="h1"
            sx={{
              mb: 1,
              fontWeight: 600,
              color: 'text.primary',
            }}
          >
            Simulador de Destinos de Exportação
          </Typography>
          <Typography
            variant="body1"
            sx={{
              color: 'text.secondary',
            }}
          >
            Descubra os melhores mercados para seu produto
          </Typography>
        </Box>

        <Autocomplete<NCMItem, false, false, true>
          options={AVAILABLE_NCMS}
          value={selectedNCM}
          onChange={(_, newValue) => {
            // Type guard: only NCMItem | null is valid
            if (typeof newValue === 'string') {
              // If freeSolo returns a string, treat as null
              setSelectedNCM(null);
              setFormData((prev) => ({ ...prev, ncm: newValue.replace(/\D/g, '') }));
              setInputDisplayValue(newValue);
            } else {
              setSelectedNCM(newValue);
              if (newValue) {
                // When selecting from dropdown, set only the NCM code for validation
                setFormData((prev) => ({ ...prev, ncm: newValue.code }));
                setInputDisplayValue(formatNCM(newValue));
                if (errors.ncm) {
                  setErrors((prev) => ({ ...prev, ncm: undefined }));
                }
              } else {
                setFormData((prev) => ({ ...prev, ncm: '' }));
                setInputDisplayValue('');
              }
            }
          }}
          inputValue={inputDisplayValue}
          onInputChange={(_, newInputValue) => {
            setInputDisplayValue(newInputValue);
            // Extract only digits for validation (in case user types freely)
            const digitsOnly = newInputValue.replace(/\D/g, '');
            setFormData((prev) => ({ ...prev, ncm: digitsOnly }));
            if (errors.ncm) {
              setErrors((prev) => ({ ...prev, ncm: undefined }));
            }
          }}
          getOptionLabel={(option) => {
            // Handle freeSolo string input
            if (typeof option === 'string') {
              return option;
            }
            return formatNCM(option);
          }}
          isOptionEqualToValue={(option, value) => option.code === value.code}
          disabled={isLoading}
          freeSolo
          renderInput={(params) => (
            <TextField
              {...params}
              label="Produto (NCM)"
              placeholder="Ex: 17011400 - Açúcar de cana, em bruto"
              error={Boolean(errors.ncm)}
              helperText={
                errors.ncm ||
                'Digite o código NCM ou busque pelo nome do produto'
              }
              required
              slotProps={{
                input: {
                  ...params.InputProps,
                  startAdornment: (
                    <>
                      <InputAdornment position="start">
                        <SearchIcon sx={{ color: 'text.secondary' }} />
                      </InputAdornment>
                      {params.InputProps.startAdornment}
                    </>
                  ),
                },
              }}
              sx={{
                '& .MuiOutlinedInput-root': {
                  backgroundColor: 'background.paper',
                },
              }}
            />
          )}
          renderOption={(props, option) => (
            <Box component="li" {...props} key={option.code}>
              <Box>
                <Typography variant="body2" sx={{ fontWeight: 600 }}>
                  {option.code}
                </Typography>
                <Typography
                  variant="caption"
                  sx={{ color: 'text.secondary', display: 'block' }}
                >
                  {option.description}
                </Typography>
                <Typography
                  variant="caption"
                  sx={{
                    color: 'primary.main',
                    fontSize: '0.7rem',
                    fontStyle: 'italic',
                  }}
                >
                  {option.chapterName}
                </Typography>
              </Box>
            </Box>
          )}
          sx={{ width: '100%' }}
        />

        <TextField
          label="Volume (kg)"
          placeholder="1000"
          type="number"
          value={formData.volume_kg}
          onChange={(e) => handleInputChange('volume_kg', e.target.value)}
          error={Boolean(errors.volume_kg)}
          helperText={
            errors.volume_kg || 'Quantidade em quilogramas para simulação'
          }
          disabled={isLoading}
          fullWidth
          required
          slotProps={{
            htmlInput: {
              min: 1,
              max: 1000000,
              step: 1,
            },
          }}
          sx={{
            '& .MuiOutlinedInput-root': {
              backgroundColor: 'background.paper',
            },
          }}
        />

        <Button
          type="submit"
          variant="contained"
          size="large"
          disabled={isLoading}
          fullWidth
          startIcon={
            isLoading ? (
              <CircularProgress size={20} color="inherit" />
            ) : (
              <SearchIcon />
            )
          }
          sx={{
            py: 1.75,
            fontSize: '1rem',
            fontWeight: 600,
            boxShadow: 3,
            '&:hover': {
              boxShadow: 6,
            },
          }}
        >
          {isLoading ? 'Simulando...' : 'Simular Destinos de Exportação'}
        </Button>

        <Typography
          variant="caption"
          sx={{
            textAlign: 'center',
            color: 'text.secondary',
            mt: 1,
          }}
        >
          A simulação leva apenas alguns segundos
        </Typography>
      </Stack>
    </Box>
  );
}
