import tinycolor from "tinycolor2";

export interface Color {
  name: string;
  hex: string;
  rgb: string;
  contrastColor: string;
}

export type MapName = "background" | "primary" | "warn" | "text" | "link";

export type ColorName =
  | "50"
  | "100"
  | "200"
  | "300"
  | "400"
  | "500"
  | "600"
  | "700"
  | "800"
  | "C900"
  | "A100"
  | "A200"
  | "A400"
  | "A700";

// export interface ColorMap {
//   background: { [key in ColorName]: Color[] };
//   primary: { [key in ColorName]: Color[] };
//   warn: { [key in ColorName]: Color[] };
//   text: { [key in ColorName]: Color[] };
//   link: { [key in ColorName]: Color[] };
// }

export type ColorMap = {
  [key in MapName]: { [key in ColorName]: Color[] };
};

export const DARK_PRIMARY = "#2073c4";
export const PRIMARY = "#5469d4";

export const DARK_WARN = "#ff3b5b";
export const WARN = "#cd3d56";

export const DARK_BACKGROUND = "#111827";
export const BACKGROUND = "#fafafa";

export const DARK_TEXT = "#ffffff";
export const TEXT = "#000000";

export class ColorService {
  dark: ColorMap;
  light: ColorMap;

  constructor() {
    const lP = {
      backgroundColor: BACKGROUND,
      backgroundColorDark: DARK_BACKGROUND,
      primaryColor: PRIMARY,
      primaryColorDark: DARK_PRIMARY,
      warnColor: WARN,
      warnColorDark: DARK_WARN,
      fontColor: TEXT,
      fontColorDark: DARK_TEXT,
      linkColor: BACKGROUND,
      linkColorDark: DARK_BACKGROUND,
    };
    this.dark = computeMap(lP, true);
    this.light = computeMap(lP, false);
  }
}

function computeColors(hex: string): Color[] {
  return [
    getColorObject(tinycolor(hex).lighten(52), "50"),
    getColorObject(tinycolor(hex).lighten(37), "100"),
    getColorObject(tinycolor(hex).lighten(26), "200"),
    getColorObject(tinycolor(hex).lighten(12), "300"),
    getColorObject(tinycolor(hex).lighten(6), "400"),
    getColorObject(tinycolor(hex), "500"),
    getColorObject(tinycolor(hex).darken(6), "600"),
    getColorObject(tinycolor(hex).darken(12), "700"),
    getColorObject(tinycolor(hex).darken(18), "800"),
    getColorObject(tinycolor(hex).darken(24), "900"),
    getColorObject(tinycolor(hex).lighten(50).saturate(30), "A100"),
    getColorObject(tinycolor(hex).lighten(30).saturate(30), "A200"),
    getColorObject(tinycolor(hex).lighten(10).saturate(15), "A400"),
    getColorObject(tinycolor(hex).lighten(5).saturate(5), "A700"),
  ];
}

function getColorObject(value: any, name: string): Color {
  const c = tinycolor(value);
  return {
    name: name,
    hex: c.toHexString(),
    rgb: c.toRgbString(),
    contrastColor: getContrast(c.toHexString()),
  };
}

function isLight(hex: string): boolean {
  const color = tinycolor(hex);
  return color.isLight();
}

function isDark(hex: string): boolean {
  const color = tinycolor(hex);
  return color.isDark();
}

function getContrast(color: string): string {
  const onBlack = tinycolor.readability("#000", color);
  const onWhite = tinycolor.readability("#fff", color);
  if (onBlack > onWhite) {
    return "hsla(0, 0%, 0%, 0.87)";
  } else {
    return "#ffffff";
  }
}

export function computeMap(labelpolicy: any, dark: boolean): ColorMap {
  const colorArray = {
    background: computeColors(
      dark ? labelpolicy.backgroundColorDark : labelpolicy.backgroundColor
    ),
    primary: computeColors(
      dark ? labelpolicy.primaryColorDark : labelpolicy.primaryColor
    ),
    warn: computeColors(
      dark ? labelpolicy.warnColorDark : labelpolicy.warnColor
    ),
    text: computeColors(
      dark ? labelpolicy.fontColorDark : labelpolicy.fontColor
    ),
    link: computeColors(
      dark ? labelpolicy.linkColorDark : labelpolicy.linkColor
    ),
  };

  let mapped: ColorMap = {} as any;
  Object.entries(colorArray).forEach(([mapname, colors]) => {
    (mapped as any)[mapname] = {};
    colors.forEach((color) => {
      (mapped as any)[mapname][`${color.name}`] = color.hex;
      (mapped as any)[mapname][`contrast-${color.name}`] = color.contrastColor;
    });
  });

  return mapped;
}
