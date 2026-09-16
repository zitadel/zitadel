import { render } from "@testing-library/react";
import { BrandingSettings, ThemeMode } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ThemeWrapper } from "./theme-wrapper";

vi.mock("next-themes", () => ({
  useTheme: () => ({ setTheme: vi.fn() }),
}));

vi.mock("@/helpers/colors", () => ({
  setTheme: vi.fn(),
}));

function colorSchemeMeta(): HTMLMetaElement | null {
  return document.head.querySelector<HTMLMetaElement>('meta[name="color-scheme"]');
}

function renderWithThemeMode(themeMode: ThemeMode) {
  const branding = { themeMode } as BrandingSettings;
  return render(
    <ThemeWrapper branding={branding}>
      <div>child</div>
    </ThemeWrapper>,
  );
}

describe("ThemeWrapper color-scheme", () => {
  afterEach(() => {
    colorSchemeMeta()?.remove();
    document.documentElement.classList.remove("dark");
  });

  it("declares only light when the theme mode forces light", () => {
    renderWithThemeMode(ThemeMode.LIGHT);
    expect(colorSchemeMeta()?.content).toBe("only light");
  });

  it("declares only dark when the theme mode forces dark", () => {
    renderWithThemeMode(ThemeMode.DARK);
    expect(colorSchemeMeta()?.content).toBe("only dark");
  });

  it("leaves the choice to the browser when the theme mode is auto", () => {
    renderWithThemeMode(ThemeMode.AUTO);
    expect(colorSchemeMeta()).toBeNull();
  });

  it("removes a stale declaration when the theme mode changes to auto", () => {
    const { rerender } = renderWithThemeMode(ThemeMode.LIGHT);
    expect(colorSchemeMeta()?.content).toBe("only light");

    rerender(
      <ThemeWrapper branding={{ themeMode: ThemeMode.AUTO } as BrandingSettings}>
        <div>child</div>
      </ThemeWrapper>,
    );
    expect(colorSchemeMeta()).toBeNull();
  });

  it("reuses an existing meta element instead of adding a second one", () => {
    const existing = document.createElement("meta");
    existing.name = "color-scheme";
    existing.content = "light dark";
    document.head.appendChild(existing);

    renderWithThemeMode(ThemeMode.DARK);
    expect(document.head.querySelectorAll('meta[name="color-scheme"]')).toHaveLength(1);
    expect(existing.content).toBe("only dark");
  });
});
