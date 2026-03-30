/**
 * @jest-environment node
 *
 * TDD — POST /api/onboarding
 */

import { NextRequest } from 'next/server';

// ── Mocks ────────────────────────────────────────────────────────────────────
// jest.mock é hoisted — não pode referenciar variáveis declaradas fora do factory.
// Usamos jest.fn() inline e acessamos via módulo mockado depois.

jest.mock('@clerk/nextjs/server', () => ({
  auth: jest.fn(),
  clerkClient: jest.fn(),
}));

// Acessa os mocks após o hoist
import { auth, clerkClient } from '@clerk/nextjs/server';
const mockAuth = auth as unknown as jest.Mock;
const mockClerkClient = clerkClient as unknown as jest.Mock;

const mockUpdateUserMetadata = jest.fn().mockResolvedValue({});

// ── Helpers ───────────────────────────────────────────────────────────────────

function makeRequest(body: unknown): NextRequest {
  return new NextRequest('http://localhost/api/onboarding', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
}

const validPayload = {
  onboardingCompleted: true,
  onboardingVersion: 1,
  onboardingCompletedAt: new Date().toISOString(),
  ncmChapter: '09',
  ncmLabel: 'Café, chá, mate e especiarias',
  ncmDefaultNcm8d: '09010110',
  volumeMonthlyKg: 1000,
  targetContinents: ['europa'],
  targetCountries: ['DEU', 'FRA'],
  onboardingSkipped: false,
};

// ── Tests ─────────────────────────────────────────────────────────────────────

describe('POST /api/onboarding', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    mockClerkClient.mockResolvedValue({
      users: { updateUserMetadata: mockUpdateUserMetadata },
    });
  });

  it('retorna 401 quando usuário não está autenticado', async () => {
    mockAuth.mockResolvedValue({ userId: null });

    const { POST } = await import('@/app/api/onboarding/route');
    const res = await POST(makeRequest(validPayload));

    expect(res.status).toBe(401);
    const body = await res.json();
    expect(body.error).toBe('Unauthorized');
  });

  it('retorna 200 e persiste metadata quando usuário está autenticado', async () => {
    mockAuth.mockResolvedValue({ userId: 'user_123' });

    const { POST } = await import('@/app/api/onboarding/route');
    const res = await POST(makeRequest(validPayload));

    expect(res.status).toBe(200);
    const body = await res.json();
    expect(body.ok).toBe(true);
  });

  it('chama updateUserMetadata com o userId correto', async () => {
    mockAuth.mockResolvedValue({ userId: 'user_abc' });

    const { POST } = await import('@/app/api/onboarding/route');
    await POST(makeRequest(validPayload));

    expect(mockUpdateUserMetadata).toHaveBeenCalledWith(
      'user_abc',
      expect.objectContaining({ publicMetadata: expect.any(Object) })
    );
  });

  it('retorna 400 quando payload não tem onboardingCompleted', async () => {
    mockAuth.mockResolvedValue({ userId: 'user_123' });

    const { POST } = await import('@/app/api/onboarding/route');
    const res = await POST(makeRequest({ ncmChapter: '09' }));

    expect(res.status).toBe(400);
    const body = await res.json();
    expect(body.error).toBeDefined();
  });

  it('persiste onboardingSkipped=true quando usuário pula', async () => {
    mockAuth.mockResolvedValue({ userId: 'user_skip' });

    const skippedPayload = {
      ...validPayload,
      onboardingCompleted: false,
      onboardingSkipped: true,
    };

    const { POST } = await import('@/app/api/onboarding/route');
    await POST(makeRequest(skippedPayload));

    expect(mockUpdateUserMetadata).toHaveBeenCalledWith(
      'user_skip',
      expect.objectContaining({
        publicMetadata: expect.objectContaining({ onboardingSkipped: true }),
      })
    );
  });

  it('persiste onboardingVersion=2 no re-onboarding', async () => {
    mockAuth.mockResolvedValue({ userId: 'user_reob' });

    const reOnboardPayload = { ...validPayload, onboardingVersion: 2 };

    const { POST } = await import('@/app/api/onboarding/route');
    await POST(makeRequest(reOnboardPayload));

    expect(mockUpdateUserMetadata).toHaveBeenCalledWith(
      'user_reob',
      expect.objectContaining({
        publicMetadata: expect.objectContaining({ onboardingVersion: 2 }),
      })
    );
  });
});
