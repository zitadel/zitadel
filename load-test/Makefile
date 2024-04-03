VUS ?= 20
DURATION ?= "200s"
ZITADEL_HOST ?=
ZITADEL_ORG_ID ?=
ADMIN_LOGIN_NAME ?=
ADMIN_PASSWORD ?=

.PHONY: human_password_login
human_password_login:
	k6 run tests/use_cases/human_password_login.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: machine_pat_login
machine_pat_login:
	k6 run tests/use_cases/machine_pat_login.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: user_info
user_info:
	k6 run tests/use_cases/user_info.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: introspect
introspect: ensure_modules
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	./../../xk6-modules/k6 run tests/use_cases/introspection.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: ensure_modules
ensure_modules:
ifeq (,$(wildcard $(PWD)/../../xk6-modules))
	@echo "cloning xk6-modules"
	cd ../.. && git clone git@github.com:zitadel/xk6-modules.git
endif
	cd ../../xk6-modules && git pull

.PHONY: transpile
transpile:
	npm i
	npm run bundle
