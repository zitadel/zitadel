/**
 * TLS Certificate Generation Utility
 *
 * Generates self-signed CA and server certificates using node-forge.
 * Useful for integration testing TLS/SSL connectivity with custom CAs.
 *
 * @module utils/tls
 */
import * as forge from "node-forge";
import * as fs from "node:fs";
import * as path from "node:path";

/**
 * Result of certificate generation containing both file paths and PEM contents.
 */
export interface GeneratedCertificates {
  /** Filesystem paths where certificates were written */
  paths: {
    /** Path to CA private key PEM file */
    caKey: string;
    /** Path to CA certificate PEM file */
    caCert: string;
    /** Path to server private key PEM file */
    serverKey: string;
    /** Path to server certificate PEM file */
    serverCert: string;
  };
  /** CA certificate and key in PEM format */
  ca: {
    /** CA certificate PEM string */
    cert: string;
    /** CA private key PEM string */
    key: string;
  };
  /** Server certificate and key in PEM format */
  server: {
    /** Server certificate PEM string */
    cert: string;
    /** Server private key PEM string */
    key: string;
  };
}

/**
 * Options for generating certificates.
 */
export interface GenerateCertificatesOptions {
  /** Directory where certificate files will be written (creates 'certs' subdirectory) */
  outputDir: string;
  /** Common Name (CN) for the CA certificate @default "Test CA" */
  caCommonName?: string;
  /** Common Name (CN) for the server certificate @default "localhost" */
  serverCommonName?: string;
  /** Subject Alternative Names for the server certificate @default ["localhost"] */
  serverAltNames?: string[];
  /** Certificate validity period in days @default 1 */
  validityDays?: number;
}

/**
 * Generates a self-signed CA certificate and a server certificate signed by that CA.
 *
 * Creates a certificate chain suitable for testing TLS connections:
 * - A root CA certificate that can be trusted by clients
 * - A server certificate signed by the CA with the specified SANs
 *
 * Both certificates are written to disk and returned as PEM strings.
 *
 * @param options - Configuration options for certificate generation
 * @returns Generated certificates with both file paths and PEM contents
 *
 * @example
 * ```typescript
 * const certs = generateCertificates({
 *   outputDir: '/tmp/test-certs',
 *   serverCommonName: 'api.example.com',
 *   serverAltNames: ['api.example.com', 'localhost', '127.0.0.1'],
 * });
 *
 * // Use in Node.js TLS server
 * const server = https.createServer({
 *   key: certs.server.key,
 *   cert: certs.server.cert,
 * });
 *
 * // Mount CA cert in Docker for client trust
 * // volumes: - ${certs.paths.caCert}:/etc/ssl/certs/custom-ca.crt
 * ```
 */
export function generateCertificates(options: GenerateCertificatesOptions): GeneratedCertificates {
  const {
    outputDir,
    caCommonName = "Test CA",
    serverCommonName = "localhost",
    serverAltNames = ["localhost"],
    validityDays = 1,
  } = options;

  const certsDir = path.join(outputDir, "certs");
  fs.mkdirSync(certsDir, { recursive: true });

  const paths = {
    caKey: path.join(certsDir, "ca.key"),
    caCert: path.join(certsDir, "ca.crt"),
    serverKey: path.join(certsDir, "server.key"),
    serverCert: path.join(certsDir, "server.crt"),
  };

  console.log("[TLS] Generating CA certificate...");
  const caKeyPair = forge.pki.rsa.generateKeyPair(2048);
  const caCert = forge.pki.createCertificate();

  caCert.publicKey = caKeyPair.publicKey;
  caCert.serialNumber = "01";
  caCert.validity.notBefore = new Date();
  caCert.validity.notAfter = new Date();
  caCert.validity.notAfter.setDate(caCert.validity.notBefore.getDate() + validityDays);

  const caAttrs = [{ name: "commonName", value: caCommonName }];
  caCert.setSubject(caAttrs);
  caCert.setIssuer(caAttrs);

  caCert.setExtensions([
    { name: "basicConstraints", cA: true },
    { name: "keyUsage", keyCertSign: true, cRLSign: true },
  ]);

  caCert.sign(caKeyPair.privateKey, forge.md.sha256.create());

  const caKeyPem = forge.pki.privateKeyToPem(caKeyPair.privateKey);
  const caCertPem = forge.pki.certificateToPem(caCert);

  fs.writeFileSync(paths.caKey, caKeyPem);
  fs.writeFileSync(paths.caCert, caCertPem);

  console.log("[TLS] Generating server certificate...");
  const serverKeyPair = forge.pki.rsa.generateKeyPair(2048);
  const serverCert = forge.pki.createCertificate();

  serverCert.publicKey = serverKeyPair.publicKey;
  serverCert.serialNumber = "02";
  serverCert.validity.notBefore = new Date();
  serverCert.validity.notAfter = new Date();
  serverCert.validity.notAfter.setDate(serverCert.validity.notBefore.getDate() + validityDays);

  const serverAttrs = [{ name: "commonName", value: serverCommonName }];
  serverCert.setSubject(serverAttrs);
  serverCert.setIssuer(caAttrs);

  const altNames = serverAltNames.map((name) => {
    if (/^\d+\.\d+\.\d+\.\d+$/.test(name)) {
      return { type: 7, ip: name };
    }
    return { type: 2, value: name };
  });

  if (!serverAltNames.includes("127.0.0.1")) {
    altNames.push({ type: 7, ip: "127.0.0.1" });
  }

  serverCert.setExtensions([
    { name: "basicConstraints", cA: false },
    { name: "keyUsage", digitalSignature: true, keyEncipherment: true },
    { name: "extKeyUsage", serverAuth: true },
    { name: "subjectAltName", altNames },
  ]);

  serverCert.sign(caKeyPair.privateKey, forge.md.sha256.create());

  const serverKeyPem = forge.pki.privateKeyToPem(serverKeyPair.privateKey);
  const serverCertPem = forge.pki.certificateToPem(serverCert);

  fs.writeFileSync(paths.serverKey, serverKeyPem);
  fs.writeFileSync(paths.serverCert, serverCertPem);

  console.log(`[TLS] CA certificate: ${paths.caCert}`);
  console.log(`[TLS] Server certificate: ${paths.serverCert}`);

  return {
    paths,
    ca: { cert: caCertPem, key: caKeyPem },
    server: { cert: serverCertPem, key: serverKeyPem },
  };
}
