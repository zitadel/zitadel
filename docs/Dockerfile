## Angular lint workspace and production build
FROM node:18 as builder
WORKDIR /docs
COPY docs/package.json docs/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY docs .
COPY proto /proto
RUN yarn build

## Final image for serving
FROM nginx as final
COPY docs/nginx.conf /etc/nginx/nginx.conf 
COPY --from=builder /docs/build /usr/share/nginx/html
