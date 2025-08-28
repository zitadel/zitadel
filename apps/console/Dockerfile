
FROM nginx:1.29.0
RUN  touch /var/run/nginx.pid && \
     chown -R nginx:nginx /var/cache/nginx /var/run/nginx.pid
USER nginx
COPY --chown=nginx:nginx nginx.conf /etc/nginx/nginx.conf
# TODO Needs work to be configured for the console app
COPY --chown=nginx:nginx dist/console /usr/share/nginx/html
EXPOSE 3001
ENTRYPOINT ["nginx", "-c", "/etc/nginx/nginx.conf"]
CMD ["-g", "daemon off;"]