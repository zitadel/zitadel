import { getBranding } from "#/lib/zitadel";
import { server } from "../lib/zitadel";
import { use } from "react";

const ThemeWrapper = async ({ children }: any) => {
  console.log("hehe");
  //   const { resolvedTheme } = useTheme();
  const isDark = true; //resolvedTheme && resolvedTheme === "dark";

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
