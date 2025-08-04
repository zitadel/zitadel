import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
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

export type ColorMap = {
  [_key in MapName]: Color[];
};

export const DARK_PRIMARY = "#2073c4";
export const PRIMARY = "#5469d4";

export const DARK_WARN = "#ff3b5b";
export const WARN = "#cd3d56";

export const DARK_BACKGROUND = "#111827";
export const BACKGROUND = "#fafafa";

export const DARK_TEXT = "#ffffff";
export const TEXT = "#000000";

export type LabelPolicyColors = {
  backgroundColor: string;
  backgroundColorDark: string;
  fontColor: string;
  fontColorDark: string;
  warnColor: string;
  warnColorDark: string;
  primaryColor: string;
  primaryColorDark: string;
};

type BrandingColors = {
  lightTheme: {
    backgroundColor: string;
    fontColor: string;
    primaryColor: string;
    warnColor: string;
  };
  darkTheme: {
    backgroundColor: string;
    fontColor: string;
    primaryColor: string;
    warnColor: string;
  };
};

export function setTheme(document: any, policy?: BrandingSettings) {
  const lP: BrandingColors = {
    lightTheme: {
      backgroundColor: policy?.lightTheme?.backgroundColor || BACKGROUND,
      fontColor: policy?.lightTheme?.fontColor || TEXT,
      primaryColor: policy?.lightTheme?.primaryColor || PRIMARY,
      warnColor: policy?.lightTheme?.warnColor || WARN,
    },
    darkTheme: {
      backgroundColor: policy?.darkTheme?.backgroundColor || DARK_BACKGROUND,
      fontColor: policy?.darkTheme?.fontColor || DARK_TEXT,
      primaryColor: policy?.darkTheme?.primaryColor || DARK_PRIMARY,
      warnColor: policy?.darkTheme?.warnColor || DARK_WARN,
    },
  };

  const dark = computeMap(lP, true);
  const light = computeMap(lP, false);

  setColorShades(dark.background, "background", "dark", document);
  setColorShades(light.background, "background", "light", document);

  setColorShades(dark.primary, "primary", "dark", document);
  setColorShades(light.primary, "primary", "light", document);

  setColorShades(dark.warn, "warn", "dark", document);
  setColorShades(light.warn, "warn", "light", document);

  setColorAlpha(dark.text, "text", "dark", document);
  setColorAlpha(light.text, "text", "light", document);

  setColorAlpha(dark.link, "link", "dark", document);
  setColorAlpha(light.link, "link", "light", document);
}

function setColorShades(
  map: Color[],
  type: string,
  theme: string,
  document: any,
) {
  map.forEach((color) => {
    document.documentElement.style.setProperty(
      `--theme-${theme}-${type}-${color.name}`,
      color.hex,
    );
    document.documentElement.style.setProperty(
      `--theme-${theme}-${type}-contrast-${color.name}`,
      color.contrastColor,
    );
  });
}

function setColorAlpha(
  map: Color[],
  type: string,
  theme: string,
  document: any,
) {
  map.forEach((color) => {
    document.documentElement.style.setProperty(
      `--theme-${theme}-${type}-${color.name}`,
      color.hex,
    );
    document.documentElement.style.setProperty(
      `--theme-${theme}-${type}-contrast-${color.name}`,
      color.contrastColor,
    );
    document.documentElement.style.setProperty(
      `--theme-${theme}-${type}-secondary-${color.name}`,
      `${color.hex}c7`,
    );
  });
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
  } as Color;
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

export function computeMap(branding: BrandingColors, dark: boolean): ColorMap {
  return {
    background: computeColors(
      dark
        ? branding.darkTheme.backgroundColor
        : branding.lightTheme.backgroundColor,
    ),
    primary: computeColors(
      dark ? branding.darkTheme.primaryColor : branding.lightTheme.primaryColor,
    ),
    warn: computeColors(
      dark ? branding.darkTheme.warnColor : branding.lightTheme.warnColor,
    ),
    text: computeColors(
      dark ? branding.darkTheme.fontColor : branding.lightTheme.fontColor,
    ),
    link: computeColors(
      dark ? branding.darkTheme.fontColor : branding.lightTheme.fontColor,
    ),
  };
}

export interface ColorShade {
  200: string;
  300: string;
  500: string;
  600: string;
  700: string;
  900: string;
}

