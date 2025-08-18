import { clsx } from "clsx";
import { getCardRoundnessClasses, getRoundnessClasses, ThemeableProps } from "@/lib/themeUtils";
import { getThemeConfig, ROUNDNESS_CLASSES } from "@/lib/theme";

interface SkeletonCardProps extends ThemeableProps {
  isLoading?: boolean;
}

// Helper function to get default card roundness from theme
function getDefaultCardRoundness(): string {
  const themeConfig = getThemeConfig();
  return ROUNDNESS_CLASSES[themeConfig.roundness].card;
}

// Helper function to get default spacing from theme
function getDefaultSpacing(): string {
  return process.env.NEXT_PUBLIC_THEME_PRESET === "minimal"
    ? "space-y-4"
    : process.env.NEXT_PUBLIC_THEME_PRESET === "modern"
      ? "space-y-8"
      : "space-y-6";
}

// Helper function to get default padding from theme
function getDefaultPadding(): string {
  return process.env.NEXT_PUBLIC_THEME_PRESET === "minimal"
    ? "p-4"
    : process.env.NEXT_PUBLIC_THEME_PRESET === "modern"
      ? "p-8"
      : "p-6";
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
        getCardRoundnessClasses(actualRoundness),
      )}
    >
      <div className={actualSpacing}>
        <div className={clsx("h-14 bg-gray-700", getRoundnessClasses(actualRoundness, ""))} />
        <div className={clsx("h-3 w-11/12 bg-gray-700", getRoundnessClasses(actualRoundness, ""))} />
        <div className={clsx("h-3 w-8/12 bg-gray-700", getRoundnessClasses(actualRoundness, ""))} />
      </div>
    </div>
  );
};
