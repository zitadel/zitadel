FROM cypress/factory AS login-integration-testsuite
WORKDIR /opt/app
COPY --from=login-dev-dependencies /app/apps/login/node_modules .
RUN npx cypress install
COPY . .
CMD ["npx", "cypress", "run"]
