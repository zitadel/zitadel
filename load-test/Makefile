VUS ?= 20
DURATION ?= "200s"
ZITADEL_HOST ?=
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
	go install go.k6.io/xk6/cmd/xk6@latest
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	./../../xk6-modules/k6 run dist/introspection.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_session
add_session: bundle
	k6 run dist/add_session.js --vus ${VUS} --duration ${DURATION}

.PHONY: update_session
update_session: bundle
	k6 run dist/update_session.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_and_update_session
add_and_update_session: bundle
	k6 run dist/add_update_session.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_org
add_org: bundle
	k6 run dist/add_org.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_project
add_project: bundle
	k6 run dist/add_project.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_project_grant
add_project_grant: bundle
	k6 run dist/add_project_grant.js --vus ${VUS}

.PHONY: add_human
add_human: bundle
	k6 run dist/add_human.js --vus ${VUS} --duration ${DURATION}

.PHONY: add_machine
add_machine: bundle
	k6 run dist/add_machine.js --vus ${VUS} --duration ${DURATION}

.PHONY: lint
lint:
	npm i
	npm run lint:fix

.PHONY: ensure_modules
ensure_modules:
ifeq (,$(wildcard $(PWD)/../../xk6-modules))
	@echo "cloning xk6-modules"
	cd ../.. && git clone https://github.com/zitadel/xk6-modules.git
endif
	cd ../../xk6-modules && git pull

.PHONY: bundle
bundle:
	npm i
	npm run bundle