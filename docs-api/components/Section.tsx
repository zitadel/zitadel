import * as React from "react";

type Props = {
  columns: number;
  children: React.ReactNode;
};

export function Section({ children }: Props) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 xl:gap-16 my-4 py-10 xl:py-20">
      {children}
    </div>
  );
}
