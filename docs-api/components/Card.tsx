import * as React from "react";

type Props = {
  title: string;
  children: React.ReactNode;
};

export function Card({ title, children }: Props) {
  return (
    <div className="my-4 bg-white dark:bg-background-dark-400 border border-border-light dark:border-border-dark rounded-md w-full">
      <div className="py-2 px-4 bg-black/10 text-sm dark:bg-white/10">
        {title}
      </div>
      <div className="px-4">{children}</div>
    </div>
  );
}
