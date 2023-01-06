FROM node:18-alpine

RUN apk update && apk --no-cache add git

COPY . .

RUN yarn install
RUN yarn build

EXPOSE 3000

ENTRYPOINT yarn serve --port 3000 --host 0.0.0.0