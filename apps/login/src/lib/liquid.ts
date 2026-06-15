import "server-only";
import { Liquid } from "liquidjs";

export { sanitizeLiquidOutput } from "./sanitize-liquid";
import { sanitizeLiquidOutput } from "./sanitize-liquid";

const engine = new Liquid({
  strictVariables: false, // Don't throw on missing variables
  strictFilters: false,
});

// ---------------------------------------------------------------------------
// Slot markers
//
// CONTENT uses an HTML comment sentinel — it must be split server-side so
// the login card can be server-rendered (SSR).
//
// THEME_SWITCHER and LANGUAGE_SWITCHER use element-based placeholders
// (<div data-liquid-slot="...">) that survive sanitization AND preserve
// the surrounding HTML structure (flex containers, etc.).  They are
// mounted into the placeholder elements via React portals on the client.
// ---------------------------------------------------------------------------

export const CONTENT_SENTINEL = "<!--__SLOT_CONTENT__-->";
export const THEME_SWITCHER_PLACEHOLDER = '<div data-liquid-slot="theme_switcher"></div>';
export const LANGUAGE_SWITCHER_PLACEHOLDER = '<div data-liquid-slot="language_switcher"></div>';

// ---------------------------------------------------------------------------
// Default template — matches the current layout structure exactly.
// ---------------------------------------------------------------------------

export const DEFAULT_LIQUID_TEMPLATE = `<div>
  {{ content }}
  <div class="mx-auto flex max-w-[440px] flex-row items-center justify-end space-x-4 px-4 py-4 md:max-w-full md:px-8">
    {{ language_switcher }}
    {{ theme_switcher }}
  </div>
</div>`;

// ---------------------------------------------------------------------------
// Template variables
// ---------------------------------------------------------------------------

export interface LiquidTemplateVars {
  content: string;
  theme_switcher: string;
  language_switcher: string;
  lang?: string;
  theme?: string;
  organization?: string;
  instance_host?: string;
  [key: string]: unknown;
}

// ---------------------------------------------------------------------------
// Rendering
// ---------------------------------------------------------------------------

/**
 * Render a Liquid template string with the given variables.
 * Returns the raw (unsanitized) output for further processing.
 * Returns undefined if rendering fails.
 */
export async function renderLiquidTemplateRaw(
  template: string,
  vars: LiquidTemplateVars,
): Promise<string | undefined> {
  try {
    return await engine.parseAndRender(template, vars);
  } catch (err) {
    console.error("[liquid] Failed to render template:", err);
    return undefined;
  }
}

/**
 * Render a Liquid template string with the given variables.
 * Sanitizes the output to prevent XSS.
 * Returns undefined if rendering fails.
 */
export async function renderLiquidTemplate(
  template: string,
  vars: LiquidTemplateVars,
): Promise<string | undefined> {
  const raw = await renderLiquidTemplateRaw(template, vars);
  if (raw === undefined) return undefined;
  return sanitizeLiquidOutput(raw);
}

// ---------------------------------------------------------------------------
// Template resolution
// ---------------------------------------------------------------------------

/**
 * Resolve which template to use:
 * LIQUID_TEMPLATE env var → brandingSettings.template → undefined.
 * Returns undefined when no custom template is configured — the caller
 * should render the default layout directly in React.
 */
export function getEffectiveTemplate(brandingTemplate?: string): string | undefined {
  const envTemplate = process.env.LIQUID_TEMPLATE;
  if (envTemplate) return envTemplate;
  if (brandingTemplate) return brandingTemplate;
  return undefined;
}

// ---------------------------------------------------------------------------
// Content splitting
// ---------------------------------------------------------------------------

/**
 * Split rendered template at the content sentinel. Returns the before/after
 * HTML parts (sanitized, with switcher placeholder elements preserved).
 *
 * If the content sentinel is not found, the entire output is returned as "before".
 */
export function splitAtContent(raw: string): { before: string; after: string } {
  const idx = raw.indexOf(CONTENT_SENTINEL);

  if (idx === -1) {
    return { before: sanitizeLiquidOutput(raw), after: "" };
  }

  const before = sanitizeLiquidOutput(raw.slice(0, idx));
  const after = sanitizeLiquidOutput(raw.slice(idx + CONTENT_SENTINEL.length));

  return { before, after };
}
