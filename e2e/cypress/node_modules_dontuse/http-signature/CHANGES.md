# node-http-signature changelog

## not yet released

(nothing yet)

## 1.3.6

* Update jsprim due to vulnerability in json-schema (#123)

## 1.3.5

* Add keyPassphrase option to signer (#115)
* Add support for created and expires values (#110)

## 1.3.4

* Fix breakage in v1.3.3 with the setting of the "algorithm" field in the
  Authorization header (#102)

## 1.3.3

**Bad release. Use 1.3.4.**

* Add support for an opaque param in the Authorization header (#101)
* Add support for adding the keyId and algorithm params into the signing string (#100)

## 1.3.2

* Allow Buffers to be used for verifyHMAC (#98)

## 1.3.1

* Fix node 0.10 usage (#90)

## 1.3.0

**Known issue:** This release broken http-signature with node 0.10.

* Bump dependency `sshpk`
* Add `Signature` header support (#83)

## 1.2.0

* Bump dependency `assert-plus`
* Add ability to pass a custom header name
* Replaced dependency `node-uuid` with `uuid`

## 1.1.1

* Version of dependency `assert-plus` updated: old version was missing
  some license information
* Corrected examples in `http_signing.md`, added auto-tests to
  automatically validate these examples

## 1.1.0

* Bump version of `sshpk` dependency, remove peerDependency on it since
  it now supports exchanging objects between multiple versions of itself
  where possible

## 1.0.2

* Bump min version of `jsprim` dependency, to include fixes for using
  http-signature with `browserify`

## 1.0.1

* Bump minimum version of `sshpk` dependency, to include fixes for
  whitespace tolerance in key parsing.

## 1.0.0

* First semver release.
* #36: Ensure verifySignature does not leak useful timing information
* #42: Bring the library up to the latest version of the spec (including the
       request-target changes)
* Support for ECDSA keys and signatures.
* Now uses `sshpk` for key parsing, validation and conversion.
* Fixes for #21, #47, #39 and compatibility with node 0.8

## 0.11.0

* Split up HMAC and Signature verification to avoid vulnerabilities where a
  key intended for use with one can be validated against the other method
  instead.

## 0.10.2

* Updated versions of most dependencies.
* Utility functions exported for PEM => SSH-RSA conversion.
* Improvements to tests and examples.
