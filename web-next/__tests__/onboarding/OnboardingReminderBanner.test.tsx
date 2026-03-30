/**
 * TDD — OnboardingReminderBanner
 *
 * Comportamento crítico: o banner só aparece uma vez — deve ser dispensado
 * persistentemente. Se a flag de localStorage falhar, o usuário vê o banner
 * a cada reload, o que é disruptivo.
 */

import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { REMINDER_BANNER_DISMISSED_KEY } from '@/lib/types/onboarding';

// Mock next/navigation
const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}));

import { OnboardingReminderBanner } from '@/components/onboarding/OnboardingReminderBanner';

describe('OnboardingReminderBanner', () => {
  beforeEach(() => {
    localStorage.clear();
    mockPush.mockClear();
  });

  it('exibe o banner quando a flag não está no localStorage', () => {
    render(<OnboardingReminderBanner />);
    expect(screen.getByRole('banner')).toBeInTheDocument();
  });

  it('não exibe o banner quando já foi dispensado', () => {
    localStorage.setItem(REMINDER_BANNER_DISMISSED_KEY, '1');
    render(<OnboardingReminderBanner />);
    expect(screen.queryByRole('banner')).not.toBeInTheDocument();
  });

  it('dispensa o banner ao clicar no botão X', async () => {
    render(<OnboardingReminderBanner />);
    const closeBtn = screen.getByRole('button', { name: /dispensar/i });
    fireEvent.click(closeBtn);
    // MUI Collapse usa transição CSS — aguarda unmount após state update
    await waitFor(() =>
      expect(screen.queryByRole('banner')).not.toBeInTheDocument()
    );
  });

  it('persiste a flag no localStorage ao dispensar', () => {
    render(<OnboardingReminderBanner />);
    const closeBtn = screen.getByRole('button', { name: /dispensar/i });
    fireEvent.click(closeBtn);
    expect(localStorage.getItem(REMINDER_BANNER_DISMISSED_KEY)).toBe('1');
  });

  it('navega para /onboarding ao clicar em "Completar agora"', () => {
    render(<OnboardingReminderBanner />);
    const cta = screen.getByRole('button', { name: /completar agora/i });
    fireEvent.click(cta);
    expect(mockPush).toHaveBeenCalledWith('/onboarding');
  });

  it('persiste a flag ao clicar em "Completar agora" (evita reexibição)', () => {
    render(<OnboardingReminderBanner />);
    const cta = screen.getByRole('button', { name: /completar agora/i });
    fireEvent.click(cta);
    expect(localStorage.getItem(REMINDER_BANNER_DISMISSED_KEY)).toBe('1');
  });
});
