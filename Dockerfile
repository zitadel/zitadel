FROM node:20-alpine

WORKDIR /app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

RUN mkdir /cfg
RUN ln -s .env /cfg/.env

USER nextjs

COPY --chown=nextjs:nodejs ./docker/apps/login/.next/standalone ./
COPY --chown=nextjs:nodejs ./docker/apps/login/.next/static ./apps/login/.next/static
COPY --chown=nextjs:nodejs ./docker/apps/login/public ./apps/login/public

ENV HOSTNAME="0.0.0.0"
CMD node apps/login/server.js