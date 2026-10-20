"use client";

import { getComponentRoundness } from "@/lib/theme";
import { EyeIcon, EyeSlashIcon } from "@heroicons/react/24/outline";
import { CheckCircleIcon } from "@heroicons/react/24/solid";
import { clsx } from "clsx";
import { useTranslations } from "next-intl";
import { ChangeEvent, DetailedHTMLProps, forwardRef, InputHTMLAttributes, ReactNode, useState } from "react";

export type TextInputProps = DetailedHTMLProps<InputHTMLAttributes<HTMLInputElement>, HTMLInputElement> & {
  label: string;
  suffix?: string;
  placeholder?: string;
  defaultValue?: string;
  error?: string | ReactNode;
  success?: string | ReactNode;
  disabled?: boolean;
  onChange?: (value: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (value: ChangeEvent<HTMLInputElement>) => void;
  roundness?: string; // Allow override via props
  showPasswordToggle?: boolean; // Escape hatch to hide the reveal button on password inputs
};

const styles = (error: boolean, disabled: boolean, roundnessClasses: string = "rounded-md") =>
  clsx(
    {
      "h-[40px] mb-[2px] p-[7px] bg-input-light-background dark:bg-input-dark-background transition-colors duration-300 grow": true,
      "border border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500": true,
      "focus:outline-none focus:ring-0 text-base text-black dark:text-white placeholder:italic placeholder-gray-700 dark:placeholder-gray-700": true,
      "border border-warn-light-500 dark:border-warn-dark-500 hover:border-warn-light-500 hover:dark:border-warn-dark-500 focus:border-warn-light-500 focus:dark:border-warn-dark-500":
        error,
      "pointer-events-none text-gray-500 dark:text-gray-800 border border-input-light-border dark:border-input-dark-border hover:border-light-hoverborder hover:dark:border-hoverborder cursor-default":
        disabled,
    },
    roundnessClasses, // Apply the full roundness classes directly
  );

// Helper function to get default input roundness from theme
function getDefaultInputRoundness(): string {
  return getComponentRoundness("input");
}

export const TextInput = forwardRef<HTMLInputElement, TextInputProps>(
  (
    {
      label,
      placeholder,
      defaultValue,
      suffix,
      required = false,
      error,
      disabled,
      success,
      onChange,
      onBlur,
      roundness,
      type,
      showPasswordToggle = true,
      ...props
    },
    ref,
  ) => {
    const t = useTranslations("common");

    // Use theme-based roundness if not explicitly provided
    const actualRoundness = roundness || getDefaultInputRoundness();

    const [passwordRevealed, setPasswordRevealed] = useState(false);
    const isPassword = type === "password";
    const hasToggle = isPassword && showPasswordToggle && !disabled;

    return (
      <label className="text-12px text-input-light-label dark:text-input-dark-label relative flex flex-col">
        <span className={`mb-1 leading-3 ${error ? "text-warn-light-500 dark:text-warn-dark-500" : ""}`}>
          {label} {required && "*"}
        </span>
        <input
          suppressHydrationWarning
          ref={ref}
          className={clsx(styles(!!error, !!disabled, actualRoundness), hasToggle && "pr-10")}
          defaultValue={defaultValue}
          required={required}
          disabled={disabled}
          placeholder={placeholder}
          autoComplete={props.autoComplete ?? "off"}
          onChange={(e) => onChange && onChange(e)}
          onBlur={(e) => onBlur && onBlur(e)}
          {...props}
          type={isPassword && passwordRevealed ? "text" : type}
        />

        {hasToggle && (
          <button
            type="button"
            data-testid="password-reveal-button"
            aria-label={passwordRevealed ? t("hidePassword") : t("showPassword")}
            aria-pressed={passwordRevealed}
            onMouseDown={(e) => e.preventDefault()}
            onClick={(e) => {
              e.preventDefault();
              setPasswordRevealed((revealed) => !revealed);
            }}
            className={clsx(
              "absolute right-[3px] bottom-[22px] z-30 translate-y-1/2 transform p-2",
              "text-gray-500 hover:text-black dark:text-gray-400 dark:hover:text-white",
              "focus-visible:ring-primary-light-500 dark:focus-visible:ring-primary-dark-500 focus:outline-none focus-visible:ring-2",
              actualRoundness.split(" ")[0],
            )}
          >
            {passwordRevealed ? <EyeSlashIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
          </button>
        )}

        {suffix && (
          <span
            className={clsx(
              "bg-background-light-500 dark:bg-background-dark-500 absolute right-[3px] bottom-[22px] z-30 translate-y-1/2 transform p-2",
              // Extract just the roundness part for the suffix (no padding)
              actualRoundness.split(" ")[0], // Take only the first part (rounded-full, rounded-md, etc.)
            )}
          >
            @{suffix}
          </span>
        )}

        <div className="leading-14.5px h-14.5px text-12px text-warn-light-500 dark:text-warn-dark-500 flex flex-row items-center">
          <span>{error ? error : " "}</span>
        </div>

        {success && (
          <div className="text-md mt-1 flex flex-row items-center text-green-500">
            <CheckCircleIcon className="h-4 w-4" />
            <span className="ml-1">{success}</span>
          </div>
        )}
      </label>
    );
  },
);
