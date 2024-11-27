import { ReactNode } from "react";

export function Skeleton({ children }: { children?: ReactNode }) {
  return (
    <div className="skeleton py-12 px-8 rounded-lg bg-background-light-600 dark:bg-background-dark-600 flex flex-row items-center justify-center">
      {children}
    </div>
  );
}
