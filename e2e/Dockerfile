FROM cypress/included:10.3.0

WORKDIR dependencies

COPY package.json package-lock.json ./
RUN npm ci --omit=dev

ENTRYPOINT [ "bash", "-c" ]

CMD [ "cp -R /dependencies/node_modules/ . && cypress run" ]
