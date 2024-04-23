import clsx from "clsx";
import { ReactNode } from "react";

export enum BadgeState {
  Info = "info",
  Error = "error",
  Success = "success",
  Alert = "alert",
}

export type StateBadgeProps = {
  state: BadgeState;
  children: ReactNode;
};

const getBadgeClasses = (state: BadgeState) =>
  clsx({
    "w-fit border-box h-18.5px flex flex-row items-center whitespace-nowrap tracking-wider leading-4 items-center justify-center px-2 py-2px text-12px rounded-full shadow-sm":
      true,
    "bg-state-success-light-background text-state-success-light-color dark:bg-state-success-dark-background dark:text-state-success-dark-color ":
      state === BadgeState.Success,
    "bg-state-neutral-light-background text-state-neutral-light-color dark:bg-state-neutral-dark-background dark:text-state-neutral-dark-color":
      state === BadgeState.Info,
    "bg-state-error-light-background text-state-error-light-color dark:bg-state-error-dark-background dark:text-state-error-dark-color":
      state === BadgeState.Error,
    "bg-state-alert-light-background text-state-alert-light-color dark:bg-state-alert-dark-background dark:text-state-alert-dark-color":
      state === BadgeState.Alert,
  });

export function StateBadge({
  state = BadgeState.Success,
  children,
}: StateBadgeProps) {
  return <span className={`${getBadgeClasses(state)}`}>{children}</span>;
}
