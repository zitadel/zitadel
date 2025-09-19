# ZITADEL API Documentation

A Next.js application that provides interactive API documentation for ZITADEL services using Scalar API Reference.

## Features

- **Interactive API Documentation**: Browse and test ZITADEL APIs directly in the browser
- **Version Management**: Switch between different API versions with a dropdown selector
- **Service Selection**: Easy navigation between different ZITADEL services (admin, auth, management, etc.)
- **Auto-generated**: Documentation is automatically generated from proto files
- **Modern UI**: Clean, responsive interface powered by Scalar
- **Versioned Artifacts**: Organized storage of API specifications by version

## Getting Started

### Prerequisites

- Node.js 18 or later
- pnpm (recommended) or npm

### Installation

1. Install dependencies:

```bash
pnpm install
```

2. Generate OpenAPI specifications from proto files:

```bash
pnpm run generate
```

This command will:
- Generate OpenAPI 3.x specifications from proto files
- Auto-generate `versions.config.json` with available versions from git tags and existing artifacts

3. Start the development server:

```bash
pnpm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Scripts

### Core Scripts
- `pnpm run dev` - Start development server
- `pnpm run build` - Build for production (includes generating OpenAPI specs)
- `pnpm run start` - Start production server
- `pnpm run generate` - Generate OpenAPI specifications from proto files AND auto-generate versions config
- `pnpm run lint` - Run ESLint

### Generation Scripts
- `pnpm run generate:openapi` - Generate only OpenAPI specifications from proto files
- `pnpm run generate:versions` - Generate/update versions.config.json based on available versions and git tags

### Version Management Scripts
- `pnpm run versions:create <version>` - Create a new version snapshot from current artifacts
- `pnpm run versions:list` - List all available versions
- `pnpm run versions:current` - Show current version info
- `pnpm run organize` - Organize current artifacts into version-specific folders

## How it works

1. **Proto Generation**: The app uses the same `plugin-download.sh` script as the main docs to download the `protoc-gen-connect-openapi` plugin
2. **OpenAPI Generation**: `buf generate` is used to convert proto files to OpenAPI 3.x specifications
3. **Version Management**: `versions.config.json` is auto-generated from git tags and available artifacts
4. **API Serving**: Next.js API routes serve the generated OpenAPI specs with version support
5. **Rendering**: Scalar API Reference renders the interactive documentation with version switching

## Version Management

The application supports multiple API versions:

### Automatic Version Detection
- Git tags matching semantic versioning (e.g., `v2.70.0`, `v4.2.1`) are automatically detected
- The `pnpm run generate:versions` script scans git tags and existing artifacts to create `versions.config.json`
- Available versions are determined by the presence of artifacts in `.artifacts/versions/`

### Creating New Versions
```bash
# Create a version snapshot from current artifacts
pnpm run versions:create v2.73.0

# This will:
# 1. Copy current artifacts to .artifacts/versions/v2.73.0/
# 2. Create metadata.json with git info
# 3. Auto-regenerate versions.config.json
```

### Version Structure
```
.artifacts/
├── openapi3/           # Current/latest artifacts
└── versions/           # Version-specific artifacts
    ├── v2.70.0/
    │   ├── openapi3/
    │   └── metadata.json
    ├── v2.71.9/
    └── v2.72.0-test/
```

## Project Structure

```
src/
├── app/
│   ├── api/
│   │   ├── openapi/[...slug]/  # API routes for serving OpenAPI specs (with version support)
│   │   └── versions/           # API route for version metadata
│   ├── layout.tsx              # Root layout
│   ├── page.tsx                # Home page with version selector
│   └── globals.css             # Global styles
├── components/
│   └── ApiReference.tsx        # Main Scalar API Reference component with version switching
└── scripts/
    ├── generate-versions-config.sh  # Auto-generate versions.config.json
    ├── manage-versions.sh           # Version creation and management
    ├── organize-artifacts.sh        # Organize artifacts by version
    └── vercel-build.sh             # Production build with version support
```

## Configuration

- `buf.gen.yaml` - Configure proto to OpenAPI generation
- `base.yaml` - Base OpenAPI configuration
- `next.config.mjs` - Next.js configuration
- `versions.config.json` - Auto-generated version configuration (created by `pnpm run generate:versions`)

## Deployment

The app can be deployed to any platform that supports Next.js applications. Make sure to run the build script which includes generating the OpenAPI specifications.

For Vercel deployment:

```bash
pnpm run build
```

The generated OpenAPI specs are included in the build output.
