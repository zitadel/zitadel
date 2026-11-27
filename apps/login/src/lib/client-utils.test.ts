import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { handleServerActionResponse, isAllowedPostErrorRedirectUrl, isExternalUrl, isSafeRedirectUri } from "./client-utils";

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

describe("isExternalUrl", () => {
  test("should return false for relative paths", () => {
    expect(isExternalUrl("/dashboard")).toBe(false);
    expect(isExternalUrl("/ui/console")).toBe(false);
    expect(isExternalUrl("/")).toBe(false);
    expect(isExternalUrl("/password?loginName=user@example.com")).toBe(false);
    expect(isExternalUrl("/signedin?sessionId=123&organization=456")).toBe(false);
  });

  test("should return true for absolute HTTPS URLs", () => {
    expect(isExternalUrl("https://example.com/callback")).toBe(true);
    expect(isExternalUrl("https://my-zitadel.com/ui/console")).toBe(true);
  });

  test("should return true for absolute HTTP URLs", () => {
    expect(isExternalUrl("http://localhost:3000/callback")).toBe(true);
    expect(isExternalUrl("http://example.com")).toBe(true);
  });

  test("should return true for custom protocol schemes (native apps)", () => {
    expect(isExternalUrl("myapp://callback")).toBe(true);
    expect(isExternalUrl("com.example.app://oauth/callback")).toBe(true);
    expect(isExternalUrl("io.zitadel.app://auth")).toBe(true);
    expect(isExternalUrl("flutter-app://callback?code=abc&state=xyz")).toBe(true);
  });

  test("should return true for protocol-relative URLs", () => {
    expect(isExternalUrl("//evil.com")).toBe(true);
    expect(isExternalUrl("//example.com/path")).toBe(true);
  });
});

describe("isAllowedPostErrorRedirectUrl", () => {
  test("accepts the known back link targets", () => {
    expect(isAllowedPostErrorRedirectUrl("/loginname")).toBe(true);
    expect(isAllowedPostErrorRedirectUrl("/idp")).toBe(true);
  });

  test.each([
    "https://evil.example",
    "//evil.example",
    "/\\evil.example",
    "javascript:alert(1)",
    "data:text/html,x",
    "",
    undefined,
    "/loginname/x",
    "/loginname?x=1",
    " /loginname",
    "/LOGINNAME",
  ])("rejects %s", (url) => {
    expect(isAllowedPostErrorRedirectUrl(url)).toBe(false);
  });
});

describe("handleServerActionResponse with a SAML response", () => {
  const samlData = {
    url: "https://sp.example.com/acs",
    fields: { SAMLResponse: "PHNhbWxwOlJlc3BvbnNlPg==", RelayState: "state-1" },
  };

  afterEach(() => {
    vi.restoreAllMocks();
    document.body.innerHTML = "";
  });

  // No component is involved: the form that started the login may already be unmounted.
  test("posts the fields to the service provider with a form on the document", () => {
    const submit = vi.spyOn(HTMLFormElement.prototype, "submit").mockImplementation(() => {});
    const router = { push: vi.fn() };
    const setError = vi.fn();

    expect(handleServerActionResponse({ samlData }, router, setError)).toBe(true);

    expect(submit).toHaveBeenCalledTimes(1);
    const form = submit.mock.contexts[0] as HTMLFormElement;
    expect(form.isConnected).toBe(true);
    expect(form.method).toBe("post");
    expect(form.action).toBe("https://sp.example.com/acs");
    expect(Object.fromEntries(new FormData(form))).toEqual(samlData.fields);
    expect(router.push).not.toHaveBeenCalled();
    expect(setError).not.toHaveBeenCalled();
  });

  test("does not post to an unsafe URL", () => {
    const submit = vi.spyOn(HTMLFormElement.prototype, "submit").mockImplementation(() => {});
    vi.spyOn(console, "warn").mockImplementation(() => {});
    const setError = vi.fn();

    expect(
      handleServerActionResponse({ samlData: { ...samlData, url: "javascript:alert(1)" } }, { push: vi.fn() }, setError),
    ).toBe(true);

    expect(submit).not.toHaveBeenCalled();
    expect(setError).toHaveBeenCalled();
  });
});
