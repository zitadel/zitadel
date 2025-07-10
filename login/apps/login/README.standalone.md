# ZITADEL Login UI - Standalone

This is the standalone version of the ZITADEL Login UI, a Next.js application that provides the authentication interface for ZITADEL.

## Quick Start

### Prerequisites

- Node.js 18+
- npm or pnpm

### Setup

1. **Prepare the standalone environment:**

   ```bash
   ./prepare-standalone.sh --install
   ```

   Or for manual control:

   ```bash
   ./prepare-standalone.sh --no-install
   npm install
   ```

2. **Start development server:**

   ```bash
   npm run dev
   ```

3. **Build for production:**
   ```bash
   npm run build:standalone
   npm run start
   ```

## Development

### Available Scripts

- `npm run dev` - Start development server with Turbopack
- `npm run build` - Build for production
- `npm run build:standalone` - Build standalone version with custom base path
- `npm run start` - Start production server
- `npm run test:unit` - Run unit tests
- `npm run lint` - Run linting
- `npm run lint:fix` - Fix linting issues

### Environment Variables

Create a `.env.local` file with your ZITADEL configuration:

```env
ZITADEL_API_URL=https://your-zitadel-instance.com
# Add other required environment variables
```

### Package Management

This standalone version automatically uses the latest published versions of:

- `@zitadel/client` - ZITADEL client library (latest)
- `@zitadel/proto` - ZITADEL protocol definitions (latest)

To update to the latest versions, simply run:

```bash
npm update @zitadel/client @zitadel/proto
```

## Differences from Monorepo Version

This standalone version includes:

- **Published packages**: Uses latest published versions of `@zitadel/client` and `@zitadel/proto`
- **Self-contained configuration**: All configuration files are standalone-ready
- **Simplified dependencies**: Removes monorepo-specific devDependencies and tooling
- **Streamlined build scripts**: Optimized scripts for standalone development
- **Independent dependency management**: No workspace or turbo dependencies

## Architecture

### Dual-Mode Design

This project supports both monorepo and standalone modes:

- **Monorepo mode**: Uses `workspace:*` dependencies for local development
- **Standalone mode**: Uses published npm packages for distribution

### Automatic Conversion

The conversion between modes is handled by intelligent scripts:

1. **`prepare-standalone.sh`**: Main conversion script for end users
2. **`scripts/prepare-standalone.js`**: Advanced preparation with latest package versions
3. **`scripts/build-standalone.js`**: Generates standalone configs without modifying current setup
4. **`scripts/config-manager.js`**: Switches between monorepo and standalone configurations

## Contributing

When contributing to this standalone version:

1. Make changes in the main monorepo first
2. Test changes in the monorepo environment
3. Update the standalone version via subtree push
4. Test the standalone version independently

## Subtree Sync

This repository is maintained as a Git subtree of the main ZITADEL repository.

To sync changes from the main repo:

```bash
# In the main ZITADEL repo
git subtree push --prefix=login/apps/login origin typescript-login-standalone
```

## License

See the main ZITADEL repository for license information.