export const COLORS = [
  {
    500: "#ef4444",
    200: "#fecaca",
    300: "#fca5a5",
    600: "#dc2626",
    700: "#b91c1c",
    900: "#7f1d1d",
  },
  {
    500: "#f97316",
    200: "#fed7aa",
    300: "#fdba74",
    600: "#ea580c",
    700: "#c2410c",
    900: "#7c2d12",
  },
  {
    500: "#f59e0b",
    200: "#fde68a",
    300: "#fcd34d",
    600: "#d97706",
    700: "#b45309",
    900: "#78350f",
  },
  {
    500: "#eab308",
    200: "#fef08a",
    300: "#fde047",
    600: "#ca8a04",
    700: "#a16207",
    900: "#713f12",
  },
  {
    500: "#84cc16",
    200: "#d9f99d",
    300: "#bef264",
    600: "#65a30d",
    700: "#4d7c0f",
    900: "#365314",
  },
  {
    500: "#22c55e",
    200: "#bbf7d0",
    300: "#86efac",
    600: "#16a34a",
    700: "#15803d",
    900: "#14532d",
  },
  {
    500: "#10b981",
    200: "#a7f3d0",
    300: "#6ee7b7",
    600: "#059669",
    700: "#047857",
    900: "#064e3b",
  },
  {
    500: "#14b8a6",
    200: "#99f6e4",
    300: "#5eead4",
    600: "#0d9488",
    700: "#0f766e",
    900: "#134e4a",
  },
  {
    500: "#06b6d4",
    200: "#a5f3fc",
    300: "#67e8f9",
    600: "#0891b2",
    700: "#0e7490",
    900: "#164e63",
  },
  {
    500: "#0ea5e9",
    200: "#bae6fd",
    300: "#7dd3fc",
    600: "#0284c7",
    700: "#0369a1",
    900: "#0c4a6e",
  },
  {
    500: "#3b82f6",
    200: "#bfdbfe",
    300: "#93c5fd",
    600: "#2563eb",
    700: "#1d4ed8",
    900: "#1e3a8a",
  },
  {
    500: "#6366f1",
    200: "#c7d2fe",
    300: "#a5b4fc",
    600: "#4f46e5",
    700: "#4338ca",
    900: "#312e81",
  },
  {
    500: "#8b5cf6",
    200: "#ddd6fe",
    300: "#c4b5fd",
    600: "#7c3aed",
    700: "#6d28d9",
    900: "#4c1d95",
  },
  {
    500: "#a855f7",
    200: "#e9d5ff",
    300: "#d8b4fe",
    600: "#9333ea",
    700: "#7e22ce",
    900: "#581c87",
  },
  {
    500: "#d946ef",
    200: "#f5d0fe",
    300: "#f0abfc",
    600: "#c026d3",
    700: "#a21caf",
    900: "#701a75",
  },
  {
    500: "#ec4899",
    200: "#fbcfe8",
    300: "#f9a8d4",
    600: "#db2777",
    700: "#be185d",
    900: "#831843",
  },
  {
    500: "#f43f5e",
    200: "#fecdd3",
    300: "#fda4af",
    600: "#e11d48",
    700: "#be123c",
    900: "#881337",
  },
];

export function getColorHash(value: string): ColorShade {
  let hash = 0;

  if (value.length === 0) {
    return COLORS[hash];
  }

  hash = hashCode(value);
  return COLORS[hash % COLORS.length];
}

export function hashCode(str: string, seed = 0): number {
  let h1 = 0xdeadbeef ^ seed,
    h2 = 0x41c6ce57 ^ seed;
  for (let i = 0, ch; i < str.length; i++) {
    ch = str.charCodeAt(i);
    h1 = Math.imul(h1 ^ ch, 2654435761);
    h2 = Math.imul(h2 ^ ch, 1597334677);
  }
  h1 =
    Math.imul(h1 ^ (h1 >>> 16), 2246822507) ^
    Math.imul(h2 ^ (h2 >>> 13), 3266489909);
  h2 =
    Math.imul(h2 ^ (h2 >>> 16), 2246822507) ^
    Math.imul(h1 ^ (h1 >>> 13), 3266489909);
  return 4294967296 * (2097151 & h2) + (h1 >>> 0);
}

export function getMembershipColor(role: string): ColorShade {
  const hash = hashCode(role);
  let color = COLORS[hash % COLORS.length];

  switch (role) {
    case "IAM_OWNER":
      color = COLORS[0];
      break;
    case "IAM_OWNER_VIEWER":
      color = COLORS[14];
      break;
    case "IAM_ORG_MANAGER":
      color = COLORS[11];
      break;
    case "IAM_USER_MANAGER":
      color = COLORS[8];
      break;

    case "ORG_OWNER":
      color = COLORS[16];
      break;
    case "ORG_USER_MANAGER":
      color = COLORS[8];
      break;
    case "ORG_OWNER_VIEWER":
      color = COLORS[14];
      break;
    case "ORG_USER_PERMISSION_EDITOR":
      color = COLORS[7];
      break;
    case "ORG_PROJECT_PERMISSION_EDITOR":
      color = COLORS[11];
      break;
    case "ORG_PROJECT_CREATOR":
      color = COLORS[12];
      break;

    case "PROJECT_OWNER":
      color = COLORS[9];
      break;
    case "PROJECT_OWNER_VIEWER":
      color = COLORS[10];
      break;
    case "PROJECT_OWNER_GLOBAL":
      color = COLORS[11];
      break;
    case "PROJECT_OWNER_VIEWER_GLOBAL":
      color = COLORS[12];
      break;

    default:
      color = COLORS[hash % COLORS.length];
      break;
  }

  return color;
}
