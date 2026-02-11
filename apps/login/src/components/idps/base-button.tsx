"use client";

import { clsx } from "clsx";
import { Loader2Icon } from "lucide-react";
import { ButtonHTMLAttributes, DetailedHTMLProps, forwardRef } from "react";
import { useFormStatus } from "react-dom";
import { getComponentRoundness, getThemeConfig, APPEARANCE_STYLES } from "@/lib/theme";

export type SignInWithIdentityProviderProps = DetailedHTMLProps<
  ButtonHTMLAttributes<HTMLButtonElement>,
  HTMLButtonElement
> & {
  name?: string;
  e2e?: string;
  /** Optional icon URL for the provider (e.g. OIDC custom icon). */
  iconUrl?: string;
  /** Optional background color for the button (e.g. hex "#1976d2"). */
  backgroundColor?: string;
};

// Helper function to get default IDP button appearance from centralized theme system
function getDefaultIdpButtonAppearance(): string {
  const themeConfig = getThemeConfig();
  const appearance = APPEARANCE_STYLES[themeConfig.appearance];
  return appearance?.["idp-button"] || "border border-divider-light dark:border-divider-dark"; // Fallback to basic border
}

export const BaseButton = forwardRef<HTMLButtonElement, SignInWithIdentityProviderProps>(function BaseButton(props, ref) {
  const formStatus = useFormStatus();
  const buttonRoundness = getComponentRoundness("button");
  const idpButtonAppearance = getDefaultIdpButtonAppearance();

  const { iconUrl, backgroundColor, className, style, ...restProps } = props;
  const hasCustomStyle = Boolean(backgroundColor);

  return (
    <button
      {...restProps}
      type="submit"
      ref={ref}
      disabled={formStatus.pending}
      className={clsx(
        `flex flex-1 cursor-pointer flex-row items-center px-4 text-sm text-text-light-500 outline-none transition-all hover:border-black focus:border-primary-light-500 dark:text-text-dark-500 hover:dark:border-white focus:dark:border-primary-dark-500`,
        buttonRoundness,
        idpButtonAppearance,
        !hasCustomStyle && `bg-background-light-400 dark:bg-background-dark-500`,
        className,
      )}
      style={
        hasCustomStyle
          ? { ...style, backgroundColor }
          : style
      }
    >
      <div className="flex flex-1 items-center justify-between gap-4">
        <div className="flex min-w-0 flex-1 flex-row items-center gap-4">
          {iconUrl && (
            <div className="flex h-12 w-12 shrink-0 items-center justify-center">
              <img
                src={iconUrl}
                alt=""
                className="h-5 w-5 object-contain"
                referrerPolicy="no-referrer"
              />
            </div>
          )}
          {iconUrl ? <span className="ml-4 shrink-0">{props.children}</span> : props.children}
        </div>
        {formStatus.pending && <Loader2Icon className="h-4 w-4 animate-spin" />}
      </div>
    </button>
  );
});
