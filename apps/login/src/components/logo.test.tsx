import { render } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { DEFAULT_THEME } from "../lib/theme";
import { Logo } from "./logo";

const LIGHT = "https://example.com/logo-light.png";
const DARK = "https://example.com/logo-dark.png";

// Queries are scoped to each render's container because the shared test setup
// does not register React Testing Library's automatic cleanup.
function renderLogo(props: Parameters<typeof Logo>[0]) {
  const { container } = render(<Logo {...props} />);
  return Array.from(container.querySelectorAll<HTMLImageElement>('img[alt="logo"]'));
}

describe("Logo", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
    delete process.env.NEXT_PUBLIC_THEME_LOGO_MAX_HEIGHT;
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("rendering", () => {
    it("should render both the light and dark variant when both are provided", () => {
      const imgs = renderLogo({ lightSrc: LIGHT, darkSrc: DARK });

      expect(imgs.map((img) => img.getAttribute("src"))).toEqual([DARK, LIGHT]);
    });

    it("should render nothing when no source is provided", () => {
      expect(renderLogo({})).toHaveLength(0);
    });

    it.each([
      ["light only", { lightSrc: LIGHT }],
      ["dark only", { darkSrc: DARK }],
    ])("should render a single image for %s", (_label, props) => {
      expect(renderLogo(props)).toHaveLength(1);
    });
  });

  describe("sizing", () => {
    it("should constrain both axes instead of pinning fixed dimensions", () => {
      const [img] = renderLogo({ lightSrc: LIGHT });

      // Without width/height attributes the asset keeps its intrinsic ratio.
      expect(img.getAttribute("width")).toBeNull();
      expect(img.getAttribute("height")).toBeNull();
      expect(img.style.maxHeight).toBe(`${DEFAULT_THEME.logoMaxHeight}px`);
      expect(img.style.maxWidth).toBe("100%");
      expect(img.className).toContain("object-contain");
    });

    it("should use the maxHeight prop when it is a positive number", () => {
      const [img] = renderLogo({ lightSrc: LIGHT, maxHeight: 140 });

      expect(img.style.maxHeight).toBe("140px");
    });

    it("should read the cap from the theme environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_LOGO_MAX_HEIGHT = "120";

      const [img] = renderLogo({ lightSrc: LIGHT });

      expect(img.style.maxHeight).toBe("120px");
    });

    it("should prefer the prop over the environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_LOGO_MAX_HEIGHT = "120";

      const [img] = renderLogo({ lightSrc: LIGHT, maxHeight: 64 });

      expect(img.style.maxHeight).toBe("64px");
    });

    // A non-finite or non-positive cap would emit something like "NaNpx", which
    // browsers discard, leaving the logo with no height constraint at all.
    it.each([
      ["NaN", Number.NaN],
      ["Infinity", Number.POSITIVE_INFINITY],
      ["zero", 0],
      ["a negative number", -40],
    ])("should fall back to the theme cap when maxHeight is %s", (_label, value) => {
      const [img] = renderLogo({ lightSrc: LIGHT, maxHeight: value });

      expect(img.style.maxHeight).toBe(`${DEFAULT_THEME.logoMaxHeight}px`);
    });
  });
});
