import { clsx } from "clsx";
import { ComponentProps } from "react";

// Look of an action that is shown as text: a secondary action below a form, or a link inside a sentence.
export const textLinkClasses =
  "hover:text-primary-light-500 dark:hover:text-primary-dark-500 cursor-pointer text-sm font-semibold transition-all disabled:cursor-default disabled:text-gray-400 dark:disabled:text-gray-700";

export function TextButton({ className, type = "button", ...props }: ComponentProps<"button">) {
  return <button type={type} className={clsx(textLinkClasses, className)} {...props} />;
}
