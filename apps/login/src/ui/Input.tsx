"use client";
import { CheckCircleIcon } from "@heroicons/react/24/solid";
import clsx from "clsx";
import React, {
  ChangeEvent,
  DetailedHTMLProps,
  forwardRef,
  InputHTMLAttributes,
  ReactNode,
} from "react";

export type TextInputProps = DetailedHTMLProps<
  InputHTMLAttributes<HTMLInputElement>,
  HTMLInputElement
> & {
  label: string;
  placeholder?: string;
  defaultValue?: string;
  error?: string | ReactNode;
  success?: string | ReactNode;
  disabled?: boolean;
  onChange?: (value: ChangeEvent<HTMLInputElement>) => void;
  onBlur?: (value: ChangeEvent<HTMLInputElement>) => void;
};

const styles = (error: boolean, disabled: boolean) =>
  clsx({
    "h-[40px] mb-[2px] rounded p-[7px] bg-input-light-background dark:bg-input-dark-background transition-colors duration-300 grow":
      true,
    "border border-input-light-border dark:border-input-dark-border hover:border-black hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500":
      true,
    "focus:outline-none focus:ring-0 text-base text-black dark:text-white placeholder:italic placeholder-gray-700 dark:placeholder-gray-700":
      true,
    "border border-warn-light-500 dark:border-warn-dark-500 hover:border-warn-light-500 hover:dark:border-warn-dark-500 focus:border-warn-light-500 focus:dark:border-warn-dark-500":
      error,
    "pointer-events-none text-gray-500 dark:text-gray-800 border border-input-light-border dark:border-input-dark-border hover:border-light-hoverborder hover:dark:border-hoverborder cursor-default":
      disabled,
  });

// eslint-disable-next-line react/display-name
export const TextInput = forwardRef<HTMLInputElement, TextInputProps>(
  (
    {
      label,
      placeholder,
      defaultValue,
      required = false,
      error,
      disabled,
      success,
      onChange,
      onBlur,
      ...props
    },
    ref,
  ) => {
    return (
      <label className="flex flex-col text-12px text-input-light-label dark:text-input-dark-label">
        <span
          className={`leading-3 mb-1 ${
            error ? "text-warn-light-500 dark:text-warn-dark-500" : ""
          }`}
        >
          {label} {required && "*"}
        </span>
        <input
          ref={ref}
          className={styles(!!error, !!disabled)}
          defaultValue={defaultValue}
          required={required}
          disabled={disabled}
          placeholder={placeholder}
          autoComplete={props.autoComplete ?? "off"}
          onChange={(e) => onChange && onChange(e)}
          onBlur={(e) => onBlur && onBlur(e)}
          {...props}
        />

        <div className="leading-14.5px h-14.5px text-warn-light-500 dark:text-warn-dark-500 flex flex-row items-center text-12px">
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
