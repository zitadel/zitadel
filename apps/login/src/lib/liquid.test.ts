import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

// Mock "server-only" so tests can import the module
vi.mock("server-only", () => ({}));

import {
  renderLiquidTemplate,
  renderLiquidTemplateRaw,
  getEffectiveTemplate,
  sanitizeLiquidOutput,
  splitAtContent,
  DEFAULT_LIQUID_TEMPLATE,
  CONTENT_SENTINEL,
  THEME_SWITCHER_PLACEHOLDER,
  LANGUAGE_SWITCHER_PLACEHOLDER,
  LiquidTemplateVars,
} from "./liquid";

// Helper to build default vars with slot markers
function defaultVars(overrides?: Partial<LiquidTemplateVars>): LiquidTemplateVars {
  return {
    content: CONTENT_SENTINEL,
    theme_switcher: THEME_SWITCHER_PLACEHOLDER,
    language_switcher: LANGUAGE_SWITCHER_PLACEHOLDER,
    ...overrides,
  };
}

describe("renderLiquidTemplate", () => {
  it("renders a valid template with variables", async () => {
    const template = "<div>Hello {{ name }}!</div>";
    const result = await renderLiquidTemplate(template, {
      ...defaultVars(),
      name: "World",
    });
    expect(result).toBe("<div>Hello World!</div>");
  });

  it("renders content variable", async () => {
    const template = "<header>Header</header>{{ content }}<footer>Footer</footer>";
    const result = await renderLiquidTemplate(template, {
      ...defaultVars(),
      content: "<main>Page</main>",
    });
    expect(result).toContain("<header>Header</header>");
    expect(result).toContain("<main>Page</main>");
    expect(result).toContain("<footer>Footer</footer>");
  });

  it("strips arbitrary data-* attributes (only data-liquid-slot is allowed)", async () => {
    const template = '<div data-lang="{{ lang }}" data-theme="{{ theme }}">{{ content }}</div>';
    const result = await renderLiquidTemplate(template, {
      ...defaultVars(),
      content: "test",
      lang: "de",
      theme: "dark",
    });
    // data-lang and data-theme are stripped — only data-liquid-slot is whitelisted
    expect(result).not.toContain("data-lang");
    expect(result).not.toContain("data-theme");
    expect(result).toContain("test");
  });

  it("renders organization and instance_host variables", async () => {
    const template = "<span>{{ organization }} - {{ instance_host }}</span>";
    const result = await renderLiquidTemplate(template, {
      ...defaultVars(),
      organization: "org-123",
      instance_host: "my.zitadel.cloud",
    });
    expect(result).toContain("org-123");
    expect(result).toContain("my.zitadel.cloud");
  });

  it("strips <script> tags from output", async () => {
    const template = '<div>{{ content }}<script>alert("xss")</script></div>';
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "safe" });
    expect(result).not.toContain("<script>");
    expect(result).not.toContain("alert");
    expect(result).toContain("safe");
  });

  it("strips <script> tags with attributes", async () => {
    const template = '<div><script type="text/javascript" src="evil.js"></script>{{ content }}</div>';
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "ok" });
    expect(result).not.toContain("<script");
    expect(result).not.toContain("evil.js");
  });

  it("strips on* event handler attributes", async () => {
    const template = '<div onclick="alert(1)" onmouseover="hack()">{{ content }}</div>';
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "safe" });
    expect(result).not.toContain("onclick");
    expect(result).not.toContain("onmouseover");
  });

  it("strips <style> tags from output", async () => {
    const template = "<style>body { display: none; }</style><div>{{ content }}</div>";
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "visible" });
    expect(result).not.toContain("<style>");
    expect(result).not.toContain("display: none");
  });

  it("strips <iframe> tags from output", async () => {
    const template = '<iframe src="https://evil.com"></iframe><div>{{ content }}</div>';
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "safe" });
    expect(result).not.toContain("<iframe");
    expect(result).not.toContain("evil.com");
  });

  it("preserves allowed structural HTML", async () => {
    const template = `
      <header class="my-header">
        <nav id="main-nav"><a href="/home">Home</a></nav>
      </header>
      <main>{{ content }}</main>
      <footer class="my-footer">
        <a href="/privacy" target="_blank" rel="noopener">Privacy</a>
      </footer>
    `;
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "Page content" });
    expect(result).toContain("<header");
    expect(result).toContain("<nav");
    expect(result).toContain("<main>");
    expect(result).toContain("<footer");
    expect(result).toContain('href="/home"');
    expect(result).toContain('href="/privacy"');
    expect(result).toContain('class="my-header"');
    expect(result).toContain('class="my-footer"');
  });

  it("preserves img tags with allowed attributes", async () => {
    const template = '<img src="/logo.png" alt="Logo" width="100" height="50" class="logo" />';
    const result = await renderLiquidTemplate(template, { ...defaultVars() });
    expect(result).toContain("<img");
    expect(result).toContain('src="/logo.png"');
    expect(result).toContain('alt="Logo"');
  });

  it("preserves inline style attributes", async () => {
    const template = '<div style="color: red; padding: 10px;">{{ content }}</div>';
    const result = await renderLiquidTemplate(template, { ...defaultVars(), content: "styled" });
    expect(result).toContain("style=");
    expect(result).toMatch(/color:\s*red/);
  });

  it("returns undefined for malformed template", async () => {
    const template = "{% for item in %}broken{% endfor %}";
    const result = await renderLiquidTemplate(template, { ...defaultVars() });
    expect(result).toBeUndefined();
  });

  it("handles missing variables gracefully without throwing", async () => {
    const template = "<div>{{ nonexistent }}</div>";
    const result = await renderLiquidTemplate(template, { ...defaultVars() });
    expect(result).toBe("<div></div>");
  });
});

