FROM login-pnpm AS login-test-acceptance-dependencies
COPY ./apps/login-test-acceptance/package.json ./apps/login-test-acceptance/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter=login-test-acceptance \
COPY ./apps/login-test-acceptance ./apps/login-test-acceptance
COPY --from=login-test-acceptance-setup / /
CMD ["pnpm", "test:acceptance"]
