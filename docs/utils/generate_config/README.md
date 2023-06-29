# Generate a markdown table from yaml

This package generates a markdown table from the zitadel configuration files (steps.yaml, defaults.yaml) to automate documentation.

## Installation

Install dependencies with `npm i`

## Test

Run tests with `npm run test` or continuously with `npm run test:watch`

## Run

Execute the script `node index.js`. Then move the files with

`mv ./output/_defaults.mdx ../../docs/self-hosting/manage/configure/env` and `mv ./output/_steps.mdx ../../docs/self-hosting/manage/configure/env`