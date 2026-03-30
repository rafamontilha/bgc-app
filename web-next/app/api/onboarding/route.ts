import { auth, clerkClient } from '@clerk/nextjs/server';
import { NextRequest, NextResponse } from 'next/server';
import type { OnboardingMetadata } from '@/lib/types/onboarding';

/**
 * POST /api/onboarding
 * Saves onboarding data to user.publicMetadata via Clerk Backend SDK.
 * Must be server-side: publicMetadata is write-protected from the browser.
 */
export async function POST(req: NextRequest) {
  const { userId } = await auth();

  if (!userId) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }

  const body = await req.json();
  const metadata: OnboardingMetadata = body;

  // Basic validation
  if (typeof metadata.onboardingCompleted !== 'boolean') {
    return NextResponse.json({ error: 'Invalid payload' }, { status: 400 });
  }

  const client = await clerkClient();
  await client.users.updateUserMetadata(userId, {
    publicMetadata: metadata as unknown as Record<string, unknown>,
  });

  return NextResponse.json({ ok: true });
}
