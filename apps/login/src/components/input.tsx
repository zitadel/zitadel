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
};

const styles = (error: boolean, disabled: boolean, hasPasswordToggle: boolean, roundnessClasses: string = "rounded-md") =>
  clsx(
    {
      "h-[40px] mb-[2px] p-[7px] bg-input-light-background dark:bg-input-dark-background transition-colors duration-300 grow": true,
      "border border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500": true,
      "focus:outline-none focus:ring-0 text-base text-black dark:text-white placeholder:italic placeholder-gray-700 dark:placeholder-gray-700": true,
      "pr-10": hasPasswordToggle,
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
      ...props
    },
    ref,
  ) => {
    // Use theme-based roundness if not explicitly provided
    const actualRoundness = roundness || getDefaultInputRoundness();

    const t = useTranslations("password");
    const isPasswordInput = props.type === "password";
    const [showPassword, setShowPassword] = useState(false);
    const toggleLabel = showPassword ? t("hidePassword") : t("showPassword");

    return (
      <label className="relative flex flex-col text-12px text-input-light-label dark:text-input-dark-label">
        <span className={`mb-1 leading-3 ${error ? "text-warn-light-500 dark:text-warn-dark-500" : ""}`}>
          {label} {required && "*"}
        </span>
        <div className="relative flex items-center">
          <input
            suppressHydrationWarning
            ref={ref}
            className={styles(!!error, !!disabled, isPasswordInput, actualRoundness)}
            defaultValue={defaultValue}
            required={required}
            disabled={disabled}
            placeholder={placeholder}
            autoComplete={props.autoComplete ?? "off"}
            onChange={(e) => onChange && onChange(e)}
            onBlur={(e) => onBlur && onBlur(e)}
            {...props}
            type={isPasswordInput && showPassword ? "text" : props.type}
          />

          {suffix && (
            <span
              className={clsx(
                "absolute bottom-[22px] right-[3px] z-30 translate-y-1/2 transform bg-background-light-500 p-2 dark:bg-background-dark-500",
                // Extract just the roundness part for the suffix (no padding)
                actualRoundness.split(" ")[0], // Take only the first part (rounded-full, rounded-md, etc.)
              )}
            >
              @{suffix}
            </span>
          )}

          {isPasswordInput && (
            <button
              type="button"
              aria-label={toggleLabel}
              aria-pressed={showPassword}
              title={toggleLabel}
              disabled={disabled}
              className={clsx(
                "absolute bottom-[22px] right-[5px] z-30 flex h-[30px] w-[30px] translate-y-1/2 transform flex-row items-center justify-center border-none bg-transparent text-gray-400 transition-colors hover:text-gray-900 disabled:cursor-default disabled:text-gray-500 dark:text-gray-400 dark:hover:text-white dark:disabled:text-gray-700",
                // Extract just the roundness part for the toggle button (no padding)
                actualRoundness.split(" ")[0], // Take only the first part (rounded-full, rounded-md, etc.)
              )}
              onClick={() => setShowPassword((value) => !value)}
            >
              {showPassword ? <EyeSlashIcon className="h-5 w-5" /> : <EyeIcon className="h-5 w-5" />}
            </button>
          )}
        </div>

        <div className="leading-14.5px h-14.5px flex flex-row items-center text-12px text-warn-light-500 dark:text-warn-dark-500">
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
