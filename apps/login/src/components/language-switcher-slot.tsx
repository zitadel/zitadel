"use client";

import { LanguageSwitcher } from "./language-switcher";
import { useLanguages } from "./languages-context";

/**
 * Thin client wrapper that reads languages from context and renders LanguageSwitcher.
 * Used as a slot component inside the Liquid template rendering pipeline.
 */
export function LanguageSwitcherSlot() {
  const languages = useLanguages();
  if (!languages.length) return null;
  return <LanguageSwitcher languages={languages} />;
}
