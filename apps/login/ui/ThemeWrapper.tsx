"use client";

import { getBranding } from "#/lib/zitadel";
import { useTheme } from "next-themes";
import { server } from "../lib/zitadel";

const ThemeWrapper = async ({ children }: any) => {
  const { resolvedTheme } = useTheme();
  const isDark = resolvedTheme && resolvedTheme === "dark";

  try {
    const policy = await getBranding(server);

    const backgroundStyle = {
      backgroundColor: `${policy?.backgroundColorDark}.`,
    };

    console.log(policy);

    return (
      <div className={`${isDark ? "ui-dark" : "ui-light"} `}>
        <div style={backgroundStyle}>{children}</div>
      </div>
    );
  } catch (error) {
    console.error(error);

    return (
      <div className={`${isDark ? "ui-dark" : "ui-light"} `}>{children}</div>
    );
  }
};

export default ThemeWrapper;
