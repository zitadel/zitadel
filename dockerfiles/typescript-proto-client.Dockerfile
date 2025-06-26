FROM login-pnpm AS typescript-proto-client
COPY ./login/packages/zitadel-proto/package.json ./packages/zitadel-proto/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter zitadel-proto
COPY --from=proto-files /buf.yaml /buf.lock /proto-files/
COPY --from=proto-files /zitadel /proto-files/zitadel
COPY ./login/packages/zitadel-proto/buf.gen.yaml ./packages/zitadel-proto/
RUN cd packages/zitadel-proto && pnpm exec buf generate /proto-files
