"use client";

import { clsx } from "clsx";
import { Loader2Icon } from "lucide-react";
import { ButtonHTMLAttributes, DetailedHTMLProps, forwardRef } from "react";
import { useFormStatus } from "react-dom";

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
  const formStatus = useFormStatus();

  return (
    <button
      {...props}
      type="submit"
      ref={ref}
      disabled={formStatus.pending}
      className={clsx(
        "flex flex-1 cursor-pointer flex-row items-center rounded-md border border-divider-light bg-background-light-400 px-4 text-sm text-text-light-500 outline-none transition-all hover:border-black focus:border-primary-light-500 dark:border-divider-dark dark:bg-background-dark-500 dark:text-text-dark-500 hover:dark:border-white focus:dark:border-primary-dark-500",
        props.className,
      )}
    >
      <div className="flex flex-1 items-center justify-between gap-4">
        <div className="flex flex-1 flex-row items-center">
          {props.children}
        </div>
        {formStatus.pending && <Loader2Icon className="h-4 w-4 animate-spin" />}
      </div>
    </button>
  );
});
