# Generate the gRPC clients for Angular
FROM bufbuild/buf as generator
COPY buf.*.yaml .
COPY proto /proto
RUN buf generate

## Angular lint workspace and production build
FROM node:18 as builder
WORKDIR /docs
COPY --from=generator .artifacts/gen/ .artifacts/gen/
COPY docs/package.json docs/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY docs .
RUN yarn prodbuild

## Final image for production serving
FROM nginx:1.23.3 as final
COPY docs/nginx.conf /etc/nginx/nginx.conf 
#COPY --from=builder /docs/build /usr/share/nginx/html
COPY docs/build /usr/share/nginx/html/docs