import { beforeEach, describe, expect, test, vi } from "vitest";
import { getShownLoginName, stripLoginNameSuffix } from "./login-name";

vi.mock("@/lib/server/zitadel", () => ({
  getDomainSettings: vi.fn(),
  getOrganizationDomains: vi.fn(),
}));

const serviceConfig = { baseUrl: "https://example.com" };

describe("stripLoginNameSuffix", () => {
  test("removes the domain suffix appended to an e-mail-style username", () => {
    expect(stripLoginNameSuffix("anna@example.com@studio.example.com", ["studio.example.com"])).toBe("anna@example.com");
  });

  test("removes the suffix of any of the organization's domains", () => {
    expect(
      stripLoginNameSuffix("anna@example.com@studio.auth.example.com", ["studio.example.com", "studio.auth.example.com"]),
    ).toBe("anna@example.com");
  });

  test("ignores case", () => {
    expect(stripLoginNameSuffix("Anna@Example.com@Studio.Example.com", ["studio.example.COM"])).toBe("Anna@Example.com");
  });

  test("keeps a login name with a single @, which may be a complete username", () => {
    expect(stripLoginNameSuffix("theodora@studio.example.com", ["studio.example.com"])).toBe("theodora@studio.example.com");
  });

  test("keeps a login name of another domain", () => {
    expect(stripLoginNameSuffix("anna@example.com@other.example.com", ["studio.example.com"])).toBe(
      "anna@example.com@other.example.com",
    );
  });

  test("keeps a login name that only ends like the domain", () => {
    expect(stripLoginNameSuffix("anna@example.com@otherstudio.example.com", ["studio.example.com"])).toBe(
      "anna@example.com@otherstudio.example.com",
    );
  });

  test("keeps the login name without domains", () => {
    expect(stripLoginNameSuffix("anna@example.com@studio.example.com", [])).toBe("anna@example.com@studio.example.com");
  });
});

describe("getShownLoginName", () => {
  let getDomainSettings: any;
  let getOrganizationDomains: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/server/zitadel");
    getDomainSettings = vi.mocked(zitadel.getDomainSettings);
    getOrganizationDomains = vi.mocked(zitadel.getOrganizationDomains);
    getDomainSettings.mockResolvedValue({ loginNameIncludesDomain: true });
    getOrganizationDomains.mockResolvedValue(["studio.example.com", "studio.auth.example.com"]);
  });

  test("strips the organization's domain when the suffix is hidden and appended", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com");
    expect(getDomainSettings).toHaveBeenCalledWith({ serviceConfig, organization: "org-1" });
    expect(getOrganizationDomains).toHaveBeenCalledWith({ serviceConfig, organizationId: "org-1" });
  });

  test("strips a login name of a secondary domain of the organization", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.auth.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com");
  });

  test("keeps a typed e-mail username that ends with the organization's domain", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "theodora@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("theodora@studio.example.com");
  });

  test("keeps the login name when the organization does not append its domain", async () => {
    getDomainSettings.mockResolvedValue({ loginNameIncludesDomain: false });

    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
  });

  test("keeps the login name and reads nothing when the suffix is not hidden", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: false,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
    expect(getDomainSettings).not.toHaveBeenCalled();
    expect(getOrganizationDomains).not.toHaveBeenCalled();
  });

  test("keeps the login name without an organization", async () => {
    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
    expect(getOrganizationDomains).not.toHaveBeenCalled();
  });

  test("keeps the login name when the organization's settings cannot be read", async () => {
    getDomainSettings.mockRejectedValue(new Error("permission denied"));

    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
  });

  test("keeps the login name when the organization's domains cannot be read", async () => {
    getOrganizationDomains.mockRejectedValue(new Error("permission denied"));

    const shown = await getShownLoginName({
      serviceConfig,
      loginName: "anna@example.com@studio.example.com",
      organizationId: "org-1",
      hideLoginNameSuffix: true,
    });

    expect(shown).toBe("anna@example.com@studio.example.com");
  });
});
