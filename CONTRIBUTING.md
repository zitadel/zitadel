# Contributing to zitadel-typescript

## Introduction

Thank you for your interest about how to contribute! As you might know there is more than code to contribute. You can find all information needed to start contributing here.

Please give us and our community the chance to get rid of security vulnerabilities by responsibly disclose this kind of issues by contacting [security@zitadel.com](mailto:security@zitadel.com).

The strongest part of a community is the possibility to share thoughts. That's why we try to react as soon as possible to your ideas, thoughts and feedback.
We love to discuss as much as possible in an open space like in the ZITADEL [issues](https://github.com/zitadel/zitadel/issues) and [discussions](https://github.com/zitadel/zitadel/discussions) section or in our [chat](https://zitadel.com/chat), but we understand your doubts and provide further contact options [here](https://zitadel.com/contact).

If you want to give an answer or be part of discussions please be kind. Treat others like you want to be treated. Read more about our code of conduct [here](CODE_OF_CONDUCT.md).

## How to contribute

Anyone can be a contributor. Either you found a typo, or you have an awesome feature request you could implement, we encourage you to create a Pull Request.

### Pull Requests

- The latest changes are always in `main`, so please make your Pull Request against that branch.
- Pull Requests should be raised for any change
- Pull Requests need approval of an ZITADEL core engineer @zitadel/engineers before merging
- We use ESLint/Prettier for linting/formatting, so please run `pnpm lint:fix` before committing to make resolving conflicts easier (VSCode users, check out [this ESLint extension](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint) and [this Prettier extension](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode) to fix lint and formatting issues in development)
- We encourage you to test your changes, and if you have the opportunity, please make those tests part of the Pull Request
- If you add new functionality, please provide the corresponding documentation as well and make it part of the Pull Request

### Setting up local environment

A quick guide on how to setup your ZITADEL app locally to work on it and test out any changes:

1. Clone the repo:

```sh
git clone https://github.com/zitadel/zitadel-typescript.git
cd zitadel-typescript
```

3. Install packages. Developing requires Node.js v16:

```sh
pnpm install
```

4. Populate `.env.local`:

Copy `/apps/login/.env` to `/apps/login/.env.local`, and add your instance env variables for each entry.

```sh
cp apps/login/.env apps/login/.env.local
```

5. Start the developer application/server:

```sh
pnpm dev
```

The application will be available on `http://localhost:3000`

That's it! 🎉
