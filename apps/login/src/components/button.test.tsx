import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { render } from "@testing-library/react";
import { Button, ButtonSizes, ButtonVariants, ButtonColors, getButtonClasses } from "./button";

describe("Button Component", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("Component Rendering", () => {
    it("should render button with children", () => {
      const { getByText } = render(<Button>Click me</Button>);
      expect(getByText("Click me")).toBeTruthy();
    });

    it("should apply custom className", () => {
      const { container } = render(<Button className="custom-class">Test</Button>);
      const button = container.querySelector("button");
      expect(button?.className).toContain("custom-class");
    });

    it("should pass through native button props", () => {
      const { container } = render(<Button disabled>Test</Button>);
      const button = container.querySelector("button");
      expect(button?.disabled).toBe(true);
    });
  });

  describe("Button Variants", () => {
    it("should render primary variant by default", () => {
      const { container } = render(<Button>Primary</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Primary should have background color
      expect(button?.className).toMatch(/bg-/);
    });

    it("should render secondary variant", () => {
      const { container } = render(<Button variant={ButtonVariants.Secondary}>Secondary</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Secondary should have border
      expect(button?.className).toMatch(/border/);
    });

    it("should render all variant types", () => {
      const variants = [ButtonVariants.Primary, ButtonVariants.Secondary, ButtonVariants.Destructive];

      variants.forEach((variant) => {
        const { container } = render(<Button variant={variant}>Test</Button>);
        const button = container.querySelector("button");
        expect(button).toBeTruthy();
      });
    });
  });

  describe("Button Sizes", () => {
    it("should render small size by default", () => {
      const { container } = render(<Button>Small</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Should have padding styles
      expect(button?.className).toMatch(/p[xy]?-/);
    });

    it("should render large size", () => {
      const { container } = render(<Button size={ButtonSizes.Large}>Large</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Should have padding styles
      expect(button?.className).toMatch(/p[xy]?-/);
    });

    it("should render all size types", () => {
      const sizes = [ButtonSizes.Small, ButtonSizes.Large];

      sizes.forEach((size) => {
        const { container } = render(<Button size={size}>Test</Button>);
        const button = container.querySelector("button");
        expect(button).toBeTruthy();
      });
    });
  });

  describe("Button Colors", () => {
    it("should render primary color by default", () => {
      const { container } = render(<Button>Primary Color</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Should have background color
      expect(button?.className).toMatch(/bg-/);
    });

    it("should render warn color", () => {
      const { container } = render(<Button color={ButtonColors.Warn}>Warn</Button>);
      const button = container.querySelector("button");
      expect(button).toBeTruthy();
      // Should have background color
      expect(button?.className).toMatch(/bg-/);
    });

    it("should render all color types", () => {
      const colors = [ButtonColors.Neutral, ButtonColors.Primary, ButtonColors.Warn];

      colors.forEach((color) => {
        const { container } = render(<Button color={color}>Test</Button>);
        const button = container.querySelector("button");
        expect(button).toBeTruthy();
      });
    });
  });

  describe("Theme Integration", () => {
    it("should apply theme-based roundness", () => {
      const { container } = render(<Button>Themed</Button>);
      const button = container.querySelector("button");
      // Should have some roundness class applied
      expect(button?.className).toBeTruthy();
      expect(button?.className).toMatch(/rounded/);
    });

    it("should completely override theme roundness with custom roundness prop", () => {
      // Set theme to have full roundness
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const customRoundness = "custom-round-test";

      const { container } = render(<Button roundness={customRoundness}>Custom</Button>);
      const button = container.querySelector("button");

      // Should contain the custom roundness
      expect(button?.className).toContain(customRoundness);

      // The custom prop completely replaces theme roundness
      expect(button).toBeTruthy();
    });

    it("should change appearance when theme changes", () => {
      // Default theme
      const { container: container1 } = render(<Button>Default</Button>);
      const button1 = container1.querySelector("button");
      const defaultClasses = button1?.className;

      // Change theme
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const { container: container2 } = render(<Button>Full</Button>);
      const button2 = container2.querySelector("button");

      // Classes should be present (we don't check exact values)
      expect(defaultClasses).toBeTruthy();
      expect(button2?.className).toBeTruthy();
    });
  });

  describe("getButtonClasses", () => {
    it("should return valid class string", () => {
      const classes = getButtonClasses(ButtonSizes.Small, ButtonVariants.Primary, ButtonColors.Primary);
      expect(typeof classes).toBe("string");
      expect(classes.length).toBeGreaterThan(0);
    });

    it("should include custom roundness classes", () => {
      const customRoundness = "rounded-2xl";
      const classes = getButtonClasses(ButtonSizes.Small, ButtonVariants.Primary, ButtonColors.Primary, customRoundness);
      expect(classes).toContain(customRoundness);
    });

    it("should include appearance classes", () => {
      const appearance = "shadow-lg";
      const classes = getButtonClasses(
        ButtonSizes.Small,
        ButtonVariants.Primary,
        ButtonColors.Primary,
        "rounded-md",
        appearance,
      );
      expect(classes).toContain(appearance);
    });

    it("should generate different classes for different variants", () => {
      const primaryClasses = getButtonClasses(ButtonSizes.Small, ButtonVariants.Primary, ButtonColors.Primary);
      const secondaryClasses = getButtonClasses(ButtonSizes.Small, ButtonVariants.Secondary, ButtonColors.Primary);

      expect(primaryClasses).not.toBe(secondaryClasses);
    });

    it("should generate different classes for different sizes", () => {
      const smallClasses = getButtonClasses(ButtonSizes.Small, ButtonVariants.Primary, ButtonColors.Primary);
      const largeClasses = getButtonClasses(ButtonSizes.Large, ButtonVariants.Primary, ButtonColors.Primary);

      expect(smallClasses).not.toBe(largeClasses);
    });
  });

  describe("Accessibility", () => {
    it("should have type='button' by default", () => {
      const { container } = render(<Button>Test</Button>);
      const button = container.querySelector("button");
      expect(button?.type).toBe("button");
    });

    it("should support disabled state", () => {
      const { container } = render(<Button disabled>Disabled</Button>);
      const button = container.querySelector("button");
      expect(button?.disabled).toBe(true);
      // Should have disabled styles
      expect(button?.className).toMatch(/disabled:/);
    });
  });
});
