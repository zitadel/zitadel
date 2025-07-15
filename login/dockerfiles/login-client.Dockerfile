FROM typescript-proto-client AS login-client
# Copy package.json first for better dependency caching
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./packages/zitadel-client
# Copy source code
COPY packages/zitadel-client ./packages/zitadel-client
# Build the client (equivalent to turbo build for zitadel-client)
RUN cd packages/zitadel-client && pnpm build
