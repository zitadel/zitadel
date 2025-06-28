FROM login-pnpm AS login-test-acceptance-dependencies
COPY ./apps/login-test-acceptance/package.json ./apps/login-test-acceptance/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter=login-test-acceptance && \
    cd apps/login-test-acceptance && \
    pnpm exec playwright install --with-deps chromium
COPY ./apps/login-test-acceptance ./apps/login-test-acceptance
CMD ["bash", "-c", "cd apps/login-test-acceptance && pnpm test:acceptance test"]
