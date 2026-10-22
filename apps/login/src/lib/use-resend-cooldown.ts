"use client";

import { useCallback, useEffect, useState } from "react";
import { RESEND_COOLDOWN_SECONDS } from "./constants";

/**
 * Client-side cooldown for "Resend code" buttons. This is a UX guard against
 * accidental repeated clicks (each resend invalidates the previous code), not
 * a security control.
 *
 * `start()` begins the countdown; `remaining` is the number of seconds left
 * (0 when the button may be used again).
 */
export function useResendCooldown(seconds: number = RESEND_COOLDOWN_SECONDS) {
  const [endsAt, setEndsAt] = useState<number | null>(null);
  const [remaining, setRemaining] = useState(0);

  useEffect(() => {
    if (endsAt === null) return;

    const tick = () => {
      const left = Math.max(0, Math.ceil((endsAt - Date.now()) / 1000));
      setRemaining(left);
      if (left === 0) setEndsAt(null);
    };

    tick();
    const interval = setInterval(tick, 250);
    return () => clearInterval(interval);
  }, [endsAt]);

  const start = useCallback(() => {
    setEndsAt(Date.now() + seconds * 1000);
    setRemaining(seconds);
  }, [seconds]);

  const isCoolingDown = remaining > 0;
  // e.g. " (24s)" while cooling down, so labels can append it
  const suffix = isCoolingDown ? ` (${remaining}s)` : "";

  return { remaining, isCoolingDown, suffix, start };
}
