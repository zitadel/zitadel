FROM login-base AS zitadel-proto
COPY packages/zitadel-proto packages/zitadel-proto
RUN pnpm generate

FROM scratch
COPY --from=zitadel-proto /build/packages/zitadel-proto /
