# Inspired by https://pnpm.io/docker#example-3-build-on-cicd
# Inspired by https://pnpm.io/docker#minimizing-docker-image-size-and-build-time

FROM node:20-slim AS base

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN apt-get update
RUN apt-get install -y git
RUN npm install -g corepack
RUN corepack enable
RUN corepack prepare pnpm@latest --activate
RUN pnpm install turbo@^2 --global

FROM base AS builder
# Set working directory
WORKDIR /app
# Replace <your-major-version> with the major version installed in your repository. For example:
RUN pnpm install turbo@^2 --global
COPY . .
 
# Generate a partial monorepo with a pruned lockfile for a target workspace.
# Assuming "web" is the name entered in the project's package.json: { name: "web" }
RUN turbo prune @zitadel/login --docker
 
# Add lockfile and package.json's of isolated subworkspace
FROM base AS installer

WORKDIR /app
 
# First install the dependencies (as they change less often)
COPY --from=builder /app/out/json/ .
RUN pnpm install --frozen-lockfile
 
# Build the project
COPY --from=builder /app/out/full/ .

RUN turbo run build

FROM base AS runner
WORKDIR /app

# Don't run production as root
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs
USER nextjs

# Automatically leverage output traces to reduce image size
# https://nextjs.org/docs/advanced-features/output-file-tracing
COPY --from=installer --chown=nextjs:nodejs /app/apps/login/.next/standalone ./
COPY --from=installer --chown=nextjs:nodejs /app/apps/login/.next/static ./apps/login/.next/static
COPY --from=installer --chown=nextjs:nodejs /app/apps/login/public ./apps/login/public

ENV HOSTNAME="0.0.0.0"
CMD node apps/login/server.js