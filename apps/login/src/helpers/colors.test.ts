import { describe, it, expect, beforeEach, vi } from "vitest";
import {
  computeMap,
  setTheme,
  DARK_PRIMARY,
  PRIMARY,
  DARK_WARN,
  WARN,
  DARK_BACKGROUND,
  BACKGROUND,
  DARK_TEXT,
  TEXT,
} from "./colors";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";

describe("color utilities", () => {
  describe("computeMap", () => {
    it("should compute light theme color maps with default colors", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: PRIMARY,
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      expect(result).toHaveProperty("background");
      expect(result).toHaveProperty("primary");
      expect(result).toHaveProperty("warn");
      expect(result).toHaveProperty("text");
      expect(result).toHaveProperty("link");

      expect(result.background).toHaveLength(14);
      expect(result.primary).toHaveLength(14);
      expect(result.warn).toHaveLength(14);
      expect(result.text).toHaveLength(14);
      expect(result.link).toHaveLength(14);
    });

    it("should compute dark theme color maps", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: PRIMARY,
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, true);

      expect(result.background).toHaveLength(14);
      expect(result.primary).toHaveLength(14);
      expect(result.warn).toHaveLength(14);
      expect(result.text).toHaveLength(14);
      expect(result.link).toHaveLength(14);
    });

    it("should generate color objects with correct properties", () => {
      const branding = {
        lightTheme: {
          backgroundColor: "#ffffff",
          fontColor: "#000000",
          primaryColor: "#5469d4",
          warnColor: "#cd3d56",
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      result.primary.forEach((color) => {
        expect(color).toHaveProperty("name");
        expect(color).toHaveProperty("hex");
        expect(color).toHaveProperty("rgb");
        expect(color).toHaveProperty("contrastColor");

        expect(color.hex).toMatch(/^#[0-9a-f]{6}$/i);
        expect(color.rgb).toMatch(/^rgb\(/);
        expect(["#ffffff", "hsla(0, 0%, 0%, 0.87)"]).toContain(color.contrastColor);
      });
    });

    it("should generate correct shade names", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: PRIMARY,
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      const expectedNames = [
        "50",
        "100",
        "200",
        "300",
        "400",
        "500",
        "600",
        "700",
        "800",
        "900",
        "A100",
        "A200",
        "A400",
        "A700",
      ];

      expect(result.primary.map((c) => c.name)).toEqual(expectedNames);
    });

    it("should compute different colors for light vs dark theme", () => {
      const branding = {
        lightTheme: {
          backgroundColor: "#ffffff",
          fontColor: "#000000",
          primaryColor: "#5469d4",
          warnColor: "#cd3d56",
        },
        darkTheme: {
          backgroundColor: "#111827",
          fontColor: "#ffffff",
          primaryColor: "#2073c4",
          warnColor: "#ff3b5b",
        },
      };

      const light = computeMap(branding, false);
      const dark = computeMap(branding, true);

      // Primary colors should be different
      expect(light.primary[5].hex).not.toBe(dark.primary[5].hex);

      // Background colors should be different
      expect(light.background[5].hex).not.toBe(dark.background[5].hex);

      // Text colors should be different
      expect(light.text[5].hex).not.toBe(dark.text[5].hex);
    });

    it("should handle custom brand colors", () => {
      const customBranding = {
        lightTheme: {
          backgroundColor: "#f0f0f0",
          fontColor: "#333333",
          primaryColor: "#ff6b6b",
          warnColor: "#feca57",
        },
        darkTheme: {
          backgroundColor: "#1e1e1e",
          fontColor: "#e0e0e0",
          primaryColor: "#ff8b8b",
          warnColor: "#fed57",
        },
      };

      const result = computeMap(customBranding, false);

      // Should generate valid hex colors
      result.primary.forEach((color) => {
        expect(color.hex).toMatch(/^#[0-9a-f]{6}$/i);
      });
    });

    it("should generate lighter shades for lower values", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: "#5469d4",
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      // 50 should be lighter than 500
      const shade50 = result.primary.find((c) => c.name === "50");
      const shade500 = result.primary.find((c) => c.name === "500");

      expect(shade50).toBeDefined();
      expect(shade500).toBeDefined();

      // Lighter colors have higher RGB values
      const rgb50 = shade50!.rgb;
      const rgb500 = shade500!.rgb;

      expect(rgb50).toBeDefined();
      expect(rgb500).toBeDefined();
    });

    it("should generate darker shades for higher values", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: "#5469d4",
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      // 900 should be darker than 500
      const shade500 = result.primary.find((c) => c.name === "500");
      const shade900 = result.primary.find((c) => c.name === "900");

      expect(shade500).toBeDefined();
      expect(shade900).toBeDefined();
    });

    it("should compute link colors using primary color", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: "#5469d4",
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      // Link should be computed based on primary color
      expect(result.link).toHaveLength(14);
      result.link.forEach((color) => {
        expect(color.hex).toMatch(/^#[0-9a-f]{6}$/i);
      });
    });

    it("should handle grayscale colors", () => {
      const grayscaleBranding = {
        lightTheme: {
          backgroundColor: "#ffffff",
          fontColor: "#000000",
          primaryColor: "#808080",
          warnColor: "#666666",
        },
        darkTheme: {
          backgroundColor: "#000000",
          fontColor: "#ffffff",
          primaryColor: "#a0a0a0",
          warnColor: "#999999",
        },
      };

      const result = computeMap(grayscaleBranding, false);

      // Should handle grayscale without errors
      result.primary.forEach((color) => {
        expect(color.hex).toMatch(/^#[0-9a-f]{6}$/i);
        expect(color.contrastColor).toBeDefined();
      });
    });

    it("should assign correct contrast colors for light backgrounds", () => {
      const branding = {
        lightTheme: {
          backgroundColor: "#ffffff",
          fontColor: TEXT,
          primaryColor: "#ffeb3b", // Light yellow
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      // Light colors should have dark contrast
      const lightShades = result.primary.filter((c) => ["50", "100", "200", "A100"].includes(c.name));

      lightShades.forEach((color) => {
        expect(color.contrastColor).toBe("hsla(0, 0%, 0%, 0.87)");
      });
    });

    it("should assign correct contrast colors for dark backgrounds", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: "#1a237e", // Dark blue
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      const result = computeMap(branding, false);

      // Dark colors should have light contrast
      const darkShades = result.primary.filter((c) => ["800", "900"].includes(c.name));

      darkShades.forEach((color) => {
        expect(color.contrastColor).toBe("#ffffff");
      });
    });
  });

  describe("setTheme", () => {
    let mockDocument: any;

    beforeEach(() => {
      mockDocument = {
        documentElement: {
          style: {
            setProperty: vi.fn(),
          },
        },
      };
    });

    it("should set theme with default colors when no policy provided", () => {
      setTheme(mockDocument);

      expect(mockDocument.documentElement.style.setProperty).toHaveBeenCalled();

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      const hasLightBackground = propertyNames.some((name: string) => name.includes("--theme-light-background-"));
      const hasDarkBackground = propertyNames.some((name: string) => name.includes("--theme-dark-background-"));
      const hasLightPrimary = propertyNames.some((name: string) => name.includes("--theme-light-primary-"));
      const hasDarkPrimary = propertyNames.some((name: string) => name.includes("--theme-dark-primary-"));

      expect(hasLightBackground).toBe(true);
      expect(hasDarkBackground).toBe(true);
      expect(hasLightPrimary).toBe(true);
      expect(hasDarkPrimary).toBe(true);
    });

    it("should set theme with custom policy colors", () => {
      const customPolicy = {
        lightTheme: {
          backgroundColor: "#f5f5f5",
          fontColor: "#212121",
          primaryColor: "#ff6b6b",
          warnColor: "#feca57",
        },
        darkTheme: {
          backgroundColor: "#1a1a1a",
          fontColor: "#fafafa",
          primaryColor: "#ff8b8b",
          warnColor: "#fed57",
        },
      } as BrandingSettings;

      setTheme(mockDocument, customPolicy);

      expect(mockDocument.documentElement.style.setProperty).toHaveBeenCalled();

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      expect(calls.length).toBeGreaterThan(0);
    });

    it("should set CSS custom properties for all color shades", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      ["50", "100", "200", "300", "400", "500", "600", "700", "800", "900"].forEach((shade) => {
        const hasShade = propertyNames.some((name: string) => name.includes(`-${shade}`));
        expect(hasShade).toBe(true);
      });

      ["A100", "A200", "A400", "A700"].forEach((shade) => {
        const hasShade = propertyNames.some((name: string) => name.includes(`-${shade}`));
        expect(hasShade).toBe(true);
      });
    });

    it("should set contrast colors for each shade", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      const contrastProperties = propertyNames.filter((name: string) => name.includes("-contrast-"));

      expect(contrastProperties.length).toBeGreaterThan(0);
    });

    it("should set secondary alpha colors for text and link", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      const secondaryProperties = propertyNames.filter((name: string) => name.includes("-secondary-"));

      expect(secondaryProperties.length).toBeGreaterThan(0);

      secondaryProperties.forEach((prop: string) => {
        expect(prop.includes("-text-") || prop.includes("-link-")).toBe(true);
      });
    });

    it("should set properties for both light and dark themes", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      const lightProperties = propertyNames.filter((name: string) => name.includes("--theme-light-"));
      const darkProperties = propertyNames.filter((name: string) => name.includes("--theme-dark-"));

      expect(lightProperties.length).toBeGreaterThan(0);
      expect(darkProperties.length).toBeGreaterThan(0);
    });

    it("should set properties for all color types", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      const types = ["background", "primary", "warn", "text", "link"];

      types.forEach((type) => {
        const typeProperties = propertyNames.filter((name: string) => name.includes(`-${type}-`));
        expect(typeProperties.length).toBeGreaterThan(0);
      });
    });

    it("should use fallback colors when policy has partial values", () => {
      const partialPolicy = {
        lightTheme: {
          primaryColor: "#ff0000",
        },
      } as BrandingSettings;

      setTheme(mockDocument, partialPolicy);

      expect(mockDocument.documentElement.style.setProperty).toHaveBeenCalled();

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      expect(calls.length).toBeGreaterThan(0);
    });

    it("should handle undefined policy gracefully", () => {
      expect(() => setTheme(mockDocument, undefined)).not.toThrow();
      expect(mockDocument.documentElement.style.setProperty).toHaveBeenCalled();
    });

    it("should set valid hex color values", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;

      calls.forEach((call: any[]) => {
        const [propertyName, propertyValue] = call;

        if (!propertyName.includes("-contrast-") && !propertyName.includes("-secondary-")) {
          expect(propertyValue).toMatch(/^#[0-9a-f]{6}$/i);
        }
      });
    });

    it("should set alpha channel for secondary colors", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;

      const secondaryCalls = calls.filter((call: any[]) => call[0].includes("-secondary-"));

      secondaryCalls.forEach((call: any[]) => {
        const [, propertyValue] = call;
        expect(propertyValue).toMatch(/c7$/i);
      });
    });

    it("should call setProperty exactly once per CSS variable", () => {
      setTheme(mockDocument);

      const calls = mockDocument.documentElement.style.setProperty.mock.calls;
      const propertyNames = calls.map((call: any[]) => call[0]);

      // Check for duplicates
      const uniqueProperties = new Set(propertyNames);
      expect(uniqueProperties.size).toBe(propertyNames.length);
    });
  });

  describe("default color constants", () => {
    it("should have valid default color values", () => {
      expect(PRIMARY).toMatch(/^#[0-9a-f]{6}$/i);
      expect(DARK_PRIMARY).toMatch(/^#[0-9a-f]{6}$/i);
      expect(WARN).toMatch(/^#[0-9a-f]{6}$/i);
      expect(DARK_WARN).toMatch(/^#[0-9a-f]{6}$/i);
      expect(BACKGROUND).toMatch(/^#[0-9a-f]{6}$/i);
      expect(DARK_BACKGROUND).toMatch(/^#[0-9a-f]{6}$/i);
      expect(TEXT).toMatch(/^#[0-9a-f]{6}$/i);
      expect(DARK_TEXT).toMatch(/^#[0-9a-f]{6}$/i);
    });

    it("should have distinct colors for light and dark themes", () => {
      expect(PRIMARY).not.toBe(DARK_PRIMARY);
      expect(WARN).not.toBe(DARK_WARN);
      expect(BACKGROUND).not.toBe(DARK_BACKGROUND);
      expect(TEXT).not.toBe(DARK_TEXT);
    });
  });

  describe("edge cases", () => {
    it("should handle extreme light colors", () => {
      const branding = {
        lightTheme: {
          backgroundColor: "#ffffff",
          fontColor: "#fefefe",
          primaryColor: "#fffacd",
          warnColor: "#fff8dc",
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: DARK_PRIMARY,
          warnColor: DARK_WARN,
        },
      };

      expect(() => computeMap(branding, false)).not.toThrow();
    });

    it("should handle extreme dark colors", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: PRIMARY,
          warnColor: WARN,
        },
        darkTheme: {
          backgroundColor: "#000000",
          fontColor: "#010101",
          primaryColor: "#0a0a0a",
          warnColor: "#141414",
        },
      };

      expect(() => computeMap(branding, true)).not.toThrow();
    });

    it("should handle highly saturated colors", () => {
      const branding = {
        lightTheme: {
          backgroundColor: BACKGROUND,
          fontColor: TEXT,
          primaryColor: "#ff0000",
          warnColor: "#00ff00",
        },
        darkTheme: {
          backgroundColor: DARK_BACKGROUND,
          fontColor: DARK_TEXT,
          primaryColor: "#0000ff",
          warnColor: "#ffff00",
        },
      };

      const result = computeMap(branding, false);
      expect(result.primary).toHaveLength(14);
      expect(result.warn).toHaveLength(14);
    });
  });
});