describe("sanitizeLiquidOutput", () => {
  it("strips script tags with newline bypass attempts", () => {
    const html = "<img\nsrc=x\nonerror=alert(1)>";
    const result = sanitizeLiquidOutput(html);
    expect(result).not.toContain("onerror");
  });

  it("preserves data-liquid-slot placeholder elements", () => {
    const html = '<div style="display:flex"><div data-liquid-slot="theme_switcher"></div></div>';
    const result = sanitizeLiquidOutput(html);
    expect(result).toContain('data-liquid-slot="theme_switcher"');
    expect(result).toContain("display:flex");
  });

  // ---------------------------------------------------------------
  // XSS attack vector tests
  // ---------------------------------------------------------------

  it("strips <script> tags", () => {
    const result = sanitizeLiquidOutput('<div><script>alert("xss")</script></div>');
    expect(result).not.toContain("<script");
    expect(result).not.toContain("alert");
  });

  it("strips <script> tags with attributes", () => {
    const result = sanitizeLiquidOutput(
      '<div><script type="text/javascript" src="evil.js"></script></div>',
    );
    expect(result).not.toContain("<script");
    expect(result).not.toContain("evil.js");
  });

  it("strips on* event handler attributes", () => {
    const result = sanitizeLiquidOutput('<div onclick="alert(1)" onmouseover="hack()">test</div>');
    expect(result).not.toContain("onclick");
    expect(result).not.toContain("onmouseover");
    expect(result).toContain("test");
  });

  it("strips javascript: URIs in href", () => {
    const result = sanitizeLiquidOutput('<a href="javascript:alert(1)">click</a>');
    expect(result).not.toContain("javascript:");
    expect(result).toContain("click");
  });

  it("strips case-varied javascript: URIs", () => {
    const result = sanitizeLiquidOutput('<a href="JAVASCRIPT:alert(1)">click</a>');
    expect(result).not.toContain("JAVASCRIPT:");
    expect(result).not.toContain("javascript:");
  });

  it("strips javascript: URIs with tab injection", () => {
    const result = sanitizeLiquidOutput('<a href="java\tscript:alert(1)">click</a>');
    expect(result).not.toContain("alert");
  });

  it("strips HTML entity encoded javascript: URIs", () => {
    const result = sanitizeLiquidOutput('<a href="&#106;avascript:alert(1)">click</a>');
    expect(result).not.toContain("alert");
  });

  it("strips data: URIs with scripts in href", () => {
    const result = sanitizeLiquidOutput(
      '<a href="data:text/html,<script>alert(1)</script>">click</a>',
    );
    expect(result).not.toContain("data:");
    expect(result).not.toContain("alert");
  });

  it("strips <svg> tags and SVG-based XSS", () => {
    const result = sanitizeLiquidOutput('<svg onload="alert(1)"><circle></circle></svg>');
    expect(result).not.toContain("<svg");
    expect(result).not.toContain("onload");
    expect(result).not.toContain("alert");
  });

  it("strips <iframe> tags", () => {
    const result = sanitizeLiquidOutput('<iframe src="https://evil.com"></iframe>');
    expect(result).not.toContain("<iframe");
    expect(result).not.toContain("evil.com");
  });

  it("strips <form> tags", () => {
    const result = sanitizeLiquidOutput(
      '<form action="https://evil.com"><input type="submit"></form>',
    );
    expect(result).not.toContain("<form");
    expect(result).not.toContain("<input");
    expect(result).not.toContain("evil.com");
  });

  it("strips <object> and <embed> tags", () => {
    const result = sanitizeLiquidOutput(
      '<object data="evil.swf"></object><embed src="evil.swf">',
    );
    expect(result).not.toContain("<object");
    expect(result).not.toContain("<embed");
    expect(result).not.toContain("evil.swf");
  });

  it("strips <meta> refresh redirect", () => {
    const result = sanitizeLiquidOutput(
      '<meta http-equiv="refresh" content="0;url=https://evil.com">',
    );
    expect(result).not.toContain("<meta");
    expect(result).not.toContain("evil.com");
  });

  it("strips <base> tag", () => {
    const result = sanitizeLiquidOutput('<base href="https://evil.com">');
    expect(result).not.toContain("<base");
    expect(result).not.toContain("evil.com");
  });

  it("strips <style> tags", () => {
    const result = sanitizeLiquidOutput("<style>body { display: none; }</style><div>ok</div>");
    expect(result).not.toContain("<style");
    expect(result).not.toContain("display: none");
  });

  it("strips ontoggle event handlers", () => {
    const result = sanitizeLiquidOutput(
      '<details open ontoggle="alert(1)"><summary>X</summary></details>',
    );
    expect(result).not.toContain("ontoggle");
    expect(result).not.toContain("alert");
  });

  it("neutralises math/mXSS mutation vectors", () => {
    const result = sanitizeLiquidOutput(
      '<math><mtext><table><mglyph><style><!--</style><img title="--></style><img src=1 onerror=alert(1)">',
    );
    expect(result).not.toContain("onerror");
    expect(result).not.toContain("alert");
  });

  it("strips arbitrary data-* attributes but preserves data-liquid-slot", () => {
    const html =
      '<div data-liquid-slot="theme_switcher" data-custom="evil" data-foo="bar">ok</div>';
    const result = sanitizeLiquidOutput(html);
    expect(result).toContain('data-liquid-slot="theme_switcher"');
    expect(result).not.toContain("data-custom");
    expect(result).not.toContain("data-foo");
    expect(result).toContain("ok");
  });

  it("allows safe href schemes (http, https, mailto, tel)", () => {
    const result = sanitizeLiquidOutput(
      '<a href="https://example.com">https</a>' +
        '<a href="http://example.com">http</a>' +
        '<a href="mailto:user@example.com">mail</a>' +
        '<a href="tel:+1234567890">phone</a>',
    );
    expect(result).toContain('href="https://example.com"');
    expect(result).toContain('href="http://example.com"');
    expect(result).toContain('href="mailto:user@example.com"');
    expect(result).toContain('href="tel:+1234567890"');
  });

  // ---------------------------------------------------------------
  // End-to-end pipeline: template → render → sanitize → split
  // ---------------------------------------------------------------

  it("blocks scripts injected through Liquid template variables", async () => {
    const template = "<div>{{ user_content }}</div>";
    const result = await renderLiquidTemplate(template, {
      ...defaultVars(),
      user_content: '<img src=x onerror="alert(document.cookie)">',
    });
    expect(result).not.toContain("onerror");
    expect(result).not.toContain("alert");
  });

  it("blocks script injection through splitAtContent pipeline", () => {
    const raw = `<script>alert("before")</script>${CONTENT_SENTINEL}<script>alert("after")</script>`;
    const { before, after } = splitAtContent(raw);
    expect(before).not.toContain("<script");
    expect(after).not.toContain("<script");
  });

  it("blocks javascript: URIs through splitAtContent pipeline", () => {
    const raw = `<a href="javascript:alert(1)">evil</a>${CONTENT_SENTINEL}<footer>ok</footer>`;
    const { before, after } = splitAtContent(raw);
    expect(before).not.toContain("javascript:");
  });
});

