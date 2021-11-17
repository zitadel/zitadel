#######################
## These steps set platform / arch type specific variables
#######################
FROM python:3 as output-copy

ARG DEEPL_AUTH_KEY

ENV CONSOLE /console/src/assets/i18n

RUN git clone https://github.com/Saigesp/json-deepl-translate.git
WORKDIR /json-deepl-translate

RUN echo DEEPL_AUTH_KEY=${DEEPL_AUTH_KEY} > .env
RUN pip install -r requirements.txt

RUN mkdir -p i18n
RUN mkdir -p output
COPY /console/src/assets/i18n ./i18n

RUN python main.py ./i18n/en.json --local ES --output ../output/es.json

FROM scratch as output
COPY --from=output-copy /json-deepl-translate/output /console/src/assets/i18n