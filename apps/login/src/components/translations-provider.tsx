"use client";

import { createInstance } from "i18next";
import React from "react";
import { I18nextProvider } from "react-i18next";
import initTranslations from "../app/i18n";

type Props = {
  locale: string;
  children: React.ReactNode;
  namespaces: string[];
  resources: Record<string, Record<string, string>>;
};

export function TranslationsProvider({
  children,
  locale,
  namespaces,
  resources,
}: Props) {
  const i18n = createInstance();

  initTranslations(locale, namespaces, i18n, resources);

  return <I18nextProvider i18n={i18n}>{children}</I18nextProvider>;
}
