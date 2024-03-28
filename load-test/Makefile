VUS ?= 20
DURATION ?= "200s"
ZITADEL_HOST ?=
ZITADEL_ORG_ID ?=
ADMIN_LOGIN_NAME ?=
ADMIN_PASSWORD ?=

.PHONY: human_password_login
human_password_login:
	k6 run tests/use_cases/password_login.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: machine_pat_login
machine_pat_login:
	k6 run tests/use_cases/machine_pat_login.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: user_info
user_info:
	k6 run tests/use_cases/user_info.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: transpile
transpile:
	npm i
	npm run bundle
