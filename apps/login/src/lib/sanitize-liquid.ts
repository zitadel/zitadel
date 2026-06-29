import sanitize from "sanitize-html";

/**
 * Sanitize rendered Liquid output: allow structural HTML but block
 * scripts, styles, iframes, and all on* event handler attributes.
 *
 * This module is shared between server (`liquid.ts`) and client
 * (`liquid-slot-renderer.tsx`) components.
 *
 * Security notes:
 * - Only http/https/mailto/tel schemes are allowed in href/src.
 * - `data-*` attributes are restricted to `data-liquid-slot` only.
 * - `on*` event handlers are blocked by sanitize-html by default
 *   (they are not in any allowedAttributes list).
 * - `<script>`, `<style>`, `<iframe>`, `<form>`, `<input>`, `<svg>`,
 *   `<math>`, `<object>`, `<embed>` are not in allowedTags and are
 *   discarded.
 */
const SANITIZE_OPTIONS: sanitize.IOptions = {
  allowedTags: sanitize.defaults.allowedTags.concat([
    "header",
    "footer",
    "nav",
    "section",
    "article",
    "aside",
    "main",
    "figure",
    "figcaption",
    "img",
    "picture",
    "source",
  ]),
  allowedAttributes: {
    ...sanitize.defaults.allowedAttributes,
    "*": ["class", "id", "style", "data-liquid-slot"],
    a: ["href", "target", "rel", "class", "id", "style"],
    img: ["src", "alt", "width", "height", "class", "id", "style"],
  },
  // Only allow safe URL schemes — blocks javascript:, data:, vbscript: etc.
  allowedSchemes: ["http", "https", "mailto", "tel"],
  // Explicitly disallow all script-related tags
  disallowedTagsMode: "discard",
};

export function sanitizeLiquidOutput(html: string): string {
  return sanitize(html, SANITIZE_OPTIONS);
}
