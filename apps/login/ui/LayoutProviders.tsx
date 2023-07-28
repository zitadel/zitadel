"use client";
import { ZitadelReactProvider } from "@zitadel/react";
import { ThemeProvider, useTheme } from "next-themes";

type Props = {
  children: React.ReactNode;
};

export function LayoutProviders({ children }: Props) {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme && resolvedTheme === "dark";

  //   useEffect(() => {
  //     console.log("layoutproviders useeffect");
  //     setTheme(document);
  //   });

  return (
    <div className={`${isDark ? "ui-dark" : "ui-light"} `}>
      <ZitadelReactProvider dark={isDark}>{children}</ZitadelReactProvider>
    </div>
  );
}
