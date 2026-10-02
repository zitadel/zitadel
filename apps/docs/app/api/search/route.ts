// app/api/search/route.ts
import { NextResponse } from 'next/server';

const RATE_LIMIT_MAX = 40; 
const MAX_CACHE_SIZE = 10000; 
const rateLimitCache = new Map<string, { count: number; resetTime: number }>();

function isRateLimited(ip: string): boolean {
  const now = Date.now();
  const record = rateLimitCache.get(ip);

  if (!record || now > record.resetTime) {
    if (rateLimitCache.size > MAX_CACHE_SIZE) {
      rateLimitCache.clear();
    }
    rateLimitCache.set(ip, { count: 1, resetTime: now + 60000 });
    return false;
  }

  if (record.count >= RATE_LIMIT_MAX) return true;

  record.count++;
  return false;
}

export async function GET(request: Request) {
  try {
    const forwardedFor = request.headers.get('x-forwarded-for');
    const ip = forwardedFor ? forwardedFor.split(',')[0].trim() : '127.0.0.1';

    if (isRateLimited(ip)) {
      return NextResponse.json({ error: 'Rate limit exceeded' }, { status: 429 });
    }

    const { searchParams } = new URL(request.url);
    const query = searchParams.get('q');
    const searchType = searchParams.get('type') || 'docs'; // Extract the type parameter

    if (!query || query.trim() === '') {
      return NextResponse.json([]);
    }
    const safeQuery = query.substring(0, 150).trim();

    const baseUrl = process.env.DOCS_SEARCH_URL || 'http://localhost:8080';
    const cleanBaseUrl = baseUrl.replace(/\/+$/, '');

    const backendUrl = new URL(`${cleanBaseUrl}/api/search/docs`);
    backendUrl.searchParams.set('q', safeQuery);
    backendUrl.searchParams.set('limit', '15');
    backendUrl.searchParams.set('type', searchType); // Forward the type parameter

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 8000);

    const response = await fetch(backendUrl.toString(), {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${process.env.DOCS_SEARCH_SECRET}`,
        'Content-Type': 'application/json'
      },
      signal: controller.signal,
      next: { revalidate: 300 }
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      throw new Error(`Express returned status: ${response.status}`);
    }

    const data = await response.json();
    return NextResponse.json(data);

  } catch (error: any) {
    console.error('Proxy Search Error:', error.name === 'AbortError' ? 'Timeout' : error.message);
    return NextResponse.json([]);
  }
}