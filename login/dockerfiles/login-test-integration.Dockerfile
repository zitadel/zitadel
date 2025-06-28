FROM login-pnpm AS login-test-integration-dependencies
COPY ./apps/login-test-integration/package.json ./apps/login-test-integration/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter=login-test-integration
FROM cypress/factory:5.10.0 AS login-test-integration
WORKDIR /opt/app
COPY --from=login-test-integration-dependencies /build/apps/login-test-integration .
RUN npm install cypress
RUN npx cypress install
COPY ./apps/login-test-integration .
CMD ["npx", "cypress", "run"]
