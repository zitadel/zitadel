import { cleanup, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, test, vi } from "vitest";
import { UserAvatar } from "./user-avatar";

vi.mock("@/lib/server/zitadel", () => ({
  getDomainSettings: vi.fn(),
  getOrganizationDomains: vi.fn(),
}));

vi.mock("@/lib/server/service-url", () => ({
  getServiceConfig: vi.fn(() => ({ serviceConfig: { baseUrl: "https://example.com" } })),
}));

vi.mock("next/headers", () => ({
  headers: vi.fn(() => new Headers()),
}));

vi.mock("@/components/avatar", () => ({ Avatar: vi.fn(() => null) }));

const loginName = "anna@example.com@studio.example.com";

describe("UserAvatar", () => {
  let getOrganizationDomains: any;

  beforeEach(async () => {
    vi.clearAllMocks();
    const zitadel = await import("@/lib/server/zitadel");
    getOrganizationDomains = vi.mocked(zitadel.getOrganizationDomains);
    getOrganizationDomains.mockResolvedValue(["studio.example.com"]);
    vi.mocked(zitadel.getDomainSettings).mockResolvedValue({ loginNameIncludesDomain: true } as never);
  });

  afterEach(() => {
    cleanup();
  });

  test("shows the login name when the suffix is not hidden", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, organizationId: "org-1", hideLoginNameSuffix: false }));

    expect(screen.getByText(loginName)).toBeDefined();
    expect(getOrganizationDomains).not.toHaveBeenCalled();
  });

  test("shows the login name without the organization domain when the branding hides the suffix", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, organizationId: "org-1", hideLoginNameSuffix: true }));

    expect(screen.getByText("anna@example.com")).toBeDefined();
    expect(screen.queryByText(loginName)).toBeNull();
  });

  test("shows the login name when the suffix is hidden but the organization is unknown", async () => {
    render(await UserAvatar({ loginName, showDropdown: false, hideLoginNameSuffix: true }));

    expect(screen.getByText(loginName)).toBeDefined();
    expect(getOrganizationDomains).not.toHaveBeenCalled();
  });

  test("shows the login name when the organization cannot be read", async () => {
    getOrganizationDomains.mockRejectedValue(new Error("permission denied"));

    render(await UserAvatar({ loginName, showDropdown: false, organizationId: "org-1", hideLoginNameSuffix: true }));

    expect(screen.getByText(loginName)).toBeDefined();
  });
});
