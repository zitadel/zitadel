# Load Tests

This package contains code for load testing specific endpoints of ZITADEL using [k6](https://k6.io).

## Prerequisite

* [npm](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
* [k6](https://k6.io/docs/get-started/installation/)
* [go](https://go.dev/doc/install)
* running ZITADEL

## Structure

The use cases under tests are defined in `src/use_cases`. The implementation of ZITADEL resources and calls are located under `src`.

## Execution

### Env vars

- `VUS`: Amount of parallel processes execute the test (default is 20)
- `DURATION`: Defines how long the tests are executed (default is `200s`)
- `ZITADEL_HOST`: URL of ZITADEL (default is `http://localhost:8080`)

To setup the tests we use the credentials of console and log in using an admin. The user must be able to create organizations and all resources inside organizations.

- `ADMIN_LOGIN_NAME`: `zitadel-admin@zitadel.localhost`
- `ADMIN_PASSWORD`: `Password1!`

### Test

Before you run the tests you need an initialized user. The tests don't implement the change password screen during login.

* `make human_password_login`  
  setup: creates human users  
  test: uses the previously created humans to sign in using the login ui
* `make machine_pat_login`  
  setup: creates machines and a pat for each machine  
  test: calls user info endpoint with the given pats
* `make user_info`  
  setup: creates human users and signs them in  
  test: calls user info endpoint using the given humans
* `make manipulate_user`  
  test: creates a human, updates its profile, locks the user and then deletes it 
* `make introspect`  
  setup: creates projects, one api per project, one key per api and generates the jwt from the given keys  
  test: calls introspection endpoint using the given JWTs