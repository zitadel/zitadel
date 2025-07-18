import { ReactNode } from "react";

export function Skeleton({ children }: { children?: ReactNode }) {
  return (
    <div className="skeleton flex flex-row items-center justify-center rounded-lg bg-background-light-600 px-8 py-12 dark:bg-background-dark-600">
      {children}
    </div>
  );
}
