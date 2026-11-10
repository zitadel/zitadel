const LANGUAGE_PLACEHOLDER = "{{.Lang}}";

export function resolveLocalizedLegalLink(link: string | undefined, locale: string | undefined): string | undefined {
  if (!link || !locale) {
    return link;
  }

  return link.replaceAll(LANGUAGE_PLACEHOLDER, locale);
}
