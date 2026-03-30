'use client';

import React, { useState } from 'react';
import {
  Box,
  Typography,
  Chip,
  Collapse,
  List,
  ListItemButton,
  ListItemText,
  ListItemIcon,
  Divider,
} from '@mui/material';
import PublicIcon from '@mui/icons-material/Public';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import ExpandLessIcon from '@mui/icons-material/ExpandLess';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import type { TradeRegion } from '@/lib/types/onboarding';
import { CONTINENTS, getCountriesByContinent } from '@/lib/data/trade-regions';

export interface OnboardingStep3RegionsProps {
  value: TradeRegion[];
  onChange: (regions: TradeRegion[]) => void;
}

export function OnboardingStep3Regions({
  value,
  onChange,
}: OnboardingStep3RegionsProps): React.ReactElement {
  const [expanded, setExpanded] = useState<string | null>(null);

  const selectedCodes = new Set(value.map((r) => r.code));

  function toggleItem(item: TradeRegion) {
    if (selectedCodes.has(item.code)) {
      onChange(value.filter((r) => r.code !== item.code));
    } else {
      onChange([...value, item]);
    }
  }

  function toggleContinent(continent: TradeRegion) {
    const countries = getCountriesByContinent(continent.code);
    const allSelected = countries.every((c) => selectedCodes.has(c.code));

    if (allSelected) {
      // deselect continent + all its countries
      const remove = new Set([continent.code, ...countries.map((c) => c.code)]);
      onChange(value.filter((r) => !remove.has(r.code)));
    } else {
      // select continent + all missing countries
      const toAdd: TradeRegion[] = [];
      if (!selectedCodes.has(continent.code)) toAdd.push(continent);
      for (const c of countries) {
        if (!selectedCodes.has(c.code)) toAdd.push(c);
      }
      onChange([...value, ...toAdd]);
    }
  }

  function toggleExpand(code: string) {
    setExpanded((prev) => (prev === code ? null : code));
  }

  function removeChip(code: string) {
    // when removing a continent chip, also remove its countries
    const continent = CONTINENTS.find((c) => c.code === code);
    if (continent) {
      const countryCodes = new Set(
        getCountriesByContinent(code).map((c) => c.code)
      );
      onChange(value.filter((r) => r.code !== code && !countryCodes.has(r.code)));
    } else {
      onChange(value.filter((r) => r.code !== code));
    }
  }

  const selectedCount = value.length;

  return (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, mb: 1 }}>
        <PublicIcon sx={{ color: 'var(--md-sys-color-primary)', fontSize: 28 }} />
        <Typography variant="h6" sx={{ fontWeight: 600, color: 'var(--md-sys-color-on-surface)' }}>
          Para onde você quer exportar?
        </Typography>
      </Box>

      <Typography
        variant="body2"
        sx={{ color: 'var(--md-sys-color-on-surface-variant)', mb: 3, pl: '44px' }}
      >
        Selecione continentes ou países específicos. Você pode escolher mais de
        um. Esta etapa é opcional.
      </Typography>

      {/* Selected chips */}
      {selectedCount > 0 && (
        <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 2 }}>
          {value.map((item) => (
            <Chip
              key={item.code}
              label={item.label}
              onDelete={() => removeChip(item.code)}
              size="small"
              sx={{
                bgcolor:
                  item.type === 'continent'
                    ? 'var(--md-sys-color-primary)'
                    : 'var(--md-sys-color-secondary-container)',
                color:
                  item.type === 'continent'
                    ? 'var(--md-sys-color-on-primary)'
                    : 'var(--md-sys-color-on-secondary-container)',
                '& .MuiChip-deleteIcon': {
                  color:
                    item.type === 'continent'
                      ? 'var(--md-sys-color-on-primary)'
                      : 'var(--md-sys-color-on-secondary-container)',
                  opacity: 0.7,
                },
              }}
            />
          ))}
        </Box>
      )}

      {/* Continent list */}
      <Box
        sx={{
          border: '1px solid var(--md-sys-color-outline-variant)',
          borderRadius: '12px',
          overflow: 'hidden',
        }}
      >
        {CONTINENTS.map((continent, idx) => {
          const countries = getCountriesByContinent(continent.code);
          const isOpen = expanded === continent.code;
          const isContinentSelected = selectedCodes.has(continent.code);
          const selectedCountriesCount = countries.filter((c) =>
            selectedCodes.has(c.code)
          ).length;

          return (
            <React.Fragment key={continent.code}>
              {idx > 0 && <Divider />}
              <ListItemButton
                onClick={() => toggleExpand(continent.code)}
                sx={{
                  px: 2,
                  py: 1.5,
                  bgcolor: isContinentSelected
                    ? 'var(--md-sys-color-primary-container)'
                    : 'transparent',
                  '&:hover': {
                    bgcolor: isContinentSelected
                      ? 'var(--md-sys-color-primary-container)'
                      : 'var(--md-sys-color-surface-variant)',
                    opacity: 0.9,
                  },
                }}
              >
                <ListItemIcon sx={{ minWidth: 36 }}>
                  {isContinentSelected ? (
                    <CheckCircleIcon
                      fontSize="small"
                      sx={{ color: 'var(--md-sys-color-primary)' }}
                    />
                  ) : (
                    <PublicIcon
                      fontSize="small"
                      sx={{ color: 'var(--md-sys-color-on-surface-variant)' }}
                    />
                  )}
                </ListItemIcon>

                <ListItemText
                  primary={continent.label}
                  secondary={
                    selectedCountriesCount > 0
                      ? `${selectedCountriesCount} país${selectedCountriesCount > 1 ? 'es' : ''} selecionado${selectedCountriesCount > 1 ? 's' : ''}`
                      : `${countries.length} países`
                  }
                  primaryTypographyProps={{ fontWeight: 600, variant: 'body2' }}
                  secondaryTypographyProps={{ variant: 'caption' }}
                />

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                  <Chip
                    label={isContinentSelected ? 'Todo continente' : 'Selecionar tudo'}
                    size="small"
                    onClick={(e) => {
                      e.stopPropagation();
                      toggleContinent(continent);
                    }}
                    sx={{
                      fontSize: '0.7rem',
                      bgcolor: isContinentSelected
                        ? 'var(--md-sys-color-primary)'
                        : 'transparent',
                      color: isContinentSelected
                        ? 'var(--md-sys-color-on-primary)'
                        : 'var(--md-sys-color-on-surface-variant)',
                      border: isContinentSelected
                        ? 'none'
                        : '1px solid var(--md-sys-color-outline)',
                    }}
                  />
                  {isOpen ? (
                    <ExpandLessIcon fontSize="small" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }} />
                  ) : (
                    <ExpandMoreIcon fontSize="small" sx={{ color: 'var(--md-sys-color-on-surface-variant)' }} />
                  )}
                </Box>
              </ListItemButton>

              <Collapse in={isOpen} unmountOnExit>
                <List dense disablePadding>
                  {countries.map((country) => {
                    const isCountrySelected = selectedCodes.has(country.code);
                    return (
                      <React.Fragment key={country.code}>
                        <Divider sx={{ ml: 8 }} />
                        <ListItemButton
                          onClick={() => toggleItem(country)}
                          sx={{
                            pl: 6,
                            py: 1,
                            bgcolor: isCountrySelected
                              ? 'var(--md-sys-color-secondary-container)'
                              : 'transparent',
                            '&:hover': {
                              bgcolor: isCountrySelected
                                ? 'var(--md-sys-color-secondary-container)'
                                : 'var(--md-sys-color-surface-variant)',
                              opacity: 0.9,
                            },
                          }}
                        >
                          <ListItemIcon sx={{ minWidth: 32 }}>
                            {isCountrySelected ? (
                              <CheckCircleIcon
                                fontSize="small"
                                sx={{ color: 'var(--md-sys-color-secondary)', fontSize: 18 }}
                              />
                            ) : (
                              <Box sx={{ width: 18, height: 18 }} />
                            )}
                          </ListItemIcon>
                          <ListItemText
                            primary={country.label}
                            primaryTypographyProps={{ variant: 'body2' }}
                          />
                          <Typography
                            variant="caption"
                            sx={{
                              fontFamily: 'monospace',
                              color: 'var(--md-sys-color-on-surface-variant)',
                            }}
                          >
                            {country.iso3}
                          </Typography>
                        </ListItemButton>
                      </React.Fragment>
                    );
                  })}
                </List>
              </Collapse>
            </React.Fragment>
          );
        })}
      </Box>

      {selectedCount === 0 && (
        <Typography
          variant="caption"
          sx={{ mt: 1.5, display: 'block', color: 'var(--md-sys-color-on-surface-variant)' }}
        >
          Nenhuma região selecionada — você pode pular este passo.
        </Typography>
      )}
    </Box>
  );
}
