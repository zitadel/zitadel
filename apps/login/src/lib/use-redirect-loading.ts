"use client";

import { useCallback, useEffect, useRef, useState } from "react";

const DEFAULT_REDIRECT_LOADING_TIMEOUT_MS = 6000;

export function useRedirectLoading(timeoutMs: number = DEFAULT_REDIRECT_LOADING_TIMEOUT_MS) {
  const [loading, setLoadingState] = useState<boolean>(false);
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const redirectingRef = useRef<boolean>(false);

  const clearTimeoutRef = useCallback(() => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
      timeoutRef.current = null;
    }
  }, []);

  const setLoading = useCallback((value: boolean) => {
      if (value) {
        clearTimeoutRef();
        redirectingRef.current = false;
        setLoadingState(true);
        return;
      }

      if (redirectingRef.current) {
        return;
      }

      clearTimeoutRef();
      setLoadingState(false);
    },
    [clearTimeoutRef],
  );

  const startLoading = useCallback(() => {
    setLoading(true);
  }, [setLoading]);

  const stopLoading = useCallback(() => {
    setLoading(false);
  }, [setLoading]);

  const startRedirectLoading = useCallback(() => {
    clearTimeoutRef();
    redirectingRef.current = true;
    setLoadingState(true);
    timeoutRef.current = setTimeout(() => {
      redirectingRef.current = false;
      setLoadingState(false);
      timeoutRef.current = null;
    }, timeoutMs);
  }, [clearTimeoutRef, timeoutMs]);

  useEffect(() => {
    return () => {
      clearTimeoutRef();
      redirectingRef.current = false;
    };
  }, [clearTimeoutRef]);

  return {
    loading,
    setLoading,
    startLoading,
    stopLoading,
    startRedirectLoading,
  };
}
