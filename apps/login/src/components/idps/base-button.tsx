"use client";

import { clsx } from "clsx";
import { ButtonHTMLAttributes, DetailedHTMLProps, forwardRef } from "react";

export type SignInWithIdentityProviderProps = DetailedHTMLProps<
  ButtonHTMLAttributes<HTMLButtonElement>,
  HTMLButtonElement
> & {
  name?: string;
  e2e?: string;
};

export const BaseButton = forwardRef<
  HTMLButtonElement,
  SignInWithIdentityProviderProps
>(function BaseButton(props, ref) {
  return (
    <button
      {...props}
      type="button"
      ref={ref}
      className={clsx(
        "transition-all cursor-pointer flex flex-row items-center bg-background-light-400 text-text-light-500 dark:bg-background-dark-500 dark:text-text-dark-500 border border-divider-light hover:border-black dark:border-divider-dark hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500 outline-none rounded-md px-4 text-sm",
        props.className,
      )}
    />
  );
});
