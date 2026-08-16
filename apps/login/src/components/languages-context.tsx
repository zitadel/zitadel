"use client";

import { Lang } from "@/lib/i18n";
import { createContext, ReactNode, useContext } from "react";

const LanguagesContext = createContext<Lang[]>([]);

export function LanguagesProvider({ languages, children }: { languages: Lang[]; children: ReactNode }) {
  return <LanguagesContext.Provider value={languages}>{children}</LanguagesContext.Provider>;
}

export function useLanguages(): Lang[] {
  return useContext(LanguagesContext);
}
