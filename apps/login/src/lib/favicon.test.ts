import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { buildIconsMetadata, resolveIconsMetadata } from "./favicon";

const BASE_PATH = "/ui/v2/login";

vi.mock("next/headers", () => ({
  headers: vi.fn(async () => new Headers()),
}));

vi.mock("@/lib/service-url", () => ({
  getServiceConfig: () => ({ serviceConfig: {} }),
}));

vi.mock("@/lib/zitadel", () => ({
  getBrandingSettings: vi.fn(),
}));

const { getBrandingSettings } = await import("@/lib/zitadel");

describe("buildIconsMetadata", () => {
  describe("bundled favicons", () => {
    it("should prefix every bundled icon with the base path", () => {
      const icons = buildIconsMetadata(undefined, BASE_PATH);

      expect(icons).toEqual({
        icon: [
          { url: `${BASE_PATH}/favicon/favicon-32x32.png`, sizes: "32x32", type: "image/png" },
          { url: `${BASE_PATH}/favicon/favicon-16x16.png`, sizes: "16x16", type: "image/png" },
        ],
        shortcut: `${BASE_PATH}/favicon/favicon.ico`,
        apple: { url: `${BASE_PATH}/favicon/apple-touch-icon.png`, sizes: "180x180" },
      });
    });

    it("should work without a base path", () => {
      const icons = buildIconsMetadata();

      expect(icons).toEqual({
        icon: [
          { url: "/favicon/favicon-32x32.png", sizes: "32x32", type: "image/png" },
          { url: "/favicon/favicon-16x16.png", sizes: "16x16", type: "image/png" },
        ],
        shortcut: "/favicon/favicon.ico",
        apple: { url: "/favicon/apple-touch-icon.png", sizes: "180x180" },
      });
    });

    it("should reference assets that are actually served below the base path", () => {
      // Regression guard: the paths must not be root relative, otherwise they
      // resolve outside the login app and return 404.
      const icons = buildIconsMetadata(undefined, BASE_PATH) as { shortcut: string };

      expect(icons.shortcut.startsWith(`${BASE_PATH}/`)).toBe(true);
    });

    // A leading "//" would be read as a protocol relative URL and resolve
    // against another host instead of this app.
    it.each([
      ["/", "/favicon/favicon.ico"],
      ["", "/favicon/favicon.ico"],
      ["   ", "/favicon/favicon.ico"],
      ["/ui/v2/login/", "/ui/v2/login/favicon/favicon.ico"],
      ["/ui/v2/login//", "/ui/v2/login/favicon/favicon.ico"],
      ["ui/v2/login", "/ui/v2/login/favicon/favicon.ico"],
    ])("should normalize base path %j into a single leading slash", (basePath, expected) => {
      const icons = buildIconsMetadata(undefined, basePath) as { shortcut: string };

      expect(icons.shortcut).toBe(expected);
      expect(icons.shortcut.startsWith("//")).toBe(false);
    });
  });

  describe("branding icon", () => {
    const brandingIcon = "https://example.com/assets/v1/123/policy/label/icon-456";

    it("should use the branding icon for all icon types", () => {
      const icons = buildIconsMetadata(brandingIcon, BASE_PATH);

      expect(icons).toEqual({
        icon: brandingIcon,
        shortcut: brandingIcon,
        apple: brandingIcon,
      });
    });

    it("should not prefix the branding icon with the base path", () => {
      const icons = buildIconsMetadata(brandingIcon, BASE_PATH) as { icon: string };

      expect(icons.icon).toBe(brandingIcon);
    });

    it.each(["", undefined])("should fall back to the bundled favicons when the branding icon is %j", (value) => {
      const icons = buildIconsMetadata(value, BASE_PATH) as { shortcut: string };

      expect(icons.shortcut).toBe(`${BASE_PATH}/favicon/favicon.ico`);
    });
  });
});

describe("resolveIconsMetadata", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv, NEXT_PUBLIC_BASE_PATH: BASE_PATH };
    vi.mocked(getBrandingSettings).mockReset();
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it("should use the branding icon of the light theme when set", async () => {
    const iconUrl = "https://id.example.com/assets/v1/1/policy/label/icon-2";
    vi.mocked(getBrandingSettings).mockResolvedValue({ lightTheme: { iconUrl } } as never);

    await expect(resolveIconsMetadata()).resolves.toEqual({
      icon: iconUrl,
      shortcut: iconUrl,
      apple: iconUrl,
    });
  });

  it("should fall back to the dark theme icon when the light theme has none", async () => {
    const iconUrl = "https://id.example.com/assets/v1/1/policy/label/icon-dark";
    vi.mocked(getBrandingSettings).mockResolvedValue({ lightTheme: { iconUrl: "" }, darkTheme: { iconUrl } } as never);

    await expect(resolveIconsMetadata()).resolves.toEqual({
      icon: iconUrl,
      shortcut: iconUrl,
      apple: iconUrl,
    });
  });

  it("should use the bundled favicons when branding defines no icon", async () => {
    vi.mocked(getBrandingSettings).mockResolvedValue({ lightTheme: { iconUrl: "" } } as never);

    const icons = (await resolveIconsMetadata()) as { shortcut: string };

    expect(icons.shortcut).toBe(`${BASE_PATH}/favicon/favicon.ico`);
  });

  // The favicon must never take the page down with it: a failing branding
  // lookup only degrades to the bundled icons.
  it("should use the bundled favicons when the branding lookup fails", async () => {
    const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});
    vi.mocked(getBrandingSettings).mockRejectedValue(new Error("fetch() returned undefined"));

    const icons = (await resolveIconsMetadata()) as { shortcut: string };

    expect(icons.shortcut).toBe(`${BASE_PATH}/favicon/favicon.ico`);
    expect(consoleError).toHaveBeenCalled();
    consoleError.mockRestore();
  });
});
