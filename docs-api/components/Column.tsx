import * as React from "react";

type Props = {
  position: "LEFT" | "RIGHT";
  children: React.ReactNode;
};

export function Column({ position, children }: Props) {
  return (
    <div className={`w-full ${position === "LEFT" ? "order-1" : "order-2"}`}>
      {children}
    </div>
  );
}
