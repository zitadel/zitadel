// Preload script that makes Node.js read CA certificates from SSL_CERT_DIR
// the same way Go's crypto/x509 does: every regular file in the directory is
// read regardless of filename.  OpenSSL (used by Node when --use-openssl-ca
// is set) only reads files named with a subject-hash (e.g. "9d66eef0.0"),
// silently ignoring plain names like "corp-ca.crt".  This script closes that
// gap so operators get identical behaviour across the Go backend and the
// Node.js login container.
//
// Usage: node --require /app/load-ssl-cert-dir.cjs ...
//
// The script is a no-op when SSL_CERT_DIR is not set.

"use strict";

const tls = require("tls");
const fs = require("fs");
const path = require("path");

const sslCertDir = process.env.SSL_CERT_DIR;
if (!sslCertDir) {
  return; // nothing to do — system trust via SSL_CERT_FILE is sufficient
}

// Go splits on ":" and walks every directory in order.
const dirs = sslCertDir.split(":").filter(Boolean);
const certs = [];

for (const dir of dirs) {
  let entries;
  try {
    entries = fs.readdirSync(dir);
  } catch {
    // Directory doesn't exist or isn't readable — skip, like Go does.
    continue;
  }

  for (const entry of entries) {
    const filePath = path.join(dir, entry);

    try {
      const stat = fs.lstatSync(filePath);

      // Go's readUniqueDirectoryEntries skips symlinks whose target is a
      // bare filename (no slash), because those are the c_rehash hash links
      // pointing at files in the same directory.  This avoids loading the
      // same certificate twice.
      if (stat.isSymbolicLink()) {
        const target = fs.readlinkSync(filePath);
        if (!target.includes("/")) {
          continue;
        }
      }

      if (!stat.isFile() && !stat.isSymbolicLink()) {
        continue;
      }

      const content = fs.readFileSync(filePath, "utf8");
      if (content.includes("-----BEGIN CERTIFICATE-----")) {
        certs.push(content);
      }
    } catch {
      // Unreadable file — skip silently, matching Go behaviour.
      continue;
    }
  }
}

if (certs.length === 0) {
  return;
}

const origCreateSecureContext = tls.createSecureContext;

tls.createSecureContext = function (options) {
  const ctx = origCreateSecureContext.call(this, options);

  // Only inject directory CAs when the caller did not pass an explicit "ca"
  // option.  This matches Go's SystemCertPool() semantics: the extra certs
  // are part of the system pool, not a forced override.
  if (!options || !options.ca) {
    for (let i = 0; i < certs.length; i++) {
      try {
        ctx.context.addCACert(certs[i]);
      } catch (err) {
        // Log once per bad cert, then remove it so subsequent contexts
        // don't pay the cost of retrying a known-bad certificate.
        console.error("load-ssl-cert-dir: skipping malformed certificate at index " + i + ": " + err.message);
        certs.splice(i, 1);
        i--;
      }
    }
  }

  return ctx;
};
