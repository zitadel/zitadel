# SAML Federated Logout Implementation Guide

This document captures the detailed technical implementation of SAML Federated Logout in Zitadel V2, specifically focusing on the fix for "Invalid requester" errors with Keycloak and the robust parameter ordering logic.

## 1. Overview

SAML Federated Logout allows a Service Provider (SP - Zitadel) to notify an Identity Provider (IdP - e.g., Keycloak, Auth0) that a user has logged out, ensuring a consistent session state across applications.

The process roughly follows:
1. User clicks logout in Zitadel.
2. Zitadel terminates its local session.
3. Zitadel checks if the session was established via an external IdP.
4. If yes, and if Federated Logout is enabled, Zitadel generates a SAML `LogoutRequest`.
5. The user's browser is redirected to the IdP's Single Logout Service (SLO) endpoint.
6. The IdP terminates its session and redirects back to Zitadel (SAML Response).
7. Zitadel finalizes the logout.

## 2. Core Components

### 2.1 Command Layer
- **File**: `internal/command/session_logout_federated.go`
- **Function**: `StartFederatedLogout`
- **Responsibility**: Orchestrates the logout process. It retrieves the session, identifies the connected IdP, fetches the IdP configuration, and generates the SAML request.

### 2.2 Data Fetching & Decoupling
To ensure testability, the command orchestration is decoupled from direct database access:
- **Interface**: `FederatedLogoutDataFetcher`
  - Abstracts `IDPUserLinks` (finding which IdP the user used)
  - Abstracts `IDPTemplateByID` (fetching IdP SAML config)

- **Interface**: `FederatedLogoutEventstore`
  - Abstracts access to the event store for finding the `SessionIndex` (via `IDPIntent` events) and pushing new logout events.

### 2.3 SAML Request Generation
- **Function**: `generateSAMLLogoutRequest`
- **Key Logic**:
  - Uses `crewjam/saml` library to create the base `LogoutRequest` XML.
  - **CRITICAL**: Manually handles the HTTP-Redirect binding signature to enforce strict parameter ordering.

## 3. Strict Parameter Ordering (The Key Fix)

### The Problem
The SAML 2.0 specification (Bindings, Section 3.4.4.1) mandates that for HTTP-Redirect binding, the signature must be calculated over the query string parameters in a strictly defined order:
1. `SAMLRequest`
2. `RelayState` (if present)
3. `SigAlg`

The standard Go `net/url` library's `Values.Encode()` method sorts query parameters alphabetically by key. This results in `RelayState` appearing *before* `SAMLRequest`, violating the specification. Many IdPs (like Keycloak) reject such requests with "Signature validation failed" or "Invalid requester".

### The Solution
We bypass `Values.Encode()` and manually construct the query string using `bytes.Buffer`:

```go
var queryBuf bytes.Buffer
queryBuf.WriteString("SAMLRequest=")
queryBuf.WriteString(url.QueryEscape(samlRequest))

if relayState != "" {
    queryBuf.WriteString("&RelayState=")
    queryBuf.WriteString(url.QueryEscape(relayState))
}

queryBuf.WriteString("&SigAlg=")
queryBuf.WriteString(url.QueryEscape(sp.ServiceProvider.SignatureMethod))

query := queryBuf.String()
```

This ensures the string to be signed (`query`) exactly matches the order expected by the IdP. We then sign this string and append the `&Signature=...` at the end.

## 4. Session Index Handling

The `SessionIndex` is a unique identifier tied to the authenticated session on the IdP side. It is required by some IdPs to process logout correctly.
- **Retrieval**: We traverse the `idpintent` events in the Eventstore to find the most recent `SAMLSucceededEvent`.
- **Decryption**: The SAML Assertion inside the event is decrypted.
- **Extraction**: The `SessionIndex` is parsed from the assertion's `AuthnStatement`.

If found, it is added to the `LogoutRequest`. If not, the logout proceeds with just the `NameID`.

## 5. Testing

### 5.1 Parameter Ordering Test
- **File**: `internal/command/session_logout_federated_ordering_test.go`
- **Purpose**: A dedicated, black-box style test that generates a `LogoutRequest` and parses the resulting URL to verify:
  1. The existence of all 4 parameters (`SAMLRequest`, `RelayState`, `SigAlg`, `Signature`).
  2. The exact position of each parameter in the raw query string.
  3. The validity of the cryptographic signature against the public key.

### 5.2 Unit Tests
- **File**: `internal/command/session_logout_federated_test.go`
- **Purpose**: Verifies usage of the decoupled interfaces, error handling (e.g., missing IdP config), and ensures the flow proceeds correctly under various data states.

## 6. How to Verify
1. Configure a SAML IdP (e.g., Keycloak).
2. Login with a user via this IdP.
3. Trigger a logout.
4. Observe the network traffic or URL. You should see a redirect to the IdP with a URL like:
   `https://idp.com/slo?SAMLRequest=...&RelayState=...&SigAlg=...&Signature=...`
5. The IdP should accept the request and show a logout confirmation or redirect back.

## 7. Common Pitfalls
- **NameID Format**: Some IdPs are strict about `NameQualifier` and `SPNameQualifier`. We explicitly clear these fields in the `LogoutRequest` unless specifically needed, to ensure broader compatibility.
- **Binding**: Currently, we default to `HTTP-Redirect` for SLO because `HTTP-POST` requires rendering an HTML form, which is complex in the current backend command structure.
