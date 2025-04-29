# ZITADEL Proto

This package provides the Protocol Buffers (proto) definitions used by ZITADEL projects. It includes the proto files and generated code for interacting with ZITADEL's gRPC APIs.

## Installation

To install the package, use npm or yarn:

```sh
npm install @zitadel/proto
# or
yarn add @zitadel/proto
# or
pnpm add @zitadel/proto
```

## Usage

This package supports both ESM and CommonJS imports. The API is organized into version-specific namespaces: `v1`, `v2`, and `v3alpha`.

### ESM (ECMAScript Modules)

```typescript
// Import the entire package
import * as zitadel from "@zitadel/proto";

// Use the version-specific namespaces
const userRequest = new zitadel.v1.user.GetUserRequest();

// Or import specific versions
import { v2 } from "@zitadel/proto";
const userServiceRequest = new v2.user_service.GetUserRequest();
```

### CommonJS

```typescript
// Import the entire package
const zitadel = require("@zitadel/proto");

// Use the version-specific namespaces
const userRequest = new zitadel.v1.user.GetUserRequest();
```

## API Structure

The package is organized into version-specific namespaces:

- `v1`: Contains the original ZITADEL API
- `v2`: Contains the newer version of the API with improved organization
- `v3alpha`: Contains the alpha version of the upcoming API

## Package Structure

The package is organized as follows:

- `index.ts`: Main entry point that exports the version-specific APIs
- `v1.ts`: Exports all v1 API modules
- `v2.ts`: Exports all v2 API modules
- `v3alpha.ts`: Exports all v3alpha API modules
- `zitadel/`: Contains the generated proto files

## Development

### Generating the proto files

The proto files are generated from the ZITADEL API definitions using [buf](https://buf.build/).

```sh
pnpm generate
```

### Building the package

```sh
pnpm build
```

### Testing

To test both ESM and CommonJS imports:

```sh
pnpm test
```

Or test them individually:

```bash
pnpm test:cjs  # Test CommonJS imports
pnpm test:esm  # Test ESM imports
```

## Documentation

For detailed documentation and API references, please visit the [ZITADEL documentation](https://zitadel.com/docs).

## Contributing

Contributions are welcome! Please read the contributing guidelines before getting started.
