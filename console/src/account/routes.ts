// The same exact allowlist is used for navigation and OIDC continuation.
export const accountPaths = [
  '/',
  '/perfil',
  '/seguranca',
  '/seguranca/senha',
  '/seguranca/passkeys',
  '/seguranca/mfa',
] as const;
export function safeAccountPath(value: unknown): string {
  return typeof value === 'string' && (accountPaths as readonly string[]).includes(value) ? value : '/';
}
