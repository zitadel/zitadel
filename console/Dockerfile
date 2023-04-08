ARG NODE_VERSION=18
FROM node:${NODE_VERSION} as npm-base
WORKDIR /console
COPY console/package.json console/package-lock.json ./
RUN npm ci
RUN npm install -g @angular/cli
COPY proto/ /proto/
COPY console/buf.gen.yaml buf.gen.yaml
RUN npm run generate
COPY console .
RUN npm run lint && npm run build