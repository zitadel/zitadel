import { beforeEach, describe, expect, it, vi } from "vitest";

import { loginMetadataFromBranding } from "./link-preview-metadata";
import { loginMetadataForService } from "./login-metadata";
import { getBrandingSettings, getDefaultOrg, type ServiceConfig } from "./zitadel";

vi.mock("next/headers", () => ({ headers: vi.fn() }));
vi.mock("./service-url", () => ({ getServiceConfig: vi.fn() }));
vi.mock("./zitadel", () => ({
  getBrandingSettings: vi.fn(),
  getDefaultOrg: vi.fn(),
}));

describe("loginMetadataFromBranding", () => {
  it("uses custom title, description, and image", () => {
    const metadata = loginMetadataFromBranding("Login with Zitadel", {
      linkPreviewTitle: "Acme",
      linkPreviewDescription: "Sign in to Acme",
      linkPreviewImageUrl: "https://auth.example.com/preview.png",
    });

    expect(metadata.title).toBe("Login with Zitadel");
    expect(metadata.description).toBe("Sign in to Acme");
    expect(metadata.openGraph).toMatchObject({
      title: "Acme",
      description: "Sign in to Acme",
      images: ["https://auth.example.com/preview.png"],
    });
    expect(metadata.twitter).toMatchObject({
      card: "summary_large_image",
      title: "Acme",
      description: "Sign in to Acme",
      images: ["https://auth.example.com/preview.png"],
    });
  });

  it("keeps the screen title and omits description and image when branding is empty", () => {
    const metadata = loginMetadataFromBranding("Welcome back!", {
      linkPreviewTitle: "",
      linkPreviewDescription: "",
      linkPreviewImageUrl: "",
    });

    expect(metadata.title).toBe("Welcome back!");
    expect(metadata.description).toBeUndefined();
    expect(metadata.openGraph).toEqual({ title: "Welcome back!" });
    expect(metadata.twitter).toEqual({ card: "summary", title: "Welcome back!" });
  });
});

describe("loginMetadataForService", () => {
  const serviceConfig = { baseUrl: "https://auth.example.com" } as unknown as ServiceConfig;
  const brandingByOrganization: Record<string, { linkPreviewTitle: string }> = {
    instance: { linkPreviewTitle: "Instance" },
    org1: { linkPreviewTitle: "Org One" },
    default: { linkPreviewTitle: "Default Org" },
  };

  beforeEach(() => {
    vi.mocked(getBrandingSettings).mockReset();
    vi.mocked(getDefaultOrg).mockReset();
    vi.mocked(getBrandingSettings).mockImplementation(
      async ({ organization }) => brandingByOrganization[organization ?? "instance"] as never,
    );
    vi.mocked(getDefaultOrg).mockResolvedValue({ id: "default" } as never);
  });

  it("uses the organization from the request", async () => {
    const metadata = await loginMetadataForService(serviceConfig, "Login", "org1");

    expect(getBrandingSettings).toHaveBeenCalledWith({ serviceConfig, organization: "org1" });
    expect(metadata.openGraph).toMatchObject({ title: "Org One" });
  });

  it("uses the instance branding without an organization", async () => {
    const metadata = await loginMetadataForService(serviceConfig, "Login");

    expect(getDefaultOrg).not.toHaveBeenCalled();
    expect(getBrandingSettings).toHaveBeenCalledWith({ serviceConfig, organization: undefined });
    expect(metadata.openGraph).toMatchObject({ title: "Instance" });
  });

  it("uses the default organization when the page falls back to it", async () => {
    const metadata = await loginMetadataForService(serviceConfig, "Login", undefined, {
      fallbackToDefaultOrganization: true,
    });

    expect(getBrandingSettings).toHaveBeenCalledWith({ serviceConfig, organization: "default" });
    expect(metadata.openGraph).toMatchObject({ title: "Default Org" });
  });

  it("prefers the request organization over the default organization", async () => {
    await loginMetadataForService(serviceConfig, "Login", "org1", { fallbackToDefaultOrganization: true });

    expect(getDefaultOrg).not.toHaveBeenCalled();
    expect(getBrandingSettings).toHaveBeenCalledWith({ serviceConfig, organization: "org1" });
  });

  it("ignores a repeated organization query parameter", async () => {
    const { getLoginPageMetadata } = await import("./login-metadata");
    const { headers } = await import("next/headers");
    const { getServiceConfig } = await import("./service-url");
    vi.mocked(headers).mockResolvedValue(new Headers() as never);
    vi.mocked(getServiceConfig).mockReturnValue({ serviceConfig } as never);

    await getLoginPageMetadata({ searchParams: Promise.resolve({ organization: ["org1", "org2"] }) }, "Login");

    expect(getBrandingSettings).toHaveBeenCalledWith({ serviceConfig, organization: undefined });
  });

  it("falls back to the screen title when branding cannot be loaded", async () => {
    vi.mocked(getBrandingSettings).mockRejectedValue(new Error("unavailable"));
    vi.spyOn(console, "error").mockImplementation(() => {});

    const metadata = await loginMetadataForService(serviceConfig, "Login");

    expect(metadata.openGraph).toEqual({ title: "Login" });
    expect(metadata.twitter).toEqual({ card: "summary", title: "Login" });
  });
});
