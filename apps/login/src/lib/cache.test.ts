import { afterEach, describe, expect, test, vi } from "vitest";
import { PromiseCache } from "./cache";

// Suppress logger output during tests
vi.mock("./logger", () => ({
  createLogger: () => ({
    warn: vi.fn(),
    info: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

/**
 * Creates a mock clock for lru-cache's `perf` option.
 * lru-cache uses `performance.now()` internally for TTL tracking, and caches
 * the reference at import time. The `perf` option lets us inject a custom clock.
 */
function mockClock(initial = 1000) {
  let now = initial;
  return {
    perf: { now: () => now },
    advance: (ms: number) => {
      now += ms;
    },
  };
}

describe("PromiseCache", () => {
  let cache: PromiseCache;

  afterEach(() => {
    cache?.clear();
  });

  describe("getOrFetch", () => {
    test("should return the fetcher result on cache miss", async () => {
      cache = new PromiseCache(10);
      const result = await cache.getOrFetch("key1", () => Promise.resolve("value1"), 60_000);
      expect(result).toBe("value1");
      expect(cache.size).toBe(1);
    });

    test("should return cached value on cache hit", async () => {
      cache = new PromiseCache(10);
      let callCount = 0;
      const fetcher = () => {
        callCount++;
        return Promise.resolve(`value-${callCount}`);
      };

      const first = await cache.getOrFetch("key1", fetcher, 60_000);
      const second = await cache.getOrFetch("key1", fetcher, 60_000);

      expect(first).toBe("value-1");
      expect(second).toBe("value-1");
      expect(callCount).toBe(1);
    });

    test("should return stale value immediately after TTL expires (SWR)", async () => {
      const clock = mockClock();
      cache = new PromiseCache(10, clock.perf);

      let callCount = 0;
      let resolveRevalidation: ((v: string) => void) | undefined;

      const fetcher = () => {
        callCount++;
        if (callCount === 1) {
          return Promise.resolve("value-1");
        }
        // Second call: return a promise we control
        return new Promise<string>((resolve) => {
          resolveRevalidation = resolve;
        });
      };

      // First call — blocks on fetch
      const first = await cache.getOrFetch("key1", fetcher, 100);
      expect(first).toBe("value-1");
      expect(callCount).toBe(1);

      // Expire the entry
      clock.advance(101);

      // Second call after expiry — should get stale value immediately
      const second = await cache.getOrFetch("key1", fetcher, 100);
      expect(second).toBe("value-1"); // stale value, not blocking
      expect(callCount).toBe(2); // revalidation was triggered

      // Third call while revalidation is in-flight — also gets stale value
      const third = await cache.getOrFetch("key1", fetcher, 100);
      expect(third).toBe("value-1"); // stale value
      expect(callCount).toBe(2); // no additional fetch (already revalidating)

      // Resolve the background revalidation and flush the .then() chain
      resolveRevalidation!("value-2");
      await new Promise((r) => setTimeout(r, 0));
      await new Promise((r) => setTimeout(r, 0));

      // Now we should get the fresh value
      const fourth = await cache.getOrFetch("key1", fetcher, 100);
      expect(fourth).toBe("value-2");
    });

    test("should keep stale value when revalidation fails", async () => {
      const clock = mockClock();
      cache = new PromiseCache(10, clock.perf);

      let callCount = 0;

      const fetcher = () => {
        callCount++;
        if (callCount === 1) {
          return Promise.resolve("value-1");
        }
        return Promise.reject(new Error("network error"));
      };

      const first = await cache.getOrFetch("key1", fetcher, 100);
      expect(first).toBe("value-1");

      clock.advance(101);

      // After TTL: returns stale, triggers failing revalidation
      const second = await cache.getOrFetch("key1", fetcher, 100);
      expect(second).toBe("value-1");
      expect(callCount).toBe(2);

      // Flush the rejected promise's handlers
      await new Promise((r) => setTimeout(r, 0));
      await new Promise((r) => setTimeout(r, 0));

      // Next call can retry
      clock.advance(1);
      const third = await cache.getOrFetch("key1", fetcher, 100);
      expect(third).toBe("value-1"); // still stale
      expect(callCount).toBe(3); // retried
    });

    test("should re-fetch after entry is evicted", async () => {
      cache = new PromiseCache(10);
      let callCount = 0;
      const fetcher = () => {
        callCount++;
        return Promise.resolve(`value-${callCount}`);
      };

      const first = await cache.getOrFetch("key1", fetcher, 60_000);
      expect(first).toBe("value-1");

      // Evict the entry manually
      cache.clear();

      const second = await cache.getOrFetch("key1", fetcher, 60_000);
      expect(second).toBe("value-2");
      expect(callCount).toBe(2);
    });

    test("should reject on first-fetch failure", async () => {
      cache = new PromiseCache(10);
      const failingFetcher = () => Promise.reject(new Error("fail"));

      await expect(cache.getOrFetch("key1", failingFetcher, 60_000)).rejects.toThrow();

      // Wait a tick for internal cleanup
      await new Promise((r) => setTimeout(r, 0));
      expect(cache.size).toBe(0);
    });

    test("should deduplicate concurrent requests for the same key", async () => {
      cache = new PromiseCache(10);
      let callCount = 0;
      const fetcher = () => {
        callCount++;
        return new Promise<string>((resolve) => setTimeout(() => resolve(`value-${callCount}`), 10));
      };

      const [a, b] = await Promise.all([
        cache.getOrFetch("key1", fetcher, 60_000),
        cache.getOrFetch("key1", fetcher, 60_000),
      ]);

      expect(a).toBe("value-1");
      expect(b).toBe("value-1");
      expect(callCount).toBe(1);
    });
  });

  describe("maxSize eviction", () => {
    test("should evict entries when maxSize is exceeded", async () => {
      cache = new PromiseCache(3);

      await cache.getOrFetch("a", () => Promise.resolve(1), 60_000);
      await cache.getOrFetch("b", () => Promise.resolve(2), 60_000);
      await cache.getOrFetch("c", () => Promise.resolve(3), 60_000);
      expect(cache.size).toBe(3);

      // Adding a 4th entry should trigger eviction of the LRU entry ("a")
      await cache.getOrFetch("d", () => Promise.resolve(4), 60_000);
      expect(cache.size).toBe(3);

      // "a" should have been evicted — re-fetching should call a new fetcher
      let refetched = false;
      await cache.getOrFetch(
        "a",
        () => {
          refetched = true;
          return Promise.resolve(10);
        },
        60_000,
      );
      expect(refetched).toBe(true);
    });

    test("should respect maxSize of 1", async () => {
      cache = new PromiseCache(1);

      await cache.getOrFetch("a", () => Promise.resolve(1), 60_000);
      expect(cache.size).toBe(1);

      await cache.getOrFetch("b", () => Promise.resolve(2), 60_000);
      expect(cache.size).toBe(1);

      // Only "b" should remain
      let aRefetched = false;
      await cache.getOrFetch(
        "a",
        () => {
          aRefetched = true;
          return Promise.resolve(10);
        },
        60_000,
      );
      expect(aRefetched).toBe(true);
    });
  });

  describe("clear", () => {
    test("should remove all entries", async () => {
      cache = new PromiseCache(10);

      await cache.getOrFetch("a", () => Promise.resolve(1), 60_000);
      await cache.getOrFetch("b", () => Promise.resolve(2), 60_000);
      expect(cache.size).toBe(2);

      cache.clear();
      expect(cache.size).toBe(0);
    });
  });
});
