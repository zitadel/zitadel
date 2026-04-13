// Theme configuration system for customizable login experience

export type ThemeRoundness = "edgy" | "mid" | "full";
export type ThemeLayout = "side-by-side" | "top-to-bottom";
export type ThemeAppearance = "flat" | "material" | "glass";
export type ThemeSpacing = "regular" | "compact";

export interface ComponentRoundnessConfig {
  card: ThemeRoundness;
  button: ThemeRoundness;
  input: ThemeRoundness;
  image: ThemeRoundness;
  avatar: ThemeRoundness;
  avatarContainer: ThemeRoundness;
  themeSwitch: ThemeRoundness;
}

export interface ThemeConfig {
  roundness: ThemeRoundness; // Global fallback
  componentRoundness?: ComponentRoundnessConfig; // Component-specific overrides
  layout: ThemeLayout;
  backgroundImage?: string;
  appearance: ThemeAppearance;
  spacing: ThemeSpacing;
}

// Default component-specific roundness configuration
export const DEFAULT_COMPONENT_ROUNDNESS: ComponentRoundnessConfig = {
  card: "mid",
  button: "mid",
  input: "mid",
  image: "mid",
  avatar: "full", // Avatars default to full roundness
  avatarContainer: "full", // Avatar containers default to full roundness
  themeSwitch: "full", // Theme switch defaults to full roundness
};

// Default theme configuration
export const DEFAULT_THEME: ThemeConfig = {
  roundness: "mid",
  componentRoundness: DEFAULT_COMPONENT_ROUNDNESS,
  layout: "top-to-bottom",
  appearance: "flat",
  spacing: "regular",
};

// Get theme configuration from environment variables
export function getThemeConfig(): ThemeConfig {
  const globalRoundness = process.env.NEXT_PUBLIC_THEME_ROUNDNESS as ThemeRoundness;

  // If global roundness is set via env var, use it for all components
  // Otherwise, use component-specific defaults
  const componentRoundness = globalRoundness
    ? {
        card: globalRoundness,
        button: globalRoundness,
        input: globalRoundness,
        image: globalRoundness,
        avatar: globalRoundness,
        avatarContainer: globalRoundness,
        themeSwitch: globalRoundness,
      }
    : DEFAULT_COMPONENT_ROUNDNESS;

  return {
    roundness: globalRoundness || DEFAULT_THEME.roundness,
    componentRoundness: componentRoundness,
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
    avatar: "rounded-none",
    avatarContainer: "rounded-none",
    themeSwitch: "rounded-none",
  },
  mid: {
    card: "rounded-lg",
    button: "rounded-md",
    input: "rounded-md",
    image: "rounded-lg",
    avatar: "rounded-lg",
    avatarContainer: "rounded-md",
    themeSwitch: "rounded-md",
  },
  full: {
    card: "rounded-3xl",
    button: "rounded-full",
    input: "rounded-full pl-4",
    image: "rounded-full",
    avatar: "rounded-full",
    avatarContainer: "rounded-full",
    themeSwitch: "rounded-full",
  },
} as const;

// Helper function to get component-specific roundness
export function getComponentRoundness(componentType: keyof ComponentRoundnessConfig): string {
  const themeConfig = getThemeConfig();

  // Use component-specific roundness if available, otherwise fall back to global roundness
  const roundnessLevel = themeConfig.componentRoundness?.[componentType] || themeConfig.roundness;

  return ROUNDNESS_CLASSES[roundnessLevel][componentType];
}

// Spacing configuration
export const SPACING_STYLES = {
  regular: {
    spacing: "space-y-6",
    padding: "p-6 py-8",
  },
  compact: {
    spacing: "space-y-4",
    padding: "p-4",
  },
} as const;

// Appearance styling (complete design philosophies)
export const APPEARANCE_STYLES = {
  flat: {
    card: "bg-background-light-400 dark:bg-background-dark-500 border border-opacity-20 border border-black/10 dark:border-white/10",
    button: "border border-button-light-border dark:border-button-dark-border", // No shadows for flat design
    "idp-button": "border border-button-light-border dark:border-button-dark-border", // No shadows for flat design
    typography: "font-normal",
    background: "bg-background-light-500 dark:bg-background-dark-500", // Same as usual background
  },
  material: {
    card: "bg-background-light-400 dark:bg-background-dark-500 shadow-sm border-0",
    button: "shadow hover:shadow-xl active:shadow-xl", // Material shadows for buttons
    "idp-button":
      "!bg-background-[#00000020] !dark:bg-background-[#ffffff50] transition shadow shadow-md hover:shadow-lg active:shadow-xl", // Material shadows for IDP buttons
    typography: "font-medium",
    background: "bg-background-light-400 dark:bg-background-dark-500", // Current system (shade 400)
  },
  glass: {
    card: "backdrop-blur-md bg-white/10 dark:bg-black/10 border border-white/20 dark:border-white/10 shadow-xl",
    button:
      "backdrop-blur-sm bg-white/20 dark:bg-black/20 border border-white/30 dark:border-white/20 shadow-lg hover:shadow-xl", // Glass effect for buttons
    "idp-button":
      "backdrop-blur-sm bg-white/20 dark:bg-black/20 border border-white/30 dark:border-white/20 shadow-lg hover:shadow-xl", // Glass effect for IDP buttons
    typography: "font-medium",
    background: "bg-transparent", // Transparent background to show blur effect
  },
} as const;
