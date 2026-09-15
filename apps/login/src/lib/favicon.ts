import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings } from "@/lib/zitadel";
import type { Metadata } from "next";
import { headers } from "next/headers";

/**
 * Bundled favicon assets, relative to the public directory.
 *
 * These are served below NEXT_PUBLIC_BASE_PATH, so every path has to be
 * prefixed with it. Next.js does not prefix URLs given in the metadata object.
 */
const BUNDLED_ICONS = {
  ico: "/favicon/favicon.ico",
  png32: "/favicon/favicon-32x32.png",
  png16: "/favicon/favicon-16x16.png",
  apple: "/favicon/apple-touch-icon.png",
} as const;

/**
 * Normalizes a base path into a prefix that can be concatenated with the icon
 * paths above.
 *
 * Next.js rejects a base path of "/" and one with a trailing slash, but this
 * helper takes the prefix as an argument, so it does not rely on that
 * validation: a stray slash would otherwise produce "//favicon/..." , which a
 * browser resolves as a protocol relative URL pointing at another host.
 */
function normalizeBasePath(basePath: string): string {
  const trimmed = basePath.trim().replace(/\/+$/, "");
  if (!trimmed) {
    return "";
  }
  return trimmed.startsWith("/") ? trimmed : `/${trimmed}`;
}

/**
 * Builds the favicon metadata.
 *
 * When the branding (label policy) defines an icon, it is used so the favicon
 * follows instance branding. That URL is already absolute — it points at the
 * assets API — and must not be prefixed. Otherwise the bundled ZITADEL
 * favicons are used, prefixed with the base path.
 */
export function buildIconsMetadata(brandingIconUrl?: string, basePath: string = ""): Metadata["icons"] {
  if (brandingIconUrl) {
    return {
      icon: brandingIconUrl,
      shortcut: brandingIconUrl,
      apple: brandingIconUrl,
    };
  }

  const prefix = normalizeBasePath(basePath);

  return {
    icon: [
      { url: `${prefix}${BUNDLED_ICONS.png32}`, sizes: "32x32", type: "image/png" },
      { url: `${prefix}${BUNDLED_ICONS.png16}`, sizes: "16x16", type: "image/png" },
    ],
    shortcut: `${prefix}${BUNDLED_ICONS.ico}`,
    apple: { url: `${prefix}${BUNDLED_ICONS.apple}`, sizes: "180x180" },
  };
}

/**
 * Resolves the favicon metadata for the login pages.
 *
 * Only the instance level icon can be resolved here: layouts do not receive
 * searchParams, so the organization is unknown where the metadata is
 * generated. A failing branding lookup falls back to the bundled favicons
 * instead of breaking the page.
 */
export async function resolveIconsMetadata(): Promise<Metadata["icons"]> {
  let brandingIconUrl: string | undefined;

  try {
    const _headers = await headers();
    const { serviceConfig } = getServiceConfig(_headers);
    const branding = await getBrandingSettings({ serviceConfig });
    brandingIconUrl = branding?.lightTheme?.iconUrl || branding?.darkTheme?.iconUrl || undefined;
  } catch (e) {
    console.error("Failed to load branding settings for the favicon", e);
  }

  return buildIconsMetadata(brandingIconUrl, process.env.NEXT_PUBLIC_BASE_PATH ?? "");
}
