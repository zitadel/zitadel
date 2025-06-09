FROM login-pnpm AS zitadel-proto
COPY packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
COPY packages/zitadel-proto packages/zitadel-proto
RUN pnpm generate

FROM scratch AS typescript-proto-client
COPY --from=zitadel-proto /build/packages/zitadel-proto /
