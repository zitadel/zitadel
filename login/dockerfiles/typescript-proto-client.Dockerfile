FROM login-pnpm AS typescript-proto-client
# Copy package.json first for better dependency caching
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
# Copy source code
COPY packages/zitadel-proto ./packages/zitadel-proto
# Generate proto files (equivalent to turbo generate)
RUN cd packages/zitadel-proto && pnpm generate
