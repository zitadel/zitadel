# ZITADEL Client

This package exports services and utilities to interact with ZITADEL

## Installation

To install the package, use npm or yarn:

```sh
npm install @zitadel/client
```

or

```sh
yarn add @zitadel/client
```

## Usage

### Importing Services

You can import and use the services provided by this package to interact with ZITADEL.

```ts
import { createSettingsServiceClient, makeReqCtx } from "@zitadel/client/v2";
import { createServerTransport } from "@zitadel/client/node";

// Example usage
const transport = createServerTransport(process.env.ZITADEL_SERVICE_USER_TOKEN!, { baseUrl: process.env.ZITADEL_API_URL! });

const settingsService = createSettingsServiceClient(transport);

settingsService.getBrandingSettings({ ctx: makeReqCtx("orgId") }, {});
```

### Utilities

This package also provides various utilities to work with ZITADEL

```ts
import { timestampMs } from "@zitadel/client";

// Example usage
console.log(`${timestampMs(session.creationDate)}`);
```

## Documentation

For detailed documentation and API references, please visit the [ZITADEL documentation](https://zitadel.com/docs).

## Contributing

Contributions are welcome! Please read the contributing guidelines before getting started.
