FROM typescript-proto-client AS login-client
COPY packages/zitadel-client/package.json ./packages/zitadel-client/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --workspace-root --filter ./packages/zitadel-client
COPY packages/zitadel-client ./packages/zitadel-client
RUN cd packages/zitadel-client && pnpm build
