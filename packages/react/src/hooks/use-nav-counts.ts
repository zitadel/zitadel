"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import { fetchNavCounts, type NavCounts } from "../api/fetch-nav-counts";

const CACHE_TTL_MS = 60_000; // 60 seconds

interface CacheEntry {
  data: NavCounts;
  timestamp: number;
  orgId: string | null;
}

// Module-level cache shared across component instances
let cache: CacheEntry | null = null;

/**
 * Hook that fetches and caches sidebar nav counts.
 * - 60s TTL cache keyed by orgId
 * - On org change, cache is invalidated and refetched
 * - Returns null while loading (no flash of zeros)
 */
export function useNavCounts(orgId?: string | null): NavCounts | null {
  const [counts, setCounts] = useState<NavCounts | null>(() => {
    // Initialize from cache if valid
    const key = orgId ?? null;
    if (cache && cache.orgId === key && Date.now() - cache.timestamp < CACHE_TTL_MS) {
      return cache.data;
    }
    return null;
  });
  const fetchingRef = useRef(false);

  const load = useCallback(async () => {
    const key = orgId ?? null;

    // Check cache first
    if (cache && cache.orgId === key && Date.now() - cache.timestamp < CACHE_TTL_MS) {
      setCounts(cache.data);
      return;
    }

    // Prevent concurrent fetches
    if (fetchingRef.current) return;
    fetchingRef.current = true;

    try {
      const data = await fetchNavCounts(orgId);
      cache = { data, timestamp: Date.now(), orgId: key };
      setCounts(data);
    } catch {
      // Silently fail — sidebar just won't show counts
    } finally {
      fetchingRef.current = false;
    }
  }, [orgId]);

  useEffect(() => {
    load();
  }, [load]);

  return counts;
}
