const ZITADEL_DOMAIN = process.env.ZITADEL_API_URL
  ? new URL(process.env.ZITADEL_API_URL).hostname
  : '*.zitadel.cloud';

const ZITADEL_PROTOCOL = process.env.ZITADEL_API_URL
  ? new URL(process.env.ZITADEL_API_URL).protocol.replace(':', '')
  : 'https';

console.log('Environment ZITADEL_API_URL:', process.env.ZITADEL_API_URL);
console.log('Using ZITADEL_DOMAIN for CSP:', ZITADEL_DOMAIN);
console.log('Using ZITADEL_PROTOCOL for CSP:', ZITADEL_PROTOCOL);

export const DEFAULT_CSP =
  `default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; connect-src 'self'; child-src; style-src 'self' 'unsafe-inline'; font-src 'self'; object-src 'none'; img-src 'self' https://${ZITADEL_DOMAIN};`;
