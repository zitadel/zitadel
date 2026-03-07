import { render, waitFor } from "@testing-library/react";
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";

import { CustomCssWrapper } from "./custom-css-wrapper";

const mockUseThemeConfig = vi.fn();

vi.mock("@/lib/theme-hooks", () => ({
  useThemeConfig: () => mockUseThemeConfig(),
}));

describe("CustomCssWrapper", () => {
  beforeEach(() => {
    mockUseThemeConfig.mockReset();
    document.querySelectorAll('link[data-custom-css="true"]').forEach((link) => {
      link.remove();
    });
  });

  afterEach(() => {
    mockUseThemeConfig.mockReset();
    document.querySelectorAll('link[data-custom-css="true"]').forEach((link) => {
      link.remove();
    });
  });

  it("should render children when no custom CSS file is provided", () => {
    mockUseThemeConfig.mockReturnValue({
      customCssFile: undefined,
    } as any);

    const { getByText } = render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    expect(getByText("Test Child")).toBeTruthy();
  });

  it("should not inject link tag when custom CSS file is empty string", () => {
    mockUseThemeConfig.mockReturnValue({
      customCssFile: "",
    } as any);

    render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    const link = document.querySelector('link[data-custom-css="true"]');
    expect(link).toBeNull();
  });

  it("should inject link tag with local file path", async () => {
    const customCssFile = "/custom-theme.css";
    mockUseThemeConfig.mockReturnValue({
      customCssFile,
    } as any);

    render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    await waitFor(() => {
      const link = document.querySelector('link[data-custom-css="true"]');
      expect(link).toBeTruthy();
      expect(link?.getAttribute("rel")).toBe("stylesheet");
      expect(link?.getAttribute("href")).toBe(customCssFile);
    });
  });

  it("should remove link tag on cleanup", async () => {
    const customCssFile = "/custom-theme.css";
    mockUseThemeConfig.mockReturnValue({
      customCssFile,
    } as any);

    const { unmount } = render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    await waitFor(() => {
      const link = document.querySelector('link[data-custom-css="true"]');
      expect(link).toBeTruthy();
    });

    unmount();

    await waitFor(() => {
      const link = document.querySelector('link[data-custom-css="true"]');
      expect(link).toBeNull();
    });
  });

  it("should only inject one link tag even with multiple renders", async () => {
    const customCssFile = "/themes/dark.css";
    mockUseThemeConfig.mockReturnValue({
      customCssFile,
    } as any);

    const { rerender } = render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    rerender(
      <CustomCssWrapper>
        <div>Test Child Updated</div>
      </CustomCssWrapper>,
    );

    await waitFor(() => {
      const links = document.querySelectorAll('link[data-custom-css="true"]');
      expect(links.length).toBe(1);
    });
  });

  it("should update href when customCssFile changes", async () => {
    mockUseThemeConfig.mockReturnValue({
      customCssFile: "/first.css",
    } as any);

    const { rerender } = render(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    await waitFor(() => {
      const link = document.querySelector('link[data-custom-css="true"]');
      expect(link?.getAttribute("href")).toBe("/first.css");
    });

    mockUseThemeConfig.mockReturnValue({
      customCssFile: "/second.css",
    } as any);

    rerender(
      <CustomCssWrapper>
        <div>Test Child</div>
      </CustomCssWrapper>,
    );

    await waitFor(() => {
      const link = document.querySelector('link[data-custom-css="true"]');
      expect(link?.getAttribute("href")).toBe("/second.css");
    });

    const links = document.querySelectorAll('link[data-custom-css="true"]');
    expect(links.length).toBe(1);
  });
});
