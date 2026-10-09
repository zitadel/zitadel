import { render } from "@testing-library/react";
import { BrandingSettings, ThemeMode } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ColorSchemeMeta, ThemeWrapper } from "./theme-wrapper";

vi.mock("next-themes", () => ({
  useTheme: () => ({ setTheme: vi.fn() }),
}));

vi.mock("@/helpers/colors", () => ({
  setTheme: vi.fn(),
}));

function colorSchemeMetas(): HTMLMetaElement[] {
  return Array.from(document.head.querySelectorAll<HTMLMetaElement>('meta[name="color-scheme"]'));
}

function wrapper(branding: BrandingSettings | undefined) {
  return (
    <ThemeWrapper branding={branding}>
      <div>child</div>
    </ThemeWrapper>
  );
}

function renderWithThemeMode(themeMode: ThemeMode) {
  return render(wrapper({ themeMode } as BrandingSettings));
}

describe("ThemeWrapper color-scheme", () => {
  afterEach(() => {
    document.documentElement.classList.remove("dark");
  });

  it("declares only light when the theme mode forces light", () => {
    const { unmount } = renderWithThemeMode(ThemeMode.LIGHT);
    expect(colorSchemeMetas().map((m) => m.content)).toEqual(["only light"]);
    unmount();
  });

  it("declares only dark when the theme mode forces dark", () => {
    const { unmount } = renderWithThemeMode(ThemeMode.DARK);
    expect(colorSchemeMetas().map((m) => m.content)).toEqual(["only dark"]);
    unmount();
  });

  const unforcedModes: [string, ThemeMode][] = [
    ["auto", ThemeMode.AUTO],
    ["unspecified", ThemeMode.UNSPECIFIED],
  ];

  it.each(unforcedModes)("leaves the choice to the browser when the theme mode is %s", (_name, themeMode) => {
    const { unmount } = renderWithThemeMode(themeMode);
    expect(colorSchemeMetas()).toHaveLength(0);
    unmount();
  });

  it("declares nothing without branding", () => {
    const { unmount } = render(wrapper(undefined));
    expect(colorSchemeMetas()).toHaveLength(0);
    unmount();
  });

  it("removes the declaration when the theme mode changes to auto", () => {
    const { rerender, unmount } = renderWithThemeMode(ThemeMode.LIGHT);
    expect(colorSchemeMetas()).toHaveLength(1);

    rerender(wrapper({ themeMode: ThemeMode.AUTO } as BrandingSettings));
    expect(colorSchemeMetas()).toHaveLength(0);
    unmount();
  });

  it("removes the declaration when branding disappears", () => {
    const { rerender, unmount } = renderWithThemeMode(ThemeMode.DARK);
    expect(colorSchemeMetas()).toHaveLength(1);

    rerender(wrapper(undefined));
    expect(colorSchemeMetas()).toHaveLength(0);
    unmount();
  });
});

describe("ColorSchemeMeta", () => {
  const cases: [string, ThemeMode | undefined, string[]][] = [
    ["light", ThemeMode.LIGHT, ["only light"]],
    ["dark", ThemeMode.DARK, ["only dark"]],
    ["auto", ThemeMode.AUTO, []],
    ["unspecified", ThemeMode.UNSPECIFIED, []],
    ["unknown", undefined, []],
  ];

  it.each(cases)("declares %s as %j", (_name, themeMode, expected) => {
    const { unmount } = render(<ColorSchemeMeta themeMode={themeMode} />);
    expect(colorSchemeMetas().map((m) => m.content)).toEqual(expected);
    unmount();
  });
});
