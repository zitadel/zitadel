import { clsx } from "clsx";
import { ThemeableProps } from "@/lib/themeUtils";
import { getThemeConfig, SPACING_STYLES, getComponentRoundness } from "@/lib/theme";

interface SkeletonCardProps extends ThemeableProps {
  isLoading?: boolean;
}

// Helper function to get default card roundness from theme
function getDefaultCardRoundness(): string {
  return getComponentRoundness("card");
}

// Helper function to get default spacing from centralized theme system
function getDefaultSpacing(): string {
  const themeConfig = getThemeConfig();
  return SPACING_STYLES[themeConfig.spacing].spacing;
}

// Helper function to get default padding from centralized theme system
function getDefaultPadding(): string {
  const themeConfig = getThemeConfig();
  return SPACING_STYLES[themeConfig.spacing].padding;
}
export const SkeletonCard = ({
  isLoading,
  roundness, // Will use theme default if not provided
  spacing, // Will use theme default if not provided
  padding, // Will use theme default if not provided
}: SkeletonCardProps) => {
  // Use theme-based values if not explicitly provided
  const actualRoundness = roundness || getDefaultCardRoundness();
  const actualSpacing = spacing || getDefaultSpacing();
  const actualPadding = padding || getDefaultPadding();

  return (
    <div
      className={clsx(
        "bg-gray-900/80",
        actualPadding,
        {
          "relative overflow-hidden before:absolute before:inset-0 before:-translate-x-full before:animate-[shimmer_1.5s_infinite] before:bg-gradient-to-r before:from-transparent before:via-white/10 before:to-transparent":
            isLoading,
        },
        actualRoundness, // Apply the full roundness classes directly
      )}
    >
      <div className={actualSpacing}>
        <div className={clsx("h-14 bg-gray-700", actualRoundness.split(" ")[0])} />
        <div className={clsx("h-3 w-11/12 bg-gray-700", actualRoundness.split(" ")[0])} />
        <div className={clsx("h-3 w-8/12 bg-gray-700", actualRoundness.split(" ")[0])} />
      </div>
    </div>
  );
};
