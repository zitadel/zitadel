import data from '@/components/ErrorReference/data.json';

// Serves the ~1MB error-reference catalog from a route handler instead of a
// static import, so it never lands in the client JS bundle for /apis/errors
// — the browser only pays for it once, on demand, and caches it after.
export function GET() {
  return Response.json(data, {
    headers: { 'Cache-Control': 'public, max-age=31536000, immutable' },
  });
}
