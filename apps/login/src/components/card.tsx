import { clsx } from "clsx";
import { HTMLAttributes, forwardRef, ReactNode } from "react";
import { getThemeConfig, APPEARANCE_STYLES, SPACING_STYLES, getComponentRoundness } from "@/lib/theme";

export interface CardProps extends HTMLAttributes<HTMLDivElement> {
  children: ReactNode;
  roundness?: string; // Allow override via props
  padding?: string; // Allow override via props
}

// Helper function to get default card roundness from theme
function getDefaultCardRoundness(): string {
  return getComponentRoundness("card");
}

// Helper function to get default padding from centralized theme system
function getDefaultCardPadding(): string {
  const themeConfig = getThemeConfig();
  return SPACING_STYLES[themeConfig.spacing].padding;
}

// Helper function to get default background from centralized theme system
function getDefaultCardBackground(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];

  // Use appearance-specific background if defined, otherwise fallback to material design (current system)
  return appearance?.background || "bg-background-light-400 dark:bg-background-dark-500";
}

// Helper function to get default card styling from centralized theme system
function getDefaultCardStyling(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];
  return appearance?.card || "shadow-sm border-0"; // Fallback to material design
}

// eslint-disable-next-line react/display-name
export const Card = forwardRef<HTMLDivElement, CardProps>(
  (
    {
      children,
      className = "",
      roundness, // Will use theme default if not provided
      padding, // Will use theme default if not provided
      ...props
    },
    ref,
  ) => {
    // Use theme-based values if not explicitly provided
    const actualRoundness = roundness || getDefaultCardRoundness();
    const actualPadding = padding || getDefaultCardPadding();
    const actualBackground = getDefaultCardBackground();
    const actualCardStyling = getDefaultCardStyling();

    return (
      <div
        ref={ref}
        className={clsx(
          actualBackground,
          actualCardStyling,
          actualPadding,
          actualRoundness, // Apply the full roundness classes directly
          className,
        )}
        {...props}
      >
        {children}
      </div>
    );
  },
);
