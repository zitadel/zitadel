import { beforeEach, describe, expect, test, vi } from "vitest";
import {
  handleServerActionResponse,
  isCrossOrigin,
  isSafeRedirectUri,
  shouldUseHardNavigation,
} from "./client-utils";

describe("isSafeRedirectUri", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  test("should reject empty uri", async () => {
    expect(await isSafeRedirectUri("")).toBe(false);
    expect(await isSafeRedirectUri(undefined as unknown as string)).toBe(false);
  });

  test("should allow relative paths", async () => {
    expect(await isSafeRedirectUri("/dashboard")).toBe(true);
    expect(await isSafeRedirectUri("/ui/console")).toBe(true);
    expect(await isSafeRedirectUri("/")).toBe(true);
  });

  test("should reject protocol-relative paths", async () => {
    expect(await isSafeRedirectUri("//evil.com")).toBe(false);
    expect(await isSafeRedirectUri("\\\\evil.com")).toBe(false);
  });

  test("should allow absolute URLs matching the current host", async () => {
    expect(await isSafeRedirectUri("https://my-zitadel.com/console")).toBe(true);
  });

  test("should allow absolute URLs to external domains", async () => {
    expect(await isSafeRedirectUri("https://evil.com")).toBe(true);
    expect(await isSafeRedirectUri("https://settings.com/dashboard")).toBe(true);
  });

  test("should reject non-http(s) protocols", async () => {
    expect(await isSafeRedirectUri("javascript:alert(1)")).toBe(false);
    expect(await isSafeRedirectUri("data:text/html;base64,PHNjcmlwdD5hbGVydCgxKTwvc2NyaXB0Pg==")).toBe(false);
    expect(await isSafeRedirectUri("file:///etc/passwd")).toBe(false);
  });

  test("should gracefully reject invalid, unparsable URLs that aren't relative", async () => {
    expect(await isSafeRedirectUri("not-a-valid-url")).toBe(false);
    expect(await isSafeRedirectUri("http://%")).toBe(false);
  });
});

describe("isCrossOrigin", () => {
  test("returns true for absolute URL with different origin", () => {
    expect(isCrossOrigin("http://client-account.ludocare.local:3002/idp-callback", "http://login.ludocare.local:3021")).toBe(
      true,
    );
  });

  test("returns false for absolute URL with same origin", () => {
    expect(isCrossOrigin("http://login.ludocare.local:3021/signedin", "http://login.ludocare.local:3021")).toBe(false);
  });

  test("returns false for relative URL", () => {
    expect(isCrossOrigin("/signedin", "http://login.ludocare.local:3021")).toBe(false);
  });
});

describe("shouldUseHardNavigation", () => {
  test("returns true for cross-origin absolute URL", () => {
    expect(
      shouldUseHardNavigation(
        "http://client-account.ludocare.local:3002/idp-callback",
        "http://login.ludocare.local:3021",
      ),
    ).toBe(true);
  });

  test("returns true for proxied route", () => {
    expect(shouldUseHardNavigation("/oauth/v2/callback?x=1", "http://login.ludocare.local:3021")).toBe(true);
  });

  test("returns false for internal route", () => {
    expect(shouldUseHardNavigation("/signedin?x=1", "http://login.ludocare.local:3021")).toBe(false);
  });
});

describe("handleServerActionResponse", () => {
  test("uses browser navigation for cross-origin redirects", () => {
    const router = { push: vi.fn() };
    const setSamlData = vi.fn();
    const setError = vi.fn();
    const hardNavigate = vi.fn();

    const handled = handleServerActionResponse(
      { redirect: "http://client-account.ludocare.local:3002/idp-callback" },
      router,
      setSamlData,
      setError,
      hardNavigate,
    );

    expect(handled).toBe(true);
    expect(hardNavigate).toHaveBeenCalledWith("http://client-account.ludocare.local:3002/idp-callback");
    expect(router.push).not.toHaveBeenCalled();
  });

  test("uses next router for internal redirects", () => {
    const router = { push: vi.fn() };
    const setSamlData = vi.fn();
    const setError = vi.fn();
    const hardNavigate = vi.fn();

    const handled = handleServerActionResponse({ redirect: "/signedin" }, router, setSamlData, setError, hardNavigate);

    expect(handled).toBe(true);
    expect(router.push).toHaveBeenCalledWith("/signedin");
    expect(hardNavigate).not.toHaveBeenCalled();
  });
});
