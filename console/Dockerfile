FROM node:18 as npm-base

WORKDIR /console

COPY third_party /third_party
COPY proto /proto

COPY console/package.json console/package-lock.json ./
RUN npm ci

COPY console .

#######################
## angular lint workspace and prod build
#######################
FROM npm-base as angular-build

RUN npm run build

#######################
## Only Copy Assets
#######################
FROM scratch as angular-export
COPY --from=angular-build /console/dist/console .

##
FROM node:18-alpine as final
RUN npm install -g http-server
COPY --from=angular-build /console/dist/console /site
EXPOSE  8080
CMD ["http-server", "--cors", "-p8080", "/site"]