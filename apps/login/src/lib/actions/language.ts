"use server";

import { LANGS } from "@/lib/i18n";
import { setLanguageCookie as setLanguageCookieValue } from "@/lib/server/cookies";

export async function setLanguageCookie(language: string) {
  if (!LANGS.some((lang) => lang.code === language)) {
    return;
  }

  await setLanguageCookieValue(language);
}
