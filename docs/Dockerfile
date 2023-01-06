FROM node:18-slim

COPY . .

RUN yarn install
RUN yarn build

EXPOSE 3000

ENTRYPOINT yarn serve --port 3000 --host 0.0.0.0