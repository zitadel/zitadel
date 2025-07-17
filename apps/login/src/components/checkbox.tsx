import classNames from "clsx";
import {
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
        <div className="flex h-5 items-center">
          <div className="box-sizing block">
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
                "form-checkbox rounded border-gray-300 text-primary-light-500 shadow-sm focus:border-indigo-300 focus:ring focus:ring-indigo-200 focus:ring-opacity-50 focus:ring-offset-0 dark:text-primary-dark-500",
                className,
              )}
              {...props}
            />
          </div>
        </div>
        {children}
      </div>
    );
  },
);
