# Host the ZITADEL Documentation 
FROM node:18-alpine AS deploy
USER node
WORKDIR /app
COPY --from=build /src/package.json ./
COPY --from=build /src/build/ ./
EXPOSE 3000
CMD ["yarn", "serve"]