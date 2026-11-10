import { TextQueryMethod } from "@zitadel/proto/zitadel/object/v2/object_pb";
import { beforeEach, describe, expect, test, vi } from "vitest";
import { startIdentityProviderFlow } from "./zitadel";

vi.mock("./service", () => ({
  createServiceForHost: vi.fn(),
}));

vi.mock("next-intl/server", () => ({
  getTranslations: vi.fn(() => (key: string) => key),
}));

vi.mock("./fingerprint", () => ({
  getUserAgent: vi.fn(),
}));

describe("zitadel queries are case insensitive", () => {
  const serviceConfig = { serviceUrl: "https://api.example.com" } as any;
  let mockCreateServiceForHost: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { createServiceForHost } = await import("./service");
    mockCreateServiceForHost = vi.mocked(createServiceForHost);
  });

  describe("getOrgsByDomain", () => {
    test("should look up the organization domain ignoring case", async () => {
      // Regression test for: https://github.com/zitadel/zitadel/issues/12025
      // Domain discovery for "jackson@Example.com" must find the org owning "example.com".
      const listOrganizations = vi.fn().mockResolvedValue({ result: [] });
      mockCreateServiceForHost.mockResolvedValue({ listOrganizations });

      const { getOrgsByDomain } = await import("./zitadel");
      await getOrgsByDomain({ serviceConfig, domain: "Example.com" });

      expect(listOrganizations).toHaveBeenCalledWith(
        {
          queries: [
            {
              query: {
                case: "domainQuery",
                value: { domain: "Example.com", method: TextQueryMethod.EQUALS_IGNORE_CASE },
              },
            },
          ],
        },
        {},
      );
    });
  });

  describe("listUsers", () => {
    test("should search login names ignoring case", async () => {
      const listUsersMock = vi.fn().mockResolvedValue({ result: [], details: {} });
      mockCreateServiceForHost.mockResolvedValue({ listUsers: listUsersMock });

      const { listUsers } = await import("./zitadel");
      await listUsers({ serviceConfig, loginName: "Jackson@Example.com" });

      const { queries } = listUsersMock.mock.calls[0][0];
      expect(queries[0].query.case).toBe("loginNameQuery");
      expect(queries[0].query.value.method).toBe(TextQueryMethod.EQUALS_IGNORE_CASE);
    });

    test("should search user names and emails ignoring case", async () => {
      const listUsersMock = vi.fn().mockResolvedValue({ result: [], details: {} });
      mockCreateServiceForHost.mockResolvedValue({ listUsers: listUsersMock });

      const { listUsers } = await import("./zitadel");
      await listUsers({ serviceConfig, userName: "Jackson", email: "Jackson@Example.com" });

      const { queries } = listUsersMock.mock.calls[0][0];
      const orQueries = queries[0].query.value.queries;
      const methods = orQueries.map((q: any) => q.query.value.method);
      expect(orQueries.map((q: any) => q.query.case)).toEqual(["userNameQuery", "emailQuery"]);
      expect(methods).toEqual([TextQueryMethod.EQUALS_IGNORE_CASE, TextQueryMethod.EQUALS_IGNORE_CASE]);
    });
  });
});

describe("startIdentityProviderFlow", () => {
  const serviceConfig = { baseUrl: "https://api.example.com" } as any;
  let startIdentityProviderIntent: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    vi.clearAllMocks();
    const { createServiceForHost } = await import("./service");
    startIdentityProviderIntent = vi.fn().mockResolvedValue({
      nextStep: { case: "authUrl", value: "https://idp.example.com/auth" },
    });
    vi.mocked(createServiceForHost).mockResolvedValue({ startIdentityProviderIntent } as any);
  });

  test("forwards the login hint to the API", async () => {
    const result = await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "user@example.com" },
    });

    expect(result).toEqual({ url: "https://idp.example.com/auth" });
    expect(startIdentityProviderIntent).toHaveBeenCalledWith({
      idpId: "idp123",
      content: {
        case: "urls",
        value: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "user@example.com" },
      },
    });
  });

  test("drops a login hint the API would reject instead of failing the IdP flow", async () => {
    const loginHint = "a".repeat(201);

    const result = await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint },
    });

    expect(result).toEqual({ url: "https://idp.example.com/auth" });
    expect(startIdentityProviderIntent).toHaveBeenCalledWith({
      idpId: "idp123",
      content: {
        case: "urls",
        value: { successUrl: "https://login/success", failureUrl: "https://login/failure" },
      },
    });
  });

  test("does not set an empty login hint", async () => {
    await startIdentityProviderFlow({
      serviceConfig,
      idpId: "idp123",
      urls: { successUrl: "https://login/success", failureUrl: "https://login/failure", loginHint: "" },
    });

    const value = startIdentityProviderIntent.mock.calls[0][0].content.value;
    expect(value).not.toHaveProperty("loginHint");
  });
});
