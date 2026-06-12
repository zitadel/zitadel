"use client";

import { useLanguages } from "./languages-context";
import { LanguageSwitcher } from "./language-switcher";

/**
 * Thin client wrapper that reads languages from context and renders LanguageSwitcher.
 * Used as a slot component inside the Liquid template rendering pipeline.
 */
export function LanguageSwitcherSlot() {
  const languages = useLanguages();
  if (!languages.length) return null;
  return <LanguageSwitcher languages={languages} />;
}
