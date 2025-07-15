# ZITADEL Proto

This package provides the Protocol Buffers (proto) definitions used by ZITADEL projects. It includes the proto files and generated code for interacting with ZITADEL's gRPC APIs.

## Installation

To install the package, use npm or yarn:

```sh
npm install @zitadel/proto
```

or

```sh
yarn add @zitadel/proto
```

## Usage

To use the proto definitions in your project, import the generated code:

```ts
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";

const org: Organization | null = await getDefaultOrg();
```

## Documentation

For detailed documentation and API references, please visit the [ZITADEL documentation](https://zitadel.com/docs).

## Contributing

Contributions are welcome! Please read the contributing guidelines before getting started.
