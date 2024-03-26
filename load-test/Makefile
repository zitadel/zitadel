VUS ?= 20
DURATION ?= "200s"

.PHONY: password_login
password_login:
	k6 run tests/use_cases/password_login.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}

.PHONY: user_info
user_info:
	k6 run tests/use_cases/user_info.js --vus ${VUS} --duration ${DURATION} --env MAX_VUS=${VUS}
