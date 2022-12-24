# Build the ZITADEL Documentation
FROM node:18 AS build
WORKDIR /src
COPY package.json ./
COPY yarn.lock ./
RUN yarn install --frozen-lockfile 
COPY . .
RUN yarn build

# Host the ZITADEL Documentation 
FROM nginx:alpine
COPY --from=build /src/build/ /usr/share/nginx/html/docs