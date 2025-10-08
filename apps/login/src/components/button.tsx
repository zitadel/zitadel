import { clsx } from "clsx";
import { ButtonHTMLAttributes, DetailedHTMLProps, forwardRef } from "react";
import { ThemeableProps } from "@/lib/themeUtils";
import { getThemeConfig, getComponentRoundness, APPEARANCE_STYLES } from "@/lib/theme";

export enum ButtonSizes {
  Small = "Small",
  Large = "Large",
}

export enum ButtonVariants {
  Primary = "Primary",
  Secondary = "Secondary",
  Destructive = "Destructive",
}

export enum ButtonColors {
  Neutral = "Neutral",
  Primary = "Primary",
  Warn = "Warn",
}

export type ButtonProps = DetailedHTMLProps<ButtonHTMLAttributes<HTMLButtonElement>, HTMLButtonElement> & {
  size?: ButtonSizes;
  variant?: ButtonVariants;
  color?: ButtonColors;
} & ThemeableProps;

export const getButtonClasses = (
  size: ButtonSizes,
  variant: ButtonVariants,
  color: ButtonColors,
  roundnessClasses: string = "rounded-md", // Default fallback
  appearance: string = "", // Theme appearance (shadows, borders, etc.)
) =>
  clsx(
    {
      "box-border leading-36px text-14px inline-flex items-center focus:outline-none transition-colors transition-shadow duration-300": true,
      "disabled:border-none disabled:bg-gray-300 disabled:text-gray-600 disabled:shadow-none disabled:cursor-not-allowed disabled:dark:bg-gray-800 disabled:dark:text-gray-900":
        variant === ButtonVariants.Primary,
      "bg-primary-light-500 dark:bg-primary-dark-500 hover:bg-primary-light-400 hover:dark:bg-primary-dark-400 text-primary-light-contrast-500 dark:text-primary-dark-contrast-500":
        variant === ButtonVariants.Primary && color !== ButtonColors.Warn,
      "bg-warn-light-500 dark:bg-warn-dark-500 hover:bg-warn-light-400 hover:dark:bg-warn-dark-400 text-white dark:text-white":
        variant === ButtonVariants.Primary && color === ButtonColors.Warn,
      "border border-button-light-border dark:border-button-dark-border text-gray-950 hover:bg-gray-500 hover:bg-opacity-20 hover:dark:bg-white hover:dark:bg-opacity-10 focus:bg-gray-500 focus:bg-opacity-20 focus:dark:bg-white focus:dark:bg-opacity-10 dark:text-white disabled:text-gray-600 disabled:hover:bg-transparent disabled:dark:hover:bg-transparent disabled:cursor-not-allowed disabled:dark:text-gray-900":
        variant === ButtonVariants.Secondary,
      "border border-button-light-border dark:border-button-dark-border text-warn-light-500 dark:text-warn-dark-500 hover:bg-warn-light-500 hover:bg-opacity-10 dark:hover:bg-warn-light-500 dark:hover:bg-opacity-10 focus:bg-warn-light-500 focus:bg-opacity-20 dark:focus:bg-warn-light-500 dark:focus:bg-opacity-20":
        color === ButtonColors.Warn && variant !== ButtonVariants.Primary,
      "px-16 py-2": size === ButtonSizes.Large,
      "px-4 h-[36px]": size === ButtonSizes.Small,
    },
    roundnessClasses, // Apply the full roundness classes directly
    appearance, // Apply appearance-specific styling (shadows, borders, etc.)
  );

// Helper function to get default button roundness from theme
function getDefaultButtonRoundness(): string {
  return getComponentRoundness("button");
}

// Helper function to get default button appearance from centralized theme system
function getDefaultButtonAppearance(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];
  return appearance?.button || "border border-button-light-border dark:border-button-dark-border"; // Fallback to flat design
}

// eslint-disable-next-line react/display-name
export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      children,
      className = "",
      variant = ButtonVariants.Primary,
      size = ButtonSizes.Small,
      color = ButtonColors.Primary,
      roundness, // Will use theme default if not provided
      ...props
    },
    ref,
  ) => {
    // Use theme-based values if not explicitly provided
    const actualRoundness = roundness || getDefaultButtonRoundness();
    const actualAppearance = getDefaultButtonAppearance();

    return (
      <button
        type="button"
        ref={ref}
        className={`${getButtonClasses(size, variant, color, actualRoundness, actualAppearance)} ${className}`}
        {...props}
      >
        {children}
      </button>
    );
  },
);
