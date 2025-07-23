FROM login-build-base AS login-client-dependencies
# Copy package.json first for better dependency caching
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./packages/zitadel-client

FROM login-zitadel-proto AS login-zitadel-client
# Copy dependencies from build base
COPY --from=login-client-dependencies /build/node_modules ./node_modules
COPY --from=login-client-dependencies /build/packages/zitadel-client/node_modules ./packages/zitadel-client/node_modules

# Copy generated proto files (@zitadel/proto package)
COPY --from=login-zitadel-proto /build/packages/zitadel-proto ./packages/zitadel-proto

# Copy source code
COPY packages/zitadel-client ./packages/zitadel-client
# Build the @zitadel/client package (equivalent to turbo build for zitadel-client)
RUN cd packages/zitadel-client && pnpm build
