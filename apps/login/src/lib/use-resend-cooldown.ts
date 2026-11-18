"use client";

import { useCallback, useEffect, useState } from "react";
import { RESEND_COOLDOWN_SECONDS } from "./constants";

/**
 * Countdown for "Resend code" buttons. Started with the default duration it is
 * a UX guard against accidental repeated clicks (each resend invalidates the
 * previous code), not a security control; the OTP page starts it with the
 * cooldown the server enforces.
 *
 * `start()` begins the countdown, for `seconds` or the duration passed to it
 * (e.g. the server's `retryAfterSeconds`); `remaining` is the number of seconds
 * left (0 when the button may be used again). `startNow` starts it whenever it
 * turns true, for a page that a code was just sent before arriving on; set on
 * the first render, the countdown already runs in that render, so the button is
 * never shown enabled first.
 */
export function useResendCooldown({
  seconds = RESEND_COOLDOWN_SECONDS,
  startNow = false,
}: { seconds?: number; startNow?: boolean } = {}) {
  const [endsAt, setEndsAt] = useState<number | null>(() => (startNow ? Date.now() + seconds * 1000 : null));
  const [remaining, setRemaining] = useState(() => (startNow ? seconds : 0));

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

  const start = useCallback(
    (duration: number = seconds) => {
      setEndsAt(Date.now() + duration * 1000);
      setRemaining(duration);
    },
    [seconds],
  );

  useEffect(() => {
    if (startNow) {
      start();
    }
  }, [startNow, start]);

  const isCoolingDown = remaining > 0;

  return { remaining, isCoolingDown, start };
}
