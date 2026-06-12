import {
  getEffectiveTemplate,
  renderLiquidTemplateRaw,
  splitAtContent,
  LiquidTemplateVars,
  CONTENT_SENTINEL,
  THEME_SWITCHER_PLACEHOLDER,
  LANGUAGE_SWITCHER_PLACEHOLDER,
} from "@/lib/liquid";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { ReactNode } from "react";
import { DynamicThemeClient } from "./dynamic-theme-client";
import { LanguageSwitcherSlot } from "./language-switcher-slot";
import { LiquidSlotRenderer } from "./liquid-slot-renderer";
import ThemeSwitch from "./theme-switch";

/**
 * Default layout rendered directly in React when no custom Liquid template
 * is configured. The footer (language switcher + theme switcher) is passed
 * as a prop to DynamicThemeClient so it renders inside the same
 * width-constrained container as the card.
 */
function DefaultLayout({
  branding,
  children,
}: {
  branding?: BrandingSettings;
  children: ReactNode | ((isSideBySide: boolean) => ReactNode);
}) {
  const footerSlot = (
    <div className="flex flex-row items-center justify-end space-x-4 py-4">
      <LanguageSwitcherSlot />
      <ThemeSwitch />
    </div>
  );

  return (
    <DynamicThemeClient branding={branding} footer={footerSlot}>
      {children}
    </DynamicThemeClient>
  );
}

/**
 * Server component wrapper around DynamicThemeClient that adds LiquidJS template support
 * with multiple named slots: {{ content }}, {{ theme_switcher }}, {{ language_switcher }}.
 *
 * **Default** (no template): Renders the card with a footer containing the switchers,
 * all inside the same width-constrained container.
 *
 * **Custom template**: Splits the rendered output at {{ content }} (SSR-compatible).
 * The {{ theme_switcher }} and {{ language_switcher }} slots use element-based
 * placeholders that survive sanitization, then React portals mount the actual
 * components into them — preserving the template's HTML structure (flex, grid, etc.).
 */
export async function DynamicTheme({
  branding,
  children,
}: {
  children: ReactNode | ((isSideBySide: boolean) => ReactNode);
  branding?: BrandingSettings;
}) {
  // Resolve template: LIQUID_TEMPLATE env var → branding.template → undefined
  const template = getEffectiveTemplate(
    (branding as Record<string, unknown> | undefined)?.template as string | undefined,
  );

  // No custom template → render default React layout directly
  if (!template) {
    return <DefaultLayout branding={branding}>{children}</DefaultLayout>;
  }

  // Portal slots: React components mounted into placeholder elements
  const switcherSlots = {
    theme_switcher: <ThemeSwitch />,
    language_switcher: <LanguageSwitcherSlot />,
  };

  // Build Liquid variables:
  // - content uses a comment sentinel (split server-side for SSR)
  // - switcher slots use element placeholders (mounted via portals)
  const vars: LiquidTemplateVars = {
    content: CONTENT_SENTINEL,
    theme_switcher: THEME_SWITCHER_PLACEHOLDER,
    language_switcher: LANGUAGE_SWITCHER_PLACEHOLDER,
    lang: "", // TODO: pass locale when available at this level
    theme: "", // Resolved client-side by ThemeWrapper
    organization: "",
    instance_host: "",
  };

  const raw = await renderLiquidTemplateRaw(template, vars);

  if (!raw) {
    // Template render failed — fallback to default layout
    return <DefaultLayout branding={branding}>{children}</DefaultLayout>;
  }

  // Split at content sentinel, sanitize each half
  // (switcher placeholders survive sanitization as real elements)
  const { before, after } = splitAtContent(raw);

  console.log("[liquid] raw length:", raw.length);
  console.log("[liquid] before:", JSON.stringify(before));
  console.log("[liquid] after:", JSON.stringify(after));

  return (
    <>
      {before && (
        <LiquidSlotRenderer html={before} slots={switcherSlots} />
      )}
      <DynamicThemeClient branding={branding}>{children}</DynamicThemeClient>
      {after && (
        <LiquidSlotRenderer html={after} slots={switcherSlots} />
      )}
    </>
  );
}
