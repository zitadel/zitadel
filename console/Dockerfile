## Generate the gRPC clients for Angular
FROM bufbuild/buf as generator
COPY buf.*.yaml .
COPY proto /proto
RUN buf generate --include-imports

## Cache step for dependencies
FROM node:18 as npm-base
WORKDIR /console
COPY console/package.json console/yarn.lock ./
RUN yarn install --frozen-lockfile

## Angular lint workspace and production build
FROM npm-base as angular-build
COPY --from=generator console/src/app/proto/generated src/app/proto/generated
COPY console .
RUN yarn prodbuild

## Extra step which allows to output the build artifacts during CI
FROM scratch as angular-export
COPY --from=angular-build /console/dist/console .

## Final image for production serving
FROM node:18-alpine as final
RUN yarn add global http-server
COPY --from=angular-build /console/dist/console /app
EXPOSE  8080
CMD ["http-server", "--cors", "-p8080", "/app"]