describe("getEffectiveTemplate", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  it("returns LIQUID_TEMPLATE env var when set", () => {
    process.env.LIQUID_TEMPLATE = "<div>{{ content }}</div>";
    const result = getEffectiveTemplate();
    expect(result).toBe("<div>{{ content }}</div>");
  });

  it("returns LIQUID_TEMPLATE env var even when branding template is provided", () => {
    process.env.LIQUID_TEMPLATE = "<div>ENV</div>";
    const result = getEffectiveTemplate("<div>BRANDING</div>");
    expect(result).toBe("<div>ENV</div>");
  });

  it("returns branding template when no env var is set", () => {
    delete process.env.LIQUID_TEMPLATE;
    const result = getEffectiveTemplate("<div>{{ content }}</div>");
    expect(result).toBe("<div>{{ content }}</div>");
  });

  it("returns undefined when neither env var nor branding template is set", () => {
    delete process.env.LIQUID_TEMPLATE;
    const result = getEffectiveTemplate();
    expect(result).toBeUndefined();
  });

  it("returns undefined when branding template is empty string", () => {
    delete process.env.LIQUID_TEMPLATE;
    const result = getEffectiveTemplate("");
    expect(result).toBeUndefined();
  });

  it("returns undefined when no overrides are set", () => {
    delete process.env.LIQUID_TEMPLATE;
    const result = getEffectiveTemplate();
    expect(result).toBeUndefined();
  });
});

