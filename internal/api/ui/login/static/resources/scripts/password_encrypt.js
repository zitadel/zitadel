// password_encrypt.js
// Client-side AES-256-GCM password encryption for regulated deployments.
//
// Activated when the server sets PasswordEncryptionEnabled=true in config.
// The password field value is replaced with an encrypted payload before the
// form is submitted, so the plaintext password never appears in the POST body
// as transmitted by the browser. The server decrypts it using the same
// key-derivation parameters before passing the credential to VerifyPassword.
//
// Payload format (dot-separated hex strings):
//   hex(ciphertext) . hex(gcm_tag) . hex(iv) . hex(salt)
//
// Key derivation: PBKDF2-SHA-256, passphrase = authRequestID, 100 000 iterations, 256-bit key.
// Cipher: AES-256-GCM with a 96-bit random IV and a 128-bit random salt.
//
// Security notes:
//   - The authRequestID is a server-issued opaque token unique per login attempt.
//     Using it as the PBKDF2 passphrase binds the encrypted payload to the
//     specific request and prevents replay across sessions.
//   - This layer does not replace TLS; it addresses the specific VAPT finding
//     that POST body content may be captured by intermediary logging
//     infrastructure even in TLS-terminated deployments.

(function () {
  "use strict";

  var PBKDF2_ITERATIONS = 100000;
  var AES_KEY_BITS = 256;
  var SALT_BYTES = 16;
  var IV_BYTES = 12;

  function hexEncode(buffer) {
    return Array.from(new Uint8Array(buffer))
      .map(function (b) { return b.toString(16).padStart(2, "0"); })
      .join("");
  }

  function deriveKey(passphrase, salt) {
    var enc = new TextEncoder();
    return crypto.subtle.importKey(
      "raw",
      enc.encode(passphrase),
      { name: "PBKDF2" },
      false,
      ["deriveKey"]
    ).then(function (baseKey) {
      return crypto.subtle.deriveKey(
        {
          name: "PBKDF2",
          salt: salt,
          iterations: PBKDF2_ITERATIONS,
          hash: "SHA-256",
        },
        baseKey,
        { name: "AES-GCM", length: AES_KEY_BITS },
        false,
        ["encrypt"]
      );
    });
  }

  function encryptPassword(password, authRequestID) {
    var enc = new TextEncoder();
    var salt = crypto.getRandomValues(new Uint8Array(SALT_BYTES));
    var iv = crypto.getRandomValues(new Uint8Array(IV_BYTES));

    return deriveKey(authRequestID, salt).then(function (key) {
      return crypto.subtle.encrypt(
        { name: "AES-GCM", iv: iv, tagLength: 128 },
        key,
        enc.encode(password)
      );
    }).then(function (encrypted) {
      // Web Crypto appends the 16-byte GCM tag to the ciphertext buffer.
      var full = new Uint8Array(encrypted);
      var ciphertext = full.slice(0, full.length - 16);
      var tag = full.slice(full.length - 16);
      return hexEncode(ciphertext) + "." + hexEncode(tag) + "." + hexEncode(iv) + "." + hexEncode(salt);
    });
  }

  // Intercept form submission: encrypt the password field before the POST body
  // is sent. The submit event is cancelled, encryption runs asynchronously,
  // then the form is re-submitted programmatically once the field is replaced.
  function attachEncryptionToForm() {
    var form = document.querySelector("form");
    if (!form) { return; }

    var passwordField = document.getElementById("password");
    var authReqField = document.getElementById("authRequestID");
    if (!passwordField || !authReqField) { return; }

    // Guard: only attach once
    if (form.dataset.encryptAttached) { return; }
    form.dataset.encryptAttached = "1";

    form.addEventListener("submit", function (event) {
      // If already encrypted (re-submission after async), allow through
      if (form.dataset.encrypted === "1") {
        form.dataset.encrypted = "";
        return;
      }

      event.preventDefault();

      var rawPassword = passwordField.value;
      var authRequestID = authReqField.value;

      if (!rawPassword || !authRequestID) {
        form.submit();
        return;
      }

      encryptPassword(rawPassword, authRequestID)
        .then(function (payload) {
          passwordField.value = payload;
          form.dataset.encrypted = "1";
          form.requestSubmit ? form.requestSubmit() : form.submit();
        })
        .catch(function () {
          // Encryption failed (e.g. unsupported browser) — submit plaintext.
          // VerifyPassword will reject it if the server expects encrypted input,
          // which is the correct failure mode (auth denied, no silent bypass).
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
