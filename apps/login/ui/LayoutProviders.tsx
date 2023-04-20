"use client";
import { ThemeProvider } from "next-themes";
import ThemeWrapper from "./ThemeWrapper";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="system"
      storageKey="cp-theme"
      value={{ dark: "dark" }}
    >
      <ThemeWrapper>{children}</ThemeWrapper>
    </ThemeProvider>
  );
}
