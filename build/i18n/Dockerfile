#######################
## These steps set platform and env
#######################
FROM python:3 as python-base

ARG DEEPL_AUTH_KEY
ARG LANGUAGE

RUN git clone https://github.com/Saigesp/json-deepl-translate.git
WORKDIR /json-deepl-translate
RUN echo DEEPL_AUTH_KEY=${DEEPL_AUTH_KEY} > .env
RUN pip install -r requirements.txt

#######################
## console translations
#######################

RUN mkdir -p i18n
RUN mkdir -p output
COPY /console/src/assets/i18n ./i18n

RUN python main.py ./i18n/en.json --local ${LANGUAGE} --output ../output/${LANGUAGE}.json

FROM scratch as output
COPY --from=python-base /json-deepl-translate/output /console/src/assets/i18n

#######################
## notification translations
#######################

FROM python-base as notification
RUN pip install PyYAML
RUN rm -rf i18n output output-yaml
RUN mkdir -p i18n output output-yaml
COPY /internal/notification/static/i18n ./i18n
RUN python -c 'import yaml, json, sys; print(json.dumps(yaml.safe_load(open(sys.argv[1])), indent=2))' ./i18n/en.yaml > ./i18n/en.json
RUN python main.py ./i18n/en.json --local ${LANGUAGE} --output ../output/${LANGUAGE}.json
RUN python -c 'import yaml, sys; print(yaml.dump(yaml.safe_load(open(sys.argv[1])), default_flow_style=False))' ./output/${LANGUAGE}.json > ./output-yaml/${LANGUAGE}.yaml
FROM scratch as notification-output
COPY --from=notification /json-deepl-translate/output-yaml /internal/notification/static/i18n

#######################
## internal translations
#######################

FROM python-base as internal
RUN pip install PyYAML
RUN rm -rf i18n output output-yaml
RUN mkdir -p i18n output output-yaml
COPY /internal/static/i18n ./i18n
RUN python -c 'import yaml, json, sys; print(json.dumps(yaml.safe_load(open(sys.argv[1])), indent=2))' ./i18n/en.yaml > ./i18n/en.json
RUN python main.py ./i18n/en.json --local ${LANGUAGE} --output ../output/${LANGUAGE}.json
RUN python -c 'import yaml, sys; print(yaml.dump(yaml.safe_load(open(sys.argv[1])), default_flow_style=False))' ./output/${LANGUAGE}.json > ./output-yaml/${LANGUAGE}.yaml
FROM scratch as internal-output
COPY --from=internal /json-deepl-translate/output-yaml /internal/static/i18n

#######################
## login translations
#######################

FROM python-base as login
RUN pip install PyYAML
RUN rm -rf i18n output output-yaml
RUN mkdir -p i18n output output-yaml
COPY /internal/ui/login/static/i18n ./i18n
RUN python -c 'import yaml, json, sys; print(json.dumps(yaml.safe_load(open(sys.argv[1])), indent=2))' ./i18n/en.yaml > ./i18n/en.json
RUN python main.py ./i18n/en.json --local ${LANGUAGE} --output ../output/${LANGUAGE}.json
RUN python -c 'import yaml, sys; print(yaml.dump(yaml.safe_load(open(sys.argv[1])), default_flow_style=False))' ./output/${LANGUAGE}.json > ./output-yaml/${LANGUAGE}.yaml
FROM scratch as login-output
COPY --from=login /json-deepl-translate/output-yaml /internal/ui/login/static/i18n