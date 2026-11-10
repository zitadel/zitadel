import { ReactNode } from "react";

export function Skeleton({ children }: { children?: ReactNode }) {
  return (
    <div className="skeleton bg-background-light-600 dark:bg-background-dark-600 flex flex-row items-center justify-center rounded-lg px-8 py-12">
      {children}
    </div>
  );
}
