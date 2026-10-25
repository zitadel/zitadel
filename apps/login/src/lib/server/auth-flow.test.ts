import { beforeEach, describe, expect, test, vi } from "vitest";
import { completeAuthFlow } from "./auth-flow";

vi.mock("next/headers", () => ({
  headers: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  redirect: vi.fn((url: string) => {
    throw Object.assign(new Error("NEXT_REDIRECT"), { digest: `NEXT_REDIRECT;replace;${url};307;` });
  }),
}));

vi.mock("@/lib/cookies", () => ({
  getAllSessions: vi.fn(),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: vi.fn(),
}));

vi.mock("@/lib/zitadel", () => ({
  listSessions: vi.fn(),
}));

vi.mock("@/lib/oidc", () => ({
  loginWithOIDCAndSession: vi.fn(),
}));

vi.mock("@/lib/saml", () => ({
  loginWithSAMLAndSession: vi.fn(),
}));

describe("completeAuthFlow", () => {
  beforeEach(async () => {
    vi.clearAllMocks();
    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("@/lib/service-url");
    const { getAllSessions } = await import("@/lib/cookies");
    vi.mocked(headers).mockResolvedValue({} as any);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig: {} } as any);
    vi.mocked(getAllSessions).mockResolvedValue([]);
  });

  test("redirects to an absolute OIDC callback URL", async () => {
    const { loginWithOIDCAndSession } = await import("@/lib/oidc");
    const { redirect } = await import("next/navigation");
    vi.mocked(loginWithOIDCAndSession).mockResolvedValue({ redirect: "https://app.example.com/callback?code=abc" });

    await expect(completeAuthFlow({ sessionId: "session", requestId: "oidc_request" })).rejects.toMatchObject({
      digest: expect.stringMatching(/^NEXT_REDIRECT;/),
    });
    expect(redirect).toHaveBeenCalledWith("https://app.example.com/callback?code=abc");
  });

  test("redirects to an absolute SAML redirect binding URL", async () => {
    const { loginWithSAMLAndSession } = await import("@/lib/saml");
    const { redirect } = await import("next/navigation");
    vi.mocked(loginWithSAMLAndSession).mockResolvedValue({ redirect: "https://sp.example.com/acs?SAMLResponse=abc" });

    await expect(completeAuthFlow({ sessionId: "session", requestId: "saml_request" })).rejects.toBeDefined();
    expect(redirect).toHaveBeenCalledWith("https://sp.example.com/acs?SAMLResponse=abc");
  });

  test("returns a relative redirect for client-side navigation", async () => {
    const { loginWithOIDCAndSession } = await import("@/lib/oidc");
    const { redirect } = await import("next/navigation");
    vi.mocked(loginWithOIDCAndSession).mockResolvedValue({ redirect: "/signedin?loginName=user" });

    await expect(completeAuthFlow({ sessionId: "session", requestId: "oidc_request" })).resolves.toEqual({
      redirect: "/signedin?loginName=user",
    });
    expect(redirect).not.toHaveBeenCalled();
  });

  test("returns SAML POST data for the auto-submit form", async () => {
    const { loginWithSAMLAndSession } = await import("@/lib/saml");
    const { redirect } = await import("next/navigation");
    const samlData = { url: "https://sp.example.com/acs", fields: { SAMLResponse: "abc" } };
    vi.mocked(loginWithSAMLAndSession).mockResolvedValue({ samlData });

    await expect(completeAuthFlow({ sessionId: "session", requestId: "saml_request" })).resolves.toEqual({ samlData });
    expect(redirect).not.toHaveBeenCalled();
  });

  test("returns a callback URL with a semicolon for client-side navigation", async () => {
    const { loginWithOIDCAndSession } = await import("@/lib/oidc");
    const { redirect } = await import("next/navigation");
    const url = "https://app.example.com/callback;tenant=foo?code=abc&state=xyz";
    vi.mocked(loginWithOIDCAndSession).mockResolvedValue({ redirect: url });

    await expect(completeAuthFlow({ sessionId: "session", requestId: "oidc_request" })).resolves.toEqual({
      redirect: url,
    });
    expect(redirect).not.toHaveBeenCalled();
  });

  test("does not redirect to an unsafe URL", async () => {
    const { loginWithOIDCAndSession } = await import("@/lib/oidc");
    const { redirect } = await import("next/navigation");
    vi.mocked(loginWithOIDCAndSession).mockResolvedValue({ redirect: "javascript:alert(1)" });

    await completeAuthFlow({ sessionId: "session", requestId: "oidc_request" });
    expect(redirect).not.toHaveBeenCalled();
  });
});
