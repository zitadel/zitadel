const ZITADEL_DOMAIN = process.env.ZITADEL_API_URL
  ? new URL(process.env.ZITADEL_API_URL).hostname
  : '*.zitadel.cloud';

export const DEFAULT_CSP =
  `default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; connect-src 'self'; child-src; style-src 'self' 'unsafe-inline'; font-src 'self'; object-src 'none'; img-src 'self' https://${ZITADEL_DOMAIN};`;
