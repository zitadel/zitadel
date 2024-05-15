import classNames from "clsx";
import React, {
  DetailedHTMLProps,
  forwardRef,
  InputHTMLAttributes,
  useEffect,
  useState,
} from "react";

export type CheckboxProps = DetailedHTMLProps<
  InputHTMLAttributes<HTMLInputElement>,
  HTMLInputElement
> & {
  checked: boolean;
  disabled?: boolean;
  onChangeVal?: (checked: boolean) => void;
};

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
  function Checkbox(
    {
      className = "",
      checked = false,
      disabled = false,
      onChangeVal,
      children,
      ...props
    },
    ref,
  ) {
    const [enabled, setEnabled] = useState<boolean>(checked);

    useEffect(() => {
      setEnabled(checked);
    }, [checked]);

    return (
      <div className="relative flex items-start">
        <div className="flex items-center h-5">
          <input
            ref={ref}
            checked={enabled}
            onChange={(event) => {
              setEnabled(event.target?.checked);
              onChangeVal && onChangeVal(event.target?.checked);
            }}
            disabled={disabled}
            type="checkbox"
            className={classNames(
              enabled
                ? "border-none text-primary-light-500 dark:text-primary-dark-500 bg-primary-light-500 active:bg-primary-light-500 dark:bg-primary-dark-500 active:dark:bg-primary-dark-500"
                : "border-2 border-gray-500 dark:border-white bg-transparent dark:bg-transparent",
              "focus:border-gray-500 focus:dark:border-white focus:ring-opacity-40 focus:dark:ring-opacity-40 focus:ring-offset-0 focus:ring-2 dark:focus:ring-offset-0 dark:focus:ring-2 focus:ring-gray-500 focus:dark:ring-white",
              "h-4 w-4 rounded-sm ring-0 outline-0 checked:ring-0 checked:dark:ring-0 active:border-none active:ring-0",
              "disabled:bg-gray-500 disabled:text-gray-500 disabled:border-gray-200 disabled:cursor-not-allowed",
              className,
            )}
            {...props}
          />
        </div>
        {children}
      </div>
    );
  },
);
