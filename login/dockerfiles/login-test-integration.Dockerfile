FROM login-pnpm AS login-test-integration-dependencies

# Install dependencies with proper caching
COPY ./apps/login-test-integration/package.json ./apps/login-test-integration/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile --filter=login-test-integration

FROM cypress/factory:5.10.0 AS login-test-integration
WORKDIR /opt/app

# Copy built dependencies
COPY --from=login-test-integration-dependencies /build/apps/login-test-integration .

# Install Cypress with caching
RUN npm install cypress
RUN npx cypress install

# Copy source code (separate layer for better caching)
COPY ./apps/login-test-integration .

# Run integration tests (equivalent to turbo test:integration)
CMD ["npx", "cypress", "run"]
