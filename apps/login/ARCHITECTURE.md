# Login Application Architecture

## Table of Contents

- [Overview](#overview)
- [Technology Stack](#technology-stack)
- [Application Structure](#application-structure)
- [Session Management](#session-management)
- [Authentication Flows](#authentication-flows)
- [Middleware Architecture](#middleware-architecture)
- [Security Considerations](#security-considerations)
- [Multi-Factor Authentication](#multi-factor-authentication)
- [Identity Provider Integration](#identity-provider-integration)
- [Deployment Modes](#deployment-modes)
- [Next.js Implementation](#nextjs-implementation)
  - [Server Actions](#server-actions)
  - [Internationalization](#internationalization)
  - [Theming System](#theming-system)
  - [Performance Optimizations](#performance-optimizations)
- [Error Handling](#error-handling)
- [Monitoring and Debugging](#monitoring-and-debugging)
- [Cross-Layer Security Analysis](#cross-layer-security-analysis)
- [Conclusion](#conclusion)

## Overview

The ZITADEL Login application is a Next.js 15-based authentication frontend that establishes secure sessions for multiple authentication protocols including OIDC (OpenID Connect), SAML, and device authorization flows. It serves as the primary user-facing authentication interface for ZITADEL's identity platform.

### Key Features

- **Multiple Authentication Methods**: Password, passkeys (WebAuthn), OTP (Email/SMS/TOTP), U2F
- **Protocol Support**: OIDC, SAML, LDAP connections, Device Code Authorization Flow
- **User Management**: Password reset, email verification, user registration, invite flows
- **Multi-Factor Authentication**: TOTP, OTP (Email/SMS), U2F, passkeys
- **Flexible Deployment**: Supports both multi-tenant and self-hosted single-instance deployments
- **Theming Support**: Environment-driven theming with responsive layouts

## Technology Stack

### Core Framework

```json
{
  "next": "15.4.0-canary.86",
  "react": "19.1.0",
  "react-dom": "19.1.0"
}
```

**Next.js 15** provides:

- **App Router**: File-system based routing with React Server Components
- **Server Actions**: Type-safe server-side mutations without API routes
- **Dynamic IO**: Improved data fetching and caching capabilities
- **Middleware**: Edge runtime for request interception and routing

### Communication Layer

```typescript
// gRPC Communication
"nice-grpc": "2.0.1"  // Type-safe gRPC client for Node.js

// Protocol Buffers
"@zitadel/proto": "workspace:*"
"@zitadel/client": "workspace:*"
```

The application communicates with ZITADEL's backend exclusively via **gRPC** using protocol buffers for type-safe, efficient communication.

### Internationalization

```typescript
"next-intl": "^3.25.1"  // i18n for Next.js App Router
```

Provides server-safe internationalization with support for organization-specific translations.

### UI Components

```typescript
"@headlessui/react": "^2.1.9"  // Unstyled, accessible UI components
"@radix-ui/react-tooltip": "^1.2.7"  // Accessible tooltip primitives
"tailwindcss": "3.4.14"  // Utility-first CSS framework
```

## Application Structure

### Directory Layout

```
apps/login/
├── src/
│   ├── app/                    # Next.js App Router pages
│   │   ├── (login)/           # Grouped authentication routes
│   │   │   ├── loginname/     # Username/email entry
│   │   │   ├── password/      # Password authentication
│   │   │   ├── passkey/       # WebAuthn passkey flows
│   │   │   ├── mfa/           # Multi-factor authentication
│   │   │   ├── otp/           # One-time password (TOTP/Email/SMS)
│   │   │   ├── u2f/           # Universal 2nd Factor
│   │   │   ├── idp/           # Identity Provider flows
│   │   │   ├── register/      # User registration
│   │   │   ├── verify/        # Email verification
│   │   │   ├── device/        # Device authorization flow
│   │   │   └── accounts/      # Account selection
│   │   ├── login/             # OIDC/SAML login initiation
│   │   ├── security/          # Security settings endpoint
│   │   └── healthy/           # Health check
│   ├── lib/                   # Core business logic
│   │   ├── server/            # Server Actions
│   │   │   ├── session.ts     # Session management
│   │   │   ├── password.ts    # Password operations
│   │   │   ├── passkeys.ts    # Passkey registration/verification
│   │   │   ├── idp.ts         # IDP integration
│   │   │   ├── device.ts      # Device authorization
│   │   │   ├── cookie.ts      # Cookie operations
│   │   │   └── auth-flow.ts   # Flow completion
│   │   ├── zitadel.ts         # gRPC service clients
│   │   ├── session.ts         # Session validation logic
│   │   ├── cookies.ts         # Cookie management
│   │   ├── client.ts          # Client-side helpers
│   │   ├── service-url.ts     # Multi-tenancy URL resolution
│   │   └── verify-helper.ts   # Verification helpers
│   ├── components/            # React components
│   ├── i18n/                  # Internationalization
│   ├── middleware.ts          # Next.js middleware
│   └── styles/                # Global styles
├── constants/
│   └── csp.js                 # Content Security Policy
├── locales/                   # Translation files
├── public/                    # Static assets
└── scripts/
    └── entrypoint.sh          # Container entrypoint
```

### Routing Strategy

The application uses Next.js 15's **App Router** with:

- **Route Groups**: `(login)` groups authentication routes without affecting URL structure
- **Server Components**: Default for all pages, enabling server-side rendering
- **Server Actions**: Direct server-side mutations marked with `"use server"`
- **Dynamic Routes**: Parameter-based routing for flexible flows

## Session Management

### Architecture Overview

Session management is the cornerstone of the login application's security model. It operates entirely **server-side** with HTTP-only cookies to prevent client-side access and XSS attacks.

### Session Flow

```
┌─────────────┐
│   Browser   │
└──────┬──────┘
       │ 1. Authentication Request
       ▼
┌─────────────────────┐
│  Server Action      │
│  (e.g., sendPassword)│
└──────┬──────────────┘
       │ 2. Create/Update Session (gRPC)
       ▼
┌─────────────────────┐
│  ZITADEL Backend    │
│  (Session Service)  │
└──────┬──────────────┘
       │ 3. Session Token + Metadata
       ▼
┌─────────────────────┐
│  Cookie Management  │
│  (HTTP-only)        │
└──────┬──────────────┘
       │ 4. Set-Cookie Header
       ▼
┌─────────────┐
│   Browser   │
│  (Cookie    │
│   Storage)  │
└─────────────┘
```

### Session Storage Architecture

#### Cookie Structure

Sessions are stored as a **JSON array** in an HTTP-only cookie named `"sessions"`:

```typescript
type SessionCookie = {
  id: string; // ZITADEL session ID
  token: string; // Session token (opaque)
  loginName: string; // User's login name
  organization?: string; // Organization ID (optional)
  creationTs: string; // Creation timestamp
  expirationTs: string; // Expiration timestamp
  changeTs: string; // Last change timestamp
  requestId?: string; // OIDC/SAML request ID (if linked)
};
```

**Example cookie value**:

```json
[
  {
    "id": "session_abc123",
    "token": "token_xyz789",
    "loginName": "user@example.com",
    "organization": "org_456",
    "creationTs": "1699564800000",
    "expirationTs": "1699651200000",
    "changeTs": "1699564800000",
    "requestId": "oidc_request_123"
  }
]
```

#### Cookie Configuration

```typescript
// From src/lib/cookies.ts
const cookieConfig = {
  name: "sessions",
  httpOnly: true, // Prevents JavaScript access
  path: "/",
  sameSite: iFrameEnabled ? "none" : "lax",
  secure: process.env.NODE_ENV === "production",
};
```

**Security Features**:

1. **HTTP-only**: Cookies cannot be accessed via `document.cookie` or JavaScript
2. **Secure flag**: Transmitted only over HTTPS in production
3. **SameSite policy**:
   - `"lax"`: Default - allows cookies with top-level navigation
   - `"none"`: For iframe embedding (requires `secure: true`)

#### Multi-Session Support

The application supports **multiple concurrent sessions** for:

- Multiple organizations
- Different users on the same device
- Session switching without re-authentication

```typescript
// Add new session to cookie array
await addSessionToCookie({ session, iFrameEnabled });

// Retrieve specific session by login name
await getSessionCookieByLoginName({
  loginName,
  organization,
});

// Get all sessions
const sessions = await getAllSessions();
```

### Session Creation

Sessions are created via gRPC calls to ZITADEL's Session Service:

```typescript
// From src/lib/server/cookie.ts
export async function createSessionAndUpdateCookie(command: {
  checks: Checks; // Authentication factors
  requestId: string | undefined;
  lifetime?: Duration;
}): Promise<Session> {
  // 1. Create session via gRPC
  const createdSession = await createSessionFromChecks({
    serviceUrl,
    checks: command.checks,
    lifetime: sessionLifetime, // Default: 24 hours
  });

  // 2. Retrieve full session details
  const response = await getSession({
    serviceUrl,
    sessionId: createdSession.sessionId,
    sessionToken: createdSession.sessionToken,
  });

  // 3. Store session in HTTP-only cookie
  const sessionCookie: CustomCookieData = {
    id: createdSession.sessionId,
    token: createdSession.sessionToken,
    loginName: response.session.factors.user.loginName,
    // ... timestamps and metadata
  };

  await addSessionToCookie({
    session: sessionCookie,
    iFrameEnabled,
  });

  return response.session;
}
```

### Session Updates

Session updates occur when authentication factors are added:

```typescript
// From src/lib/server/session.ts
export async function updateSession(options: {
  loginName?: string;
  sessionId?: string;
  checks?: Checks; // New authentication factor
  challenges?: RequestChallenges;
  lifetime?: Duration;
}) {
  // 1. Find existing session cookie
  const recentSession = sessionId
    ? await getSessionCookieById({ sessionId })
    : await getSessionCookieByLoginName({ loginName, organization });

  // 2. Update session via gRPC
  const session = await setSessionAndUpdateCookie({
    recentCookie: recentSession,
    checks,
    challenges,
    lifetime,
  });

  // 3. Return updated session metadata
  return {
    sessionId: session.id,
    factors: session.factors,
    challenges: session.challenges,
  };
}
```

### Session Validation

Session validation ensures that sessions meet security policies:

```typescript
// From src/lib/session.ts
export async function isSessionValid({ serviceUrl, session }: { serviceUrl: string; session: Session }): Promise<boolean> {
  // 1. Check session has user
  if (!session.factors?.user) return false;

  // 2. Check expiration
  const stillValid = session.expirationDate ? timestampDate(session.expirationDate) > new Date() : true;
  if (!stillValid) return false;

  // 3. Verify authentication factors
  const validPassword = !!session.factors.password?.verifiedAt;
  const validPasskey = !!session.factors.webAuthN?.verifiedAt;
  const validIDP = !!session.factors.intent?.verifiedAt;

  if (!(validPassword || validPasskey || validIDP)) return false;

  // 4. Check MFA requirements
  const loginSettings = await getLoginSettings({
    serviceUrl,
    organization: session.factors.user.organizationId,
  });

  const isMfaRequired = shouldEnforceMFA(session, loginSettings);

  if (isMfaRequired) {
    const authMethods = await listAuthenticationMethodTypes({
      serviceUrl,
      userId: session.factors.user.id,
    });

    // Verify MFA factor is present and verified
    const mfaValid = checkMfaFactors(session, authMethods);
    if (!mfaValid) return false;
  }

  // 5. Optional: Email verification check
  if (process.env.EMAIL_VERIFICATION === "true") {
    const user = await getUserByID({
      serviceUrl,
      userId: session.factors.user.id,
    });

    const humanUser = user?.user?.type.case === "human" ? user.user.type.value : undefined;

    if (humanUser && !humanUser.email?.isVerified) {
      return false;
    }
  }

  return true;
}
```

### Session Factors

Sessions track multiple authentication factors:

```typescript
type SessionFactors = {
  user?: {
    id: string;
    loginName: string;
    organizationId?: string;
  };
  password?: {
    verifiedAt?: Timestamp;
  };
  webAuthN?: {
    verifiedAt?: Timestamp;
    userVerified?: boolean; // For passkeys
  };
  intent?: {
    verifiedAt?: Timestamp; // For IDP authentication
  };
  totp?: {
    verifiedAt?: Timestamp; // Time-based OTP
  };
  otpEmail?: {
    verifiedAt?: Timestamp;
  };
  otpSms?: {
    verifiedAt?: Timestamp;
  };
};
```

### Session Lifetime Management

Different authentication methods have different lifetimes:

```typescript
// Password check
lifetime: loginSettings?.passwordCheckLifetime

// Passkey/U2F check
lifetime: loginSettings?.multiFactorCheckLifetime

// OTP (Email/SMS) check
lifetime: loginSettings?.secondFactorCheckLifetime

// IDP check
lifetime: loginSettings?.externalLoginCheckLifetime

// Default fallback
lifetime: { seconds: BigInt(60 * 60 * 24), nanos: 0 }  // 24 hours
```

### Cookie Size Management

The cookie has a **maximum size limit** of 2048 bytes:

```typescript
const MAX_COOKIE_SIZE = 2048;

// When adding session, check size
const temp = [...currentSessions, session];
if (JSON.stringify(temp).length >= MAX_COOKIE_SIZE) {
  // Replace oldest session with new one
  currentSessions = [session].concat(currentSessions.slice(1));
}
```

### Session Cleanup

Expired sessions are automatically cleaned up:

```typescript
export async function getAllSessions<T>(cleanup: boolean = false): Promise<SessionCookie<T>[]> {
  const sessions = JSON.parse(cookieValue);

  if (cleanup) {
    const now = new Date();
    return sessions.filter((session) =>
      session.expirationTs ? timestampDate(timestampFromMs(Number(session.expirationTs))) > now : true,
    );
  }

  return sessions;
}
```

## Authentication Flows

### 1. Password Authentication

```
User → /loginname → /password → Session Created → [MFA Check] → Redirect
```

**Implementation**:

```typescript
// src/lib/server/password.ts
export async function sendPassword(command: {
  loginName: string;
  organization?: string;
  checks: Checks;
  requestId?: string;
}): Promise<{ error: string } | { redirect: string }> {
  // 1. Find or create session
  let sessionCookie = await getSessionCookieByLoginName({
    loginName: command.loginName,
    organization: command.organization,
  });

  // 2. Validate password via session update
  try {
    session = await setSessionAndUpdateCookie({
      recentCookie: sessionCookie,
      checks: command.checks, // Contains password
      lifetime: loginSettings.passwordCheckLifetime,
    });
  } catch (error) {
    // Handle lockout attempts
    if ("failedAttempts" in error) {
      const lockoutSettings = await getLockoutSettings({
        serviceUrl,
        orgId: command.organization,
      });
      // Return error with attempt count
    }
  }

  // 3. Check password expiry
  const passwordChangedCheck = checkPasswordChangeRequired(
    expirySettings,
    session,
    humanUser,
    command.organization,
    command.requestId,
  );

  // 4. Check email verification
  const emailVerificationCheck = checkEmailVerification(session, humanUser, command.organization, command.requestId);

  // 5. Check MFA requirements
  const authMethods = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: session.factors.user.id,
  });

  const mfaFactorCheck = await checkMFAFactors(
    serviceUrl,
    session,
    loginSettings,
    authMethods,
    command.organization,
    command.requestId,
  );

  // 6. Complete flow or redirect
  if (command.requestId) {
    // OIDC/SAML flow
    return await completeFlowOrGetUrl(
      {
        sessionId: session.id,
        requestId: command.requestId,
        organization: command.organization,
      },
      loginSettings?.defaultRedirectUri,
    );
  }

  // Regular flow
  return await completeFlowOrGetUrl(
    {
      loginName: session.factors.user.loginName,
      organization: session.factors.user.organizationId,
    },
    loginSettings?.defaultRedirectUri,
  );
}
```

### 2. Passkey Authentication

```
User → /loginname → /passkey → WebAuthn Challenge → Verification → Session Created
```

**Passkey Flow**:

```typescript
// src/lib/server/passkeys.ts

// 1. Request WebAuthn challenge
export async function updateSession(options: {
  loginName: string;
  challenges?: RequestChallenges;  // WebAuthn challenge
}) {
  const session = await setSessionAndUpdateCookie({
    recentCookie,
    challenges: {
      webAuthN: {
        domain: hostname,
        userVerificationRequirement: "required"
      }
    },
  });

  return {
    sessionId: session.id,
    challenges: session.challenges,  // Return challenge for client
  };
}

// 2. Verify passkey response
export async function sendPasskey(command: {
  loginName: string;
  checks: Checks;  // Contains WebAuthn credential
  requestId?: string;
}) {
  const session = await setSessionAndUpdateCookie({
    recentCookie,
    checks: command.checks,  // WebAuthn verification
    lifetime: loginSettings.multiFactorCheckLifetime,
  });

  // Passkeys satisfy MFA requirements automatically
  return await completeFlowOrGetUrl({...});
}
```

### 3. Identity Provider (IDP) Authentication

```
User → /idp → External IDP → Callback → Create Session → [MFA Check] → Redirect
```

**IDP Flow**:

```typescript
// src/lib/server/idp.ts

// 1. Redirect to IDP
export async function redirectToIdp(formData: FormData) {
  const idpId = formData.get("id") as string;
  const provider = formData.get("provider") as string;

  // For LDAP, collect credentials first
  if (provider === "ldap") {
    redirect(`/idp/ldap?idpId=${idpId}&...`);
  }

  // For OAuth/OIDC IDPs
  const response = await startIdentityProviderFlow({
    serviceUrl,
    idpId,
    successUrl: `/idp/${provider}/process?...`,
    failureUrl: `/idp/${provider}/failure?...`,
  });

  redirect(response.redirect);  // Redirect to IDP
}

// 2. Handle IDP callback
export async function createNewSessionFromIdpIntent(command: {
  userId: string;
  idpIntent: {
    idpIntentId: string;
    idpIntentToken: string;
  };
  requestId?: string;
}) {
  // Create session with IDP intent factor
  const session = await createSessionForIdpAndUpdateCookie({
    userId: command.userId,
    idpIntent: command.idpIntent,
    lifetime: loginSettings.externalLoginCheckLifetime,
  });

  // Check email verification
  const emailCheck = checkEmailVerification(session, humanUser, ...);

  // Check MFA (if forceMfa is enabled, not forceMfaLocalOnly)
  const authMethods = await listAuthenticationMethodTypes({
    serviceUrl,
    userId: session.factors.user.id,
  });

  const mfaCheck = await checkMFAFactors(
    serviceUrl,
    session,
    loginSettings,
    authMethods || [],
    command.organization,
    command.requestId,
  );

  return completeFlowOrGetUrl({...});
}
```

### 4. OIDC/SAML Flow Initiation

OIDC and SAML authentication flows are initiated when external applications redirect users to ZITADEL's authorization endpoints. The login UI then handles the authentication process and completes the flow.

#### Flow Initiation Overview

```
┌─────────────────────┐
│   Application       │
│   (Relying Party)   │
└──────────┬──────────┘
           │ 1. Redirect to authorization endpoint
           ▼
┌─────────────────────────────────────────────┐
│   ZITADEL Backend                           │
│   /oauth/v2/authorize?client_id=...        │
│   /saml/v2/SSO                              │
│                                             │
│   • Validate client credentials             │
│   • Create AuthRequest (V2_xxx)             │
│   • Store request with parameters           │
└──────────┬──────────────────────────────────┘
           │ 2. Redirect to login UI
           ▼
┌─────────────────────────────────────────────┐
│   Login UI: /login Route                    │
│   /login?authRequest=V2_xxx (OIDC)          │
│   /login?samlRequest=xxx (SAML)             │
│                                             │
│   • Validate request ID                     │
│   • Load existing sessions                  │
│   • Fetch request details from backend      │
│   • Determine authentication path           │
└──────────┬──────────────────────────────────┘
           │ 3. Route based on state
           ├──────────┬───────────┬────────────┐
           │          │           │            │
    No session  Has session  Force login  Silent auth
           │          │           │            │
           ▼          ▼           ▼            ▼
    /loginname   Complete or  /loginname   Complete or
                 /accounts                   error
```

#### OIDC Authorization Endpoint

Applications initiate OIDC flows by redirecting to:

```
https://${ZITADEL_DOMAIN}/oauth/v2/authorize?
  client_id=xxx
  &redirect_uri=https://app.example.com/callback
  &response_type=code
  &scope=openid profile email
  &state=random_state
  &nonce=random_nonce
  &prompt=login        # Optional: force login
  &login_hint=user@example.com  # Optional: pre-fill username
```

#### SAML SSO Endpoint

SAML Service Providers initiate flows by sending SAML AuthnRequests to:

```
https://${ZITADEL_DOMAIN}/saml/v2/SSO
```

The SAML request contains the service provider's entity ID and assertion consumer URL. ZITADEL processes this similarly to OIDC, creating a SAML request ID and redirecting to:

```
https://login.example.com/login?samlRequest=${SAML_REQUEST_ID}
```

#### Login UI Entry Point

The login UI receives auth requests at `/src/app/login/route.ts`:

```typescript
export async function GET(request: NextRequest) {
  const searchParams = request.nextUrl.searchParams;

  // 1. Block React Server Component requests (internal Next.js)
  if (isRSCRequest(searchParams)) {
    return NextResponse.json({ error: "RSC requests not supported" }, { status: 400 });
  }

  // 2. Validate and extract request ID
  const requestId = validateAuthRequest(searchParams);
  // Converts: authRequest=V2_xxx → "oidc_V2_xxx"
  // Converts: samlRequest=xxx → "saml_xxx"

  if (!requestId) {
    return NextResponse.json({ error: "No valid authentication request found" }, { status: 400 });
  }

  // 3. Load existing sessions from HTTP-only cookies
  const sessionCookies = await getAllSessions();
  const sessions = await loadSessions({
    serviceUrl,
    ids: sessionCookies.map((s) => s.id),
  });

  // 4. Delegate to appropriate handler
  if (requestId.startsWith("oidc_")) {
    return handleOIDCFlowInitiation({
      serviceUrl,
      requestId,
      sessions,
      sessionCookies,
      request,
    });
  } else if (requestId.startsWith("saml_")) {
    return handleSAMLFlowInitiation({
      serviceUrl,
      requestId,
      sessions,
      sessionCookies,
      request,
    });
  }
}
```

#### OIDC Flow Initiation Handler

The OIDC handler (`/src/lib/server/flow-initiation.ts`) determines the authentication path:

```typescript
export async function handleOIDCFlowInitiation(params: FlowInitiationParams) {
  const { serviceUrl, requestId, sessions, sessionCookies, request } = params;

  // 1. Fetch auth request details from ZITADEL backend
  const { authRequest } = await getAuthRequest({
    serviceUrl,
    authRequestId: requestId.replace("oidc_", ""),
  });

  // 2. Extract organization and IDP from OIDC scopes
  let organization = "";
  let idpId = "";

  if (authRequest?.scope) {
    // Organization ID scope: urn:zitadel:iam:org:id:123456789
    const orgScope = authRequest.scope.find((s) => /urn:zitadel:iam:org:id:([0-9]+)/.test(s));
    if (orgScope) {
      organization = /urn:zitadel:iam:org:id:([0-9]+)/.exec(orgScope)?.[1] ?? "";
    }

    // Organization domain scope: urn:zitadel:iam:org:domain:primary:example.com
    const orgDomainScope = authRequest.scope.find((s) => /urn:zitadel:iam:org:domain:primary:(.+)/.test(s));
    if (orgDomainScope) {
      const orgDomain = /urn:zitadel:iam:org:domain:primary:(.+)/.exec(orgDomainScope)?.[1];
      // Resolve organization by domain
      const orgs = await getOrgsByDomain({ serviceUrl, domain: orgDomain });
      if (orgs.result?.length === 1) {
        organization = orgs.result[0].id;
      }
    }

    // IDP scope: urn:zitadel:iam:org:idp:id:xxx (bypass login UI, direct IDP redirect)
    const idpScope = authRequest.scope.find((s) => /urn:zitadel:iam:org:idp:id:(.+)/.test(s));
    if (idpScope) {
      idpId = /urn:zitadel:iam:org:idp:id:(.+)/.exec(idpScope)?.[1] ?? "";

      // Start IDP flow immediately
      const url = await startIdentityProviderFlow({
        serviceUrl,
        idpId,
        successUrl: `${origin}/idp/${provider}/process?requestId=${requestId}&organization=${organization}`,
        failureUrl: `${origin}/idp/${provider}/failure?requestId=${requestId}&organization=${organization}`,
      });

      return NextResponse.redirect(url);
    }
  }

  // 3. Handle OIDC prompt parameter
  if (authRequest.prompt.includes(Prompt.CREATE)) {
    // Registration flow requested
    const registerUrl = constructUrl(request, "/register");
    registerUrl.searchParams.set("requestId", requestId);
    if (organization) {
      registerUrl.searchParams.set("organization", organization);
    }
    return NextResponse.redirect(registerUrl);
  }

  if (authRequest.prompt.includes(Prompt.LOGIN)) {
    // Force login (ignore existing sessions)
    const loginNameUrl = constructUrl(request, "/loginname");
    loginNameUrl.searchParams.set("requestId", requestId);

    if (authRequest.loginHint) {
      loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
      // Optionally auto-submit if configured
    }
    if (organization) {
      loginNameUrl.searchParams.set("organization", organization);
    }

    return NextResponse.redirect(loginNameUrl);
  }

  if (authRequest.prompt.includes(Prompt.SELECT_ACCOUNT)) {
    // Show account selection page
    return NextResponse.redirect(constructUrl(request, "/accounts"));
  }

  if (authRequest.prompt.includes(Prompt.NONE)) {
    // Silent authentication - must use existing session or fail
    const selectedSession = await findValidSession({
      serviceUrl,
      sessions,
      authRequest,
    });

    if (!selectedSession) {
      return NextResponse.json({ error: "No active session found" }, { status: 400 });
    }

    const cookie = sessionCookies.find((c) => c.id === selectedSession.id);

    // Complete flow immediately without user interaction
    const { callbackUrl } = await createCallback({
      serviceUrl,
      req: create(CreateCallbackRequestSchema, {
        authRequestId: authRequest.id,
        callbackKind: {
          case: "session",
          value: { sessionId: cookie.id, sessionToken: cookie.token },
        },
      }),
    });

    return NextResponse.redirect(callbackUrl);
  }

  // 4. Default behavior: check for existing valid session
  if (sessions.length > 0) {
    const selectedSession = await findValidSession({
      serviceUrl,
      sessions,
      authRequest,
    });

    if (selectedSession) {
      // Valid session exists - complete flow
      const cookie = sessionCookies.find((c) => c.id === selectedSession.id);

      try {
        const { callbackUrl } = await createCallback({
          serviceUrl,
          req: create(CreateCallbackRequestSchema, {
            authRequestId: authRequest.id,
            callbackKind: {
              case: "session",
              value: { sessionId: cookie.id, sessionToken: cookie.token },
            },
          }),
        });

        return NextResponse.redirect(callbackUrl);
      } catch (error) {
        // Callback failed (session expired, etc.) - show account selection
        return NextResponse.redirect(constructUrl(request, "/accounts"));
      }
    } else {
      // Sessions exist but none are valid - show account selection
      return NextResponse.redirect(constructUrl(request, "/accounts"));
    }
  }

  // 5. No sessions - start fresh authentication
  const loginNameUrl = constructUrl(request, "/loginname");
  loginNameUrl.searchParams.set("requestId", requestId);

  if (authRequest?.loginHint) {
    loginNameUrl.searchParams.set("loginName", authRequest.loginHint);
  }
  if (organization) {
    loginNameUrl.searchParams.set("organization", organization);
  }

  return NextResponse.redirect(loginNameUrl);
}
```

#### SAML Flow Initiation Handler

SAML flow initiation is simpler (no prompt parameter):

```typescript
export async function handleSAMLFlowInitiation(params: FlowInitiationParams) {
  const { serviceUrl, requestId, sessions, sessionCookies, request } = params;

  // 1. Fetch SAML request details
  const { samlRequest } = await getSAMLRequest({
    serviceUrl,
    samlRequestId: requestId.replace("saml_", ""),
  });

  if (!samlRequest) {
    return NextResponse.json({ error: "No samlRequest found" }, { status: 400 });
  }

  // 2. No sessions - start authentication
  if (sessions.length === 0) {
    const loginNameUrl = constructUrl(request, "/loginname");
    loginNameUrl.searchParams.set("requestId", requestId);
    return NextResponse.redirect(loginNameUrl);
  }

  // 3. Find valid session
  const selectedSession = await findValidSession({
    serviceUrl,
    sessions,
    samlRequest,
  });

  if (!selectedSession) {
    // No valid session - show account selection
    return NextResponse.redirect(constructUrl(request, "/accounts"));
  }

  const cookie = sessionCookies.find((c) => c.id === selectedSession.id);

  if (!cookie) {
    return NextResponse.redirect(constructUrl(request, "/accounts"));
  }

  // 4. Complete SAML flow
  try {
    const { url, binding } = await createResponse({
      serviceUrl,
      req: create(CreateResponseRequestSchema, {
        samlRequestId: samlRequest.id,
        responseKind: {
          case: "session",
          value: {
            sessionId: cookie.id,
            sessionToken: cookie.token,
          },
        },
      }),
    });

    if (binding.case === "redirect") {
      // HTTP-Redirect binding
      return NextResponse.redirect(url);
    } else if (binding.case === "post") {
      // HTTP-POST binding - return auto-submit form
      const html = `
        <html>
          <body onload="document.forms[0].submit()">
            <form action="${url}" method="post">
              <input type="hidden" name="RelayState" value="${binding.value.relayState}" />
              <input type="hidden" name="SAMLResponse" value="${binding.value.samlResponse}" />
              <noscript>
                <button type="submit">Continue</button>
              </noscript>
            </form>
          </body>
        </html>
      `;

      return new NextResponse(html, {
        headers: { "Content-Type": "text/html" },
      });
    }
  } catch (error) {
    console.error("SAML createResponse failed:", error);
    return NextResponse.redirect(constructUrl(request, "/accounts"));
  }
}
```

#### OIDC Prompt Parameter Behavior

The `prompt` parameter controls authentication behavior:

| Prompt           | Behavior                                                                                       |
| ---------------- | ---------------------------------------------------------------------------------------------- |
| `none`           | Silent authentication - use existing session or fail immediately. No user interaction allowed. |
| `login`          | Force re-authentication even if valid session exists. Ignores existing sessions.               |
| `consent`        | Not yet implemented - treated as default behavior.                                             |
| `select_account` | Show account selection page even if only one session exists.                                   |
| `create`         | Redirect to registration flow instead of login.                                                |

**Prompt Priority**: If multiple prompt values are provided (space-separated), they are processed in priority order: `create` > `select_account` > `login` > `none`. Only the first matching prompt is acted upon.

#### Request ID Format

Request IDs are prefixed to identify the protocol:

- **OIDC**: `oidc_V2_123456789` (V2 prefix indicates Login Version 2)
- **SAML**: `saml_456789`
- **Device**: `device_789012` (handled at `/device` endpoint, not `/login`)

The prefix allows the login UI to:

1. Route to the correct handler
2. Call the appropriate backend API (getAuthRequest vs getSAMLRequest)
3. Complete with the correct callback method

#### Middleware Proxy Behavior

The middleware proxies certain paths to ZITADEL backend:

```typescript
export async function middleware(request: NextRequest) {
  const proxyPaths = [
    "/.well-known/", // OpenID configuration
    "/oauth/", // Token, introspection, revocation
    "/oidc/", // Userinfo, end_session
    "/idps/callback/", // IDP OAuth/SAML callbacks
    "/saml/", // SAML SSO/SLO endpoints
  ];

  // Only proxy in self-hosted mode
  if (process.env.ZITADEL_API_URL && process.env.ZITADEL_SERVICE_USER_TOKEN) {
    // Rewrite to ZITADEL backend
    request.nextUrl.href = `${serviceUrl}${request.nextUrl.pathname}`;

    // Add headers for backend
    requestHeaders.set("x-zitadel-public-host", request.nextUrl.host);
    requestHeaders.set("x-zitadel-instance-host", instanceHost);

    return NextResponse.rewrite(request.nextUrl, {
      request: { headers: requestHeaders },
    });
  }

  // Multi-tenant mode: no proxy (integrated stack)
  return NextResponse.next();
}
```

**Why proxy these paths?**

In self-hosted deployments, the login UI and ZITADEL backend are separate services. The middleware proxies protocol endpoints so:

1. **Token endpoint** (`/oauth/v2/token`): Applications exchange authorization codes for tokens
2. **Userinfo endpoint** (`/oidc/v1/userinfo`): Applications fetch user information
3. **IDP callbacks** (`/idps/callback/*`): External IDPs return to this path after authentication
4. **SAML endpoints** (`/saml/v2/*`): SAML assertions and metadata

Without proxying, applications would need to know two URLs (login UI + ZITADEL backend), and CORS would block requests.

#### Deployment Mode Differences

**Multi-Tenant (ZITADEL Cloud)**:

- ZITADEL backend and login UI run in the same context
- Authorization endpoint directly redirects to integrated login UI
- No middleware proxy needed
- Single domain for all endpoints

**Self-Hosted**:

- ZITADEL backend and login UI are separate services
- Authorization endpoint redirects to external login UI domain
- Middleware proxies protocol endpoints to ZITADEL backend
- Login UI must be configured as trusted domain in ZITADEL

#### Security Considerations

1. **Auth Request Storage**: Auth requests are stored server-side with all parameters. The login UI never trusts client-provided values - it always fetches from ZITADEL backend.

2. **Request ID Validation**: The login UI validates request IDs exist and are accessible before processing.

3. **Session Binding**: Sessions are validated against auth request requirements (organization, authentication level, etc.) before completion.

4. **Request ID Validation**: The `/login` endpoint (GET) doesn't require CSRF protection because it only reads the request ID from query parameters and immediately validates it by fetching the full auth request from ZITADEL backend. Invalid or malicious request IDs are rejected before any authentication occurs.

5. **Proxy Security**: Self-hosted deployments must configure the login UI domain as a trusted domain in ZITADEL to prevent unauthorized proxying.

### 5. OIDC/SAML Flow Completion

After successful authentication, the login UI completes the OIDC/SAML flow by creating a callback to ZITADEL with the session.

**OIDC Flow Completion**:

```typescript
// src/lib/oidc.ts
export async function loginWithOIDCAndSession({
  serviceUrl,
  authRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithOIDCAndSession): Promise<{ error: string } | { redirect: string }> {
  const selectedSession = sessions.find((s) => s.id === sessionId);

  // Validate session is still valid
  const isValid = await isSessionValid({
    serviceUrl,
    session: selectedSession,
  });

  if (!isValid) {
    // Re-authenticate user
    return sendLoginname({
      loginName: selectedSession.factors.user?.loginName,
      organization: selectedSession.factors.user?.organizationId,
      requestId: `oidc_${authRequest}`,
    });
  }

  const cookie = sessionCookies.find((cookie) => cookie.id === selectedSession?.id);

  // Create OIDC callback
  const { callbackUrl } = await createCallback({
    serviceUrl,
    req: create(CreateCallbackRequestSchema, {
      authRequestId: authRequest,
      callbackKind: {
        case: "session",
        value: {
          sessionId: cookie.id,
          sessionToken: cookie.token,
        },
      },
    }),
  });

  return { redirect: callbackUrl };
}
```

**SAML Flow Completion**:

```typescript
// src/lib/saml.ts
export async function loginWithSAMLAndSession({
  serviceUrl,
  samlRequest,
  sessionId,
  sessions,
  sessionCookies,
}: LoginWithSAMLAndSession): Promise<{ error: string } | { redirect: string }> {
  // Similar to OIDC but creates SAML response
  const { url } = await createResponse({
    serviceUrl,
    req: create(CreateResponseRequestSchema, {
      samlRequestId: samlRequest,
      responseKind: {
        case: "session",
        value: {
          sessionId: cookie.id,
          sessionToken: cookie.token,
        },
      },
    }),
  });

  return { redirect: url };
}
```

**Flow Completion Orchestration**:

```typescript
// src/lib/server/auth-flow.ts
export async function completeAuthFlow(command: {
  sessionId: string;
  requestId: string;
}): Promise<{ error: string } | { redirect: string }> {
  const sessionCookies = await getAllSessions();
  const sessions = await loadSessions({
    serviceUrl,
    ids: sessionCookies.map((s) => s.id),
  });

  if (requestId.startsWith("oidc_")) {
    return await loginWithOIDCAndSession({
      serviceUrl,
      authRequest: requestId.replace("oidc_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });
  } else if (requestId.startsWith("saml_")) {
    return await loginWithSAMLAndSession({
      serviceUrl,
      samlRequest: requestId.replace("saml_", ""),
      sessionId,
      sessions,
      sessionCookies,
    });
  }

  return { error: "Invalid request ID format" };
}
```

### 6. Device Authorization Flow

```
Device → Display Code → User → /device?user_code=XXX → Authenticate → Authorize Device
```

**Device Flow**:

```typescript
// src/lib/server/device.ts
export async function completeDeviceAuthorization(
  deviceAuthorizationId: string,
  session?: { sessionId: string; sessionToken: string },
) {
  // With session: approve device
  // Without session: deny device
  return authorizeOrDenyDeviceAuthorization({
    serviceUrl,
    deviceAuthorizationId,
    session, // undefined = deny
  });
}
```

## Middleware Architecture

The middleware intercepts all requests to handle:

1. Multi-tenancy routing
2. Proxy requests to ZITADEL backend
3. Security header injection
4. Organization header propagation

### Middleware Implementation

```typescript
// src/middleware.ts
export const config = {
  matcher: ["/.well-known/:path*", "/oauth/:path*", "/oidc/:path*", "/idps/callback/:path*", "/saml/:path*", "/:path*"],
};

export async function middleware(request: NextRequest) {
  const requestHeaders = new Headers(request.headers);

  // 1. Extract organization from query parameter
  const organization = request.nextUrl.searchParams.get("organization");
  if (organization) {
    requestHeaders.set("x-zitadel-i18n-organization", organization);
  }

  // 2. Determine if this is a proxy path
  const proxyPaths = ["/.well-known/", "/oauth/", "/oidc/", "/idps/callback/", "/saml/"];

  const isMatched = proxyPaths.some((prefix) => request.nextUrl.pathname.startsWith(prefix));

  // 3. Check if proxy is configured (self-hosted mode only)
  // Multi-tenant mode does NOT set these variables - runs integrated stack
  const isProxyConfigured = !!(process.env.ZITADEL_API_URL && process.env.ZITADEL_SERVICE_USER_TOKEN);

  // 4. For non-proxy paths or missing configuration, just add headers
  if (!isMatched || !isProxyConfigured) {
    return NextResponse.next({
      request: { headers: requestHeaders },
    });
  }

  // 5. Proxy mode (self-hosted only - login UI separate from backend)
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  const instanceHost = serviceUrl.replace("https://", "").replace("http://", "");

  // Add headers for ZITADEL backend
  requestHeaders.set("x-zitadel-public-host", request.nextUrl.host);
  requestHeaders.set("x-zitadel-instance-host", instanceHost);

  // 6. Set security headers
  const responseHeaders = new Headers();
  responseHeaders.set("Access-Control-Allow-Origin", "*");
  responseHeaders.set("Access-Control-Allow-Headers", "*");

  // 7. Apply iframe embedding settings if enabled
  const securitySettings = await loadSecuritySettings(request);
  if (securitySettings?.embeddedIframe?.enabled) {
    responseHeaders.set(
      "Content-Security-Policy",
      `${DEFAULT_CSP} frame-ancestors ${securitySettings.embeddedIframe.allowedOrigins.join(" ")};`,
    );
    responseHeaders.delete("X-Frame-Options");
  }

  // 8. Rewrite request to ZITADEL backend
  request.nextUrl.href = `${serviceUrl}${request.nextUrl.pathname}${request.nextUrl.search}`;

  return NextResponse.rewrite(request.nextUrl, {
    request: { headers: requestHeaders },
    headers: responseHeaders,
  });
}
```

### Service URL Resolution

The application determines the ZITADEL backend URL based on deployment mode:

```typescript
// src/lib/service-url.ts
export function getServiceUrlFromHeaders(headers: ReadonlyHeaders): {
  serviceUrl: string;
} {
  let instanceUrl;

  // 1. Multi-tenant: Use forwarded host
  const forwardedHost = headers.get("x-zitadel-forward-host");
  if (forwardedHost) {
    instanceUrl = forwardedHost.startsWith("http://") ? forwardedHost : `https://${forwardedHost}`;
  }
  // 2. Self-hosted with fixed URL: Use ZITADEL_API_URL
  else if (process.env.ZITADEL_API_URL) {
    instanceUrl = process.env.ZITADEL_API_URL;
  }
  // 3. Self-hosted with custom domain: Use host header
  else {
    const host = headers.get("host");
    if (host) {
      const [hostname] = host.split(":");
      if (hostname !== "localhost") {
        instanceUrl = host.startsWith("http") ? host : `https://${host}`;
      }
    }
  }

  if (!instanceUrl) {
    throw new Error("Service URL could not be determined");
  }

  return { serviceUrl: instanceUrl };
}
```

## Security Considerations

### HTTP-Only Cookies

All session cookies are **HTTP-only**, preventing JavaScript access:

```typescript
cookiesList.set({
  name: "sessions",
  httpOnly: true, // ← Cannot be accessed via JavaScript
  secure: process.env.NODE_ENV === "production",
  sameSite: iFrameEnabled ? "none" : "lax",
});
```

**Benefits**:

- Protection against XSS attacks
- Sessions cannot be stolen via client-side scripts
- Cookie manipulation requires server-side access

### SameSite Cookie Policy

The `sameSite` attribute prevents CSRF attacks:

```typescript
sameSite: iFrameEnabled ? "none" : "lax";
```

- **`"lax"`** (default): Cookies sent with top-level navigation (GET requests)
  - Blocks CSRF on state-changing requests (POST, etc.)
  - Allows cookies when user clicks links
- **`"none"`** (iframe mode): Required for cross-origin iframe embedding
  - Must be combined with `secure: true`
  - Only works over HTTPS

### Content Security Policy

The application enforces a strict CSP:

```javascript
// constants/csp.js
export const DEFAULT_CSP = `
  default-src 'self';
  script-src 'self' 'unsafe-inline' 'unsafe-eval';
  style-src 'self' 'unsafe-inline';
  img-src 'self' data: https:;
  font-src 'self' data:;
  connect-src 'self';
  media-src 'self';
  object-src 'none';
  base-uri 'self';
  form-action 'self';
`;
```

**Iframe Embedding**:

When iframe embedding is enabled via security settings:

```typescript
if (securitySettings?.embeddedIframe?.enabled) {
  responseHeaders.set(
    "Content-Security-Policy",
    `${DEFAULT_CSP} frame-ancestors ${securitySettings.embeddedIframe.allowedOrigins.join(" ")};`,
  );
  responseHeaders.delete("X-Frame-Options");
}
```

### Lockout Protection

Failed authentication attempts trigger account lockout:

```typescript
// Password authentication error handling
catch (error) {
  if ("failedAttempts" in error) {
    const lockoutSettings = await getLockoutSettings({
      serviceUrl,
      orgId: command.organization,
    });

    const hasLimit = lockoutSettings?.maxPasswordAttempts > BigInt(0);
    const locked = hasLimit &&
      error.failedAttempts >= lockoutSettings?.maxPasswordAttempts;

    return {
      error: t("errors.failedToAuthenticate", {
        failedAttempts: error.failedAttempts,
        maxPasswordAttempts: lockoutSettings?.maxPasswordAttempts,
        lockoutMessage: locked
          ? t("errors.accountLockedContactAdmin")
          : "",
      }),
    };
  }
}
```

### CSRF Protection

Next.js Server Actions provide built-in Cross-Site Request Forgery (CSRF) protection through multiple defense mechanisms. This protection is automatic and does not require explicit CSRF tokens.

#### How CSRF Protection Works

**1. POST-Only Invocation**

All Server Actions exclusively use the POST HTTP method:

```typescript
// Server Actions are always invoked via POST
"use server";
export async function loginAction(formData: FormData) {
  // This can only be called via POST
}
```

**Benefits**:

- GET requests cannot invoke Server Actions
- SameSite cookies (default `"lax"`) block cross-site POST requests
- Traditional CSRF attacks via `<img>`, `<link>`, or simple form GET submissions are prevented

**2. Origin Validation**

Next.js 14+ validates that requests originate from the same host:

```typescript
// Automatic validation performed by Next.js
// Compares Origin header with Host/X-Forwarded-Host header
if (requestOrigin !== expectedHost) {
  // Request rejected
  throw new Error("Origin validation failed");
}
```

**What gets validated**:

- `Origin` header (sent by browser for POST requests)
- Compared against `Host` or `X-Forwarded-Host` headers
- Must match exactly, including protocol and port

**Why this matters**:

- Prevents attacks from external domains
- Ensures Server Actions only callable from pages served by the same application
- Works even in old browsers that don't support `Origin` header (request is rejected)

**3. SameSite Cookie Policy**

Cookies are configured with appropriate SameSite attributes:

```typescript
// src/lib/cookies.ts
sameSite: iFrameEnabled ? "none" : "lax";
```

- **`"lax"`** (default): Cookies only sent with same-site requests or top-level navigation
  - Blocks cross-site POST/PUT/DELETE requests
  - Default CSRF protection for modern browsers
- **`"none"`** (iframe mode): Requires `secure: true` (HTTPS only)
  - Used only when explicitly needed for iframe embedding
  - Still protected by Origin validation

**4. Proxy-Aware Validation**

Next.js's Origin validation is proxy-aware and checks against:

- `Host` header (direct connections)
- `X-Forwarded-Host` header (when behind a proxy/load balancer)

This ensures CSRF protection works correctly in both:

- Direct deployments (Next.js exposed directly)
- Proxied deployments (behind nginx, load balancer, etc.)

**Note**: The application's `getOriginalHost()` function (in `src/lib/server/host.ts`) is used for URL construction (password reset links, etc.) but is **not** part of Next.js's CSRF validation. Next.js performs its own internal Origin validation.

#### Defense in Depth

These mechanisms work together to provide comprehensive CSRF protection:

1. **SameSite cookies** prevent cross-site request credentials from being sent
2. **POST-only actions** prevent simple GET-based attacks
3. **Origin validation** ensures requests come from legitimate sources
4. **Host validation** verifies proxy configurations don't expose attacks

#### No CSRF Tokens Required

Unlike traditional approaches, Next.js Server Actions do **not** use:

- Hidden CSRF tokens in forms
- Double-submit cookies
- Custom anti-CSRF headers

This simplification is possible because:

- Modern browsers support `Origin` headers and SameSite cookies
- Server Actions are framework-controlled (not arbitrary endpoints)
- Origin validation happens automatically at the framework level

#### Legacy Browser Considerations

For browsers that don't support the `Origin` header:

- The action is rejected (fail-secure approach)
- No silent fallback to potentially unsafe behavior
- Only modern, supported browsers can invoke Server Actions

### Search Parameter Security

The login application accepts numerous search parameters (URL query parameters) across different routes. Most are safe because they're validated by the backend or used for display only.

#### Search Parameter Inventory

**🔴 High Risk: Redirect URIs**

- `post_logout_redirect` / `post_logout_redirect_uri` (logout page)
- Requires multi-layer validation to prevent open redirect attacks

**🟡 Medium Risk: Pre-validated by Backend**

- `authRequest`, `samlRequest`, `requestId` - Authentication flow identifiers validated by backend

**🟢 Low Risk: Safe Parameters**

- Entity identifiers: `organization`, `userId`, `sessionId`, `user_code`, `code`
- Display/pre-fill: `firstname`, `lastname`, `email`, `loginName`, `suffix` (XSS-protected by React)
- Flags: `submit`, `invite`, `initial`, `send`
- Backend-controlled: `loginSettings.defaultRedirectUri` (configured in ZITADEL Console)

#### Redirect URI Validation Implementation

The login application implements **defense-in-depth validation** for redirect URIs:

**Validation Architecture**:

1. **Backend Validation** (Primary): ZITADEL validates URIs against registered client configurations
2. **Frontend Validation** (Defense-in-Depth): Login UI validates using trusted domains from ZITADEL API

**Frontend Implementation** (`/src/lib/url-validation.ts`):

- Fetches trusted domains from ZITADEL API via `listTrustedDomains()`
- Multi-layer validation: protocol checking, HTTPS enforcement, domain allowlist
- Automatic fallback to safe default (`/logout/done`) if validation fails
- Security logging for monitoring

**Security Layers**:

1. ✅ Protocol Validation - Blocks `javascript:`, `data:`, `file:`, etc.
2. ✅ HTTPS Enforcement - Requires HTTPS in production
3. ✅ Domain Allowlist - Only allows configured trusted domains
4. ✅ Security Logging - Logs blocked attempts for monitoring
5. 🔄 Rate Limiting - TODO: Prevent automated abuse

**Configuration**:
Trusted domains are configured in the ZITADEL Console by administrators and fetched dynamically from the backend API, ensuring consistent validation without environment variable synchronization.

#### OIDC/SAML Backend Validation

For OIDC/SAML flows, the ZITADEL backend validates `post_logout_redirect_uri` against pre-registered client URIs:

- Exact matching required (scheme, host, port, path)
- Client binding via `id_token_hint` or `client_id` parameter
- Supports glob patterns for dynamic deployments
- Development mode allows relaxed validation for `localhost`

#### Security Best Practices

**For ZITADEL Operators**:

1. Register URIs strictly - only necessary redirect URIs
2. Avoid wildcards unless required for dynamic deployments
3. Enforce HTTPS for all production redirect URIs
4. Regular audits of registered URIs
5. Monitor and alert on validation failures

**For Application Developers**:

1. Pre-register all redirect URIs in ZITADEL Console
2. Use specific paths, not just domains
3. Match registered URIs exactly (case-sensitive)
4. Test that unregistered URIs are rejected
5. Handle rejection cases gracefully

#### Other Search Parameters

The login application accepts various search parameters across different pages. Most are safe as they're used for lookups or validated by the backend:

**Safe Parameters** (validated by backend or used for lookups):

- `authRequest`, `samlRequest`, `requestId`: Auth request IDs
- `organization`, `userId`, `sessionId`: Entity identifiers
- `user_code`: Device authorization code
- `loginName`: Username for session lookup
- `code`: Verification/reset codes
- `logout_hint`: Username hint for logout

**Pre-fill Parameters** (safe, used for UX):

- `firstname`, `lastname`, `email`: Registration pre-fill
- `loginName`, `suffix`: Username hints
- `submit`, `invite`, `initial`: Boolean flags

**Backend-Controlled** (not from search params):

- `loginSettings.defaultRedirectUri`: Configured in ZITADEL Console
  - Used as fallback redirect after authentication
  - ✅ Safe - controlled by administrators, not user input
  - Retrieved from backend settings, never from URL parameters

**No Additional Validation Needed**: Parameters like `loginName`, `email`, `firstname` are used for display/pre-fill only. They don't trigger redirects or execute code. Input sanitization for display purposes is handled by React's automatic XSS protection.

#### Custom Route Handlers

If using custom Route Handlers (`route.tsx`), manual CSRF protection is required:

```typescript
// Custom routes do NOT have automatic CSRF protection
// Must implement traditional CSRF tokens or other validation
export async function POST(request: Request) {
  // WARNING: No automatic Origin validation here
  // Must implement CSRF protection manually
}
```

**Security note**: The login application uses Server Actions exclusively and does not rely on custom Route Handlers for mutations.

### gRPC Authentication

All backend communication uses authenticated gRPC with automatic token management based on deployment mode.

#### Service Client Creation

The `createServiceForHost` function is the central factory for creating authenticated gRPC service clients:

```typescript
// src/lib/service.ts
export async function createServiceForHost<T extends ServiceClass>(service: T, serviceUrl: string) {
  let token;

  // 1. Multi-tenancy mode: Generate JWT from service account
  if (process.env.AUDIENCE && process.env.SYSTEM_USER_ID && process.env.SYSTEM_USER_PRIVATE_KEY) {
    token = await systemAPIToken();
  }
  // 2. Self-hosted mode: Use static PAT
  else if (process.env.ZITADEL_SERVICE_USER_TOKEN) {
    token = process.env.ZITADEL_SERVICE_USER_TOKEN;
  }

  if (!serviceUrl) {
    throw new Error("No instance url found");
  }

  if (!token) {
    throw new Error("No token found");
  }

  // 3. Create gRPC transport with authentication
  const transport = createServerTransport(token, serviceUrl);

  // 4. Return typed service client
  return createClientFor<T>(service)(transport);
}
```

#### Token Authentication Strategies

**Strategy 1: Multi-Tenancy with JWT Service Account** (Priority 1)

Used when all three environment variables are present:

- `AUDIENCE`: Target ZITADEL API audience
- `SYSTEM_USER_ID`: Service account user ID
- `SYSTEM_USER_PRIVATE_KEY`: Base64-encoded private key

```typescript
// src/lib/api.ts
export async function systemAPIToken() {
  const token = {
    audience: process.env.AUDIENCE,
    userID: process.env.SYSTEM_USER_ID,
    token: Buffer.from(process.env.SYSTEM_USER_PRIVATE_KEY, "base64").toString("utf-8"),
  };

  // Generate short-lived JWT (typically 1 hour)
  return newSystemToken({
    audience: token.audience,
    subject: token.userID,
    key: token.token, // Private key for signing
  });
}
```

**Benefits**:

- **Dynamic JWT generation**: Fresh tokens for each request
- **Short-lived tokens**: Enhanced security (typically 1-hour expiry)
- **Region-wide access**: Single service account can access all instances in a region
- **Key rotation**: Private key can be rotated without code changes

**Strategy 2: Static Personal Access Token (PAT)** (Priority 2)

Used when `ZITADEL_SERVICE_USER_TOKEN` is provided:

```bash
# Environment variable
ZITADEL_SERVICE_USER_TOKEN=<long-lived-token>
```

**Benefits**:

- **Simple setup**: Single token, no key management
- **Self-hosted friendly**: Easy for single-instance deployments
- **File-based loading**: Can load from file for Kubernetes secrets

```typescript
// Token file support (checked during startup)
if (process.env.ZITADEL_SERVICE_USER_TOKEN_FILE) {
  // Block until file exists and read token
  const token = await readTokenFile(process.env.ZITADEL_SERVICE_USER_TOKEN_FILE);
  process.env.ZITADEL_SERVICE_USER_TOKEN = token;
}
```

**Required Role: `IAM_LOGIN_CLIENT`**

The service user must have the `IAM_LOGIN_CLIENT` role, which provides:

- Necessary permissions for authentication flows (session management, settings, OIDC/SAML)
- Restricted from user/organization management operations
- Follows principle of least privilege for login operations only

#### gRPC Transport Creation

```typescript
// src/lib/zitadel.ts
export function createServerTransport(token: string, baseUrl: string) {
  return libCreateServerTransport(token, {
    baseUrl,
    httpVersion: "2", // HTTP/2 required for gRPC
  });
}
```

The transport automatically:

- Adds `Authorization: Bearer <token>` header to all requests
- Uses HTTP/2 for multiplexed streams
- Handles connection pooling and keep-alives
- Provides automatic reconnection on failures

#### Service Client Usage

Once created, service clients provide type-safe gRPC methods:

```typescript
// Get login settings
const settingsService = await createServiceForHost(SettingsService, serviceUrl);

const loginSettings = await settingsService.getLoginSettings(
  { ctx: makeReqCtx(organization) },
  {}, // gRPC call options
);

// Create session
const sessionService = await createServiceForHost(SessionService, serviceUrl);

const session = await sessionService.createSession(
  {
    checks: {
      /* ... */
    },
    lifetime: { seconds: BigInt(3600), nanos: 0 },
  },
  {},
);
```

#### Authentication Context Propagation

The token determines the authentication context for all gRPC calls:

```typescript
// Service account context (via createServiceForHost)
const settingsService = await createServiceForHost(SettingsService, serviceUrl);
// ↑ Uses service account token (system-level access)

// User session context (via createServerTransport directly)
const transport = createServerTransport(
  sessionToken, // ← User's session token
  serviceUrl,
);
const userService = createUserServiceClient(transport);
// ↑ Uses user's session token (user-level access)
```

**Service Account vs User Context**:

- **Service Account** (from `createServiceForHost`):
  - System-level operations
  - Read settings, list users, create sessions
  - Used for most operations in the login flow

- **User Context** (from user's session token):
  - Self-service operations
  - Change own password, set MFA, update profile
  - Used when user has active session and sufficient permissions

#### Token Priority and Fallback

```
┌─────────────────────────────────────────┐
│ Token Selection Priority                │
├─────────────────────────────────────────┤
│ 1. JWT Service Account                  │
│    ✓ AUDIENCE                           │
│    ✓ SYSTEM_USER_ID                     │
│    ✓ SYSTEM_USER_PRIVATE_KEY            │
│    → Generate fresh JWT via systemAPIToken() │
├─────────────────────────────────────────┤
│ 2. Static PAT                           │
│    ✓ ZITADEL_SERVICE_USER_TOKEN         │
│    → Use token directly                 │
├─────────────────────────────────────────┤
│ 3. No token configured                  │
│    ✗ Throw error: "No token found"      │
└─────────────────────────────────────────┘
```

#### Example: Complete Flow

```typescript
// 1. User submits password
export async function sendPassword(command: {...}) {
  // 2. Determine service URL (multi-tenant or self-hosted)
  const { serviceUrl } = getServiceUrlFromHeaders(headers);

  // 3. Create authenticated session service
  //    (automatically selects JWT or PAT based on env vars)
  const sessionService = await createServiceForHost(
    SessionService,
    serviceUrl
  );

  // 4. Make authenticated gRPC call
  const session = await sessionService.setSession({
    sessionId: existingSession.id,
    sessionToken: existingSession.token,
    checks: {
      password: { password: command.password }
    }
  });

  // 5. Session created/updated with service account permissions
  return { sessionId: session.sessionId };
}
```

#### Security Considerations

**JWT Service Account Mode**:

- Private key must be kept secure
- Key should be stored in secure secret management (Kubernetes secrets, AWS Secrets Manager, etc.)
- Short-lived JWTs reduce impact of token compromise
- Each request gets a fresh token (no token reuse)

**Static PAT Mode**:

- Token has longer lifetime (typically months/years)
- Should be rotated periodically
- If compromised, must be revoked and regenerated
- Suitable for trusted environments (self-hosted)

**Transport Security**:

- All gRPC communication over TLS (HTTPS)
- HTTP/2 provides multiplexing and header compression
- Tokens transmitted in HTTP headers (not URL)
- Connection pooling minimizes overhead

### Password Security

Password operations have additional safeguards:

```typescript
// Password reset requires verification code
export async function changePassword(command: {
  code?: string; // Verification code from email
  userId: string;
  password: string;
}) {
  // Without code, check if user was verified recently
  if (!command.code) {
    const authmethods = await listAuthenticationMethodTypes({
      serviceUrl,
      userId,
    });

    // Require verification if user has existing auth methods
    if (authmethods.authMethodTypes.length !== 0) {
      return {
        error: "Code or verification required",
      };
    }

    // Check for valid user verification check
    const hasValidUserVerificationCheck = await checkUserVerification(userId);

    if (!hasValidUserVerificationCheck) {
      return { error: "Verification required" };
    }
  }

  // Set password with code or verification
  return setUserPassword({
    serviceUrl,
    userId,
    password: command.password,
    code: command.code,
  });
}
```

### Host Header Validation for URL Templates

**⚠️ CRITICAL SECURITY REQUIREMENT**:

The login application uses the **host header** from incoming requests to construct URLs for:

- Password reset emails
- Email verification links
- User invite links
- Other notification templates

**How it works**:

```typescript
// Login application reads host header
const host = headers.get("host");
const loginUrl = `https://${host}/reset-password?token=xyz`;

// ZITADEL backend receives this URL for email templates
await sendPasswordResetEmail({
  resetUrl: loginUrl, // Used in email template
});
```

**Security Risk**:

If the host header is not validated, an attacker could:

1. Send a request with a malicious host header (e.g., `Host: evil.com`)
2. Trigger a password reset, email verification, or invite flow
3. The victim receives an email with a link to the attacker's domain
4. The link contains a valid reset token/verification code
5. When the victim clicks the link, their credentials/tokens are sent to the attacker

**Required Protection**:

ZITADEL **MUST** validate the host header against its configured trusted domains before:

- Generating email templates with URLs
- Sending password reset emails
- Sending verification emails
- Sending invite emails
- Any other flow that includes user-clickable links

**Implementation Requirements**:

ZITADEL backend MUST:

- ✅ Maintain a list of trusted/allowed domains for the instance
- ✅ Validate the host header from login application requests
- ✅ Reject requests with untrusted host headers before sending emails
- ✅ Log suspicious host header values for security monitoring
- ✅ Use only validated host values in URL templates

The login application:

- ✅ Passes the host header value in its requests to ZITADEL
- ❌ Does NOT validate the host header itself (trusts the request)
- ❌ Cannot prevent header manipulation by clients
- ✅ Relies on ZITADEL to perform validation before using URLs

**Attack Scenario Example**:

```bash
# Attacker sends malicious request
curl -H "Host: attacker.com" \
     https://login.mycompany.com/reset-password \
     -d "email=victim@company.com"

# Without validation, victim receives:
# "Reset your password: https://attacker.com/reset-password?token=valid-token-123"

# With validation, ZITADEL rejects the request:
# HTTP 400 Bad Request: Untrusted host header
```

**Why This Matters**:

This is a **critical phishing vector** because:

- Links in emails appear to come from legitimate ZITADEL notifications
- Users trust password reset and verification emails
- Valid tokens make the malicious links fully functional
- Attackers can harvest credentials and session tokens
- The attack is simple to execute and hard for users to detect

**Next.js CSRF Protection Helps Here**:

Next.js's built-in Origin validation (described in the [CSRF Protection](#csrf-protection) section) provides an additional security layer:

- When Server Actions are invoked to trigger password reset/verification flows, Next.js validates that the `Origin` header matches the `Host`/`X-Forwarded-Host` header
- This means an attacker **cannot** simply send a forged request with `Host: attacker.com` from their own domain
- The browser will send `Origin: https://attacker.com` but `Host: login.mycompany.com`, causing Next.js to reject the request
- Attackers would need to either:
  - Compromise the legitimate login domain itself (much harder)
  - Find a vulnerability that bypasses Origin validation (framework-level exploit)

**Defense in Depth**:

The combination of Next.js Origin validation and ZITADEL's trusted domain validation provides defense in depth:

1. **Next.js Layer** (Login Application):
   - Rejects forged host headers in Server Action requests
   - Prevents cross-origin Server Action invocation
   - Stops most simple host header manipulation attacks

2. **ZITADEL Layer** (Backend):
   - Validates host header against trusted domains before sending emails
   - Protects against attacks that bypass the Next.js layer
   - Ensures only validated URLs appear in email templates
   - Prevents internal API calls with malicious hosts

**Remaining Attack Vectors**:

Even with Next.js protection, ZITADEL's validation is still critical because:

- Direct API calls to ZITADEL that bypass the Next.js frontend
- Internal requests between services that don't go through Next.js
- Custom integrations that use ZITADEL's email APIs directly
- Future changes to the login application architecture

**Best Practice**: Both layers should validate independently - never rely solely on frontend protection for security-critical operations like email template generation.

### User Agent Verification for Invite Flows

**Security Mechanism**: To protect against session hijacking and token replay attacks during user invite and verification flows, the application uses a **verification check cookie** that binds the verification flow to the specific browser that initiated it.

**How it works**:

```typescript
// When user verifies email/invite without existing authentication methods
const userAgentId = await getOrSetFingerprintId();

// Create SHA-256 hash binding userId to browser fingerprint
const verificationCheck = crypto.createHash("sha256").update(`${user.userId}:${userAgentId}`).digest("hex");

// Set short-lived verification cookie
await cookiesList.set({
  name: "verificationCheck",
  value: verificationCheck,
  httpOnly: true,
  path: "/",
  maxAge: 300, // 5 minutes
});
```

**Browser Fingerprinting**:

The application generates a persistent fingerprint ID for each browser:

```typescript
// Generate unique fingerprint ID on first visit
export async function getFingerprintId() {
  return uuidv4(); // Generates unique UUID
}

// Store in long-lived HTTP-only cookie
await cookiesList.set({
  name: "fingerprintId",
  value: fingerprintId,
  httpOnly: true,
  path: "/",
  maxAge: 31536000, // 1 year
});
```

The fingerprint is included in all session creation requests along with user agent details:

```typescript
const userAgentData: UserAgent = {
  ip: headers.get("x-forwarded-for") ?? headers.get("remoteAddress"),
  header: { "user-agent": { values: userAgentHeaderValues } },
  description: `${browserDescription}, ${deviceDescription}, ${engineDescription}, ${osDescription}`,
  fingerprintId: fingerprintId,
};
```

**Attack Prevention**:

This mechanism prevents several attack vectors:

1. **Verification Link Interception**:

   ```
   Attacker intercepts: /verify?code=xyz&userId=123
   → When attacker visits link, verification check fails
   → Different browser fingerprint = different hash
   → Flow rejected
   ```

2. **Session Hijacking During Invite**:
   - Verification check ensures only the browser that received the invite email can complete setup
   - Even with valid verification code, attacker needs matching fingerprint

3. **Token Replay Attacks**:
   - Verification check cookie expires after 5 minutes
   - Hash binds specific user to specific browser
   - Prevents reuse of verification URLs across browsers

**When This Protection Applies**:

This verification check is used specifically when:

- User verifies email/invite AND has **no existing authentication methods**
- User is being redirected to set up their first authenticator
- Session doesn't exist or user is not authenticated

**Implementation Details**:

```typescript
// Set verification check during email/invite verification
if (authMethodTypes.length === 0) {
  // No auth methods - user needs to set up authenticator
  const userAgentId = await getOrSetFingerprintId();
  const verificationCheck = crypto.createHash("sha256").update(`${user.userId}:${userAgentId}`).digest("hex");

  await cookies().set({
    name: "verificationCheck",
    value: verificationCheck,
    httpOnly: true,
    path: "/",
    maxAge: 300, // 5 minutes
  });

  return { redirect: `/authenticator/set?sessionId=${session.id}` };
}
```

**Security Properties**:

- ✅ **Browser-bound**: Verification only succeeds in the same browser
- ✅ **Time-limited**: 5-minute expiration window
- ✅ **HTTP-only**: Cannot be accessed or modified by JavaScript
- ✅ **Cryptographic binding**: SHA-256 hash prevents tampering
- ✅ **Automatic cleanup**: Cookie expires after 5 minutes
- ✅ **Session isolation**: Each verification flow gets unique hash

**Limitations**:

- ❌ Protection only applies to users with no authentication methods
- ❌ Users with existing authenticators don't get verification check
- ❌ Browser fingerprint can change (incognito mode, browser updates)
- ❌ Shared devices may retain fingerprint across users

**Why 5 Minutes?**:

The short expiration balances security and usability:

- Long enough for legitimate users to complete authenticator setup
- Short enough to limit attack window
- Forces attackers to act immediately, increasing detection likelihood

## Multi-Factor Authentication

### MFA Enforcement Policy

MFA enforcement is determined by login settings:

```typescript
// src/lib/verify-helper.ts
export function shouldEnforceMFA(session: Session, loginSettings: LoginSettings | undefined): boolean {
  if (!loginSettings) return false;

  // Passkeys are inherently multi-factor
  const authenticatedWithPasskey = session.factors?.webAuthN?.verifiedAt && session.factors?.webAuthN?.userVerified;

  if (authenticatedWithPasskey) return false;

  // forceMfa: MFA required for ALL auth methods
  if (loginSettings.forceMfa) return true;

  // forceMfaLocalOnly: MFA only for password auth
  if (loginSettings.forceMfaLocalOnly) {
    const authenticatedWithPassword = !!session.factors?.password?.verifiedAt;
    const authenticatedWithIDP = !!session.factors?.intent?.verifiedAt;

    // IDP users skip MFA with forceMfaLocalOnly
    if (authenticatedWithIDP) return false;

    // Password users require MFA with forceMfaLocalOnly
    if (authenticatedWithPassword) return true;
  }

  return false;
}
```

### MFA Flow

```typescript
export async function checkMFAFactors(
  serviceUrl: string,
  session: Session,
  loginSettings: LoginSettings | undefined,
  authMethods: AuthenticationMethodType[],
  organization?: string,
  requestId?: string,
) {
  // Filter to MFA methods only
  const availableMultiFactors = authMethods?.filter(
    (m) =>
      m === AuthenticationMethodType.TOTP ||
      m === AuthenticationMethodType.OTP_SMS ||
      m === AuthenticationMethodType.OTP_EMAIL ||
      m === AuthenticationMethodType.U2F,
  );

  // Passkeys satisfy MFA automatically
  const hasAuthenticatedWithPasskey = session.factors?.webAuthN?.verifiedAt && session.factors?.webAuthN?.userVerified;

  if (hasAuthenticatedWithPasskey) return;

  // Single MFA method: redirect directly
  if (availableMultiFactors?.length === 1) {
    const factor = availableMultiFactors[0];
    const params = new URLSearchParams({
      loginName: session.factors?.user?.loginName,
      organization,
      requestId,
    });

    switch (factor) {
      case AuthenticationMethodType.TOTP:
        return { redirect: `/otp/time-based?` + params };
      case AuthenticationMethodType.OTP_SMS:
        return { redirect: `/otp/sms?` + params };
      case AuthenticationMethodType.OTP_EMAIL:
        return { redirect: `/otp/email?` + params };
      case AuthenticationMethodType.U2F:
        return { redirect: `/u2f?` + params };
    }
  }
  // Multiple MFA methods: show selection page
  else if (availableMultiFactors?.length > 1) {
    return { redirect: `/mfa?` + params };
  }
  // MFA enforced but no methods: force setup
  else if (shouldEnforceMFA(session, loginSettings) && !availableMultiFactors.length) {
    return {
      redirect: `/mfa/set?` + params + `&force=true&checkAfter=true`,
    };
  }
  // MFA skip lifetime: allow temporary skip
  else if (loginSettings?.mfaInitSkipLifetime && shouldEnforceMFA(session, loginSettings)) {
    const user = await getUserByID({
      serviceUrl,
      userId: session.factors.user.id,
    });

    const humanUser = user?.user?.type.case === "human" ? user.user.type.value : undefined;

    // Check if skip is still valid
    if (humanUser?.mfaInitSkipped) {
      const skipTime = timestampDate(humanUser.mfaInitSkipped).getTime();
      const skipLifetime = Number(loginSettings.mfaInitSkipLifetime.seconds) * 1000;
      const elapsed = Date.now() - skipTime;

      if (elapsed <= skipLifetime) return; // Skip still valid
    }

    // Offer optional MFA setup
    return {
      redirect: `/mfa/set?` + params + `&force=false&checkAfter=true`,
    };
  }
}
```

### Supported MFA Methods

1. **TOTP (Time-based One-Time Password)**
   - RFC 6238 compliant
   - Works with Google Authenticator, Authy, etc.

2. **OTP Email**
   - Code sent via email
   - Short-lived verification codes

3. **OTP SMS**
   - Code sent via SMS
   - Requires phone number verification

4. **U2F (Universal 2nd Factor)**
   - Hardware security keys
   - FIDO U2F protocol

5. **Passkeys (WebAuthn)**
   - Inherently multi-factor (possession + biometric/PIN)
   - Satisfies MFA requirements automatically

## Identity Provider Integration

### Supported IDP Types

1. **OAuth/OIDC Providers**
   - Google, GitHub, Azure AD, etc.
   - Standard OAuth 2.0 authorization code flow

2. **LDAP**
   - Username/password collection in UI
   - Backend LDAP bind authentication

3. **SAML**
   - Enterprise SSO
   - SAML 2.0 Web Browser SSO Profile

### IDP Flow Architecture

```
┌─────────┐
│  User   │
└────┬────┘
     │ 1. Select IDP
     ▼
┌─────────────────┐
│  /idp page      │
└────┬────────────┘
     │ 2. Start IDP Flow (Server Action)
     ▼
┌─────────────────┐
│  startIDPFlow   │
│  - Registers    │
│    success/fail │
│    URLs         │
└────┬────────────┘
     │ 3. Redirect to IDP
     ▼
┌─────────────────┐
│ External IDP    │
│ (OAuth/SAML)    │
└────┬────────────┘
     │ 4. User authenticates
     ▼
┌─────────────────┐
│ IDP Callback    │
│ /idp/{provider}/│
│ process         │
└────┬────────────┘
     │ 5. Exchange code/assertion
     ▼
┌─────────────────┐
│ ZITADEL Backend │
│ - Validates     │
│ - Creates intent│
└────┬────────────┘
     │ 6. IDP Intent
     ▼
┌─────────────────┐
│ Create Session  │
│ with IDP factor │
└────┬────────────┘
     │ 7. [Optional: Link/Register]
     │ 8. [Optional: MFA Check]
     ▼
┌─────────────────┐
│   Redirect      │
└─────────────────┘
```

### IDP Linking

Users can link external IDPs to existing accounts:

```typescript
// Link mode: User must be authenticated
const linkOnly = formData.get("linkOnly") === "true";

if (linkOnly) {
  // Require existing session
  const existingSession = await getMostRecentSessionCookie();

  if (!existingSession) {
    return { error: "Authentication required to link IDP" };
  }
}
```

### LDAP Special Handling

LDAP requires credential collection before redirect:

```typescript
export async function redirectToIdp(formData: FormData) {
  const provider = formData.get("provider") as string;

  if (provider === "ldap") {
    // Redirect to LDAP credential page
    redirect(`/idp/ldap?idpId=${idpId}&...`);
  }

  // Other IDPs redirect directly
  const response = await startIdentityProviderFlow({...});
  redirect(response.redirect);
}

// LDAP authentication
export async function createNewSessionForLDAP(command: {
  username: string;
  password: string;
  idpId: string;
}) {
  const response = await startLDAPIdentityProviderFlow({
    serviceUrl,
    idpId: command.idpId,
    username: command.username,
    password: command.password,
  });

  // Returns IDP intent on success
  const { userId, idpIntentId, idpIntentToken } =
    response.nextStep.value;

  // Continue with IDP intent flow
  return {
    redirect: `/idp/ldap/success?userId=${userId}&id=${idpIntentId}&token=${idpIntentToken}`,
  };
}
```

## Deployment Modes

### 1. Multi-Tenant Deployment

**Configuration**:

```bash
# Required environment variables for multi-tenancy
AUDIENCE=https://api.zitadel.cloud  # Target ZITADEL API audience
SYSTEM_USER_ID=<service-account-user-id>
SYSTEM_USER_PRIVATE_KEY=<base64-encoded-private-key>

# Optional: Override API URL (defaults to AUDIENCE)
ZITADEL_API_URL=https://zitadel.cloud
```

**Characteristics**:

- Centralized login UI serving multiple instances in a region
- Instance context passed via `x-zitadel-forward-host` header
- JWT service account with region-wide access
- **No middleware proxy** - ZITADEL runs the entire application stack together
- Organization-specific branding and translations

**Host Resolution**:

```typescript
// Multi-tenant: Use forwarded host header to determine target instance
const forwardedHost = headers.get("x-zitadel-forward-host");
if (forwardedHost) {
  instanceUrl = `https://${forwardedHost}`; // e.g., custom-abc123.zitadel.cloud
}
```

**⚠️ CRITICAL SECURITY REQUIREMENT**:

The multi-tenant deployment model relies on a **trustworthy proxy/gateway** that:

1. **Validates the original customer request** against the target instance
2. **Sets the `x-zitadel-forward-host` header** based on validated routing rules
3. **Prevents header injection** - ensures clients cannot override or forge this header
4. **Strips any existing `x-zitadel-forward-host` headers** from incoming requests before adding the validated one

**Security Model**:

```
┌──────────────────────────────────────────────────────────────┐
│                    Trusted Proxy/Gateway                     │
│                                                              │
│  1. Receives: https://custom-abc123.zitadel.cloud/login      │
│  2. Validates: custom-abc123 exists and is authorized        │
│  3. Strips: Any x-zitadel-forward-host from client           │
│  4. Sets: x-zitadel-forward-host: custom-abc123.zitadel.cloud│
│  5. Forwards to: Login App                                   │
└──────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌──────────────────────────────────────────────────────────────┐
│                      Login Application                       │
│                                                              │
│  - Trusts x-zitadel-forward-host header implicitly           │
│  - Uses header to determine target ZITADEL instance          │
│  - Creates gRPC clients for that specific instance           │
└──────────────────────────────────────────────────────────────┘
```

**Why This Matters**:

Without a trusted proxy that validates requests and controls header injection:

- **Instance impersonation**: Attackers could access any instance by forging headers
- **Unauthorized access**: Users could authenticate against wrong instances
- **Cross-tenant attacks**: Session/data leakage between different customers

**Implementation Requirements**:

The proxy MUST:

- ✅ Validate the original request URL/domain against allowed instances
- ✅ Strip all `x-zitadel-forward-host` headers from client requests
- ✅ Set `x-zitadel-forward-host` based ONLY on validated routing logic
- ✅ Use TLS to communicate with the login application
- ✅ Not expose the login application directly to clients

The login application:

- ❌ Does NOT validate the `x-zitadel-forward-host` header
- ❌ Does NOT verify instance authorization
- ✅ Assumes the header is set by a trusted component
- ✅ Uses the header value directly for instance resolution

**Authentication**:

```typescript
// Multi-tenant uses JWT service account
if (process.env.AUDIENCE && process.env.SYSTEM_USER_ID && process.env.SYSTEM_USER_PRIVATE_KEY) {
  // Generate fresh JWT for each request
  token = await systemAPIToken();
}
```

### 2. Self-Hosted Single-Instance

**Configuration**:

```bash
# Required environment variables for self-hosted
ZITADEL_API_URL=https://login-abc123.zitadel.cloud
ZITADEL_SERVICE_USER_TOKEN=<personal-access-token>
```

**Characteristics**:

- Single ZITADEL instance
- Direct communication with local backend
- Static PAT (Personal Access Token) for authentication
- **Middleware proxies OIDC/SAML endpoints** to ZITADEL backend (same as multi-tenant)
- Simplified deployment (single token, no JWT generation)

**⚠️ SERVICE USER ROLE REQUIREMENT**:

The service user that generates the `ZITADEL_SERVICE_USER_TOKEN` must be assigned the **`IAM_LOGIN_CLIENT`** role. This role provides:

- ✅ **Sufficient permissions** for all login operations:
  - Create and manage user sessions
  - Read login/branding/security settings
  - Validate authentication factors
  - Handle OIDC/SAML flows
  - Manage MFA verification

- ✅ **Principle of least privilege**:
  - No user management permissions (create/delete users)
  - No organization management permissions
  - No system administration access
  - Scoped specifically to authentication flows

**Setting up the service user**:

1. Create a service user in ZITADEL Console
2. Assign the `IAM_LOGIN_CLIENT` role to this user
3. Generate a Personal Access Token (PAT) for this user
4. Use the PAT as `ZITADEL_SERVICE_USER_TOKEN`

```bash
# The service user should have exactly this role - no more, no less
# Role: IAM_LOGIN_CLIENT
# Scope: Instance-level (for self-hosted single instance)
```

**Authentication**:

```typescript
// Self-hosted uses static PAT
if (process.env.ZITADEL_SERVICE_USER_TOKEN) {
  token = process.env.ZITADEL_SERVICE_USER_TOKEN;
}
```

**Host Resolution**:

```typescript
// Self-hosted: Use ZITADEL_API_URL
if (process.env.ZITADEL_API_URL) {
  instanceUrl = process.env.ZITADEL_API_URL;
}
// Fallback: Use host header (only for non-localhost)
// Note: This is rarely used and requires login UI on same domain as ZITADEL
else {
  const host = headers.get("host");
  const [hostname] = host.split(":");
  if (hostname !== "localhost") {
    instanceUrl = `https://${host}`;
  }
}
```

**⚠️ TRUSTED DOMAIN REQUIREMENT**:

The domain where the login application runs must be configured as a **trusted domain** in the ZITADEL instance:

1. Add the login UI domain to the instance's allowed domains via API (e.g., `https://login.mycompany.com`)
2. Use the ZITADEL Management API or Admin API to configure trusted domains
3. This allows the login UI to proxy requests to ZITADEL

**How the middleware forwards requests**:

When proxying OIDC/SAML endpoints, the middleware adds these headers to identify the request origin:

```typescript
requestHeaders.set("x-zitadel-public-host", request.nextUrl.host); // Login UI domain
requestHeaders.set("x-zitadel-instance-host", instanceHost); // Target ZITADEL instance
```

ZITADEL validates the `x-zitadel-public-host` against its trusted domains configuration.

Without this configuration:

- **ZITADEL will reject proxied requests** from the login UI domain
- CORS errors will prevent the login UI from communicating with ZITADEL
- Session cookies may be blocked by browser security policies
- Authentication flows will fail

### Deployment Differences

- **Multi-tenant (ZITADEL Cloud)**:
  - ZITADEL runs the entire application stack together (backend + login UI)
  - No proxy needed - everything runs in the same context
  - `ZITADEL_API_URL` and `ZITADEL_SERVICE_USER_TOKEN` are **not set**
  - OIDC/SAML endpoints handled directly by the integrated stack

- **Self-hosted**:
  - Login UI and ZITADEL backend are separate services
  - Proxy required to forward protocol endpoints to backend
  - `ZITADEL_API_URL` and `ZITADEL_SERVICE_USER_TOKEN` **must be set**
  - Middleware rewrites requests to ZITADEL backend

**Why Proxy Only in Self-Hosted?**

In self-hosted deployments, the login UI needs to proxy OIDC/SAML endpoints because:

1. **Separate services**: Login UI (Next.js) and ZITADEL backend run independently
2. **Protocol endpoints** (`/.well-known/`, `/oauth/`, `/oidc/`, `/saml/`) must be handled by ZITADEL backend
3. **IDP callbacks** (`/idps/callback/`) return to the same domain and must reach ZITADEL
4. **Unified domain**: Clients interact with a single domain for both UI and protocol endpoints

**Proxy Behavior (Self-Hosted Only)**:

```typescript
// Middleware rewrites proxy paths to backend
// This ONLY applies to self-hosted deployments where login UI is separate
const isProxyConfigured = !!(process.env.ZITADEL_API_URL && process.env.ZITADEL_SERVICE_USER_TOKEN);

if (isProxyConfigured && isProxyPath) {
  request.nextUrl.href = `${serviceUrl}${request.nextUrl.pathname}`;
  return NextResponse.rewrite(request.nextUrl, {
    headers: {
      "x-zitadel-public-host": request.nextUrl.host,
      "x-zitadel-instance-host": instanceHost,
    },
  });
}
```

### 3. Development Mode

**Configuration**:

```bash
# .env.local
ZITADEL_API_URL=http://localhost:8080
DEBUG=true  # Disables caching
```

**Characteristics**:

- Local ZITADEL instance
- Hot reloading
- Disabled caching for immediate updates
- HTTP allowed for localhost

### Deployment Comparison

| Feature                 | Multi-Tenant                                  | Self-Hosted                             | Development                            |
| ----------------------- | --------------------------------------------- | --------------------------------------- | -------------------------------------- |
| **Authentication**      | JWT (AUDIENCE + SYSTEM_USER_ID + PRIVATE_KEY) | Static PAT (ZITADEL_SERVICE_USER_TOKEN) | JWT OR Static PAT (same as production) |
| **ZITADEL_API_URL**     | Optional (defaults to AUDIENCE)               | Required                                | Required                               |
| **Token Generation**    | Dynamic JWT per request                       | Static long-lived token                 | Dynamic JWT or Static PAT              |
| **Middleware Proxy**    | No (integrated stack)                         | Yes (for OIDC/SAML endpoints)           | No (direct access)                     |
| **Host Resolution**     | `x-zitadel-forward-host` header               | `ZITADEL_API_URL`                       | `ZITADEL_API_URL`                      |
| **Access Scope**        | Region-wide (all instances)                   | Single instance                         | Single instance                        |
| **Caching**             | Enabled                                       | Enabled                                 | Disabled (DEBUG=true)                  |
| **HTTPS**               | Required                                      | Recommended                             | Not enabled (HTTP)                     |
| **Organization Header** | Propagated                                    | Propagated                              | Propagated                             |

## Next.js Implementation

### Server Actions

Server Actions provide type-safe, server-side mutations without API routes:

### Server Action Pattern

```typescript
// src/lib/server/password.ts
"use server";  // ← Marks file as server-only

import { headers } from "next/headers";

export async function sendPassword(command: {
  loginName: string;
  checks: Checks;
}): Promise<{ error: string } | { redirect: string }> {
  // Server-side only code
  const _headers = await headers();
  const { serviceUrl } = getServiceUrlFromHeaders(_headers);

  // Perform server-side operations
  const session = await setSessionAndUpdateCookie({...});

  // Return serializable result
  return { redirect: "/signedin" };
}
```

### Client-Side Usage

```typescript
// src/components/password-form.tsx
"use client";

import { sendPassword } from "@/lib/server/password";
import { useFormState } from "react-dom";

export function PasswordForm() {
  const [state, formAction] = useFormState(sendPassword, null);

  return (
    <form action={formAction}>
      <input name="password" type="password" />
      <button type="submit">Sign In</button>
      {state?.error && <p>{state.error}</p>}
    </form>
  );
}
```

### Benefits

1. **Type Safety**: Full TypeScript support across client/server boundary
2. **No API Routes**: Direct function calls from client to server
3. **Progressive Enhancement**: Forms work without JavaScript
4. **Automatic Serialization**: Return values are automatically serialized
5. **Server-Only Code**: Sensitive operations never exposed to client

### Internationalization

#### Translation Architecture

```typescript
// Server-side translations
import { getTranslations } from "next-intl/server";

export async function sendPassword(command: {...}) {
  const t = await getTranslations("password");

  return {
    error: t("errors.couldNotVerifyPassword")
  };
}
```

### Organization-Specific Translations

```typescript
// Fetch translations from ZITADEL
export async function getHostedLoginTranslation({
  serviceUrl,
  organization,
  locale,
}: {
  serviceUrl: string;
  organization?: string;
  locale?: string;
}) {
  const settingsService = await createServiceForHost(SettingsService, serviceUrl);

  return settingsService.getHostedLoginTranslation({
    level: organization ? { case: "organizationId", value: organization } : { case: "instance", value: true },
    locale,
  });
}
```

### Translation Propagation

The organization parameter is propagated through middleware:

```typescript
// middleware.ts
const organization = request.nextUrl.searchParams.get("organization");
if (organization) {
  requestHeaders.set("x-zitadel-i18n-organization", organization);
}
```

### Theming System

The application supports environment-variable driven theming:

#### Theme Configuration

```bash
# .env.local
NEXT_PUBLIC_THEME_ROUNDNESS=mid          # edgy | mid | full
NEXT_PUBLIC_THEME_LAYOUT=side-by-side    # side-by-side | top-to-bottom
NEXT_PUBLIC_THEME_APPEARANCE=material    # flat | material
NEXT_PUBLIC_THEME_SPACING=regular        # regular | compact
```

### Server-Safe Theme Functions

```typescript
// src/lib/theme.ts
export function getThemeConfig() {
  return {
    roundness: process.env.NEXT_PUBLIC_THEME_ROUNDNESS || "mid",
    layout: process.env.NEXT_PUBLIC_THEME_LAYOUT || "side-by-side",
    appearance: process.env.NEXT_PUBLIC_THEME_APPEARANCE || "material",
    spacing: process.env.NEXT_PUBLIC_THEME_SPACING || "regular",
  };
}

export function getComponentRoundness(component: keyof ComponentRoundnessConfig): string {
  const config = getThemeConfig();
  const roundness = config.roundness;

  const mapping = {
    button: {
      edgy: "rounded-none",
      mid: "rounded-md",
      full: "rounded-full",
    },
    // ... other components
  };

  return mapping[component][roundness];
}
```

### Responsive Layout Hook

```typescript
// src/lib/theme-hooks.ts
"use client";

export function useResponsiveLayout() {
  const [isMdOrSmaller, setIsMdOrSmaller] = useState(false);
  const themeConfig = getThemeConfig();

  useEffect(() => {
    const mediaQuery = window.matchMedia("(max-width: 768px)");
    setIsMdOrSmaller(mediaQuery.matches);

    const handler = (e: MediaQueryListEvent) => setIsMdOrSmaller(e.matches);
    mediaQuery.addEventListener("change", handler);

    return () => mediaQuery.removeEventListener("change", handler);
  }, []);

  const isSideBySide = themeConfig.layout === "side-by-side" && !isMdOrSmaller;

  return { isSideBySide, isResponsiveOverride: isMdOrSmaller };
}
```

### Branding Integration

```typescript
// Fetch organization-specific branding
export async function getBrandingSettings({ serviceUrl, organization }: { serviceUrl: string; organization?: string }) {
  const settingsService = await createServiceForHost(SettingsService, serviceUrl);

  return settingsService.getBrandingSettings({
    ctx: makeReqCtx(organization),
  });
}
```

### Performance Optimizations

#### Caching Strategy

```typescript
// src/lib/zitadel.ts
const useCache = process.env.DEBUG !== "true";

async function cacheWrapper<T>(callback: Promise<T>) {
  "use cache";  // Next.js 15 cache directive
  cacheLife("hours");  // Cache for 1 hour

  return callback;
}

export async function getLoginSettings({
  serviceUrl,
  organization,
}: {
  serviceUrl: string;
  organization?: string;
}) {
  const callback = settingsService.getLoginSettings({...});

  return useCache ? cacheWrapper(callback) : callback;
}
```

**Cached Resources**:

- Login settings
- Branding settings
- Security settings
- Translations

**Cache Invalidation**:

- Disabled in development (`DEBUG=true`)
- Automatic TTL-based invalidation (1 hour)
- Organization-specific cache keys

### gRPC Connection Pooling

```typescript
// Reuse gRPC transport connections
let transportCache: Map<string, Transport> = new Map();

export function createServerTransport(token: string, baseUrl: string) {
  const cacheKey = `${token}:${baseUrl}`;

  if (transportCache.has(cacheKey)) {
    return transportCache.get(cacheKey);
  }

  const transport = libCreateServerTransport(token, {
    baseUrl,
    httpVersion: "2",
  });

  transportCache.set(cacheKey, transport);
  return transport;
}
```

### Server Component Strategy

```typescript
// Default: Server Components (no "use client")
export default async function LoginPage() {
  // Rendered on server, no JavaScript shipped to client
  const branding = await getBrandingSettings({...});

  return (
    <div>
      {/* Static HTML rendered on server */}
    </div>
  );
}

// Client Components only when necessary
"use client";  // ← Only for interactive components
export function PasswordForm() {
  // Client-side interactivity
}
```

## Error Handling

### Server Action Error Pattern

```typescript
// Always return serializable error/success objects
export async function sendPassword(
  command: {...}
): Promise<{ error: string } | { redirect: string }> {
  try {
    const session = await setSessionAndUpdateCookie({...});
    return { redirect: "/signedin" };
  } catch (error) {
    if ("failedAttempts" in error) {
      return {
        error: t("errors.failedToAuthenticate", {
          failedAttempts: error.failedAttempts,
        }),
      };
    }
    return { error: t("errors.unknownError") };
  }
}
```

### gRPC Error Handling

```typescript
import { ConnectError } from "@zitadel/client";

try {
  await setPassword({ serviceUrl, payload });
} catch (error: ConnectError) {
  // gRPC error codes
  if (error.code === 7) {
    // PERMISSION_DENIED
    return { error: t("errors.sessionNotValid") };
  }
  if (error.code === 9) {
    // FAILED_PRECONDITION
    return { error: t("errors.failedPrecondition") };
  }

  // Handle error details
  const details = error.findDetails(CredentialsCheckErrorSchema);
  if (details[0]?.failedAttempts) {
    // Handle lockout
  }
}
```

### Fallback Error Pages

```typescript
// app/global-error.tsx
export default function GlobalError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <html>
      <body>
        <h2>Something went wrong!</h2>
        <button onClick={() => reset()}>Try again</button>
      </body>
    </html>
  );
}
```

## Monitoring and Debugging

### Logging Strategy

```typescript
// Structured logging for debugging
console.log("Password auth: OIDC/SAML flow with requestId:", requestId);
console.log("Session is valid:", isValid);
console.warn("No session cookie found, returning empty array");
console.error("Error getting session:", error);
```

### Health Check Endpoint

```typescript
// app/healthy/page.tsx
export default function HealthyPage() {
  return new Response("OK", { status: 200 });
}
```

### Debug Mode

```bash
# .env.local
DEBUG=true  # Disables caching, enables verbose logging
```

## Cross-Layer Security Analysis

### Overview

The ZITADEL Login application operates across three architectural layers that interact to provide authentication services. Security vulnerabilities can arise not just within individual layers, but in the interactions and assumptions between them. This section defines security boundaries, attack vectors, and acceptance criteria to ensure changes in any layer don't compromise the login application's security.

### Architectural Layers

```
┌───────────────────────────────────────┐       ┌───────────────────────────────────────┐
│        Presentation Layer             │       │         Service Layer                 │
│    (Login UI - Next.js App)           │       │     (ZITADEL Backend APIs)            │
│                                       │       │                                       │
│  • User Interface & Forms             │       │  • Session Service                    │
│  • Server Actions (auth logic)        │       │  • User Service                       │
│  • Session Cookie Management          │◄─────►│  • Settings Service                   │
│  • WebAuthn/Passkey flows             │       │  • Auth Service (OIDC/SAML)           │
│  • URL parameter handling             │       │  • Notification Service               │
└───────────────┬───────────────────────┘       └───────────────┬───────────────────────┘
                │                                               │
                │ gRPC/HTTPS                                    │ SQL/gRPC
                │ (JWT/PAT Auth)                                │ (Queries)
                │                                               │
                ▼                                               ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           Infrastructure Layer                                  │
│                        (Shared Platform Services)                               │
│                                                                                 │
│  ┌──────────────────────────┐  ┌──────────────────────────┐                    │
│  │   Network Services       │  │   Secret Management      │                    │
│  │                          │  │                          │                    │
│  │  • TLS/HTTPS             │  │  • Key Vault             │                    │
│  │  • Load Balancers        │  │  • Token Storage         │                    │
│  │  • Reverse Proxies       │  │  • Certificates          │                    │
│  │  • Firewalls             │  │  • Secret Rotation       │                    │
│  │  • DDoS Protection       │  │                          │                    │
│  └──────────────────────────┘  └──────────────────────────┘                    │
│                                                                                 │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │                      Logging & Monitoring                                 │  │
│  │         • Security Event Logging  • Metrics  • Alerting                   │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────────┘

Communication Flows:
━━━━━━━━━━━━━━━━

1. Presentation ◄────► Service
   • Protocol: gRPC over HTTPS
   • Authentication: JWT (multi-tenant) or PAT (self-hosted)
   • Purpose: API calls for sessions, users, settings, auth flows

2. Presentation ──► Infrastructure
   • Uses: TLS termination, load balancing, proxy routing
   • Gets: Secrets (JWT/PAT keys from secret management)
   • Protected by: Network security, firewall rules

3. Service ──► Infrastructure
   • Uses: Database connections (SQL), gRPC infrastructure
   • Gets: Secrets (database credentials, encryption keys)
   • Protected by: Network segmentation, access controls

Architecture Notes:
━━━━━━━━━━━━━━━━━

• Presentation and Service are PEER layers (same level)
• Infrastructure is a SHARED FOUNDATION supporting both
• All external traffic flows through Infrastructure first
• No direct database access from Presentation layer
• All cross-layer communication is authenticated and encrypted
```

### Layer-Specific Security Responsibilities

#### Presentation Layer (Login UI)

**Security Responsibilities:**

1. **Input Validation & Sanitization**
   - Validate all user inputs before processing
   - Sanitize display data (React handles most XSS automatically)
   - Validate URL parameters (requestId, organization, loginName, etc.)

2. **Session Management**
   - Store sessions only in HTTP-only cookies
   - Never expose session tokens to client-side JavaScript
   - Validate session expiration timestamps
   - Clean up expired sessions

3. **CSRF Protection**
   - Rely on Next.js Server Actions' built-in Origin validation
   - Ensure all state-changing operations use POST
   - Maintain SameSite cookie policy

4. **Redirect Validation**
   - Validate all redirect URIs against trusted domains
   - Never blindly redirect based on user input
   - Use backend-validated redirect URIs for OIDC/SAML

5. **Authentication Context**
   - Never trust client-provided authentication state
   - Always verify session with backend before granting access
   - Validate MFA requirements server-side

6. **Error Handling**
   - Never expose internal error details to users
   - Log security-relevant errors server-side
   - Provide generic error messages to prevent information leakage

**Trust Assumptions:**

- ✅ Service Layer correctly validates authentication factors
- ✅ Service Layer enforces security policies (MFA, password complexity)
- ✅ Service Layer validates session tokens
- ✅ Infrastructure provides TLS/HTTPS
- ❌ NEVER trust client-provided data without validation
- ❌ NEVER trust cookies without backend verification

#### Service Layer (ZITADEL Backend)

**Security Responsibilities:**

1. **Authentication & Authorization**
   - Validate all gRPC authentication tokens (JWT/PAT)
   - Enforce role-based access control (e.g., IAM_LOGIN_CLIENT)
   - Validate session tokens before operations
   - Verify user permissions for requested operations

2. **Session Validation**
   - Validate session existence and expiration
   - Enforce authentication factor requirements
   - Check MFA/2FA policies
   - Validate password expiry settings

3. **Security Policy Enforcement**
   - Enforce password complexity requirements
   - Apply account lockout policies after failed attempts
   - Validate MFA enrollment requirements
   - Check email verification requirements

4. **Request Validation**
   - Validate OIDC/SAML auth requests
   - Verify client credentials and redirect URIs
   - Validate callback URLs against registered clients
   - Sanitize and validate all input data

5. **Host Header Validation**
   - Validate host headers in URL templates for emails
   - Maintain trusted domain allowlist
   - Log suspicious host header attempts

6. **Rate Limiting**
   - Implement rate limiting for authentication attempts
   - Throttle password reset requests
   - Limit session creation rate per user/IP

**Trust Assumptions:**

- ✅ Presentation Layer uses authenticated gRPC (JWT/PAT)
- ✅ Infrastructure provides secure database connections
- ✅ Infrastructure protects secrets (private keys, tokens)
- ❌ NEVER trust presentation layer to enforce security policies
- ❌ NEVER trust presentation layer's session validation
- ❌ NEVER trust host headers without validation

#### Infrastructure Layer

**Security Responsibilities:**

1. **Network Security**
   - Enforce TLS/HTTPS for all external communications
   - Implement network segmentation between layers
   - Configure firewall rules (only necessary ports open)
   - Protect against DDoS attacks

2. **Header Validation & Sanitization**
   - **Strip untrusted headers from client requests** (e.g., `X-Forwarded-*`, `Host`)
   - **Set headers based on validated routing logic only** (never trust client values)
   - Validate Origin headers for cross-origin requests
   - Configure proxy to add authenticated headers (e.g., `x-zitadel-forward-host` in multi-tenant)
   - Prevent header injection attacks

3. **Proxy/Load Balancer Configuration**
   - Set `X-Forwarded-For` and `X-Forwarded-Host` based on actual client IP and routing
   - Configure request size limits and timeouts
   - Implement rate limiting at proxy level
   - Ensure headers are forwarded correctly to application layers

4. **Secret Management**
   - Securely store private keys and tokens
   - Use secret management systems (Kubernetes secrets, AWS Secrets Manager)
   - Never log secrets

5. **Logging & Monitoring**
   - Log security events (failed logins, suspicious activity)
   - Monitor for anomalous patterns
   - Alert on security-relevant events
   - Securely store and protect logs

**Trust Assumptions:**

- ✅ Service Layer implements authentication/authorization
- ✅ Presentation Layer handles user interaction securely
- ❌ NEVER expose internal services directly to internet
- ❌ NEVER trust external network traffic
- ❌ NEVER trust client-provided headers (always strip and set based on validated routing)

### Cross-Layer Attack Vectors

#### 1. Session Fixation/Hijacking

**Attack:** Attacker attempts to reuse or steal session tokens across layers.

**Layer Interactions:**

- Presentation Layer: Stores session tokens in HTTP-only cookies
- Service Layer: Validates session tokens and enforces expiration
- Infrastructure Layer: Protects cookies with TLS

**Mitigations:**

- ✅ HTTP-only cookies (Presentation)
- ✅ Session token validation on every request (Service)
- ✅ Session expiration enforcement (Service)
- ✅ TLS/HTTPS for all communications (Infrastructure)
- ✅ Session rotation on authentication level changes (Service)

**Checklist for Changes:**

- [ ] Does this change expose session tokens to client-side JavaScript?
- [ ] Does session validation still occur on every privileged operation?
- [ ] Are session expiration timestamps checked?
- [ ] Is TLS/HTTPS enforced for cookie transmission?

#### 2. Open Redirect Vulnerabilities

**Attack:** Attacker manipulates redirect parameters to send users to malicious sites.

**Layer Interactions:**

- Presentation Layer: Accepts redirect URIs from query parameters
- Service Layer: Validates redirect URIs against client registrations
- Infrastructure Layer: May modify or forward redirect headers

**Mitigations:**

- ✅ Backend validates OIDC/SAML redirect URIs against client config (Service)
- ✅ Frontend validates post-logout redirects against trusted domains (Presentation)
- ✅ Never blindly redirect based on user input (Presentation)
- ✅ Host header validation for email links (Service)

**Checklist for Changes:**

- [ ] Are redirect URIs validated against an allowlist?
- [ ] Does the change introduce new redirect parameters?
- [ ] Are redirect URIs sanitized (no javascript:, data:, file: protocols)?
- [ ] Is the host header validated before constructing URLs?
- [ ] Are email/notification links constructed with validated hosts?

#### 3. Authentication Bypass

**Attack:** Attacker bypasses authentication checks by exploiting layer assumptions.

**Layer Interactions:**

- Presentation Layer: Creates/validates sessions via gRPC
- Service Layer: Enforces authentication policies
- Infrastructure Layer: Routes requests between layers

**Mitigations:**

- ✅ Backend validates all authentication factors (Service)
- ✅ Frontend never trusts client-provided auth state (Presentation)
- ✅ MFA requirements enforced server-side (Service)
- ✅ Session validation on every privileged operation (Service)
- ✅ Network segmentation prevents direct database access (Infrastructure)

**Checklist for Changes:**

- [ ] Are authentication checks performed server-side?
- [ ] Can client-side code bypass authentication logic?
- [ ] Are MFA requirements enforced by the backend?
- [ ] Is session validation still required for this operation?
- [ ] Does the change introduce new authentication paths?

#### 4. Authorization Bypass

**Attack:** Attacker accesses resources without proper authorization.

**Layer Interactions:**

- Presentation Layer: Uses service account token (IAM_LOGIN_CLIENT) for operations
- Service Layer: Validates token permissions and user context
- Infrastructure Layer: Doesn't perform authorization

**Mitigations:**

- ✅ Service account has minimal permissions (IAM_LOGIN_CLIENT only)
- ✅ User context validated for self-service operations (Service)
- ✅ Organization/resource ownership checked (Service)
- ✅ Session token used for user-context operations (Presentation)

**Checklist for Changes:**

- [ ] Does this operation require user context or system context?
- [ ] Are organization/resource ownership checks in place?
- [ ] Is the correct token used (service account vs user session)?
- [ ] Can users access resources belonging to other users/organizations?
- [ ] Are permissions checked at the service layer?

#### 5. Injection Attacks (SQL, Command, LDAP)

**Attack:** Attacker injects malicious code through untrusted inputs.

**Layer Interactions:**

- Presentation Layer: Accepts user input, sanitizes for display
- Service Layer: Validates and sanitizes inputs, uses parameterized queries
- Infrastructure Layer: Database enforces query structure

**Mitigations:**

- ✅ React automatically escapes JSX content (Presentation)
- ✅ gRPC protocol buffers prevent injection (Presentation → Service)
- ✅ Backend uses parameterized queries (Service)
- ✅ Input validation on all fields (Service)

**Checklist for Changes:**

- [ ] Are all user inputs validated and sanitized?
- [ ] Does the code use parameterized queries (no string concatenation)?
- [ ] Are inputs properly escaped for display?
- [ ] Does the change introduce new input fields?
- [ ] Are search/filter parameters properly validated?

#### 6. Host Header Poisoning

**Attack:** Attacker manipulates host header to poison cached content or generate malicious links.

**Layer Interactions:**

- Presentation Layer: Reads host header for URL construction
- Service Layer: Validates host header against trusted domains
- Infrastructure Layer: Forwards/modifies host headers (proxy)

**Mitigations:**

- ✅ Service layer validates host headers for email templates (Service)
- ✅ Trusted domains maintained in backend configuration (Service)
- ✅ Proxy configuration sets X-Forwarded-Host correctly (Infrastructure)
- ⚠️ Password reset/verification links use validated hosts (Service)

**Checklist for Changes:**

- [ ] Does this change construct URLs using host headers?
- [ ] Is the host header validated against trusted domains?
- [ ] Are email/notification templates affected?
- [ ] Does the change trust X-Forwarded-Host headers?
- [ ] Is proxy configuration documented for this change?

#### 7. CSRF Attacks

**Attack:** Attacker tricks user into performing unwanted actions while authenticated.

**Layer Interactions:**

- Presentation Layer: Uses Next.js Server Actions (automatic Origin validation)
- Service Layer: Doesn't perform CSRF checks (trusts presentation layer)
- Infrastructure Layer: Forwards headers correctly

**Mitigations:**

- ✅ Server Actions use POST-only (Presentation)
- ✅ Automatic Origin validation (Presentation)
- ✅ SameSite cookie policy (Presentation)
- ✅ Proxy forwards Origin/Host headers correctly (Infrastructure)

**Checklist for Changes:**

- [ ] Are state-changing operations POST-only?
- [ ] Does the change use Server Actions or custom Route Handlers?
- [ ] If custom Route Handlers, is CSRF protection implemented?
- [ ] Are SameSite cookie settings maintained?
- [ ] Does proxy configuration preserve Origin headers?

#### 8. Information Disclosure

**Attack:** Attacker gains sensitive information through error messages, logs, or responses.

**Layer Interactions:**

- Presentation Layer: Displays error messages to users
- Service Layer: Generates errors with internal details
- Infrastructure Layer: Logs requests and responses

**Mitigations:**

- ✅ Generic error messages shown to users (Presentation)
- ✅ Detailed errors logged server-side only (Service)
- ✅ Session tokens never logged (All layers)
- ✅ Secrets masked in logs (Infrastructure)

**Checklist for Changes:**

- [ ] Are error messages generic and safe for users?
- [ ] Does the change log sensitive information?
- [ ] Are session tokens/passwords excluded from logs?
- [ ] Does the response include unnecessary internal details?
- [ ] Are stack traces hidden from users in production?

### PR Acceptance Criteria for Security

When implementing features or changes in any layer, the following security acceptance criteria must be met:

#### For All Changes

- [ ] **Threat Model Updated**: Document how this change affects the threat model
- [ ] **Layer Boundaries Respected**: Change doesn't violate trust assumptions between layers
- [ ] **Security Review Completed**: Change reviewed by security-conscious team member
- [ ] **Tests Include Security Cases**: Tests cover security-relevant edge cases (auth bypass, injection, etc.)
- [ ] **Documentation Updated**: ARCHITECTURE.md and relevant docs reflect security implications

#### Presentation Layer Changes (Login UI)

- [ ] **Input Validation**: All user inputs validated and sanitized
- [ ] **Session Handling**: Sessions remain in HTTP-only cookies, never exposed to JavaScript
- [ ] **Authentication Checks**: Authentication verified server-side (never client-side only)
- [ ] **CSRF Protection**: State-changing operations use Server Actions or have CSRF protection
- [ ] **Redirect Validation**: All redirects validated against allowlist or backend configuration
- [ ] **Error Handling**: Error messages don't leak sensitive information
- [ ] **gRPC Authentication**: All backend calls use authenticated transport (JWT/PAT)
- [ ] **No Client-Side Secrets**: No tokens, keys, or secrets in client-side code
- [ ] **Security Headers**: CSP, X-Frame-Options, etc. remain intact
- [ ] **Cookie Configuration**: HTTP-only, Secure, SameSite attributes maintained

**Security Checklist Questions:**

- Does this change trust any client-provided data without backend validation?
- Can an attacker manipulate this flow to bypass authentication?
- Are sessions or tokens exposed where they shouldn't be?
- Could this change enable open redirect attacks?
- Does this change handle errors securely?

#### Service Layer Changes (ZITADEL Backend)

- [ ] **Authentication Required**: All operations validate authentication token (JWT/PAT)
- [ ] **Authorization Enforced**: User permissions checked for requested operations
- [ ] **Input Validation**: All inputs validated against schema and business rules
- [ ] **Session Validation**: Session tokens validated before accepting session operations
- [ ] **Security Policy Enforcement**: MFA, password complexity, lockout policies enforced
- [ ] **Host Header Validation**: Host headers validated if used in URL construction
- [ ] **Rate Limiting**: Rate limits applied to sensitive operations
- [ ] **Audit Logging**: Security-relevant actions logged with sufficient detail
- [ ] **Error Responses**: Error messages don't leak internal details
- [ ] **Parameterized Queries**: No SQL injection vulnerabilities (use parameterized queries)

**Security Checklist Questions:**

- Does this change trust the presentation layer to enforce security policies?
- Are authentication and authorization checks in place?
- Could this change be used to enumerate users or resources?
- Does this change validate redirect URIs or host headers?
- Are rate limits applied to prevent abuse?

#### Infrastructure Layer Changes

- [ ] **TLS/HTTPS Enforced**: All external communications use TLS
- [ ] **Secrets Protected**: Secrets stored securely, never logged
- [ ] **Network Segmentation**: Services properly isolated by network policies
- [ ] **Header Sanitization**: Client-provided headers stripped before setting routing/instance headers
- [ ] **Proxy Configuration**: Headers (X-Forwarded-\*, Origin) set by proxy based on validated routing only
- [ ] **No Client Header Trust**: Instance/routing headers (e.g., x-zitadel-forward-host) never taken from client
- [ ] **Logging Secure**: Logs don't contain secrets, stored securely
- [ ] **Monitoring Configured**: Security events monitored and alerted
- [ ] **Firewall Rules**: Only necessary ports exposed
- [ ] **DDoS Protection**: Anti-DDoS measures in place

**Security Checklist Questions:**

- Does this change expose internal services to the internet?
- Are secrets properly managed and rotated?
- Could this change leak sensitive information through logs?
- Are client-provided headers stripped before setting routing headers?
- Does the proxy set headers based on validated routing logic (not client input)?
- Is TLS/HTTPS maintained throughout the request path?

#### Cross-Layer Changes

When a change spans multiple layers:

- [ ] **Layer Communication Documented**: How layers interact is clearly documented
- [ ] **Trust Boundaries Defined**: What each layer trusts/validates is explicit
- [ ] **Defense in Depth**: Multiple layers provide overlapping security controls
- [ ] **Failure Modes Secure**: If one layer fails, others maintain security
- [ ] **End-to-End Testing**: Security tests cover the full request path
- [ ] **Attack Vector Analysis**: Cross-layer attack vectors considered and mitigated

### Security Review Process

#### Before Implementation

1. **Threat Modeling**: Identify potential attack vectors for the feature
2. **Layer Analysis**: Determine which layers are affected and their security responsibilities
3. **Trust Boundary Review**: Ensure the feature respects layer trust assumptions
4. **Acceptance Criteria**: Define security acceptance criteria from the relevant sections above

#### During Implementation

1. **Secure Coding**: Follow secure coding practices (input validation, error handling, etc.)
2. **Security Tests**: Write tests that cover security edge cases
3. **Code Review**: Include security-focused review by another team member
4. **Documentation**: Update ARCHITECTURE.md with security implications

#### Before Merge

1. **Security Checklist**: Complete all relevant security checklists above
2. **Attack Vector Analysis**: Verify identified attack vectors are mitigated
3. **Layer Boundary Validation**: Confirm trust assumptions still hold
4. **Test Coverage**: Ensure security tests are included and passing
5. **Documentation Review**: Verify security documentation is accurate and complete

### Security Testing Guidelines

#### Unit Tests

- Test input validation (valid, invalid, edge cases, malicious inputs)
- Test authentication/authorization checks
- Test error handling (no information leakage)
- Test session validation logic

#### Integration Tests

- Test cross-layer authentication flows
- Test OIDC/SAML flows with invalid redirect URIs
- Test session hijacking prevention
- Test CSRF protection
- Test host header validation

#### Security Tests

- Test authentication bypass attempts
- Test authorization bypass attempts
- Test injection attacks (SQL, XSS, command)
- Test open redirect vulnerabilities
- Test session fixation/hijacking
- Test CSRF attacks
- Test rate limiting

#### Manual Security Review

- Review for sensitive information in logs
- Review for secrets in code or environment variables
- Review redirect URI validation
- Review error messages shown to users
- Review authentication/authorization flows

### Security Incident Response

If a security vulnerability is discovered:

1. **Identify Affected Layers**: Determine which layers are impacted
2. **Assess Cross-Layer Impact**: Check if other layers have compensating controls
3. **Implement Defense in Depth**: Add security controls in multiple layers
4. **Update Acceptance Criteria**: Add new checklist items based on the vulnerability
5. **Improve Testing**: Add security tests to prevent regression
6. **Document Lessons Learned**: Update ARCHITECTURE.md with new attack vectors

### Key Principles

1. **Defense in Depth**: Multiple layers should have overlapping security controls
2. **Least Privilege**: Each layer should have minimal necessary permissions
3. **Never Trust, Always Verify**: Each layer validates inputs, even from trusted layers
4. **Fail Securely**: When errors occur, fail in a secure state (deny access)
5. **Secure by Default**: Security features enabled by default
6. **Layer Responsibility**: Each layer enforces its own security requirements
7. **Explicit Trust Boundaries**: What each layer trusts is explicitly documented
8. **Audit Everything**: Security-relevant actions are logged and monitored

## Conclusion

The ZITADEL Login application is a Next.js 15 application that:

1. **Establishes secure sessions** using HTTP-only cookies and server-side validation
2. **Supports multiple authentication flows** (password, passkeys, IDPs, MFA)
3. **Handles multiple protocols** (OIDC, SAML, Device Authorization)
4. **Works in multiple deployment modes** (multi-tenant, self-hosted)
5. **Provides strong security** (HTTP-only cookies, CSP, lockout protection)
6. **Scales efficiently** (caching, connection pooling, Server Components)

The architecture leverages Next.js 15's latest features (App Router, Server Actions, Dynamic IO) while maintaining compatibility with ZITADEL's gRPC backend through type-safe protocol buffer communication.

**Security Posture**: This architecture implements defense-in-depth across three layers (Presentation, Service, Infrastructure) with explicit trust boundaries, attack vector analysis, and strict acceptance criteria for all changes. Each layer maintains its security responsibilities while providing overlapping protections against common attack vectors.
