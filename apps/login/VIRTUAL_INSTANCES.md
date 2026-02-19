# Virtual Instances Configuration

This document explains how the ZITADEL login application implements virtual instances support for multi-tenant deployments.

## Overview

The login app supports virtual instances through a system token approach that allows it to dynamically serve multiple ZITADEL instances from a single deployment. This is achieved by combining environment variables for system authentication with the `x-zitadel-forward-host` header to determine the target instance.

## How Virtual Instances Work

### System Token Authentication

The login app uses a system user with a private key to generate JWT tokens for authenticating with ZITADEL instances. This allows the login app to:

1. Authenticate as a system user across multiple instances
2. Dynamically determine which instance to communicate with based on incoming requests
3. Route requests to the appropriate ZITADEL instance

### Instance Detection

The target ZITADEL instance is determined by the following priority order:

1. **`x-zitadel-forward-host` header** - Primary method for virtual instances
2. **`host` header** - Fallback for custom domain deployments

## Required Environment Variables

For virtual instances support, you must configure these three environment variables:

```bash
# The system API audience URL (ZITADEL instance URL)
AUDIENCE="https://your-zitadel-instance.example.com"

# The system user ID that has the necessary permissions
SYSTEM_USER_ID="your-system-user-id"

# Base64-encoded private key for the system user (PEM format)
SYSTEM_USER_PRIVATE_KEY="LS0tLS1CRUdJTi..."
```

## Configuration Detection Logic

The login app detects virtual instances mode when all three environment variables are present:

```typescript
// Virtual instances mode - requires all three variables
if (process.env.AUDIENCE && process.env.SYSTEM_USER_ID && process.env.SYSTEM_USER_PRIVATE_KEY) {
  // Use system token for multi-tenant virtual instances
  token = await systemAPIToken();
}
```

## Virtual Instance Routing

### The `x-zitadel-forward-host` Header

This is the key header that enables virtual instances. When present, it tells the login app which ZITADEL instance the request should be routed to:

```http
x-zitadel-forward-host: customer1.zitadel.app
```

The service URL determination logic:

```typescript
const forwardedHost = headers.get("x-zitadel-forward-host");
if (forwardedHost) {
  instanceUrl = forwardedHost.startsWith("http://") ? forwardedHost : `https://${forwardedHost}`;
}
```

### URL Construction

The login app also uses this header for constructing redirect URLs and callbacks:

```typescript
const forwardedHost =
  request.headers.get("x-zitadel-forward-host") ?? request.headers.get("x-forwarded-host") ?? request.headers.get("host");
```

## System User Setup

### System Token Setup

System tokens are configured at the ZITADEL runtime level and work superordinate over all instances.

1. **Generate RSA Key Pair**:

   ```bash
   openssl genrsa -traditional -out system-user-1.pem 2048
   openssl rsa -in system-user-1.pem -outform PEM -pubout -out system-user-1.pub
   ```

2. **Configure Runtime Configuration**:
   Add the system user to ZITADEL's runtime configuration:

   ```yaml
   SystemAPIUsers:
     - system-user-1:
         Path: /system-user-1.pub
         # OR use base64 encoded key directly:
         # KeyData: <base64 encoded value of system-user-1.pub>
         Memberships:
           # If no memberships are specified, the user has a membership of type System with the role "SYSTEM_OWNER"
           - MemberType: System
             Roles:
               - "SYSTEM_OWNER"
           # Optional: Add specific instance or organization restrictions
           # - MemberType: IAM
           #   Roles: "IAM_OWNER"
           #   AggregateID: "123456789012345678"
   ```

3. **Environment Variables**:
   ```bash
   export AUDIENCE="https://your-zitadel-instance.example.com"
   export SYSTEM_USER_ID="system-user-1"
   # The private key must be base64 encoded for the login app
   export SYSTEM_USER_PRIVATE_KEY=$(base64 -i system-user-1.pem)
   ```

### System User Capabilities

System users operate across all instances and organizations. If no memberships are specified, the system user automatically gets System membership with "SYSTEM_OWNER" role.

### JWT Token Generation

The login app uses the ZITADEL client library to create system tokens:

```typescript
// From src/lib/api.ts
export async function systemAPIToken() {
  const token = {
    audience: process.env.AUDIENCE,
    userID: process.env.SYSTEM_USER_ID,
    token: Buffer.from(process.env.SYSTEM_USER_PRIVATE_KEY, "base64").toString("utf-8"),
  };

  return newSystemToken({
    audience: token.audience,
    subject: token.userID,
    key: token.token,
  });
}
```

### Security Considerations

1. **Private Key Protection**: Ensure the private key is securely stored and not exposed in logs
2. **System User Permissions**: Follow the principle of least privilege for the system user
3. **Network Security**: Ensure secure communication between the login app and ZITADEL instances
4. **Header Validation**: Validate the `x-zitadel-forward-host` header to prevent header injection attacks

## Troubleshooting

### Common Issues

1. **"No token found" Error**
   - Check that either virtual instances variables or service token is set
   - Verify private key is properly base64 encoded

2. **"Service URL could not be determined" Error**
   - Ensure `x-zitadel-forward-host` header is set for virtual instances

3. **Authentication Failures**
   - Verify the system user exists and has proper permissions
   - Check that the private key matches the system user
   - Ensure the audience URL is correct

### Debug Information

The login app logs service detection decisions. Check your application logs for:

- Service URL determination
- Token type selection (system vs service token)
- Header processing information

## Example Configuration

### Virtual Instances

```bash
# .env file for virtual instances
AUDIENCE=https://api.zitadel.cloud
SYSTEM_USER_ID=system-user-1
SYSTEM_USER_PRIVATE_KEY=LS0tLS1CRUdJTi...
```

This configuration enables the ZITADEL login application to seamlessly handle multiple virtual instances while maintaining security and performance.
