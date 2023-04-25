"use client";

import { setTheme, LabelPolicyColors } from "#/utils/colors";
import { useEffect } from "react";

type Props = {
  branding: LabelPolicyColors | undefined;
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
