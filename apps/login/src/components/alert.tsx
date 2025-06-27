import {
  ExclamationTriangleIcon,
  InformationCircleIcon,
} from "@heroicons/react/24/outline";
import { clsx } from "clsx";
import { ReactNode } from "react";

type Props = {
  children: ReactNode;
  type?: AlertType;
};

export enum AlertType {
  ALERT,
  INFO,
}

const yellow =
  "border-yellow-600/40 dark:border-yellow-500/20 bg-yellow-200/30 text-yellow-600 dark:bg-yellow-700/20 dark:text-yellow-200";
const red =
  "border-red-600/40 dark:border-red-500/20 bg-red-200/30 text-red-600 dark:bg-red-700/20 dark:text-red-200";
const neutral =
  "border-divider-light dark:border-divider-dark bg-black/5 text-gray-600 dark:bg-white/10 dark:text-gray-200";

export function Alert({ children, type = AlertType.ALERT }: Props) {
  return (
    <div
      className={clsx(
        "flex flex-row items-center justify-center border rounded-md py-2 pr-2 scroll-px-40",
        {
          [yellow]: type === AlertType.ALERT,
          [neutral]: type === AlertType.INFO,
        },
      )}
    >
      {type === AlertType.ALERT && (
        <ExclamationTriangleIcon className="flex-shrink-0 h-5 w-5 mr-2 ml-2" />
      )}
      {type === AlertType.INFO && (
        <InformationCircleIcon className="flex-shrink-0 h-5 w-5 mr-2 ml-2" />
      )}
      <span className="text-sm w-full ">{children}</span>
    </div>
  );
}
