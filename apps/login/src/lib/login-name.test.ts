import { beforeEach, describe, expect, test, vi } from "vitest";
import { getShownLoginName, stripLoginNameSuffix } from "./login-name";

vi.mock("@/lib/zitadel", () => ({
  getOrganizationPrimaryDomain: vi.fn(),
}));

const serviceConfig = { baseUrl: "https://example.com" };

describe("stripLoginNameSuffix", () => {
  test("removes the domain suffix", () => {
    expect(stripLoginNameSuffix("anna@studio.example.com", "studio.example.com")).toBe("anna");
  });

  test("keeps an e-mail username when the domain suffix is removed", () => {
    expect(stripLoginNameSuffix("anna@example.com@studio.example.com", "studio.example.com")).toBe("anna@example.com");
  });

  test("ignores case", () => {
    expect(stripLoginNameSuffix("Anna@Studio.Example.com", "studio.example.COM")).toBe("Anna");
  });

  test("keeps a login name of another domain", () => {
    expect(stripLoginNameSuffix("anna@example.com", "studio.example.com")).toBe("anna@example.com");
  });

  test("keeps a login name that only ends like the domain", () => {
    expect(stripLoginNameSuffix("anna@otherstudio.example.com", "studio.example.com")).toBe("anna@otherstudio.example.com");
  });

  test("keeps a login name that is nothing but the suffix", () => {
    expect(stripLoginNameSuffix("@studio.example.com", "studio.example.com")).toBe("@studio.example.com");
  });

  test("keeps the login name without a domain", () => {
    expect(stripLoginNameSuffix("anna@studio.example.com", undefined)).toBe("anna@studio.example.com");
  });
});

describe("getShownLoginName", () => {
  let getOrganizationPrimaryDomain: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/zitadel");
    getOrganizationPrimaryDomain = vi.mocked(zitadel.getOrganizationPrimaryDomain);
    getOrganizationPrimaryDomain.mockResolvedValue("studio.example.com");
  });

  test("strips the organization's domain when the suffix is hidden", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com");
    expect(getOrganizationPrimaryDomain).toHaveBeenCalledWith({ serviceConfig, organizationId: "org-1" });
  });

  test("keeps the login name and reads nothing when the suffix is not hidden", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: false,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
    expect(getOrganizationPrimaryDomain).not.toHaveBeenCalled();
  });

  test("keeps the login name without an organization", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
    expect(getOrganizationPrimaryDomain).not.toHaveBeenCalled();
  });

  test("keeps the login name when the organization cannot be read", async () => {
    getOrganizationPrimaryDomain.mockRejectedValue(new Error("permission denied"));

    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
  });
});
