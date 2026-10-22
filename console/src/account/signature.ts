import { Injectable } from '@angular/core';
import { AbstractValidationHandler, ValidationParams } from 'angular-oauth2-oidc';

function bytes(value: string): Uint8Array<ArrayBuffer> {
  return Uint8Array.from(atob(value.replace(/-/g, '+').replace(/_/g, '/')), (c) => c.charCodeAt(0));
}
/** angular-oauth2-oidc's default NullValidationHandler does not verify JWT signatures. */
@Injectable()
export class AccountSignatureHandler extends AbstractValidationHandler {
  async validateSignature(params: ValidationParams): Promise<void> {
    const [header, payload, signature] = params.idToken.split('.');
    const parsed = JSON.parse(new TextDecoder().decode(bytes(header))) as { alg: string; kid: string };
    if (parsed.alg !== 'RS256' || !parsed.kid) throw new Error('Unsupported ID token signature');
    let keys = (params.jwks as { keys?: (JsonWebKey & { kid?: string })[] })?.keys ?? [];
    let jwk = keys.find((k: JsonWebKey & { kid?: string }) => k.kid === parsed.kid && k.kty === 'RSA');
    if (!jwk && params.loadKeys) {
      keys = ((await params.loadKeys()) as { keys?: (JsonWebKey & { kid?: string })[] })?.keys ?? [];
      jwk = keys.find((k: JsonWebKey & { kid?: string }) => k.kid === parsed.kid && k.kty === 'RSA');
    }
    if (!jwk || (jwk.use && jwk.use !== 'sig') || (jwk.alg && jwk.alg !== 'RS256')) throw new Error('Signing key missing');
    const key = await crypto.subtle.importKey('jwk', jwk, { name: 'RSASSA-PKCS1-v1_5', hash: 'SHA-256' }, false, ['verify']);
    if (
      !(await crypto.subtle.verify(
        'RSASSA-PKCS1-v1_5',
        key,
        bytes(signature),
        new TextEncoder().encode(`${header}.${payload}`),
      ))
    )
      throw new Error('Invalid ID token signature');
  }
  protected async calcHash(value: string, algorithm: string): Promise<string> {
    const digest = await crypto.subtle.digest(algorithm.toUpperCase(), new TextEncoder().encode(value));
    return String.fromCharCode(...new Uint8Array(digest));
  }
}
