FROM cypress/factory AS login-integration-testsuite
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
WORKDIR /opt/app
COPY \
  pnpm-lock.yaml \
  pnpm-workspace.yaml \
  ./
COPY ./apps/login-integration-testsuite/package.json ./apps/login-integration-testsuite/package.json
RUN --mount=type=cache,id=pnpm,target=/pnpm/store \
    pnpm install --frozen-lockfile
RUN pnpm exec cypress install
COPY ./apps/login-integration-testsuite/ .
CMD ["pnpm", "exec", "cypress", "run"]
