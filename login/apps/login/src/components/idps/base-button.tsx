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
        "flex-1 transition-all cursor-pointer flex flex-row items-center bg-background-light-400 text-text-light-500 dark:bg-background-dark-500 dark:text-text-dark-500 border border-divider-light hover:border-black dark:border-divider-dark hover:dark:border-white focus:border-primary-light-500 focus:dark:border-primary-dark-500 outline-none rounded-md px-4 text-sm",
        props.className,
      )}
    >
      <div className="flex-1 justify-between flex items-center gap-4">
        <div className="flex-1 flex flex-row items-center">
          {props.children}
        </div>
        {formStatus.pending && <Loader2Icon className="w-4 h-4 animate-spin" />}
      </div>
    </button>
  );
});
