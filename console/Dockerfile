## Generate the gRPC clients for Angular
FROM bufbuild/buf as generator
COPY buf.*.yaml .
COPY proto /proto
RUN buf generate --include-imports

## Angular lint workspace and production build
FROM node:18 as builder
WORKDIR /console
COPY --from=generator console/src/app/proto/generated src/app/proto/generated
COPY console/package.json console/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY console .
RUN NODE_OPTIONS=--max-old-space-size=4096
RUN yarn prodbuild

## Final image for production serving
FROM node:18-alpine as final
RUN yarn add global http-server
COPY --from=builder /console/dist/console /app
EXPOSE  8080
CMD ["http-server", "--cors", "-p8080", "/app"]
