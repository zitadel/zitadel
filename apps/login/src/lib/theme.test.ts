import { describe, it, expect, beforeEach, afterEach } from "vitest";
import {
  ThemeRoundness,
  ThemeLayout,
  ThemeAppearance,
  ThemeSpacing,
  ComponentRoundnessConfig,
  DEFAULT_COMPONENT_ROUNDNESS,
  DEFAULT_THEME,
  getThemeConfig,
  ROUNDNESS_CLASSES,
  getComponentRoundness,
  SPACING_STYLES,
  APPEARANCE_STYLES,
} from "./theme";

describe("Theme Configuration", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    // Reset environment variables before each test
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    // Restore original environment variables
    process.env = originalEnv;
  });

  describe("DEFAULT_COMPONENT_ROUNDNESS", () => {
    it("should have correct default values for all components", () => {
      expect(DEFAULT_COMPONENT_ROUNDNESS).toEqual({
        card: "mid",
        button: "mid",
        input: "mid",
        image: "mid",
        avatar: "full",
        avatarContainer: "full",
        themeSwitch: "full",
      });
    });

    it("should have all required component types", () => {
      const requiredComponents: (keyof ComponentRoundnessConfig)[] = [
        "card",
        "button",
        "input",
        "image",
        "avatar",
        "avatarContainer",
        "themeSwitch",
      ];

      requiredComponents.forEach((component) => {
        expect(DEFAULT_COMPONENT_ROUNDNESS).toHaveProperty(component);
      });
    });
  });

  describe("DEFAULT_THEME", () => {
    it("should have all required properties", () => {
      expect(DEFAULT_THEME).toEqual({
        roundness: "mid",
        componentRoundness: DEFAULT_COMPONENT_ROUNDNESS,
        layout: "top-to-bottom",
        appearance: "flat",
        spacing: "regular",
      });
    });

    it("should have valid default values", () => {
      expect(DEFAULT_THEME.roundness).toBe("mid");
      expect(DEFAULT_THEME.layout).toBe("top-to-bottom");
      expect(DEFAULT_THEME.appearance).toBe("flat");
      expect(DEFAULT_THEME.spacing).toBe("regular");
    });
  });

  describe("getThemeConfig", () => {
    it("should return default theme when no environment variables are set", () => {
      const config = getThemeConfig();

      expect(config.roundness).toBe(DEFAULT_THEME.roundness);
      expect(config.layout).toBe(DEFAULT_THEME.layout);
      expect(config.appearance).toBe(DEFAULT_THEME.appearance);
      expect(config.spacing).toBe(DEFAULT_THEME.spacing);
      expect(config.componentRoundness).toEqual(DEFAULT_COMPONENT_ROUNDNESS);
    });

    it("should use global roundness from environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";

      const config = getThemeConfig();

      expect(config.roundness).toBe("full");
    });

    it("should apply global roundness to all components when env var is set", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "edgy";

      const config = getThemeConfig();

      expect(config.componentRoundness?.card).toBe("edgy");
      expect(config.componentRoundness?.button).toBe("edgy");
      expect(config.componentRoundness?.input).toBe("edgy");
      expect(config.componentRoundness?.image).toBe("edgy");
      expect(config.componentRoundness?.avatar).toBe("edgy");
      expect(config.componentRoundness?.avatarContainer).toBe("edgy");
      expect(config.componentRoundness?.themeSwitch).toBe("edgy");
    });

    it("should use component-specific defaults when global roundness is not set", () => {
      const config = getThemeConfig();

      expect(config.componentRoundness).toEqual(DEFAULT_COMPONENT_ROUNDNESS);
      expect(config.componentRoundness?.avatar).toBe("full");
      expect(config.componentRoundness?.avatarContainer).toBe("full");
      expect(config.componentRoundness?.card).toBe("mid");
    });

    it("should use layout from environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_LAYOUT = "side-by-side";

      const config = getThemeConfig();

      expect(config.layout).toBe("side-by-side");
    });

    it("should use appearance from environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "material";

      const config = getThemeConfig();

      expect(config.appearance).toBe("material");
    });

    it("should use spacing from environment variable", () => {
      process.env.NEXT_PUBLIC_THEME_SPACING = "compact";

      const config = getThemeConfig();

      expect(config.spacing).toBe("compact");
    });

    it("should use background image from environment variable", () => {
      const backgroundUrl = "https://example.com/image.jpg";
      process.env.NEXT_PUBLIC_THEME_BACKGROUND_IMAGE = backgroundUrl;

      const config = getThemeConfig();

      expect(config.backgroundImage).toBe(backgroundUrl);
    });

    it("should have undefined background image when not set", () => {
      const config = getThemeConfig();

      expect(config.backgroundImage).toBeUndefined();
    });

    it("should combine multiple environment variables correctly", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      process.env.NEXT_PUBLIC_THEME_LAYOUT = "side-by-side";
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "glass";
      process.env.NEXT_PUBLIC_THEME_SPACING = "compact";
      process.env.NEXT_PUBLIC_THEME_BACKGROUND_IMAGE = "https://example.com/bg.png";

      const config = getThemeConfig();

      expect(config.roundness).toBe("full");
      expect(config.layout).toBe("side-by-side");
      expect(config.appearance).toBe("glass");
      expect(config.spacing).toBe("compact");
      expect(config.backgroundImage).toBe("https://example.com/bg.png");
    });
  });

  describe("ROUNDNESS_CLASSES", () => {
    it("should have classes for all roundness levels", () => {
      expect(ROUNDNESS_CLASSES).toHaveProperty("edgy");
      expect(ROUNDNESS_CLASSES).toHaveProperty("mid");
      expect(ROUNDNESS_CLASSES).toHaveProperty("full");
    });

    it("should have classes for all component types in each roundness level", () => {
      const components: (keyof ComponentRoundnessConfig)[] = [
        "card",
        "button",
        "input",
        "image",
        "avatar",
        "avatarContainer",
        "themeSwitch",
      ];

      const roundnessLevels: ThemeRoundness[] = ["edgy", "mid", "full"];

      roundnessLevels.forEach((level) => {
        components.forEach((component) => {
          expect(ROUNDNESS_CLASSES[level]).toHaveProperty(component);
          expect(typeof ROUNDNESS_CLASSES[level][component]).toBe("string");
          expect(ROUNDNESS_CLASSES[level][component].length).toBeGreaterThan(0);
        });
      });
    });

    it("should have distinct classes for different roundness levels", () => {
      const components: (keyof ComponentRoundnessConfig)[] = ["card", "button", "input"];

      components.forEach((component) => {
        // Each roundness level should have different classes for the same component
        expect(ROUNDNESS_CLASSES.edgy[component]).not.toBe(ROUNDNESS_CLASSES.mid[component]);
        expect(ROUNDNESS_CLASSES.mid[component]).not.toBe(ROUNDNESS_CLASSES.full[component]);
        expect(ROUNDNESS_CLASSES.edgy[component]).not.toBe(ROUNDNESS_CLASSES.full[component]);
      });
    });
  });

  describe("getComponentRoundness", () => {
    it("should return a valid CSS class string for any component", () => {
      const cardClass = getComponentRoundness("card");
      expect(typeof cardClass).toBe("string");
      expect(cardClass.length).toBeGreaterThan(0);
    });

    it("should return different classes for different components with default config", () => {
      const cardClass = getComponentRoundness("card");
      const avatarClass = getComponentRoundness("avatar");

      // Avatar defaults to full roundness, card to mid - they should be different
      expect(cardClass).not.toBe(avatarClass);
    });

    it("should change output when global roundness environment variable changes", () => {
      // Get default
      const defaultCardClass = getComponentRoundness("card");

      // Set to different roundness
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "edgy";
      const edgyCardClass = getComponentRoundness("card");

      // Should be different
      expect(edgyCardClass).not.toBe(defaultCardClass);

      // Try another roundness level
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const fullCardClass = getComponentRoundness("card");

      expect(fullCardClass).not.toBe(defaultCardClass);
      expect(fullCardClass).not.toBe(edgyCardClass);
    });

    it("should return classes for all component types", () => {
      const components: (keyof ComponentRoundnessConfig)[] = [
        "card",
        "button",
        "input",
        "image",
        "avatar",
        "avatarContainer",
        "themeSwitch",
      ];

      components.forEach((component) => {
        const cssClass = getComponentRoundness(component);
        expect(typeof cssClass).toBe("string");
        expect(cssClass.length).toBeGreaterThan(0);
      });
    });

    it("should apply global roundness to all components when set", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";

      const cardClass = getComponentRoundness("card");
      const buttonClass = getComponentRoundness("button");
      const avatarClass = getComponentRoundness("avatar");

      // All should get classes from the "full" roundness level
      expect(cardClass).toBe(ROUNDNESS_CLASSES.full.card);
      expect(buttonClass).toBe(ROUNDNESS_CLASSES.full.button);
      expect(avatarClass).toBe(ROUNDNESS_CLASSES.full.avatar);
    });

    it("should respect component-specific defaults when no global roundness is set", () => {
      delete process.env.NEXT_PUBLIC_THEME_ROUNDNESS;

      const avatarClass = getComponentRoundness("avatar");
      const cardClass = getComponentRoundness("card");

      // Avatar should use its default (full), card should use its default (mid)
      expect(avatarClass).toBe(ROUNDNESS_CLASSES.full.avatar);
      expect(cardClass).toBe(ROUNDNESS_CLASSES.mid.card);
    });
  });

  describe("SPACING_STYLES", () => {
    it("should have regular and compact spacing options", () => {
      expect(SPACING_STYLES).toHaveProperty("regular");
      expect(SPACING_STYLES).toHaveProperty("compact");
    });

    it("should have spacing and padding properties for each option", () => {
      expect(SPACING_STYLES.regular).toHaveProperty("spacing");
      expect(SPACING_STYLES.regular).toHaveProperty("padding");
      expect(SPACING_STYLES.compact).toHaveProperty("spacing");
      expect(SPACING_STYLES.compact).toHaveProperty("padding");
    });

    it("should have non-empty string values for all properties", () => {
      expect(typeof SPACING_STYLES.regular.spacing).toBe("string");
      expect(SPACING_STYLES.regular.spacing.length).toBeGreaterThan(0);
      expect(typeof SPACING_STYLES.regular.padding).toBe("string");
      expect(SPACING_STYLES.regular.padding.length).toBeGreaterThan(0);

      expect(typeof SPACING_STYLES.compact.spacing).toBe("string");
      expect(SPACING_STYLES.compact.spacing.length).toBeGreaterThan(0);
      expect(typeof SPACING_STYLES.compact.padding).toBe("string");
      expect(SPACING_STYLES.compact.padding.length).toBeGreaterThan(0);
    });

    it("should have different values between regular and compact", () => {
      expect(SPACING_STYLES.regular.spacing).not.toBe(SPACING_STYLES.compact.spacing);
      expect(SPACING_STYLES.regular.padding).not.toBe(SPACING_STYLES.compact.padding);
    });
  });

  describe("APPEARANCE_STYLES", () => {
    it("should have flat, material, and glass appearance options", () => {
      expect(APPEARANCE_STYLES).toHaveProperty("flat");
      expect(APPEARANCE_STYLES).toHaveProperty("material");
      expect(APPEARANCE_STYLES).toHaveProperty("glass");
    });

    it("should have required properties for each appearance", () => {
      const requiredProperties = ["card", "button", "idp-button", "typography", "background"];

      Object.values(APPEARANCE_STYLES).forEach((style) => {
        requiredProperties.forEach((prop) => {
          expect(style).toHaveProperty(prop);
          // @ts-ignore - dynamic property access
          expect(typeof style[prop]).toBe("string");
          // @ts-ignore - dynamic property access
          expect(style[prop].length).toBeGreaterThan(0);
        });
      });
    });

    it("should have different styles for different appearances", () => {
      // Each appearance should have distinct card styles
      expect(APPEARANCE_STYLES.flat.card).not.toBe(APPEARANCE_STYLES.material.card);
      expect(APPEARANCE_STYLES.material.card).not.toBe(APPEARANCE_STYLES.glass.card);
      expect(APPEARANCE_STYLES.flat.card).not.toBe(APPEARANCE_STYLES.glass.card);

      // Each appearance should have distinct button styles
      expect(APPEARANCE_STYLES.flat.button).not.toBe(APPEARANCE_STYLES.material.button);
      expect(APPEARANCE_STYLES.material.button).not.toBe(APPEARANCE_STYLES.glass.button);
      expect(APPEARANCE_STYLES.flat.button).not.toBe(APPEARANCE_STYLES.glass.button);
    });

    it("should have idp-button styles for all appearances", () => {
      expect(APPEARANCE_STYLES.flat["idp-button"]).toBeDefined();
      expect(typeof APPEARANCE_STYLES.flat["idp-button"]).toBe("string");

      expect(APPEARANCE_STYLES.material["idp-button"]).toBeDefined();
      expect(typeof APPEARANCE_STYLES.material["idp-button"]).toBe("string");

      expect(APPEARANCE_STYLES.glass["idp-button"]).toBeDefined();
      expect(typeof APPEARANCE_STYLES.glass["idp-button"]).toBe("string");
    });
  });

  describe("Type Safety", () => {
    it("should accept valid ThemeRoundness values", () => {
      const validValues: ThemeRoundness[] = ["edgy", "mid", "full"];
      validValues.forEach((value) => {
        expect(["edgy", "mid", "full"]).toContain(value);
      });
    });

    it("should accept valid ThemeLayout values", () => {
      const validValues: ThemeLayout[] = ["side-by-side", "top-to-bottom"];
      validValues.forEach((value) => {
        expect(["side-by-side", "top-to-bottom"]).toContain(value);
      });
    });

    it("should accept valid ThemeAppearance values", () => {
      const validValues: ThemeAppearance[] = ["flat", "material", "glass"];
      validValues.forEach((value) => {
        expect(["flat", "material", "glass"]).toContain(value);
      });
    });

    it("should accept valid ThemeSpacing values", () => {
      const validValues: ThemeSpacing[] = ["regular", "compact"];
      validValues.forEach((value) => {
        expect(["regular", "compact"]).toContain(value);
      });
    });
  });

  describe("Integration Tests", () => {
    it("should work correctly when switching between different roundness levels", () => {
      // Get classes for different roundness levels
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "edgy";
      const edgyClass = getComponentRoundness("card");

      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "mid";
      const midClass = getComponentRoundness("card");

      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const fullClass = getComponentRoundness("card");

      // All should be different
      expect(edgyClass).not.toBe(midClass);
      expect(midClass).not.toBe(fullClass);
      expect(edgyClass).not.toBe(fullClass);

      // All should be valid strings
      expect(typeof edgyClass).toBe("string");
      expect(typeof midClass).toBe("string");
      expect(typeof fullClass).toBe("string");
    });

    it("should maintain consistency across all theme properties", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      process.env.NEXT_PUBLIC_THEME_LAYOUT = "side-by-side";
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "material";
      process.env.NEXT_PUBLIC_THEME_SPACING = "compact";

      const config = getThemeConfig();

      // Verify all properties are set correctly
      expect(config.roundness).toBe("full");
      expect(config.layout).toBe("side-by-side");
      expect(config.appearance).toBe("material");
      expect(config.spacing).toBe("compact");

      // Verify component roundness is applied globally
      expect(config.componentRoundness?.card).toBe("full");
      expect(config.componentRoundness?.button).toBe("full");
    });

    it("should handle empty string environment variables by using defaults", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "";
      process.env.NEXT_PUBLIC_THEME_LAYOUT = "";

      const config = getThemeConfig();

      expect(config.roundness).toBe(DEFAULT_THEME.roundness);
      expect(config.layout).toBe(DEFAULT_THEME.layout);
    });
  });
});
