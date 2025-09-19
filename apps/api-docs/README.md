# ZITADEL API Documentation

A Next.js application that provides interactive API documentation for ZITADEL services using Scalar API Reference.

## Features

- **Interactive API Documentation**: Browse and test ZITADEL APIs directly in the browser
- **Version Management**: Simple manual configuration of available API versions
- **Service Selection**: Easy navigation between different ZITADEL services (admin, auth, management, etc.)
- **Manual Version Control**: Define exactly which versions to show in `versions.config.simple.json`
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

2. Generate OpenAPI specifications from proto files (current branch only):

```bash
pnpm run generate
```

This command will:

- Generate OpenAPI 3.x specifications from proto files for the current branch

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
- `pnpm run generate` - Generate OpenAPI specifications from proto files (current branch only)
- `pnpm run lint` - Run ESLint

### Generation Scripts

- `pnpm run generate:openapi` - Generate only OpenAPI specifications from proto files
- `pnpm run versions:create` - Generate OpenAPI spec for current branch and save to artifacts
- `pnpm run versions:create <version>` - Add a specific version to the config and generate artifacts

### Version Management

Versions are manually configured in `versions.config.json`. You can:

1. **Add versions manually** by editing `versions.config.json`
2. **Add versions via script** using `pnpm run versions:create <version>`

#### Adding Versions via Script

```bash
# Add a new version to config and generate artifacts
pnpm run versions:create v4.3.0

# This will:
# 1. Try to generate OpenAPI specs for current branch
# 2. Create artifact directory for v4.3.0
# 3. Add v4.3.0 to versions.config.json automatically
```

Example `versions.config.json`:

```json
{
  "versions": [
    {
      "id": "main",
      "name": "Latest (Main)",
      "gitRef": "main",
      "enabled": true,
      "isStable": false
    },
    {
      "id": "v4.2.2",
      "name": "v4.2.2",
      "gitRef": "v4.2.2",
      "enabled": true,
      "isStable": true
    }
  ],
  "settings": {
    "defaultVersion": "v4.2.2"
  }
}
```

## How it works

1. **Proto Generation**: The app uses the same `plugin-download.sh` script as the main docs to download the `protoc-gen-connect-openapi` plugin
2. **OpenAPI Generation**: `buf generate` is used to convert proto files to OpenAPI 3.x specifications
3. **Version Management**: Versions are manually defined in `versions.config.json`
4. **API Serving**: Next.js API routes serve the generated OpenAPI specs with version support
5. **Rendering**: Scalar API Reference renders the interactive documentation with version switching

## Version Management

The application supports multiple API versions through an automated system:

### Automatic Version Detection

- Git tags matching semantic versioning (e.g., `v4.2.2`, `v4.2.1`) are automatically detected
- The `pnpm run generate:versions` script scans git tags and existing artifacts to create `versions.config.json`
- Available versions are determined by the presence of artifacts in `.artifacts/versions/`

### Current Workflow

```bash
# Generate everything (OpenAPI specs + versions config)
pnpm run generate

# Create a version snapshot from current artifacts
pnpm run versions:create v4.2.0

# This will:
# 1. Copy current artifacts to .artifacts/versions/v4.2.0/
# 2. Create metadata.json with git info
# 3. Auto-regenerate versions.config.json
```

### Version Structure

```
.artifacts/
├── openapi3/           # Current/latest artifacts
└── versions/           # Version-specific artifacts
    ├── main/
    │   ├── openapi3/
    │   └── metadata.json
    ├── v4.2.2/
    ├── v4.2.1/
    └── v4.2.0/
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
- `versions.config.json` - Version configuration with available versions (auto-generated by `pnpm run generate:versions`)

## Deployment

The app can be deployed to any platform that supports Next.js applications. Make sure to run the build script which includes generating the OpenAPI specifications.

For Vercel deployment:

```bash
pnpm run build
```

The generated OpenAPI specs are included in the build output.
