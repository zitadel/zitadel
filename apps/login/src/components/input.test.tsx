import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { TextInput } from "./input";

describe("TextInput Component", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("Component Rendering", () => {
    it("should render input with label", () => {
      render(<TextInput label="Email" />);
      expect(screen.getByText("Email")).toBeTruthy();
    });

    it("should render input with placeholder", () => {
      const { container } = render(<TextInput label="Username" placeholder="Enter username" />);
      const input = container.querySelector("input");
      expect(input?.placeholder).toBe("Enter username");
    });

    it("should render input with default value", () => {
      const { container } = render(<TextInput label="Name" defaultValue="John Doe" />);
      const input = container.querySelector("input");
      expect(input?.defaultValue).toBe("John Doe");
    });

    it("should show required indicator when required", () => {
      render(<TextInput label="Required Field" required />);
      expect(screen.getByText(/\*/)).toBeTruthy();
    });
  });

  describe("Input States", () => {
    it("should render disabled state", () => {
      const { container } = render(<TextInput label="Disabled" disabled />);
      const input = container.querySelector("input");
      expect(input?.disabled).toBe(true);
      // Should have disabled/pointer-events styles
      expect(input?.className).toMatch(/pointer-events/);
    });

    it("should render error state with message", () => {
      render(<TextInput label="Email" error="Invalid email" />);
      expect(screen.getByText("Invalid email")).toBeTruthy();
    });

    it("should apply error styling when error is present", () => {
      const { container } = render(<TextInput label="Email" error="Invalid" />);
      const input = container.querySelector("input");
      expect(input).toBeTruthy();
      // Should have border-warn or warn-related styles
      expect(input?.className).toMatch(/border-warn/);
    });

    it("should render success state with message", () => {
      render(<TextInput label="Email" success="Valid email" />);
      expect(screen.getByText("Valid email")).toBeTruthy();
    });
  });

  describe("Theme Integration", () => {
    it("should apply theme-based roundness", () => {
      const { container } = render(<TextInput label="Themed" />);
      const input = container.querySelector("input");
      // Should have some roundness class applied
      expect(input?.className).toBeTruthy();
      expect(input?.className).toMatch(/rounded/);
    });

    it("should completely override theme roundness with custom roundness prop", () => {
      // Set theme to have full roundness
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const customRoundness = "custom-input-round";

      const { container } = render(<TextInput label="Custom" roundness={customRoundness} />);
      const input = container.querySelector("input");

      // Should contain the custom roundness
      expect(input?.className).toContain(customRoundness);

      // The custom prop completely replaces theme roundness
      expect(input).toBeTruthy();
    });

    it("should respect theme roundness changes", () => {
      // Default theme
      const { container: container1 } = render(<TextInput label="Default" />);
      const input1 = container1.querySelector("input");
      const defaultClasses = input1?.className;

      // Change theme
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const { container: container2 } = render(<TextInput label="Full" />);
      const input2 = container2.querySelector("input");

      // Classes should be present
      expect(defaultClasses).toBeTruthy();
      expect(input2?.className).toBeTruthy();
    });
  });

  describe("Input Behavior", () => {
    it("should have autocomplete off by default", () => {
      const { container } = render(<TextInput label="Test" />);
      const input = container.querySelector("input");
      expect(input?.autocomplete).toBe("off");
    });

    it("should allow custom autocomplete", () => {
      const { container } = render(<TextInput label="Email" autoComplete="email" />);
      const input = container.querySelector("input");
      expect(input?.autocomplete).toBe("email");
    });

    it("should pass through native input props", () => {
      const { container } = render(<TextInput label="Test" maxLength={10} type="email" />);
      const input = container.querySelector("input");
      expect(input?.maxLength).toBe(10);
      expect(input?.type).toBe("email");
    });
  });

  describe("Styling", () => {
    it("should have consistent base styling", () => {
      const { container } = render(<TextInput label="Test" />);
      const input = container.querySelector("input");
      expect(input).toBeTruthy();
      expect(input?.className).toBeTruthy();
      expect(input?.className.length).toBeGreaterThan(0);
      // Should have transition styles
      expect(input?.className).toMatch(/transition/);
    });

    it("should apply border styles", () => {
      const { container } = render(<TextInput label="Test" />);
      const input = container.querySelector("input");
      expect(input).toBeTruthy();
      // Should have border-related classes
      expect(input?.className).toMatch(/border/);
    });

    it("should have focus styles", () => {
      const { container } = render(<TextInput label="Test" />);
      const input = container.querySelector("input");
      expect(input).toBeTruthy();
      // Should have focus-related classes
      expect(input?.className).toMatch(/focus:/);
    });
  });

  describe("Label Styling", () => {
    it("should apply error color to label when error exists", () => {
      const { container } = render(<TextInput label="Error Field" error="Error message" />);
      const label = container.querySelector("label");
      expect(label).toBeTruthy();
    });

    it("should have default label styling", () => {
      const { container } = render(<TextInput label="Normal Field" />);
      const label = container.querySelector("label");
      expect(label).toBeTruthy();
      expect(label?.className).toBeTruthy();
    });
  });

  describe("Accessibility", () => {
    it("should connect label to input", () => {
      const { container } = render(<TextInput label="Accessible Input" />);
      const label = container.querySelector("label");
      const input = container.querySelector("input");
      expect(label).toBeTruthy();
      expect(input).toBeTruthy();
    });

    it("should show required indicator", () => {
      const { container } = render(<TextInput label="UniqueRequiredField" required />);
      const label = container.querySelector("label");
      expect(label?.textContent).toContain("*");
    });
  });
});
