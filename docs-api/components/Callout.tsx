import * as React from "react";

export function Callout({ title, children }) {
  return (
    <div className="flex flex-col rounded-sm px-4 py-3 border border-border-light dark:border-border-dark bg-white dark:bg-background-dark-400">
      <strong>{title}</strong>
      <span>{children}</span>
    </div>
  );
}
