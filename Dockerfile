FROM alpine:latest

RUN ls -la /github/workspace/
COPY /github/workspace/angular /app/console
COPY /github/workspace/go /app