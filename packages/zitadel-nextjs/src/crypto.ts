/**
 * Encrypted cookie helpers for session storage.
 *
 * Uses the Web Crypto API (available in Next.js edge runtime and Node.js ≥ 20)
 * to AES-GCM encrypt/decrypt session data stored in HTTP-only cookies.
 */

const ALGORITHM = "AES-GCM";
const KEY_LENGTH = 256;
const IV_LENGTH = 12;
const ENCODER = new TextEncoder();
const DECODER = new TextDecoder();

/**
 * Derives a CryptoKey from a secret string using PBKDF2.
 * The secret should be at least 32 characters long.
 */
async function deriveKey(secret: string): Promise<CryptoKey> {
  const keyMaterial = await crypto.subtle.importKey(
    "raw",
    ENCODER.encode(secret),
    "PBKDF2",
    false,
    ["deriveKey"],
  );

  return crypto.subtle.deriveKey(
    {
      name: "PBKDF2",
      // Static salt — acceptable because the secret itself should be high-entropy.
      salt: ENCODER.encode("@zitadel/nextjs"),
      iterations: 100_000,
      hash: "SHA-256",
    },
    keyMaterial,
    { name: ALGORITHM, length: KEY_LENGTH },
    false,
    ["encrypt", "decrypt"],
  );
}

/**
 * Encrypts a plaintext string using AES-256-GCM with the given secret.
 * Returns a base64url-encoded string of `iv || ciphertext`.
 */
export async function encrypt(
  plaintext: string,
  secret: string,
): Promise<string> {
  const key = await deriveKey(secret);
  const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH));
  const ciphertext = await crypto.subtle.encrypt(
    { name: ALGORITHM, iv },
    key,
    ENCODER.encode(plaintext),
  );

  // Concatenate iv + ciphertext
  const combined = new Uint8Array(iv.length + ciphertext.byteLength);
  combined.set(iv);
  combined.set(new Uint8Array(ciphertext), iv.length);

  return base64UrlEncode(combined);
}

/**
 * Decrypts a base64url-encoded `iv || ciphertext` string using AES-256-GCM.
 * Returns the plaintext string, or `null` if decryption fails.
 */
export async function decrypt(
  encrypted: string,
  secret: string,
): Promise<string | null> {
  try {
    const key = await deriveKey(secret);
    const combined = base64UrlDecode(encrypted);
    const iv = combined.slice(0, IV_LENGTH);
    const ciphertext = combined.slice(IV_LENGTH);

    const plaintext = await crypto.subtle.decrypt(
      { name: ALGORITHM, iv },
      key,
      ciphertext,
    );

    return DECODER.decode(plaintext);
  } catch {
    return null;
  }
}

function base64UrlEncode(data: Uint8Array): string {
  let binary = "";
  for (const byte of data) {
    binary += String.fromCharCode(byte);
  }
  return btoa(binary).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function base64UrlDecode(str: string): Uint8Array {
  const padded = str.replace(/-/g, "+").replace(/_/g, "/");
  const binary = atob(padded);
  const bytes = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i);
  }
  return bytes;
}
