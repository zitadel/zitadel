// Theme configuration system for customizable login experience

export type ThemeRoundness = "edgy" | "mid" | "full";
export type ThemeLayout = "side-by-side" | "top-to-bottom";
export type ThemePreset = "professional" | "modern" | "minimal" | "corporate";

export interface ThemeConfig {
  roundness: ThemeRoundness;
  layout: ThemeLayout;
  backgroundImage?: string;
  preset: ThemePreset;
}

// Default theme configuration
export const DEFAULT_THEME: ThemeConfig = {
  roundness: "mid",
  layout: "side-by-side",
  preset: "professional",
};

// Get theme configuration from environment variables
export function getThemeConfig(): ThemeConfig {
  return {
    roundness: (process.env.NEXT_PUBLIC_THEME_ROUNDNESS as ThemeRoundness) || DEFAULT_THEME.roundness,
    layout: (process.env.NEXT_PUBLIC_THEME_LAYOUT as ThemeLayout) || DEFAULT_THEME.layout,
    backgroundImage: process.env.NEXT_PUBLIC_THEME_BACKGROUND_IMAGE || undefined,
    preset: (process.env.NEXT_PUBLIC_THEME_PRESET as ThemePreset) || DEFAULT_THEME.preset,
  };
}

// Roundness CSS classes
export const ROUNDNESS_CLASSES = {
  edgy: {
    card: "rounded-none",
    button: "rounded-none",
    input: "rounded-none",
    image: "rounded-none",
  },
  mid: {
    card: "rounded-lg",
    button: "rounded-md",
    input: "rounded-md",
    image: "rounded-lg",
  },
  full: {
    card: "rounded-3xl",
    button: "rounded-full",
    input: "rounded-full",
    image: "rounded-full",
  },
} as const;

// Layout CSS classes
export const LAYOUT_CLASSES = {
  "side-by-side": {
    container: "lg:grid lg:grid-cols-2 lg:gap-8 min-h-screen",
    brandSection: "hidden lg:flex lg:items-center lg:justify-center",
    formSection: "flex items-center justify-center px-4 py-12 sm:px-6 lg:px-8",
    formContainer: "w-full max-w-md space-y-8",
  },
  "top-to-bottom": {
    container: "min-h-screen flex flex-col",
    brandSection: "flex-shrink-0 py-16 text-center",
    formSection: "flex-1 flex items-center justify-center px-4 py-12 sm:px-6 lg:px-8",
    formContainer: "w-full max-w-md space-y-8",
  },
} as const;

// Preset styling (non-color related)
export const PRESET_STYLES = {
  professional: {
    card: "shadow-lg border",
    spacing: "space-y-6",
    padding: "p-6",
    typography: "font-medium",
  },
  modern: {
    card: "shadow-xl border-0",
    spacing: "space-y-8",
    padding: "p-8",
    typography: "font-semibold",
  },
  minimal: {
    card: "shadow-sm border",
    spacing: "space-y-4",
    padding: "p-4",
    typography: "font-normal",
  },
  corporate: {
    card: "shadow-md border",
    spacing: "space-y-6",
    padding: "p-6",
    typography: "font-medium",
  },
} as const;

// Helper function to get CSS classes for current theme
export function getThemeClasses(theme: ThemeConfig) {
  const roundness = ROUNDNESS_CLASSES[theme.roundness];
  const layout = LAYOUT_CLASSES[theme.layout];
  const preset = PRESET_STYLES[theme.preset];

  return {
    roundness,
    layout,
    preset,
    backgroundImage: theme.backgroundImage
      ? {
          container: "relative",
          overlay: "absolute inset-0 bg-black bg-opacity-50",
          background: `bg-cover bg-center bg-no-repeat`,
        }
      : null,
  };
}
