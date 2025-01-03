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
import { createUserServiceClient } from "@zitadel/client";

const userService = createUserServiceClient({
  // Configuration options
});
```

### Utilities

This package also provides various utilities to work with ZITADEL

```ts
import { timestampMs } from "@zitadel/client";

// Example usage
console.log(`${timestampMs(session.creationDate)}`);
```
