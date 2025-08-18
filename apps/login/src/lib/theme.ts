// Theme configuration system for customizable login experience

export type ThemeRoundness = "edgy" | "mid" | "full";
export type ThemeLayout = "side-by-side" | "top-to-bottom";
export type ThemeAppearance = "flat" | "material";
export type ThemeSpacing = "regular" | "compact";

export interface ThemeConfig {
  roundness: ThemeRoundness;
  layout: ThemeLayout;
  backgroundImage?: string;
  appearance: ThemeAppearance;
  spacing: ThemeSpacing;
}

// Default theme configuration
export const DEFAULT_THEME: ThemeConfig = {
  roundness: "mid",
  layout: "top-to-bottom",
  appearance: "material",
  spacing: "regular",
};

// Get theme configuration from environment variables
export function getThemeConfig(): ThemeConfig {
  return {
    roundness: (process.env.NEXT_PUBLIC_THEME_ROUNDNESS as ThemeRoundness) || DEFAULT_THEME.roundness,
    layout: (process.env.NEXT_PUBLIC_THEME_LAYOUT as ThemeLayout) || DEFAULT_THEME.layout,
    backgroundImage: process.env.NEXT_PUBLIC_THEME_BACKGROUND_IMAGE || undefined,
    appearance: (process.env.NEXT_PUBLIC_THEME_APPEARANCE as ThemeAppearance) || DEFAULT_THEME.appearance,
    spacing: (process.env.NEXT_PUBLIC_THEME_SPACING as ThemeSpacing) || DEFAULT_THEME.spacing,
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
    input: "rounded-full pl-4",
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

// Spacing configuration
export const SPACING_STYLES = {
  regular: {
    spacing: "space-y-6",
    padding: "p-6",
  },
  compact: {
    spacing: "space-y-4",
    padding: "p-4",
  },
} as const;

// Appearance styling (complete design philosophies)
export const APPEARANCE_STYLES = {
  flat: {
    card: "border border-opacity-20 border border-black/10 dark:border-white/10",
    button: "border border-button-light-border dark:border-button-dark-border", // No shadows for flat design
    typography: "font-normal",
    background: "bg-background-light-500 dark:bg-background-dark-500", // Same as usual background
  },
  material: {
    card: "shadow-sm border-0",
    button: "shadow hover:shadow-xl active:shadow-xl", // Material shadows for buttons
    typography: "font-medium",
    background: "bg-background-light-400 dark:bg-background-dark-500", // Current system (shade 400)
  },
} as const;

// Helper function to get CSS classes for current theme
export function getThemeClasses(theme: ThemeConfig) {
  const roundness = ROUNDNESS_CLASSES[theme.roundness];
  const layout = LAYOUT_CLASSES[theme.layout];
  const spacing = SPACING_STYLES[theme.spacing];
  const appearance = APPEARANCE_STYLES[theme.appearance];

  return {
    roundness,
    layout,
    spacing,
    appearance,
    backgroundImage: theme.backgroundImage
      ? {
          container: "relative",
          overlay: "absolute inset-0 bg-black bg-opacity-50",
          background: `bg-cover bg-center bg-no-repeat`,
        }
      : null,
  };
}
