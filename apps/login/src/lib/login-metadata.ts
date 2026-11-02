import { loginMetadataFromBranding, type LinkPreviewBranding } from "@/lib/link-preview-metadata";
import { getServiceConfig } from "@/lib/service-url";
import { getBrandingSettings, getDefaultOrg, type ServiceConfig } from "@/lib/zitadel";
import type { Metadata } from "next";
import { headers } from "next/headers";

export { loginMetadataFromBranding, type LinkPreviewBranding } from "@/lib/link-preview-metadata";

export type LoginMetadataOptions = {
  // Use the default organization when no organization is given,
  // matching pages that render the default organization's branding.
  fallbackToDefaultOrganization?: boolean;
};

export type LoginPageMetadataProps = {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
};

export async function getLoginPageMetadata(
  props: LoginPageMetadataProps,
  title: string,
  options?: LoginMetadataOptions,
): Promise<Metadata> {
  const searchParams = await props.searchParams;
  const organization = typeof searchParams.organization === "string" ? searchParams.organization : undefined;
  return getLoginMetadata(title, organization, options);
}

export async function getLoginMetadata(
  title: string,
  organization?: string,
  options?: LoginMetadataOptions,
): Promise<Metadata> {
  const requestHeaders = await headers();
  const { serviceConfig } = getServiceConfig(requestHeaders);
  return loginMetadataForService(serviceConfig, title, organization, options);
}

export async function loginMetadataForService(
  serviceConfig: ServiceConfig,
  title: string,
  organization?: string,
  options?: LoginMetadataOptions,
): Promise<Metadata> {
  let branding: LinkPreviewBranding | undefined;
  try {
    let brandingOrganization = organization;
    if (!brandingOrganization && options?.fallbackToDefaultOrganization) {
      brandingOrganization = (await getDefaultOrg({ serviceConfig }))?.id;
    }
    branding = await getBrandingSettings({ serviceConfig, organization: brandingOrganization });
  } catch (error) {
    console.error("Failed to load branding for link preview", error);
  }
  return loginMetadataFromBranding(title, branding);
}
