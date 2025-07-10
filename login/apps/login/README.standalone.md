# ZITADEL Login UI - Standalone

This is the standalone version of the ZITADEL Login UI, a Next.js application that provides the authentication interface for ZITADEL.

## Quick Start

### Prerequisites

- Node.js 18+ 
- npm or pnpm

### Setup

1. **Prepare the standalone environment:**
   ```bash
   ./prepare-standalone.sh
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

## Differences from Monorepo Version

This standalone version includes:

- Self-contained configuration files
- Published versions of `@zitadel/client` and `@zitadel/proto` packages
- Standalone build scripts
- Independent dependency management

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
