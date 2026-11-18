import { act, renderHook } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { useResendCooldown } from "./use-resend-cooldown";

describe("useResendCooldown", () => {
  beforeEach(() => vi.useFakeTimers());
  afterEach(() => vi.useRealTimers());

  test("is idle until started", () => {
    const { result } = renderHook(() => useResendCooldown({ seconds: 30 }));
    expect(result.current.remaining).toBe(0);
    expect(result.current.isCoolingDown).toBe(false);
  });

  test("counts down and re-enables after the cooldown", () => {
    const { result } = renderHook(() => useResendCooldown({ seconds: 30 }));

    act(() => result.current.start());
    expect(result.current.remaining).toBe(30);
    expect(result.current.isCoolingDown).toBe(true);

    act(() => vi.advanceTimersByTime(10_000));
    expect(result.current.remaining).toBe(20);

    act(() => vi.advanceTimersByTime(20_000));
    expect(result.current.remaining).toBe(0);
    expect(result.current.isCoolingDown).toBe(false);
  });

  test("counts down from the duration passed to start", () => {
    const { result } = renderHook(() => useResendCooldown({ seconds: 30 }));

    act(() => result.current.start(5));
    expect(result.current.remaining).toBe(5);

    act(() => vi.advanceTimersByTime(5_000));
    expect(result.current.isCoolingDown).toBe(false);
  });

  test("starts when startNow turns true", () => {
    const { result, rerender } = renderHook(({ startNow }) => useResendCooldown({ seconds: 30, startNow }), {
      initialProps: { startNow: false },
    });
    expect(result.current.isCoolingDown).toBe(false);

    rerender({ startNow: true });
    expect(result.current.remaining).toBe(30);
  });

  test("is already counting down in the first render when startNow is set", () => {
    const rendered: number[] = [];
    renderHook(() => {
      const cooldown = useResendCooldown({ seconds: 30, startNow: true });
      rendered.push(cooldown.remaining);
      return cooldown;
    });

    expect(rendered[0]).toBe(30);
  });
});
