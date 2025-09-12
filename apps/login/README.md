# ZITADEL Login Application

A Next.js-based login application implementing ZITADEL's authentication flow with comprehensive multi-factor authentication support.

## Quick Start

### Prerequisites

- Node.js 18+ and pnpm
- ZITADEL instance running (can be local or remote)

### Development Setup

1. **Install dependencies:**

   ```bash
   pnpm install
   ```

2. **Set environment variables:**
   Create a `.env.local` file in the `apps/login` directory:

   ```env
   # Required: ZITADEL API endpoint
   ZITADEL_API_URL=https://your-zitadel-instance.com
   # For local development:
   # ZITADEL_API_URL=http://localhost:8080

   # Required: Service user token for API authentication
   ZITADEL_SERVICE_USER_TOKEN="your-service-user-token"

   # Optional: Enable email verification check
   EMAIL_VERIFICATION=true

   # Optional: Enable debug mode (prevents auto-redirect from root)
   DEBUG=true

   # Optional: Custom port (default: 3000)
   PORT=3001

   # Optional: Base path for the application (useful for reverse proxies)
   NEXT_PUBLIC_BASE_PATH=/ui/v2/login
   ```

3. **Run the development server:**

   ```bash
   # From project root
   pnpm dev

   # Or specifically for login app
   cd apps/login
   pnpm dev
   ```

4. **Access the application:**
   - Local: `http://localhost:3000`
   - The app will auto-redirect to `/loginname` unless `DEBUG=true`

### Production Build

```bash
# Build the application
pnpm build

# Start production server
pnpm start
```

### Testing

```bash
# Run unit tests
pnpm test

# Run unit tests in watch mode (for development)
pnpm test:watch

# Run unit tests with coverage
pnpm test:coverage

# From project root, run login app tests specifically
pnpm --filter=@zitadel/login test
```

## Application Architecture

This login application is implemented as a **state machine** with multiple pages handling different authentication steps. Each page represents a state in the authentication flow, with transitions based on user actions and system policies.

## Complete Login Flow Diagram

```mermaid
flowchart TD
    Start([User arrives]) --> Root["/"]
    Root --> |Auto-redirect| LoginName["/loginname"]

    LoginName --> |Username submitted| UsernameCheck{Username valid?}
    UsernameCheck --> |Invalid| LoginName
    UsernameCheck --> |Valid + Session exists| SessionValid{Session valid?}
    UsernameCheck --> |Valid + No session| Password["/password"]

    SessionValid --> |Valid| CheckMFA[Check MFA Requirements]
    SessionValid --> |Invalid| Password

    Password --> |Password correct| PasswordSuccess[Password verified]
    Password --> |Password incorrect| Password
    Password --> |Forgot password| ResetFlow[Password reset flow]

    PasswordSuccess --> CheckMFA

    CheckMFA --> |IDP authenticated| Complete[Authentication complete]
    CheckMFA --> |Passkey verified| Complete
    CheckMFA --> |MFA not required| Complete
    CheckMFA --> |Single MFA method| DirectMFA[Direct to specific MFA]
    CheckMFA --> |Multiple MFA methods| MFASelect["/mfa"]
    CheckMFA --> |No methods + forced| MFASetup["/mfa/set"]
    CheckMFA --> |No methods + optional| MFASetupOptional["/mfa/set?force=false"]

    DirectMFA --> |TOTP| TOTP["/otp/time-based"]
    DirectMFA --> |SMS| SMS["/otp/sms"]
    DirectMFA --> |Email| Email["/otp/email"]
    DirectMFA --> |U2F| U2F["/u2f"]

    MFASelect --> |Method selected| DirectMFA

    TOTP --> |Code verified| Complete
    SMS --> |Code verified| Complete
    Email --> |Code verified| Complete
    U2F --> |U2F verified| Complete

    MFASetup --> |Setup completed| MFAVerify[Verify new MFA]
    MFASetupOptional --> |Setup completed| MFAVerify
    MFASetupOptional --> |Skipped| Complete

    MFAVerify --> |Verified| Complete

    Complete --> |Auth request exists| Redirect[Redirect to app]
    Complete --> |No auth request| SignedIn["/signedin"]

    %% Alternative flows
    LoginName --> |Register clicked| Register["/register"]
    LoginName --> |IDP selected| IDP["/idp/[provider]"]

    Register --> |Registration complete| Password
    IDP --> |IDP success| CheckMFA

    %% Account selection
    Start --> |Multiple sessions| Accounts["/accounts"]
    Accounts --> |Session selected| SessionValid
    Accounts --> |Add account| LoginName

    %% Passkey flow
    LoginName --> |Passkey available| Passkey["/passkey"]
    Password --> |Use passkey| Passkey
    Passkey --> |Passkey verified| Complete

    %% Error states
    Password --> |Account locked| Error[Error state]
    TOTP --> |Too many failures| Error
    SMS --> |Too many failures| Error
    Email --> |Too many failures| Error
    U2F --> |Too many failures| Error

    style Complete fill:#90EE90
    style Error fill:#FFB6C1
    style CheckMFA fill:#87CEEB
```

## MFA Enforcement Logic

The MFA enforcement is controlled by the `isSessionValid()` function and `checkMFAFactors()` helper:

### Key Decision Points

