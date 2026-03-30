/**
 * TDD — OnboardingTutorial
 *
 * O tutorial deve aparecer exatamente uma vez após o primeiro onboarding.
 * A flag TUTORIAL_SEEN_KEY é o contrato: se ela estiver ausente = novo usuário.
 */

import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { TUTORIAL_SEEN_KEY } from '@/lib/types/onboarding';

import { OnboardingTutorial } from '@/components/onboarding/OnboardingTutorial';

describe('OnboardingTutorial', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('abre o modal quando TUTORIAL_SEEN_KEY não está no localStorage', () => {
    render(<OnboardingTutorial />);
    expect(screen.getByRole('dialog')).toBeInTheDocument();
  });

  it('não abre o modal quando TUTORIAL_SEEN_KEY já está definido', () => {
    localStorage.setItem(TUTORIAL_SEEN_KEY, '1');
    render(<OnboardingTutorial />);
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('exibe o primeiro slide ao abrir', () => {
    render(<OnboardingTutorial />);
    expect(screen.getByText('Bem-vindo ao BGC!')).toBeInTheDocument();
  });

  it('avança para o próximo slide ao clicar em "Próximo"', () => {
    render(<OnboardingTutorial />);
    fireEvent.click(screen.getByRole('button', { name: /próximo/i }));
    expect(screen.getByText('Simulador de Destinos')).toBeInTheDocument();
  });

  it('volta ao slide anterior ao clicar em "Anterior"', () => {
    render(<OnboardingTutorial />);
    fireEvent.click(screen.getByRole('button', { name: /próximo/i }));
    fireEvent.click(screen.getByRole('button', { name: /anterior/i }));
    expect(screen.getByText('Bem-vindo ao BGC!')).toBeInTheDocument();
  });

  it('persiste a flag e fecha ao clicar em "Começar" no último slide', () => {
    render(<OnboardingTutorial />);
    // Avança até o último slide (4 slides no total)
    fireEvent.click(screen.getByRole('button', { name: /próximo/i }));
    fireEvent.click(screen.getByRole('button', { name: /próximo/i }));
    fireEvent.click(screen.getByRole('button', { name: /próximo/i }));
    // Agora está no último slide — botão diz "Começar"
    fireEvent.click(screen.getByRole('button', { name: /começar/i }));

    expect(localStorage.getItem(TUTORIAL_SEEN_KEY)).toBe('1');
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
  });

  it('persiste a flag ao fechar com o botão X', () => {
    render(<OnboardingTutorial />);
    fireEvent.click(screen.getByRole('button', { name: /fechar tutorial/i }));
    expect(localStorage.getItem(TUTORIAL_SEEN_KEY)).toBe('1');
  });

  it('chama onClose callback ao fechar', () => {
    const onClose = jest.fn();
    render(<OnboardingTutorial onClose={onClose} />);
    fireEvent.click(screen.getByRole('button', { name: /fechar tutorial/i }));
    expect(onClose).toHaveBeenCalledTimes(1);
  });
});
