// img-src includes https: so OIDC/IdP provider icon URLs (e.g. from CDNs) can load
export const DEFAULT_CSP =
  "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; connect-src 'self'; child-src; style-src 'self' 'unsafe-inline'; font-src 'self'; object-src 'none'; img-src 'self' https:;";
