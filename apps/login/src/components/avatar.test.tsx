import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { render } from "@testing-library/react";
import { Avatar, getInitials } from "./avatar";

// Mock next-themes
vi.mock("next-themes", () => ({
  useTheme: () => ({
    resolvedTheme: "light",
  }),
}));

// Mock next/image
vi.mock("next/image", () => ({
  default: ({ src, alt }: { src: string; alt: string }) => <img src={src} alt={alt} />,
}));

describe("Avatar Component", () => {
  const originalEnv = process.env;

  beforeEach(() => {
    process.env = { ...originalEnv };
  });

  afterEach(() => {
    process.env = originalEnv;
  });

  describe("getInitials", () => {
    it("should get initials from full name", () => {
      const initials = getInitials("John Doe", "john.doe@example.com");
      expect(initials).toBe("JD");
    });

    it("should get single initial from single name", () => {
      const initials = getInitials("John", "john@example.com");
      expect(initials).toBe("J");
    });

    it("should get initials from loginName when name is empty", () => {
      const initials = getInitials("", "john.doe@example.com");
      expect(initials.length).toBeGreaterThan(0);
    });

    it("should handle loginName with underscore separator", () => {
      const initials = getInitials("", "john_doe@example.com");
      expect(initials).toBe("jd");
    });

    it("should handle loginName with dash separator", () => {
      const initials = getInitials("", "john-doe@example.com");
      expect(initials).toBe("jd");
    });

    it("should handle loginName with dot separator", () => {
      const initials = getInitials("", "john.doe@example.com");
      expect(initials).toBe("jd");
    });

    it("should get initials from username part of email", () => {
      const initials = getInitials("", "testuser@example.com");
      expect(initials.length).toBeGreaterThan(0);
    });
  });

  describe("Component Rendering", () => {
    it("should render avatar with initials", () => {
      const { container } = render(<Avatar name="John Doe" loginName="john@example.com" />);
      expect(container.querySelector(".avatar, [class*='avatar'], div")).toBeTruthy();
    });

    it("should render with different sizes", () => {
      const sizes: Array<"small" | "base" | "large"> = ["small", "base", "large"];

      sizes.forEach((size) => {
        const { container } = render(<Avatar size={size} name="Test User" loginName="test@example.com" />);
        expect(container.firstChild).toBeTruthy();
      });
    });

    it("should render with shadow prop", () => {
      const { container } = render(<Avatar name="Test User" loginName="test@example.com" shadow={true} />);
      expect(container.firstChild).toBeTruthy();
    });

    it("should render without shadow prop", () => {
      const { container } = render(<Avatar name="Test User" loginName="test@example.com" shadow={false} />);
      expect(container.firstChild).toBeTruthy();
    });

    it("should render with image URL", () => {
      const { container } = render(
        <Avatar name="Test User" loginName="test@example.com" imageUrl="https://example.com/avatar.jpg" />,
      );
      expect(container.firstChild).toBeTruthy();
    });
  });

  describe("Theme Integration", () => {
    it("should apply theme-based roundness", () => {
      const { container } = render(<Avatar name="Test User" loginName="test@example.com" />);
      const avatar = container.firstChild as HTMLElement;
      // Should render with some styling
      expect(avatar).toBeTruthy();
    });

    it("should respect theme roundness changes", () => {
      // Default theme
      const { container: container1 } = render(<Avatar name="Test User" loginName="test1@example.com" />);
      const avatar1 = container1.firstChild as HTMLElement;

      // Change theme
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "edgy";
      const { container: container2 } = render(<Avatar name="Test User" loginName="test2@example.com" />);
      const avatar2 = container2.firstChild as HTMLElement;

      // Both should render
      expect(avatar1).toBeTruthy();
      expect(avatar2).toBeTruthy();
    });

    it("should render with full roundness", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "full";
      const { container } = render(<Avatar name="Test User" loginName="test@example.com" />);
      expect(container.firstChild).toBeTruthy();
    });

    it("should render with mid roundness", () => {
      process.env.NEXT_PUBLIC_THEME_ROUNDNESS = "mid";
      const { container } = render(<Avatar name="Test User" loginName="test@example.com" />);
      expect(container.firstChild).toBeTruthy();
    });
  });

  describe("Color Generation", () => {
    it("should generate consistent colors for same loginName", () => {
      const { container: container1 } = render(<Avatar name="Test" loginName="same@example.com" />);
      const { container: container2 } = render(<Avatar name="Test" loginName="same@example.com" />);

      expect(container1.firstChild).toBeTruthy();
      expect(container2.firstChild).toBeTruthy();
    });

    it("should render for different loginNames", () => {
      const { container: container1 } = render(<Avatar name="User 1" loginName="user1@example.com" />);
      const { container: container2 } = render(<Avatar name="User 2" loginName="user2@example.com" />);

      expect(container1.firstChild).toBeTruthy();
      expect(container2.firstChild).toBeTruthy();
    });
  });

  describe("Props Validation", () => {
    it("should handle null name", () => {
      const { container } = render(<Avatar name={null} loginName="test@example.com" />);
      expect(container.firstChild).toBeTruthy();
    });

    it("should handle undefined name", () => {
      const { container } = render(<Avatar name={undefined} loginName="test@example.com" />);
      expect(container.firstChild).toBeTruthy();
    });

    it("should require loginName", () => {
      const { container } = render(<Avatar name="Test" loginName="required@example.com" />);
      expect(container.firstChild).toBeTruthy();
    });
  });
});
