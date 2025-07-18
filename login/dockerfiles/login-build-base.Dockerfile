# Base image for building login components with proper dependency caching
FROM login-pnpm AS login-build-base

# Install root workspace dependencies first (best caching)
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter .

# Copy all package.json files for dependency resolution
COPY packages/*/package.json ./packages/*/
COPY apps/*/package.json ./apps/*/

# Install all dependencies in one layer for better caching
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
