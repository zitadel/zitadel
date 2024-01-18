---
title: SAML Endpoints in ZITADEL
---

## SAML 2.0 metadata

The SAML Metadata is located within the issuer domain. This would give us $CUSTOM-DOMAIN/saml/v2/metadata.

This metadata contains all the information defined in the spec.

**Link to
spec.** [Metadata for the OASIS Security Assertion Markup Language (SAML) V2.0 – Errata Composite](https://www.oasis-open.org/committees/download.php/35391/sstc-saml-metadata-errata-2.0-wd-04-diff.pdf)

## Certificate endpoint

$CUSTOM-DOMAIN/saml/v2/certificate

The certificate endpoint provides the certificate which is used to sign the responses for download, for easier use with
different service providers which want the certificate separately instead of inside the metadata.

## SSO endpoint

$CUSTOM-DOMAIN/saml/v2/SSO

The SSO endpoint is the starting point for all initial user authentications. The user agent (browser) will be redirected
to this endpoint to authenticate the user.

Supported on this endpoint or currently `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect`
or `urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST` bindings.

**Link to
spec.** [Bindings for the OASIS Security Assertion Markup Language (SAML) V2.0 – Errata Composite](https://www.oasis-open.org/committees/download.php/35387/sstc-saml-bindings-errata-2.0-wd-05-diff.pdf)

### Required request parameters

| Parameter | Description                                                                                                                                                                         |
|---------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RelayState | (Optional) ID to associate the exchange with the original request.                                                                                                                  |
| SAMLRequest | The request made to the SAML IDP.  (base64 encoded)                                                                                                                                 |
| SigAlg | Algorithm used to sign the request, only if binding is 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect' as signature has to be provided es separate parameter. (base64 encoded) |
| Signature | Signature of the request as parameter with 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect' binding.  (base64 encoded)                                                          |

### Successful response

Depending on the content of the request the response comes back in the requested binding, but the content is the same.

| Parameter | Description                                                                                                                                                          |
|---------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| RelayState | ID to associate the exchange with the original request.                                                                                                              |
| SAMLResponse | The response form the SAML IDP.  (base64 encoded)                                                                                                                    |
| SigAlg | Algorithm used to sign the response, only if binding is 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect' as signature has to be provided es separate parameter.  (base64 encoded)  |
| Signature | Signature of the response as parameter with 'urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect' binding.  (base64 encoded)                                                            |

### Error response

Regardless of the error, the used http error code will be '200', which represents a successful request. Whereas the
response will contain a StatusCode include a message which provides more information if an error occurred.

**Link to
spec** [Assertions and Protocols for the OASIS Security Assertion Markup Language (SAML) V2.0 – Errata Composite](https://www.oasis-open.org/committees/download.php/35711/sstc-saml-core-errata-2.0-wd-06-diff.pdf)

## Custom attributes

Custom attributes are being inserted into SAML response if not already present.
Your app can use custom claims to handle more complex scenarios, such as restricting access based on these claims.

You can add custom attributes using the [complement SAMLresponse](/docs/apis/actions/customize-samlresponse) of the [actions feature](/docs/apis/actions/introduction).

Examples of Actions that result in custom attributes can be found in our [Marketplace for ZITADEL Actions](https://github.com/zitadel/actions).
