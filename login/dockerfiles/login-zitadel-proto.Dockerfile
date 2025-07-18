FROM login-pnpm AS login-zitadel-proto
# Copy package.json first for better dependency caching
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
# Copy source code
COPY packages/zitadel-proto ./packages/zitadel-proto
# Generate @zitadel/proto package - equivalent to turbo generate
RUN cd packages/zitadel-proto && pnpm generate
