/**
 * API Route: /v1/routes/compare
 * Proxy requests to Go API backend
 */

import { NextRequest, NextResponse } from 'next/server';

export const dynamic = 'force-dynamic';

export async function GET(request: NextRequest) {
  try {
    // Get API URL from environment or default to internal service
    const apiUrl = process.env.API_URL || 'http://bgc-api:8080';

    // Get all query parameters
    const searchParams = request.nextUrl.searchParams;
    const queryString = searchParams.toString();

    // Build target URL
    const targetUrl = `${apiUrl}/v1/routes/compare${queryString ? `?${queryString}` : ''}`;

    console.log(`[Proxy] GET /v1/routes/compare → ${targetUrl}`);

    // Forward request to API
    const response = await fetch(targetUrl, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const error = await response.text();
      console.error(`[Proxy] Error ${response.status}: ${error}`);
      return NextResponse.json(
        { error: `API error: ${response.statusText}` },
        { status: response.status }
      );
    }

    const data = await response.json();

    // Return data with CORS headers
    return NextResponse.json(data, {
      status: 200,
      headers: {
        'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
        'Access-Control-Allow-Headers': 'Content-Type',
      },
    });
  } catch (error) {
    console.error('[Proxy] Exception:', error);
    return NextResponse.json(
      { error: 'Failed to fetch from API' },
      { status: 500 }
    );
  }
}
