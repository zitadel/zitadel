VUS ?= 20
DURATION ?= "200s"
ZITADEL_HOST ?=
ZITADEL_ORG_ID ?=
ADMIN_LOGIN_NAME ?=
ADMIN_PASSWORD ?=

.PHONY: human_password_login
human_password_login: bundle
	k6 run dist/human_password_login.js --vus ${VUS} --duration ${DURATION}

.PHONY: machine_pat_login
machine_pat_login: bundle
	k6 run dist/machine_pat_login.js --vus ${VUS} --duration ${DURATION}

.PHONY: user_info
user_info: bundle
	k6 run dist/user_info.js --vus ${VUS} --duration ${DURATION}

.PHONY: manipulate_user
manipulate_user: bundle
	k6 run dist/manipulate_user.js --vus ${VUS} --duration ${DURATION}

.PHONY: introspect
introspect: ensure_modules bundle
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	./../../xk6-modules/k6 run dist/introspection.js --vus ${VUS} --duration ${DURATION}

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

.PHONY: bundle
bundle:
	npm i
	npm run bundle