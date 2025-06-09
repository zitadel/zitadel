FROM login-pnpm AS zitadel-test-unit-build
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
COPY packages/zitadel-tsconfig packages/zitadel-tsconfig
WORKDIR /build/packages/zitadel-client
COPY packages/zitadel-client .
COPY --from=typescript-proto-client / /build/packages/zitadel-proto
RUN pnpm build

FROM login-dev-base AS zitadel-test-unit
COPY packages/zitadel-tsconfig packages/zitadel-tsconfig
COPY --from=zitadel-test-unit-build /build/packages/zitadel-client/dist /build/packages/zitadel-client/dist
COPY apps/login apps/login
