FROM node:18 as npm-base

COPY console/package.json console/package-lock.json ./
RUN npm ci

COPY console .

RUN npm run generate

#######################
## copy for local dev
#######################
FROM scratch as npm-copy
COPY --from=npm-base /console/src/app/proto/generated /console/src/app/proto/generated

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as angular-build

#RUN npm run lint
RUN npm run prodbuild

#######################
## Only Copy Assets
#######################
FROM scratch as angular-export
COPY --from=angular-build /console/dist/console .
