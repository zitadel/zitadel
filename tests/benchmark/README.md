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

* `VUS`: Amount of parallel processes execute the test (default is 20)
* `DURATION`: Defines how long the tests are executed (default is `200s`)
* `ZITADEL_HOST`: URL of ZITADEL (default is `http://localhost:8080`)
* `ADMIN_LOGIN_NAME`: Loginanme of a human user with `IAM_OWNER`-role
* `ADMIN_PASSWORD`: password of the human user

To setup the tests we use the credentials of console and log in using an admin. The user must be able to create organizations and all resources inside organizations.

* `ADMIN_LOGIN_NAME`: `zitadel-admin@zitadel.localhost`
* `ADMIN_PASSWORD`: `Password1!`

### Test

Before you run the tests you need an initialized user. The tests don't implement the change password screen during login.

* `make human_password_login`  
  setup: creates human users  
  test: uses the previously created humans to sign in using the login ui
* `make machine_pat_login`  
  setup: creates machines and a pat for each machine  
  test: calls user info endpoint with the given pats
* `make machine_client_credentials_login`  
  setup: creates machines and a client credential secret for each machine  
  test: calls token endpoint with the `client_credentials` grant type.
* `make user_info`  
  setup: creates human users and signs them in  
  test: calls user info endpoint using the given humans
* `make manipulate_user`  
  test: creates a human, updates its profile, locks the user and then deletes it 
* `make introspect`  
  setup: creates projects, one api per project, one key per api and generates the jwt from the given keys  
  test: calls introspection endpoint using the given JWTs
* `make add_session`  
  setup: creates human users  
  test: creates new sessions with user id check
* `make oidc_session`  
  setup: creates a machine user to create the auth request and session.  
  test: creates an auth request, a session and links the session to the auth request. Implementation of [this flow](https://zitadel.com/docs/guides/integrate/login-ui/oidc-standard).
* `make otp_session`  
  setup: creates 1 human user for each VU and adds email OTP to it  
  test: creates a session based on the login name of the user, sets the email OTP challenge to the session and afterwards checks the OTP code
* `make password_session`  
  setup: creates 1 human user for each VU and adds email OTP to it  
  test: creates a session based on the login name of the user and checks for the password on a second step
* `make machine_jwt_profile_grant`  
  setup: generates private/public key, creates machine users, adds a key  
  test: creates a token and calls user info 
* `make machine_jwt_profile_grant_single_user`  
  setup: generates private/public key, creates machine user, adds a key  
  test: creates a token and calls user info in parallel for the same user
* `make users_by_metadata_key`  
  setup: creates for half of the VUS a human user and a machine for the other half, adds 3 metadata to each user
  test: calls the list users endpoint and filters by a metadata key
* `make users_by_metadata_value`  
  setup: creates for half of the VUS a human user and a machine for the other half, adds 3 metadata to each user
  test: calls the list users endpoint and filters by a metadata value
* `make verify_all_user_grants_exists`  
  setup: creates 50 projects, 1 machine per VU
  test: creates a machine and grants all projects to the machine
  teardown: the organization is not removed to verify the data of the projections are correct. You can find additional information [at the bottom of this file](./src/use_cases/verify_all_user_grants_exist.ts)
