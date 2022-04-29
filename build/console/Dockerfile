ARG NODE_VERSION=16

#######################
## With this step we prepare all node_modules, this helps caching the build
## Speed up this step by mounting your local node_modules directory
#######################
FROM node:${NODE_VERSION} as npm-base
WORKDIR /console

COPY console .
COPY --from=zitadel-base:local /proto /proto
COPY --from=zitadel-base:local /usr/local/bin /usr/local/bin/.
COPY build/console build/console/
RUN build/console/generate-grpc.sh

#######################
## copy for local dev
#######################
FROM scratch as npm-copy
COPY --from=npm-base /console/src/app/proto/generated /console/src/app/proto/generated

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as angular-build

COPY console/package.json console/package-lock.json ./
RUN npm ci

RUN npm run lint
RUN npm run prodbuild
RUN ls -la /console/dist/console

#######################
## Only Copy Assets
#######################
FROM scratch as angular-export
COPY --from=angular-build /console/dist/console .
