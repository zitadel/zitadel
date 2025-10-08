import { clsx } from "clsx";
import { ReactNode } from "react";

const Label = ({
  children,
  animateRerendering,
  color,
}: {
  children: ReactNode;
  animateRerendering?: boolean;
  color?: "default" | "pink" | "blue" | "violet" | "cyan" | "orange" | "red";
}) => {
  return (
    <div
      className={clsx("rounded-full px-1.5", {
        "bg-gray-800 text-gray-500": color === "default",
        "bg-pink-500 text-pink-100": color === "pink",
        "bg-blue-500 text-blue-100": color === "blue",
        "bg-cyan-500 text-cyan-100": color === "cyan",
        "bg-red-500 text-red-100": color === "red",
        "bg-violet-500 text-violet-100": color === "violet",
        "bg-orange-500 text-orange-100": color === "orange",
        "animate-[highlight_1s_ease-in-out_1]": animateRerendering,
      })}
    >
      {children}
    </div>
  );
};
export const Boundary = ({
  children,
  labels = ["children"],
  size = "default",
  color = "default",
  animateRerendering = true,
}: {
  children: ReactNode;
  labels?: string[];
  size?: "small" | "default";
  color?: "default" | "pink" | "blue" | "violet" | "cyan" | "orange" | "red";
  animateRerendering?: boolean;
}) => {
  return (
    <div
      className={clsx("relative rounded-lg border border-dashed", {
        "p-3 lg:p-5": size === "small",
        "p-4 lg:p-9 lg:pb-6": size === "default",
        "border-divider-light dark:border-divider-dark": color === "default",
        "border-pink-500": color === "pink",
        "border-blue-500": color === "blue",
        "border-cyan-500": color === "cyan",
        "border-red-500": color === "red",
        "border-violet-500": color === "violet",
        "border-orange-500": color === "orange",
        "animate-[rerender_1s_ease-in-out_1] text-pink-500": animateRerendering,
      })}
    >
      <div
        className={clsx(
          "absolute -top-2 flex space-x-1 text-[9px] uppercase leading-4 tracking-widest",
          {
            "left-3 lg:left-5": size === "small",
            "left-4 lg:left-9": size === "default",
          },
        )}
      >
        {labels.map((label) => {
          return (
            <Label
              key={label}
              color={color}
              animateRerendering={animateRerendering}
            >
              {label}
            </Label>
          );
        })}
      </div>

      {children}
    </div>
  );
};
