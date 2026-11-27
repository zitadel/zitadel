import { clsx } from "clsx";
import { ComponentProps } from "react";

const shared = "cursor-pointer text-sm font-semibold underline-offset-2 transition-all focus-visible:underline";

// A link inside a sentence or next to other text: the primary colour sets it apart without hovering.
export const textLinkClasses = clsx(
  shared,
  "text-primary-light-500 hover:text-primary-light-400 dark:text-primary-dark-500 hover:dark:text-primary-dark-400",
);

// An action on a line of its own, such as "Reset Password" below the password field.
const textActionClasses = clsx(shared, "hover:text-primary-light-500 dark:hover:text-primary-dark-500");

const disabledClasses = "disabled:cursor-default disabled:text-gray-400 dark:disabled:text-gray-700";

type Props = ComponentProps<"button"> & {
  // Set for a button that sits next to other text, so it is coloured like a link.
  inline?: boolean;
};

export function TextButton({ className, inline = false, type = "button", ...props }: Props) {
  return (
    <button
      type={type}
      className={clsx(inline ? textLinkClasses : textActionClasses, disabledClasses, className)}
      {...props}
    />
  );
}
