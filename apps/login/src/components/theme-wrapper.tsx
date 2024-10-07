"use client";

import { setTheme } from "@/helpers/colors";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2/branding_settings_pb";
import { ReactNode, useEffect } from "react";

type Props = {
  branding: BrandingSettings | undefined;
  children: ReactNode;
};

export const ThemeWrapper = ({ children, branding }: Props) => {
  useEffect(() => {
    setTheme(document, branding);
  }, []);

  const defaultClasses = "bg-background-light-600 dark:bg-background-dark-600";

  return <div className={defaultClasses}>{children}</div>;
};
