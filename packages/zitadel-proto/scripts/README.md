# Export Generation Script

This directory contains the automated export generation script for the `@zitadel/proto` package.

## Overview

The `generate-exports.mjs` script automatically scans all generated proto files and creates explicit exports in `package.json`. This ensures compatibility with older Node.js module resolution strategies while maintaining modern export functionality.

## How it works

1. Scans the `types/` directory for all `.d.ts` files
2. For each file, generates corresponding exports pointing to:
   - `types/` directory for TypeScript definitions
   - `es/` directory for ES modules
   - `cjs/` directory for CommonJS modules
3. Updates the `package.json` exports field
4. Preserves wildcard exports as fallbacks

## Usage

The script runs automatically as part of the proto generation:

```bash
pnpm run generate
```

Or run it separately:

```bash
pnpm run generate:exports
```

## Why this is needed

Modern Node.js supports wildcard exports like `./zitadel/*`, but older module resolution (used in CI environments) requires explicit export paths. This script bridges that gap by automatically generating explicit exports for every proto file.
