import { LRUCache } from "lru-cache";

interface FetchContext {
  fetcher: () => Promise<any>;
}

/**
 * A bounded, stale-while-revalidate in-memory promise cache backed by lru-cache.
 *
 * Features:
 * - True LRU eviction
 * - Deduplicates concurrent requests (built-in to lru-cache's fetchMethod)
 * - Serves stale data immediately while revalidating in the background
 * - Keeps stale value on fetch rejection
 * - Bounded to `maxSize` entries to prevent unbounded memory growth
 */
export class PromiseCache {
  private readonly cache: LRUCache<string, any, FetchContext>;

  constructor(maxSize = 100_000, perf?: { now: () => number }) {
    this.cache = new LRUCache<string, any, FetchContext>({
      max: Math.max(1, maxSize),
      // A global TTL is required to initialize lru-cache's TTL tracking internals.
      // Per-entry TTLs passed to fetch() will override this default.
      ttl: 1,
      allowStale: true,
      noDeleteOnStaleGet: true,
      noDeleteOnFetchRejection: true,
      allowStaleOnFetchRejection: true,
      fetchMethod: async (_key, _staleValue, { context }) => {
        return context.fetcher();
      },
      ...(perf ? { perf, ttlResolution: 0 } : {}),
    });
  }

  /**
   * Get a cached value or execute the fetcher and cache its result.
   *
   * After the first successful fetch, expired entries return the stale
   * value immediately and trigger a background revalidation (SWR).
   * Only the very first call for a key (or after eviction) blocks
   * on the fetch.
   */
  getOrFetch<T>(key: string, fetcher: () => Promise<T>, ttlMs: number): Promise<T> {
    return this.cache.forceFetch(key, {
      ttl: ttlMs,
      context: { fetcher },
    }) as Promise<T>;
  }

  /** Current number of entries (including stale). */
  get size(): number {
    return this.cache.size;
  }

  /** Clear all entries. */
  clear(): void {
    this.cache.clear();
  }
}
