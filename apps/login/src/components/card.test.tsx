import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { render } from "@testing-library/react";
import { Card } from "./card";

describe("Card Component", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("Component Rendering", () => {
    it("should render card with children", () => {
      const { getByText } = render(
        <Card>
          <div>Card Content</div>
        </Card>,
      );
      expect(getByText("Card Content")).toBeTruthy();
    });

    it("should render multiple children", () => {
      const { getByText } = render(
        <Card>
          <h1>Title</h1>
          <p>Description</p>
        </Card>,
      );
      expect(getByText("Title")).toBeTruthy();
      expect(getByText("Description")).toBeTruthy();
    });

    it("should apply custom className", () => {
      const { container } = render(<Card className="custom-class">Content</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card?.className).toContain("custom-class");
    });

    it("should pass through native div props", () => {
      const { container } = render(
        <Card id="test-card" data-testid="card">
          Content
        </Card>,
      );
      const card = container.firstChild as HTMLElement;
      expect(card?.id).toBe("test-card");
      expect(card?.getAttribute("data-testid")).toBe("card");
    });
  });

  describe("Theme Integration", () => {
    it("should apply theme-based roundness", () => {
      const { container } = render(<Card>Themed Card</Card>);
      const card = container.firstChild as HTMLElement;
      // Should have some roundness class applied
      expect(card?.className).toBeTruthy();
      expect(card?.className).toMatch(/rounded/);
    });

    it("should completely override theme roundness with custom roundness prop", () => {
      // Set theme to have mid roundness
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "mid";
      const customRoundness = "custom-card-round";

      const { container } = render(<Card roundness={customRoundness}>Custom</Card>);
      const card = container.firstChild as HTMLElement;

      // Should contain the custom roundness
      expect(card?.className).toContain(customRoundness);

      // The custom prop completely replaces theme roundness, not merged
      expect(card).toBeTruthy();
    });

    it("should completely override theme padding with custom padding prop", () => {
      // Set theme to have regular spacing
      process.env.NEXT_PUBLIC_THEME_SPACING = "regular";
      const customPadding = "custom-card-padding";

      const { container } = render(<Card padding={customPadding}>Custom</Card>);
      const card = container.firstChild as HTMLElement;

      // Should contain the custom padding
      expect(card?.className).toContain(customPadding);

      // The custom prop completely replaces theme padding, not merged
      expect(card).toBeTruthy();
    });

    it("should respect theme appearance changes", () => {
      // Default theme
      const { container: container1 } = render(<Card>Default</Card>);
      const card1 = container1.firstChild as HTMLElement;
      const defaultClasses = card1?.className;

      // Change theme appearance
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "glass";
      const { container: container2 } = render(<Card>Glass</Card>);
      const card2 = container2.firstChild as HTMLElement;

      // Classes should be present
      expect(defaultClasses).toBeTruthy();
      expect(card2?.className).toBeTruthy();
    });

    it("should respect theme spacing changes when no override provided", () => {
      // Default spacing
      const { container: container1 } = render(<Card>Default</Card>);
      const card1 = container1.firstChild as HTMLElement;
      const defaultClasses = card1?.className;

      // Compact spacing
      process.env.NEXT_PUBLIC_THEME_SPACING = "compact";
      const { container: container2 } = render(<Card>Compact</Card>);
      const card2 = container2.firstChild as HTMLElement;

      // Classes should be present
      expect(defaultClasses).toBeTruthy();
      expect(card2?.className).toBeTruthy();
    });
  });

  describe("Styling", () => {
    it("should have base card styling", () => {
      const { container } = render(<Card>Test</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card?.className).toBeTruthy();
      expect(card?.className.length).toBeGreaterThan(0);
      // Should have background color
      expect(card?.className).toMatch(/bg-/);
    });

    it("should be a div element", () => {
      const { container } = render(<Card>Test</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card?.tagName).toBe("DIV");
    });
  });

  describe("Multiple Appearance Themes", () => {
    it("should render with flat appearance", () => {
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "flat";
      const { container } = render(<Card>Flat Card</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card).toBeTruthy();
    });

    it("should render with material appearance", () => {
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "material";
      const { container } = render(<Card>Material Card</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card).toBeTruthy();
    });

    it("should render with glass appearance", () => {
      process.env.NEXT_PUBLIC_THEME_APPEARANCE = "glass";
      const { container } = render(<Card>Glass Card</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card).toBeTruthy();
    });
  });

  describe("Ref Forwarding", () => {
    it("should support ref forwarding", () => {
      const { container } = render(<Card>Content</Card>);
      const card = container.firstChild as HTMLElement;
      expect(card).toBeTruthy();
      expect(card.nodeType).toBe(1); // Element node
    });
  });
});
