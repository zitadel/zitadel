/**
 * Versioned docs are downloaded into `content/v{major}.{minor}/` by
 * `scripts/fetch-remote-content.mjs` and served under `/docs/v{major}.{minor}/...`.
 * The search index tags every page with its version so results can be scoped to
 * the version currently being viewed.
 */
export const LATEST_VERSION = 'latest';

const VERSION_SEGMENT = /^v\d+\.\d+$/;

/** `['v4.17', 'guides', ...]` -> `'v4.17'`; anything else -> `'latest'`. */
export function getVersionFromSlug(slug: string[] | undefined): string {
  const first = slug?.[0];
  return first && VERSION_SEGMENT.test(first) ? first : LATEST_VERSION;
}

/** `/v4.17/guides/...` -> `'v4.17'`; `/guides/...` -> `'latest'`. */
export function getVersionFromUrl(url: string): string {
  return getVersionFromSlug(url.split('/').filter(Boolean));
}
