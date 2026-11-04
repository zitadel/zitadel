import type { Metadata } from "next";

export type LinkPreviewBranding = {
  linkPreviewTitle?: string;
  linkPreviewDescription?: string;
  linkPreviewImageUrl?: string;
};

// loginMetadataFromBranding maps branding to document metadata.
// The browser tab title stays the translated screen title. Open Graph and Twitter
// tags use the configured link preview values, and fall back to that screen title
// when no custom title is set.
export function loginMetadataFromBranding(title: string, branding?: LinkPreviewBranding): Metadata {
  const previewTitle = branding?.linkPreviewTitle || title;
  const description = branding?.linkPreviewDescription || undefined;
  const image = branding?.linkPreviewImageUrl || undefined;

  return {
    title,
    description: description ?? null,
    openGraph: {
      title: previewTitle,
      ...(description ? { description } : {}),
      ...(image ? { images: [image] } : {}),
    },
    twitter: {
      card: image ? "summary_large_image" : "summary",
      title: previewTitle,
      ...(description ? { description } : {}),
      ...(image ? { images: [image] } : {}),
    },
  };
}
