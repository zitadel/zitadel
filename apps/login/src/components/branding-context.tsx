"use client";

import { createContext, ReactNode, useContext } from "react";

interface BrandingContextValue {
  themeMode: number; // 0=UNSPECIFIED, 1=AUTO, 2=LIGHT, 3=DARK
}

const BrandingContext = createContext<BrandingContextValue>({
  themeMode: 0, // Default to UNSPECIFIED (toggle visible with system option)
});

export function BrandingProvider({ themeMode, children }: { themeMode: number; children: ReactNode }) {
  return <BrandingContext.Provider value={{ themeMode }}>{children}</BrandingContext.Provider>;
}

export function useBrandingContext() {
  return useContext(BrandingContext);
}
