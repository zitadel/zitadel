import { Context, Hash, Liquid, Tag, TagToken, TopLevelToken } from "liquidjs";
import "server-only";
import { sanitizeLiquidOutput } from "./sanitize-liquid";

export { sanitizeLiquidOutput } from "./sanitize-liquid";

// ---------------------------------------------------------------------------
// Translation function type
//
// Matches the signature returned by next-intl's `getTranslations()` without
// a namespace: `t("namespace.key", { param: value })`.
// ---------------------------------------------------------------------------

/**
 * A translation function that resolves a dotted key and optional interpolation
 * values into a translated string. Matches next-intl's `getTranslations()`.
 */
export type TranslationFn = (key: string, values?: Record<string, unknown>) => string;

const engine = new Liquid({
  strictVariables: false, // Don't throw on missing variables
  strictFilters: false,
});

// ---------------------------------------------------------------------------
// Custom "t" tag — {% t "key" param: "value" %}
//
// Looks up a translation key via the `__t` function injected into the Liquid
// context at render time. The key is a string literal (quoted).
// Named parameters use Liquid's standard hash syntax (key: value) and are
// forwarded as ICU interpolation values.
//
// Examples:
//   {% t "loginname.title" %}
//   {% t "signedin.title" user: "John" %}
//   {% t "password.complexity.length" minLength: "8" %}
//   {% t "signedin.title" user: username %}  ← variable reference
// ---------------------------------------------------------------------------

engine.registerTag(
  "t",
  class TranslateTag extends Tag {
    private key: string;
    private hash: Hash;

    constructor(tagToken: TagToken, remainTokens: TopLevelToken[], liquid: Liquid) {
      super(tagToken, remainTokens, liquid);

      const args = tagToken.args.trim();

      // Extract the quoted key (first argument)
      const keyMatch = args.match(/^(["'])(.+?)\1/);
      if (!keyMatch) {
        throw new Error(`{% t %} tag requires a quoted translation key, e.g. {% t "loginname.title" %}. Got: ${args}`);
      }
      this.key = keyMatch[2];

      // Parse remaining args as named hash parameters using Liquid's
      // standard "key: value" syntax.
      const remaining = args.slice(keyMatch[0].length).trim();
      this.hash = new Hash(remaining);
    }

    *render(ctx: Context): Generator<unknown, string, unknown> {
      // Access the translation function from the top-level Liquid scope.
      // We use `ctx.environments` directly because `ctx.get(["__t"])` uses
      // LiquidJS's property-path traversal which doesn't find `__t`.
      const envs = ctx.environments as Record<string, unknown>;
      const t = envs["__t"] as TranslationFn | undefined;
      if (!t) {
        // No translation function available — return the key as-is
        return this.key;
      }

      // Resolve hash parameters (handles both literals and variable references)
      const params = (yield this.hash.render(ctx)) as Record<string, unknown>;

      try {
        const hasParams = Object.keys(params).length > 0;
        return t(this.key, hasParams ? params : undefined);
      } catch {
        // Translation key not found or format error — return the key
        return this.key;
      }
    }
  },
);

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
  instance_host?: string;
  /** Translation function injected into the context for the {% t %} tag. */
  __t?: TranslationFn;
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
export async function renderLiquidTemplateRaw(template: string, vars: LiquidTemplateVars): Promise<string | undefined> {
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
export async function renderLiquidTemplate(template: string, vars: LiquidTemplateVars): Promise<string | undefined> {
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
