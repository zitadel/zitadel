import data from '@/components/EndpointErrors/data.json';

// Serves the per-endpoint error tables from a route handler instead of a
// static import, so the whole catalog (every category, every operation)
// never lands in the client JS bundle for every single endpoint page,
// the browser only pays for it once, on demand, and caches it after.
//
// The cache header here has to allow for the fact that this same URL
// serves different content across deployments, this file gets rebuilt
// with fresh data every time a category's tracing JSON changes. A long
// max-age plus immutable, the pattern that's correct for a hashed asset
// URL, would be wrong here: a browser or CDN that cached the response
// before a deploy could keep serving it for up to a year afterward,
// since nothing about this URL itself changes to bust that cache. A
// short max-age with revalidation in the background keeps repeat views
// fast without risking that kind of long-lived staleness.
export function GET() {
  return Response.json(data, {
    headers: { 'Cache-Control': 'public, max-age=300, stale-while-revalidate=3600' },
  });
}
