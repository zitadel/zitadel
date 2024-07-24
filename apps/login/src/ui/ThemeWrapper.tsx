"use client";

import { setTheme } from "@/utils/colors";
import { useEffect } from "react";
import { BrandingSettings } from "@zitadel/proto/zitadel/settings/v2beta/branding_settings_pb";
import { PartialMessage } from "@zitadel/client2";

type Props = {
  branding: PartialMessage<BrandingSettings> | undefined;
  children: React.ReactNode;
};

const ThemeWrapper = ({ children, branding }: Props) => {
  useEffect(() => {
    setTheme(document, branding);
  }, []);

  const defaultClasses = "bg-background-light-600 dark:bg-background-dark-600";

  return <div className={defaultClasses}>{children}</div>;
};

export default ThemeWrapper;
