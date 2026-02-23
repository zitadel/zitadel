import forge from "node-forge";

/**
 * Generates a self-signed CA certificate and a server certificate signed by
 * that CA, returning both as PEM-encoded strings alongside their private keys.
 *
 * This is intended for test and development environments where a trusted
 * certificate chain is needed without relying on the openssl binary. The
 * server certificate includes Subject Alternative Names so that TLS clients
 * can verify the connection against the expected hostnames and IP addresses.
 */
export function generateCertificates({
  caCN = "Test CA",
  serverCN = "mock-zitadel",
  dns = ["mock-zitadel", "localhost"],
  ips = ["127.0.0.1"],
  days = 1,
} = {}) {
  function createCert({
    subject,
    issuer,
    publicKey,
    signingKey,
    extensions,
    serial = "01",
  }: {
    subject: forge.pki.CertificateField[];
    issuer: forge.pki.CertificateField[];
    publicKey: forge.pki.rsa.PublicKey;
    signingKey: forge.pki.rsa.PrivateKey;
    extensions?: forge.pki.CertificateField[];
    serial?: string;
  }) {
    const cert = forge.pki.createCertificate();
    cert.publicKey = publicKey;
    cert.serialNumber = serial;
    cert.validity.notBefore = new Date();
    cert.validity.notAfter = new Date(Date.now() + days * 86400000);
    cert.setSubject(subject);
    cert.setIssuer(issuer);
    if (extensions) {
      cert.setExtensions(extensions);
    }
    cert.sign(signingKey, forge.md.sha256.create());
    return forge.pki.certificateToPem(cert);
  }

  const caKeys = forge.pki.rsa.generateKeyPair(2048);
  const serverKeys = forge.pki.rsa.generateKeyPair(2048);
  const caSubject = [{ name: "commonName", value: caCN }];

  return {
    ca: {
      cert: createCert({
        subject: caSubject,
        issuer: caSubject,
        publicKey: caKeys.publicKey,
        signingKey: caKeys.privateKey,
        extensions: [{ name: "basicConstraints", cA: true }],
      }),
      key: forge.pki.privateKeyToPem(caKeys.privateKey),
    },
    server: {
      cert: createCert({
        subject: [{ name: "commonName", value: serverCN }],
        issuer: caSubject,
        publicKey: serverKeys.publicKey,
        signingKey: caKeys.privateKey,
        serial: "02",
        extensions: [
          {
            name: "subjectAltName",
            altNames: [
              ...dns.map(function (value) {
                return { type: 2, value };
              }),
              ...ips.map(function (ip) {
                return { type: 7, ip };
              }),
            ],
          },
        ],
      }),
      key: forge.pki.privateKeyToPem(serverKeys.privateKey),
    },
  };
}
