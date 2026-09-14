import { renderToString } from "react-dom/server";
import { describe, expect, it, vi } from "vitest";
import { Avatar } from "./avatar";

// The server never has a resolved theme, while the pre-hydration client does.
// Avatar markup must not depend on it, or React emits a hydration mismatch
// when the system theme resolves to light.
let mockResolvedTheme: string | undefined;

vi.mock("next-themes", () => ({
  useTheme: () => ({ resolvedTheme: mockResolvedTheme }),
}));

describe("Avatar hydration safety", () => {
  it("produces identical markup regardless of the resolved theme", () => {
    mockResolvedTheme = undefined; // server render
    const serverHtml = renderToString(<Avatar name="Jane Doe" loginName="jane@example.com" />);

    for (const theme of ["light", "dark", "system"]) {
      mockResolvedTheme = theme;
      const clientHtml = renderToString(<Avatar name="Jane Doe" loginName="jane@example.com" />);
      expect(clientHtml).toBe(serverHtml);
    }
  });

  it("switches palettes through CSS variables and dark: variants", () => {
    const html = renderToString(<Avatar name="Jane Doe" loginName="jane@example.com" />);
    // Trailing `:` keeps `--avatar-bg` from matching `--avatar-bg-dark`.
    expect(html).toContain("--avatar-bg:");
    expect(html).toContain("--avatar-fg:");
    expect(html).toContain("--avatar-bg-dark:");
    expect(html).toContain("--avatar-fg-dark:");
    // The dark: classes are what actually swap the palette client-side.
    expect(html).toContain("bg-[var(--avatar-bg)]");
    expect(html).toContain("dark:bg-[var(--avatar-bg-dark)]");
    expect(html).toContain("text-[var(--avatar-fg)]");
    expect(html).toContain("dark:text-[var(--avatar-fg-dark)]");
  });
});
