# ZITADEL API Documentation

A Next.js application that provides interactive API documentation for ZITADEL services using Scalar API Reference.

## Features

- **Interactive API Documentation**: Browse and test ZITADEL APIs directly in the browser
- **Service Selection**: Easy navigation between different ZITADEL services
- **Auto-generated**: Documentation is automatically generated from proto files
- **Modern UI**: Clean, responsive interface powered by Scalar

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

3. Start the development server:
```bash
pnpm run dev
```

4. Open [http://localhost:3000](http://localhost:3000) in your browser.

## Scripts

- `pnpm run dev` - Start development server
- `pnpm run build` - Build for production (includes generating OpenAPI specs)
- `pnpm run start` - Start production server
- `pnpm run generate` - Generate OpenAPI specifications from proto files
- `pnpm run lint` - Run ESLint

## How it works

1. **Proto Generation**: The app uses the same `plugin-download.sh` script as the main docs to download the `protoc-gen-connect-openapi` plugin
2. **OpenAPI Generation**: `buf generate` is used to convert proto files to OpenAPI 3.x specifications
3. **API Serving**: Next.js API routes serve the generated OpenAPI specs
4. **Rendering**: Scalar API Reference renders the interactive documentation

## Project Structure

```
src/
├── app/
│   ├── api/openapi/          # API routes for serving OpenAPI specs
│   ├── layout.tsx            # Root layout
│   ├── page.tsx              # Home page
│   └── globals.css           # Global styles
└── components/
    └── ApiReference.tsx      # Main Scalar API Reference component
```

## Configuration

- `buf.gen.yaml` - Configure proto to OpenAPI generation
- `base.yaml` - Base OpenAPI configuration
- `next.config.mjs` - Next.js configuration

## Deployment

The app can be deployed to any platform that supports Next.js applications. Make sure to run the build script which includes generating the OpenAPI specifications.

For Vercel deployment:
```bash
pnpm run build
```

The generated OpenAPI specs are included in the build output.
