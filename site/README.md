# Getting started

CAOS Site is a github action that generates a static page out of markdown files. It uses marked.js in combination with highlight.js to compile and style markdown.
The documentation is built according to the structure of a docs `folder`[Folder](https://github.com/caos/site/tree/master/site/docs) located at root of the targeted repository.

## Running locally

You can simply run the static site by using the docker-compose command below.

```Bash
COMPOSE_DOCKER_CLI_BUILD=1 DOCKER_BUILDKIT=1 docker-compose -f site/docker-compose.yml up --build
```

## Honorable Mentions

This project was created with the help of some components from [svelte](https://github.com/sveltejs/svelte)([MIT](https://github.com/sveltejs/svelte/blob/master/LICENSE)) as well as [site-kit](https://github.com/sveltejs/site-kit)([MIT](https://github.com/sveltejs/site-kit/blob/master/LICENSE)).