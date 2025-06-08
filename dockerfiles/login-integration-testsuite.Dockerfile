FROM login-base AS integration-dependencies
COPY \
  pnpm-lock.yaml \
  pnpm-workspace.yaml \
  ./
COPY ./apps/login-integration-testsuite/package.json ./apps/login-integration-testsuite/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --no-frozen-lockfile --filter=login-integration-testsuite

FROM cypress/factory AS login-integration-testsuite
WORKDIR /opt/app
COPY --from=integration-dependencies /build/apps/login-integration-testsuite .
RUN npm install cypress
RUN npx cypress install
COPY ./apps/login-integration-testsuite .
CMD ["npx", "cypress", "run"]
