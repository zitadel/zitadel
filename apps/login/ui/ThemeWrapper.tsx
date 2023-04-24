"use client";

import { LabelPolicy } from "#/../../packages/zitadel-server/dist";
import { ColorService } from "#/utils/colors";
import { useEffect } from "react";

type Props = {
  branding: LabelPolicy | undefined;
  children: React.ReactNode;
};

const ThemeWrapper = ({ children, branding }: Props) => {
  useEffect(() => {
    const colorService = new ColorService(document, branding);
  }, []);

  const defaultClasses = "bg-background-light-600 dark:bg-background-dark-600";

  //   console.log(branding);
  //   useEffect(() => {
  //     if (branding) {
  //       document.documentElement.style.setProperty(
  //         "--background-color",
  //         branding?.backgroundColor
  //       );
  //       document.documentElement.style.setProperty(
  //         "--dark-background-color",
  //         branding?.backgroundColorDark
  //       );
  //     }
  //   }, []);

  return <div className={defaultClasses}>{children}</div>;
};

export default ThemeWrapper;
