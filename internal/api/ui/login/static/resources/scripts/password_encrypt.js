// password_encrypt.js
// Client-side ECDH + AES-256-GCM password encryption for regulated deployments.
//
// Activated when the server sets PasswordEncryption.Enabled=true in config.
// The server generates an ephemeral P-256 ECDH keypair per password page render,
// embeds the public key in the page, and stores the private key server-side.
// This script generates its own ephemeral P-256 keypair, performs ECDH to derive
// a shared AES-256 key, and encrypts the password before the POST body is sent.
//
// Security guarantee: the POST body contains only the client public key and
// ciphertext. The server private key is never transmitted. A captured POST body
// alone is insufficient to recover the plaintext password.
//
// Payload format (dot-separated):
//   base64(clientUncompressedP256PubKey) . hex(ciphertext) . hex(gcm_tag) . hex(iv)
//
// Key derivation: ECDH shared secret → SHA-256 → 256-bit AES-GCM key.
// This matches the Web Crypto SubtleCrypto.deriveBits("ECDH") + SHA-256 digest path.
//
// Browser support: Web Crypto API (window.crypto.subtle) is available in all
// modern browsers: Chrome 37+, Firefox 34+, Safari 11+, Edge 12+.

(function () {
  "use strict";

  const IV_BYTES = 12;
  const GCM_TAG_BYTES = 16;

  function hexEncode(buffer) {
    return Array.from(new Uint8Array(buffer))
      .map((b) => b.toString(16).padStart(2, "0"))
      .join("");
  }

  // Derive a 256-bit AES-GCM key from the ECDH shared secret via SHA-256.
  function deriveAESKey(ecdhSharedSecret) {
    return crypto.subtle.digest("SHA-256", ecdhSharedSecret).then((keyMaterial) =>
      crypto.subtle.importKey(
        "raw",
        keyMaterial,
        { name: "AES-GCM", length: 256 },
        false,
        ["encrypt"]
      )
    );
  }

  function encryptPassword(password, serverPublicKeyBase64) {
    const enc = new TextEncoder();

    // Decode the server's ephemeral public key (raw uncompressed P-256 bytes).
    const serverPubKeyBytes = Uint8Array.from(atob(serverPublicKeyBase64), (c) =>
      c.charCodeAt(0)
    );

    // Import server public key.
    const serverPubKeyPromise = crypto.subtle.importKey(
      "raw",
      serverPubKeyBytes,
      { name: "ECDH", namedCurve: "P-256" },
      false,
      []
    );

    // Generate client ephemeral keypair.
    const clientKeyPairPromise = crypto.subtle.generateKey(
      { name: "ECDH", namedCurve: "P-256" },
      true, // extractable — we need to send the public key to the server
      ["deriveKey", "deriveBits"]
    );

    return Promise.all([serverPubKeyPromise, clientKeyPairPromise])
      .then(([serverPubKey, clientKeyPair]) => {
        // Perform ECDH: derive raw shared secret bits.
        const ecdhBitsPromise = crypto.subtle.deriveBits(
          { name: "ECDH", public: serverPubKey },
          clientKeyPair.privateKey,
          256
        );

        // Export client public key (raw uncompressed, 65 bytes for P-256).
        const clientPubKeyPromise = crypto.subtle.exportKey(
          "raw",
          clientKeyPair.publicKey
        );

        return Promise.all([ecdhBitsPromise, clientPubKeyPromise]);
      })
      .then(([sharedSecret, clientPubKeyBytes]) => {
        return deriveAESKey(sharedSecret).then((aesKey) => {
          const iv = crypto.getRandomValues(new Uint8Array(IV_BYTES));

          return crypto.subtle
            .encrypt({ name: "AES-GCM", iv, tagLength: 128 }, aesKey, enc.encode(password))
            .then((encrypted) => {
              // Web Crypto appends the 16-byte GCM tag to the ciphertext buffer.
              const full = new Uint8Array(encrypted);
              const ciphertext = full.slice(0, -GCM_TAG_BYTES);
              const tag = full.slice(-GCM_TAG_BYTES);

              // Encode client public key as base64.
              const clientPubKeyB64 = btoa(
                String.fromCharCode(...new Uint8Array(clientPubKeyBytes))
              );

              return (
                clientPubKeyB64 +
                "." +
                hexEncode(ciphertext) +
                "." +
                hexEncode(tag) +
                "." +
                hexEncode(iv)
              );
            });
        });
      });
  }

  function attachEncryptionToForm() {
    const form = document.querySelector("form");
    if (!form) { return; }

    const passwordField = document.getElementById("password");
    // serverPubKey is embedded by the template as a data attribute on the form
    // so it cannot be confused with authRequestID and requires no extra hidden field.
    const serverPubKey = form.dataset.serverPubKey;
    if (!passwordField || !serverPubKey) { return; }

    if (form.dataset.encryptAttached) { return; }
    form.dataset.encryptAttached = "1";

    form.addEventListener("submit", (event) => {
      if (form.dataset.encrypted === "1") {
        form.dataset.encrypted = "";
        return;
      }

      event.preventDefault();

      const rawPassword = passwordField.value;
      if (!rawPassword) {
        form.submit();
        return;
      }

      encryptPassword(rawPassword, serverPubKey)
        .then((payload) => {
          passwordField.value = payload;
          form.dataset.encrypted = "1";
          form.requestSubmit ? form.requestSubmit() : form.submit();
        })
        .catch(() => {
          // Encryption failed (e.g. key import error). Submit plaintext.
          // If AllowPlaintextFallback=false server-side, this submission is
          // rejected there. No silent bypass of the encryption requirement.
          form.submit();
        });
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", attachEncryptionToForm);
  } else {
    attachEncryptionToForm();
  }
}());
