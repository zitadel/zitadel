import data from '@/components/EndpointErrors/data.json';

// Serves the per-endpoint error tables from a route handler instead of a
// static import, so the whole catalog (every category, every operation)
// never lands in the client JS bundle for every single endpoint page —
// the browser only pays for it once, on demand, and caches it after.
export function GET() {
  return Response.json(data, {
    headers: { 'Cache-Control': 'public, max-age=31536000, immutable' },
  });
}
