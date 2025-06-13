FROM login-pnpm AS typescript-proto-client
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
COPY packages/zitadel-proto packages/zitadel-proto
RUN pnpm generate
