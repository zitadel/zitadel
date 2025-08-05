# syntax=docker.io/docker/dockerfile:1

FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --chown=nextjs:nodejs public ./public
COPY --chown=nextjs:nodejs .next/standalone/ ./
COPY --chown=nextjs:nodejs .next/static ./.next/static 

USER nextjs

EXPOSE 3002

ENV PORT=3002

ENV HOSTNAME="0.0.0.0"
CMD ["node", "./apps/login/server.js"]