---
title: Vue
---

This is our Zitadel [Vue.js](https://vuejs.org/) template. It shows how to authenticate as a user and show user information.

![Vue Screenshot](/img/vue/app-screen.png)

## Getting Started

First, we start by creating a new Vue app with `create-vue`, the official project scaffolding tool, which sets up everything automatically for you. To create a project, run:

```bash
npm init vue@latest
```

## Install Authentication library

To keep the template as easy as possible we use [vue-oidc-client](https://github.com/soukoku/vue-oidc-client) as our main authentication library. To install, run:

```bash
npm i vue-oidc-client
```

To run the app, type:

```bash
npm run dev
```

then open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Configuration

To setup your configuration, create a file called [auth].ts in `src/auth`.

```ts
import { createOidcAuth, LogLevel, SignInType } from 'vue-oidc-client/vue3';

const ZITADEL_ISSUER = "https:/[your-domain]-[random-string].zitadel.cloud";
const ZITADEL_CLIENT_ID = "YOUR-CLIENT-ID";
const appUrl = "http://localhost:3000/";

const mainOidc = createOidcAuth(
  "main",
  SignInType.Window,
  appUrl,
  {
    authority: ZITADEL_ISSUER,
    client_id: ZITADEL_CLIENT_ID,
    response_type: "code",
    scope: "openid profile email",
  },
  console,
  LogLevel.Debug
);

... events

export default mainOidc;

```

We recommend using the Authentication Code flow secured by PKCE for the Authentication flow.
To be able to connect to ZITADEL, navigate to your Instances Console, create or select an existing project and add your app selecting WEB, then PKCE, and then add `http://localhost:3000/auth/signinwin/main` as redirect url to your app.

Hit Create, then in the detail view of your application make sure to enable dev mode. Dev mode ensures that you can start an auth flow from a **non** https endpoint for testing.

> Note that we get a clientId but no clientSecret because it is not needed for our authentication flow.

Now go to Token settings and check the checkbox for **User Info inside ID Token** to get your users detail directly on authentication.

## Environment

Set your environment variables in the auth file.
You can find your Issuer Url on the application detail page in console.

```
ZITADEL_ISSUER=[yourIssuerUrl]
ZITADEL_CLIENT_ID=[yourClientId]
```

## Protected Routes

We use the /about page to display user information. To ensure the /about page is protected add a `meta: {authName: auth.authName }` to the routes array.

```ts
import auth from "../auth/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/about",
      name: "about",
      meta: {
        authName: auth.authName,
      },
      component: () => import("@/views/AboutView.vue"),
    },
  ],
});

auth.useRouter(router);

export default router;
```

## Modify main.ts

To make the example app function properly, you have to modify the main.ts and call `startup()` of the previously created auth file.
Note that to access auth information in your components, you have to add the line `app.config.globalProperties.$oidc = mainOidc;`.
You can then later access the information from the view `AboutView.vue`

### main.ts

```ts
import "./assets/main.css";

import { createPinia } from "pinia";
import { createApp } from "vue";

import App from "./App.vue";
import mainOidc from "./auth/auth";
import router from "./router";

mainOidc.startup().then((ok) => {
  if (ok) {
    const app = createApp(App);

    app.use(createPinia());
    app.use(router);

    app.config.globalProperties.$oidc = mainOidc;

    app.mount("#app");
  } else {
    console.log("Startup was not ok");
  }
});
```

### AboutView.vue

```vue
<template>
  <div class="about">
    <h1>This is an authenticated user page</h1>
    <div class="about" v-if="$oidc.isAuthenticated">
      <p class="username">
        <strong>{{ user.name }}</strong>
      </p>

      <button v-on:click="$oidc.signOut">Signout</button>

      <ul class="claims">
        <li v-for="c in claims" :key="c.key">
          <strong>{{ c.key }}</strong
          >: {{ c.value }}
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts">
export default {
  computed: {
    user(): any {
      return { ...this.$oidc.userProfile, accessToken: this.$oidc.accessToken };
    },
    claims() {
      if (this.user) {
        return Object.keys(this.user).map((key) => ({
          key,
          value: this.user[key],
        }));
      }
      return [];
    },
  },
};
</script>
```

You can now check wheter the user is authenticated with `$oidc.isAuthenticated` and call the `$oidc.signOut` function to log the user out.

## Completion

You have successfully integrated your Vue application with ZITADEL!

If you get stuck, consider checking out our [example](https://github.com/zitadel/zitadel-examples/tree/main/vue) application. It includes all the mentioned functionality of this quickstart. You can simply start by cloning the repository and replacing the `ZITADEL_ISSUER` and `CLIENT_ID` in the auth.ts by your own configuration. If you run into issues, contact us or raise an issue on [GitHub](https://github.com/zitadel/zitadel).
