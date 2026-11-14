import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { useResendCooldown } from "./use-resend-cooldown";

describe("useResendCooldown", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  test("is idle until started", () => {
    const { result } = renderHook(() => useResendCooldown(30));
    expect(result.current.remaining).toBe(0);
    expect(result.current.isCoolingDown).toBe(false);
    expect(result.current.suffix).toBe("");
  });

  test("counts down and re-enables after the cooldown", () => {
    const { result } = renderHook(() => useResendCooldown(30));

    act(() => result.current.start());
    expect(result.current.remaining).toBe(30);
    expect(result.current.isCoolingDown).toBe(true);

    act(() => vi.advanceTimersByTime(10_000));
    expect(result.current.remaining).toBe(20);
    expect(result.current.suffix).toBe(" (20s)");

    act(() => vi.advanceTimersByTime(20_000));
    expect(result.current.remaining).toBe(0);
    expect(result.current.isCoolingDown).toBe(false);
  });
});
