VUS ?= 20
DURATION ?= "200s"
ZITADEL_HOST ?=
ADMIN_LOGIN_NAME ?=
ADMIN_PASSWORD ?=
DATE := $(shell date '+%d-%H:%M:%S')

K6 := ./../../xk6-modules/k6

.PHONY: human_password_login
human_password_login: bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/human_password_login.js --vus ${VUS} --duration ${DURATION} --out csv=output/human_password_login_${DATE}.csv

.PHONY: machine_pat_login
machine_pat_login: bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/machine_pat_login.js --vus ${VUS} --duration ${DURATION} --out csv=output/machine_pat_login_${DATE}.csv

.PHONY: machine_client_credentials_login
machine_client_credentials_login: bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/machine_client_credentials_login.js --vus ${VUS} --duration ${DURATION} --out csv=output/machine_client_credentials_login_${DATE}.csv

.PHONY: user_info
user_info: bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/user_info.js --vus ${VUS} --duration ${DURATION} --out csv=output/user_info_${DATE}.csv

.PHONY: manipulate_user
manipulate_user: bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/manipulate_user.js --vus ${VUS} --duration ${DURATION} --out csv=output/manipulate_user_${DATE}.csv

.PHONY: introspect
introspect: ensure_modules bundle
	go install go.k6.io/xk6/cmd/xk6@latest
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/introspection.js --vus ${VUS} --duration ${DURATION} --out csv=output/introspect_${DATE}.csv

.PHONY: oidc_session
oidc_session: ensure_key_pair ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/oidc_session.js --vus ${VUS} --duration ${DURATION} --out csv=output/oidc_session_${DATE}.csv

.PHONY: otp_session
otp_session: ensure_key_pair ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/otp_session.js --vus ${VUS} --duration ${DURATION} --out csv=output/otp_session_${DATE}.csv

.PHONY: password_session
password_session: ensure_key_pair ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/password_session.js --vus ${VUS} --duration ${DURATION} --out csv=output/otp_session_${DATE}.csv

.PHONY: add_session
add_session: ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/add_session.js --vus ${VUS} --duration ${DURATION} --out csv=output/add_session_${DATE}.csv

.PHONY: machine_jwt_profile_grant
machine_jwt_profile_grant: ensure_modules ensure_key_pair bundle
	go install go.k6.io/xk6/cmd/xk6@latest
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/machine_jwt_profile_grant.js --vus ${VUS} --duration ${DURATION} --out csv=output/machine_jwt_profile_grant_${DATE}.csv

.PHONY: machine_jwt_profile_grant_single_user
machine_jwt_profile_grant_single_user: ensure_modules ensure_key_pair bundle
	go install go.k6.io/xk6/cmd/xk6@latest
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/machine_jwt_profile_grant_single_user.js --vus ${VUS} --duration ${DURATION} --out csv=output/machine_jwt_profile_grant_single_user_${DATE}.csv

.PHONY: verify_all_user_grants_exist
verify_all_user_grants_exist: ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/verify_all_user_grants_exist.js --vus ${VUS} --duration ${DURATION}
# 	--out csv=output/verify_all_user_grants_exist_${DATE}.csv

.PHONY: users_by_metadata_key
users_by_metadata_key: ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/users_by_metadata_key.js --vus ${VUS} --duration ${DURATION} --out csv=output/users_by_metadata_${DATE}.csv

.PHONY: users_by_metadata_value
users_by_metadata_value: ensure_modules bundle
	${K6} run --summary-trend-stats "min,avg,max,p(50),p(95),p(99)" dist/users_by_metadata_value.js --vus ${VUS} --duration ${DURATION} --out csv=output/users_by_metadata_${DATE}.csv

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
	mkdir -p output
	npm i
	npm run bundle
	go install go.k6.io/xk6/cmd/xk6@latest
	cd ../../xk6-modules && xk6 build --with xk6-zitadel=.

.PHONY: ensure_key_pair
ensure_key_pair:
ifeq (,$(wildcard $(PWD)/.keys))
	mkdir .keys
endif
ifeq (,$(wildcard $(PWD)/.keys/key.pem))
	openssl genrsa -out .keys/key.pem 2048
endif
ifeq (,$(wildcard $(PWD)/.keys/key.pem.pub))
	openssl rsa -in .keys/key.pem -outform PEM -pubout -out .keys/key.pem.pub
endif
