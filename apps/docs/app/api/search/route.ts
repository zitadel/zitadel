// app/api/search/route.ts
import { NextResponse } from 'next/server';

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url);
    const query = searchParams.get('q');
    const searchType = searchParams.get('type') || 'docs';

    if (!query || query.trim() === '') {
      return NextResponse.json([]);
    }
    const safeQuery = query.substring(0, 150).trim();

    // Extract original client IP to pass down to the Express rate limiter
    const forwardedFor = request.headers.get('x-forwarded-for');
    const clientIp = forwardedFor ? forwardedFor.split(',')[0].trim() : '127.0.0.1';

    const baseUrl = process.env.DOCS_SEARCH_URL || 'http://localhost:8080';
    const cleanBaseUrl = baseUrl.replace(/\/+$/, '');

    const backendUrl = new URL(`${cleanBaseUrl}/api/search/docs`);
    backendUrl.searchParams.set('q', safeQuery);
    backendUrl.searchParams.set('limit', '15');
    backendUrl.searchParams.set('type', searchType);

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 8000);

    const response = await fetch(backendUrl.toString(), {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${process.env.DOCS_SEARCH_SECRET}`,
        'Content-Type': 'application/json',
        'X-Forwarded-For': clientIp // Forwarding IP to Express
      },
      signal: controller.signal,
      next: { revalidate: 300 }
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      if (response.status === 429) {
        return NextResponse.json({ error: 'Rate limit exceeded' }, { status: 429 });
      }
      return NextResponse.json({ error: 'Search service unavailable' }, { status: response.status });
    }

    const data = await response.json();
    return NextResponse.json(data);

  } catch (error: any) {
    console.error('Proxy Search Error:', error.name === 'AbortError' ? 'Timeout' : error.message);
    const status = error.name === 'AbortError' ? 504 : 500;
    return NextResponse.json({ error: 'Search service connection failed' }, { status });
  }
}
