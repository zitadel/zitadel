---
title: "Opaque Tokens in Zitadel: Enhancing Application Security"
sidebar_label: Opaque tokens
---

In the context of application security, robust authentication mechanisms are essential for safeguarding sensitive data and ensuring user trust.
Opaque tokens, the default token type within the ZITADEL platform, play a crucial role in bolstering security measures.
This documentation elucidates the principles behind opaque tokens, their implementation within ZITADEL, and their advantages over alternative token types.

## What are Opaque Tokens?

Opaque tokens are a type of access token utilized in authentication processes, particularly within OAuth 2.0 and OpenID Connect (OIDC) frameworks. Unlike self-contained tokens like [JSON Web Tokens (JWT)](https://datatracker.ietf.org/doc/html/rfc7519), opaque tokens do not divulge user information directly.
Instead, they serve as opaque references to session data stored securely on the authorization server.

## Authentication Workflow with Opaque Tokens

1. **Token Generation**: When a user initiates an authentication process within an application integrated with ZITADEL, the authentication server generates a unique opaque token associated with the user's session.

2. **Token Presentation**: The generated opaque token is provided to the client, which subsequently presents it during requests to access protected resources within the application.

3. **Token Verification**: Upon receiving the opaque token, the application server interacts with the authorization server to validate its authenticity and retrieve detailed information about the user's session. This process ensures the integrity of the authentication flow and verifies the user's permissions to access requested resources.

## Benefits of Opaque Tokens in ZITADEL

1. **Reduced Token Exposure**: Opaque tokens mitigate the risk of token exposure since they do not contain sensitive user information directly. This reduces the likelihood of token-based attacks and enhances overall security posture.

2. **Enhanced Server-side Control**: With opaque tokens, validation occurs server-side, granting administrators greater control over authentication flows and access policies. This centralized approach facilitates comprehensive monitoring and enforcement of security measures, including server-side single-logout across all applications.

3. **Protection Against Token Tampering**: Opaque tokens prevent unauthorized manipulation of token contents, thereby ensuring the integrity and authenticity of authentication processes. This protection against token tampering further strengthens the security of applications integrated with ZITADEL.

## Opaque Tokens vs. JWT Tokens

When it comes to implementing authentication and authorization mechanisms within applications, developers often face the choice between different types of tokens, each with its own set of characteristics and advantages.
Two common types of tokens used in authentication protocols are opaque tokens and JSON Web Tokens (JWT).

### Structure

- **Opaque Tokens**: Opaque tokens are essentially references or pointers to information stored on the authorization server. They do not contain any meaningful user data within the token itself. Instead, they typically consist of a unique identifier (e.g., a session ID or database key) that allows the server to look up the associated session information.

- **JWT Tokens**: JSON Web Tokens, on the other hand, are self-contained tokens that contain user information in a JSON format. JWTs consist of three base64-encoded sections: header, payload, and signature. The payload contains claims or assertions about the user (e.g., user ID, roles, expiration time) that are digitally signed to ensure integrity.

### Token verification

- **Opaque Tokens**: Verifying opaque tokens requires interaction with the authorization server. When a client presents an opaque token to a resource server, the resource server sends the token to the authorization server for validation. The authorization server then checks the token's validity and returns the associated user information if the token is valid.

- **JWT Tokens**: JWT tokens can be verified locally by clients without needing to communicate with the authorization server. Clients can validate JWT signatures using public keys or shared secrets obtained during the token issuance process. This decentralized verification process can be faster and more scalable but requires securely distributing and managing keys.

### Token security and size

- **Opaque Tokens**: Since opaque tokens do not contain user information, they are inherently more secure in terms of protecting sensitive data. However, the reliance on server-side validation means there is an overhead associated with each token verification request, which can impact performance in high-throughput scenarios.

- **JWT Tokens**: JWT tokens contain user information within the token itself, which can be convenient for clients as it eliminates the need for frequent interactions with the authorization server. However, this also means that JWT tokens can potentially expose sensitive information if not handled and secured properly. Additionally, JWT tokens tend to be larger compared to opaque tokens due to the encoded payload.

### Use cases and trade-offs

- **Opaque Tokens**: Opaque tokens are well-suited for scenarios where security and confidentiality are top priorities, such as handling highly sensitive user data or complying with strict privacy regulations. They are particularly advantageous in distributed systems where centralized control over authentication and access policies is desired.

- **JWT Tokens**: JWT tokens are often preferred in scenarios where performance and scalability are critical, such as microservices architectures or API-based applications. The ability to verify tokens locally can reduce latency and minimize dependencies on external services. However, developers must carefully consider the implications of including sensitive information in JWT payloads and implement appropriate security measures.

## Conclusion

In conclusion, opaque tokens represent a foundational component in fortifying application security within ZITADEL.
By leveraging opaque tokens, organizations can establish robust authentication mechanisms, mitigate security risks, and maintain stringent control over access policies.
As organizations navigate the complex landscape of application security, integrating technologies such as opaque tokens becomes imperative for safeguarding sensitive data and fostering user trust.

Notes: 

- Read more about the differences in our [blog on JWT vs. Opaque tokens](https://zitadel.com/blog/jwt-vs-opaque-tokens)
- Learn how to use [token introspection](/docs/guides/integrate/token-introspection) to validate access tokens
- Decode, verify and generate valid JWT tokens with [jwt.io](https://jwt.io/)
