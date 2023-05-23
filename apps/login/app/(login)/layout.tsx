import { BrandingSettings } from "@zitadel/server";
import { ZitadelLogo } from "#/ui/ZitadelLogo";
import React from "react";
import { getBrandingSettings, server } from "#/lib/zitadel";
import { Logo } from "#/ui/Logo";

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const branding: BrandingSettings = await getBrandingSettings(server);
  let partial: Partial<BrandingSettings> | undefined;
  if (branding) {
    partial = {
      lightTheme: branding?.lightTheme,
      darkTheme: branding?.darkTheme,
    };
  }
  return (
    <div className="mx-auto flex flex-col items-center space-y-4">
      <div className="relative">
        <Logo
          lightSrc={branding.lightTheme?.logoUrl ?? "/zitadel-logo-dark.svg"}
          darkSrc={branding.darkTheme?.logoUrl ?? "/zitadel-logo-light.svg"}
          height={150}
          width={150}
        />
      </div>

      <div className="w-full">{children}</div>
      <div className="flex flex-row justify-between"></div>
    </div>
  );
}