| Condition                                | Result                | Next Action                    |
| ---------------------------------------- | --------------------- | ------------------------------ |
| IDP authenticated                        | Bypass all MFA        | Complete authentication        |
| Passkey verified (userVerified=true)     | Bypass additional MFA | Complete authentication        |
| `forceMfa` or `forceMfaLocalOnly` = true | Enforce MFA           | Check available methods        |
| MFA not required by policy               | Skip MFA validation   | Complete authentication        |
| Single MFA method available              | Auto-route            | Direct to specific method page |
| Multiple MFA methods                     | User choice           | Show MFA selection page        |
| No methods + MFA forced                  | Setup required        | Force MFA setup                |
| No methods + MFA optional                | Setup optional        | Allow skip with lifetime       |

### MFA Method Routing

| Available Method | Route                  | Description              |
| ---------------- | ---------------------- | ------------------------ |
| TOTP only        | `/otp/time-based`      | Time-based authenticator |
| SMS only         | `/otp/sms`             | SMS verification         |
| Email only       | `/otp/email`           | Email verification       |
| U2F only         | `/u2f`                 | Hardware security key    |
| Multiple methods | `/mfa`                 | User selection page      |
| None (forced)    | `/mfa/set?force=true`  | Required setup           |
| None (optional)  | `/mfa/set?force=false` | Optional setup with skip |

### Session Validation Rules

```typescript
// From /lib/session.ts - isSessionValid()
const isValid =
  sessionNotExpired && (validPassword || validPasskey || validIDP) && mfaRequirementsMet && emailVerifiedIfRequired;
```

#### MFA Requirements Logic:

1. **IDP Bypass**: If `session.factors.intent.verifiedAt` exists → Skip all MFA checks
2. **Policy Check**: Get `forceMfa` or `forceMfaLocalOnly` from login settings
3. **Method Check**: If MFA required → verify at least one configured method is verified
4. **Verification**: Check verification timestamps for: TOTP, OTP_EMAIL, OTP_SMS, U2F

## Page Responsibilities

### Core Authentication Pages

- **`/loginname`**: Username/email input, user discovery, IDP options
- **`/password`**: Password verification, reset password option
- **`/passkey`**: WebAuthn/passkey authentication

### MFA Pages

- **`/mfa`**: MFA method selection (when multiple methods available)
- **`/otp/[method]`**: OTP verification (time-based, SMS, email)
- **`/u2f`**: U2F/WebAuthn second factor verification
- **`/mfa/set`**: MFA setup/registration flow

### Account Management

- **`/accounts`**: Multiple session selection
- **`/register`**: New user registration
- **`/signedin`**: Post-authentication success page

### Utility Pages

- **`/verify`**: Email/phone verification
- **`/device`**: Device authorization flow
- **`/logout`**: Session termination

## Configuration

### Environment Variables

| Variable                     | Description                        | Default | Required |
| ---------------------------- | ---------------------------------- | ------- | -------- |
| `ZITADEL_API_URL`            | ZITADEL API endpoint               | -       | ✅       |
| `ZITADEL_SERVICE_USER_TOKEN` | Service user token for API auth    | -       | ✅       |
| `EMAIL_VERIFICATION`         | Enforce email verification         | `false` | ❌       |
| `DEBUG`                      | Enable debug mode                  | `false` | ❌       |
| `PORT`                       | Custom port for the application    | `3000`  | ❌       |
| `NEXT_PUBLIC_BASE_PATH`      | Base path for reverse proxy setups | -       | ❌       |

### Login Settings (from ZITADEL)

Key settings that affect the flow:

- `forceMfa`: Require MFA for all users
- `forceMfaLocalOnly`: Require MFA for non-IDP users
- `ignoreUnknownUsernames`: Show generic errors
- `allowRegister`: Enable registration
- `allowExternalIdp`: Enable IDP login
- `mfaInitSkipLifetime`: How long MFA setup can be skipped

## Development

### Project Structure

```
apps/login/
├── src/
│   ├── app/(login)/          # Next.js app router pages
│   │   ├── loginname/        # Username input
│   │   ├── password/         # Password verification
│   │   ├── mfa/              # MFA selection & setup
│   │   ├── otp/[method]/     # OTP verification
│   │   ├── passkey/          # Passkey auth
│   │   └── ...               # Other auth pages
│   ├── components/           # Reusable UI components
│   ├── lib/                  # Utilities & helpers
│   │   ├── session.ts        # Session validation logic
│   │   ├── verify-helper.ts  # MFA routing logic
│   │   └── server/           # Server actions
│   └── messages/             # i18n translations
├── __tests__/                # Unit tests
└── public/                   # Static assets
```

### Key Files

- **`lib/session.ts`**: Contains `isSessionValid()` and session management
- **`lib/verify-helper.ts`**: Contains `checkMFAFactors()` routing logic
- **`lib/server/`**: Server actions for form submissions
- **`components/`**: Reusable UI components for forms and layout

### Testing

The application includes:

- **Unit tests** for session validation logic (`lib/session.test.ts`)
- **Component tests** for UI elements and forms
- **Integration tests** for server actions and API calls

Key test files:

- `lib/session.test.ts` - Tests for `isSessionValid()` and MFA enforcement logic
- `components/**/*.test.tsx` - Component rendering and interaction tests

Run tests with:

```bash
# All tests
pnpm test

# Watch mode for development
pnpm test:watch

# With coverage report
pnpm test:coverage
```

### Debugging

Set `DEBUG=true` to:

- Prevent auto-redirect from root page
- Enable additional console logging
- Show detailed error information

## API Integration

The login app integrates with ZITADEL via:

- **Session API**: Create, update, and validate sessions
- **User API**: User discovery, registration, password management
- **Settings API**: Login policies, branding, security settings
- **Auth API**: OAuth/OIDC flows, device authorization

All API calls are made server-side using the ZITADEL client libraries for security.