describe("splitAtContent", () => {
  it("splits template at content sentinel", async () => {
    const raw = `<header>Header</header>${CONTENT_SENTINEL}<footer>Footer</footer>`;
    const { before, after } = splitAtContent(raw);

    expect(before).toContain("Header");
    expect(after).toContain("Footer");
  });

  it("returns entire output as before when no content sentinel", () => {
    const raw = "<div>No content slot</div>";
    const { before, after } = splitAtContent(raw);

    expect(before).toContain("No content slot");
    expect(after).toBe("");
  });

  it("sanitizes both parts", () => {
    const raw = `<script>evil</script>${CONTENT_SENTINEL}<script>bad</script>`;
    const { before, after } = splitAtContent(raw);

    expect(before).not.toContain("<script>");
    expect(after).not.toContain("<script>");
  });

  it("preserves switcher placeholder elements through sanitization", async () => {
    const raw = await renderLiquidTemplateRaw(
      '{{ content }}<div style="display:flex">{{ language_switcher }}{{ theme_switcher }}</div>',
      defaultVars(),
    );
    expect(raw).toBeDefined();

    const { after } = splitAtContent(raw!);

    // The placeholder elements should survive sanitization
    expect(after).toContain('data-liquid-slot="language_switcher"');
    expect(after).toContain('data-liquid-slot="theme_switcher"');
    // The flex container should be preserved around them
    expect(after).toContain("display:flex");
  });

  it("preserves HTML structure around switcher placeholders", async () => {
    const template = `{{ content }}<div style="display:flex;justify-content:space-between;max-width:440px;margin:0 auto;"><div><a href="/privacy">Privacy</a></div><div style="display:flex;gap:0.5rem;">{{ language_switcher }}{{ theme_switcher }}</div></div>`;
    const raw = await renderLiquidTemplateRaw(template, defaultVars());
    expect(raw).toBeDefined();

    const { after } = splitAtContent(raw!);

    // Structure should be intact
    expect(after).toContain("Privacy");
    expect(after).toContain('data-liquid-slot="language_switcher"');
    expect(after).toContain('data-liquid-slot="theme_switcher"');
    expect(after).toContain("justify-content:space-between");
  });
});

describe("DEFAULT_LIQUID_TEMPLATE", () => {
  it("contains all three slot variables", () => {
    expect(DEFAULT_LIQUID_TEMPLATE).toContain("{{ content }}");
    expect(DEFAULT_LIQUID_TEMPLATE).toContain("{{ theme_switcher }}");
    expect(DEFAULT_LIQUID_TEMPLATE).toContain("{{ language_switcher }}");
  });

  it("renders successfully with slot markers", async () => {
    const raw = await renderLiquidTemplateRaw(DEFAULT_LIQUID_TEMPLATE, defaultVars());
    expect(raw).toBeDefined();
    expect(raw).toContain(CONTENT_SENTINEL);
    expect(raw).toContain('data-liquid-slot="theme_switcher"');
    expect(raw).toContain('data-liquid-slot="language_switcher"');
  });
});
