ARG NODE_VERSION=14

#######################
## With this step we prepare all node_modules, this helps caching the build
## Speed up this step by mounting your local node_modules directory
#######################
FROM node:${NODE_VERSION} as npm-base
WORKDIR /console

COPY console/package.json console/package-lock.json ./
RUN npm ci

COPY console .
COPY --from=base /proto /proto
COPY --from=base /usr/local/bin /usr/local/bin/.
COPY build/console build/console/
RUN build/console/generate-grpc.sh

#######################
## copy for local dev
#######################
FROM scratch as npm-copy
COPY --from=npm-base /console/src/app/proto/generated ./console/src/app/proto/generated

#######################
## angular dev build
#######################
FROM npm-base as dev-angular-build
RUN npm install -g @angular/cli

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as prod-angular-build
RUN npm run lint
RUN npm run prodbuild

#######################
## Only Copy Assets
#######################
FROM scratch as prod-angular-export
COPY --from=prod-angular-build /console/dist/console .