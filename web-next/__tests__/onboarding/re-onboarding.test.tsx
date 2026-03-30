/**
 * TDD — Re-onboarding
 *
 * Usuário retorna ao wizard via perfil para atualizar NCM/destinos.
 * Critérios de aceitação:
 *   1. Perfil exibe seção "Configurações de Exportação" com NCM atual
 *   2. Botão "Alterar" navega para /onboarding?re=1
 *   3. Wizard com ?re=1 pré-preenche Step1 com o NCM salvo
 *   4. onboardingVersion é incrementado ao re-salvar
 *
 * Os testes de (3) e (4) validam a lógica de pré-fill no wizard,
 * que será implementada neste Day 7.
 */

import React from 'react';
import { render, screen } from '@testing-library/react';

// ── Mocks globais ─────────────────────────────────────────────────────────────

const mockPush = jest.fn();
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
  useSearchParams: () => ({ get: (key: string) => key === 're' ? '1' : null }),
}));

const existingMeta = {
  onboardingCompleted: true,
  onboardingVersion: 1,
  onboardingCompletedAt: '2026-01-01T00:00:00.000Z',
  ncmChapter: '09',
  ncmLabel: 'Café, chá, mate e especiarias',
  ncmDefaultNcm8d: '09010110',
  volumeMonthlyKg: 1000,
  targetContinents: ['europa'],
  targetCountries: ['DEU'],
  onboardingSkipped: false,
};

const mockUserObj = {
  firstName: 'Rafael',
  lastName: 'Silva',
  fullName: 'Rafael Silva',
  imageUrl: '',
  createdAt: new Date().getTime(),
  primaryEmailAddress: { emailAddress: 'rafael@test.com', verification: { status: 'verified' } },
  emailAddresses: [{ emailAddress: 'rafael@test.com', verification: { status: 'verified' } }],
  externalAccounts: [],
  passwordEnabled: false,
  publicMetadata: existingMeta,
  reload: jest.fn(),
  update: jest.fn(),
  delete: jest.fn(),
};

jest.mock('@clerk/nextjs', () => ({
  useUser: () => ({ user: mockUserObj, isLoaded: true }),
  useClerk: () => ({ signOut: jest.fn() }),
  useSignIn: () => ({}),
  useSignUp: () => ({}),
}));

// ── Tests: Profile page ───────────────────────────────────────────────────────

describe('Profile page — seção de configurações de exportação', () => {
  it('exibe o NCM salvo no perfil do usuário', async () => {
    const { default: ProfilePage } = await import('@/app/profile/page');
    render(<ProfilePage />);

    // Deve mostrar o código e a descrição do NCM
    expect(screen.getByText(/09/)).toBeInTheDocument();
    expect(screen.getByText(/café/i)).toBeInTheDocument();
  });

  it('exibe botão "Alterar" que navega para /onboarding?re=1', async () => {
    const { default: ProfilePage } = await import('@/app/profile/page');
    render(<ProfilePage />);

    const alterarBtn = screen.getByRole('button', { name: /alterar/i });
    expect(alterarBtn).toBeInTheDocument();

    const { fireEvent } = await import('@testing-library/react');
    fireEvent.click(alterarBtn);
    expect(mockPush).toHaveBeenCalledWith('/onboarding?re=1');
  });
});

// ── Tests: Onboarding wizard pré-preenchido ───────────────────────────────────

describe('Onboarding wizard com ?re=1 — pré-fill do NCM existente', () => {
  it('exibe o NCM previamente salvo no Step 1', async () => {
    const { default: OnboardingPage } = await import('@/app/onboarding/page');
    render(<OnboardingPage />);

    // O wizard deve mostrar o NCM salvo no Step 1
    expect(screen.getByText(/09 — Café/i)).toBeInTheDocument();
  });
});

// ── Tests: onboardingVersion ──────────────────────────────────────────────────

describe('saveMetadata — incremento de onboardingVersion no re-onboarding', () => {
  it('envia onboardingVersion=2 quando versão atual é 1', async () => {
    const fetchMock = jest.fn().mockResolvedValue({ ok: true, json: jest.fn() });
    global.fetch = fetchMock as unknown as typeof fetch;

    // Simula clique em "Concluir" no wizard
    const { default: OnboardingPage } = await import('@/app/onboarding/page');
    const { fireEvent: fe } = await import('@testing-library/react');
    render(<OnboardingPage />);

    // Navega até o último passo e conclui
    // Step 1 já está pré-preenchido com NCM do ?re=1
    const nextBtn = screen.getByRole('button', { name: /próximo/i });
    fe.click(nextBtn); // → step 2
    fe.click(screen.getByRole('button', { name: /próximo/i })); // → step 3
    fe.click(screen.getByRole('button', { name: /concluir/i })); // → save

    // Valida que o fetch foi chamado com onboardingVersion incrementado
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/onboarding',
      expect.objectContaining({
        method: 'POST',
        body: expect.stringContaining('"onboardingVersion":2'),
      })
    );
  });
});
