FROM nginx:alpine

COPY build/nginx.conf /etc/nginx/nginx.conf

RUN rm -rf /usr/share/nginx/html/*

COPY dist/console /usr/share/nginx/html/

EXPOSE 8081

ENTRYPOINT ["nginx", "-g", "daemon off;"]