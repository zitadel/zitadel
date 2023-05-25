export type Item = {
  name: string;
  slug: string;
  description?: string;
};

export const demos: { name: string; items: Item[] }[] = [
  {
    name: "Login",
    items: [
      {
        name: "Username",
        slug: "username",
        description: "The entrypoint of the application",
      },
      {
        name: "Password",
        slug: "password",
        description: "The page to request a users password",
      },
      {
        name: "Accounts",
        slug: "accounts",
        description: "List active and inactive sessions",
      },
      //   {
      //     name: "Set Password",
      //     slug: "password/set",
      //     description: "The page to set a users password",
      //   },
      //   {
      //     name: "MFA",
      //     slug: "mfa",
      //     description: "The page to request a users mfa method",
      //   },
      //   {
      //     name: "MFA Set",
      //     slug: "mfa/set",
      //     description: "The page to set a users mfa method",
      //   },
      //   {
      //     name: "MFA Create",
      //     slug: "mfa/create",
      //     description: "The page to create a users mfa method",
      //   },
      //   {
      //     name: "Passwordless",
      //     slug: "passwordless",
      //     description: "The page to login a user with his passwordless device",
      //   },
      //   {
      //     name: "Passwordless Create",
      //     slug: "passwordless/create",
      //     description: "The page to add a users passwordless device",
      //   },
    ],
  },
  {
    name: "Register",
    items: [
      {
        name: "Register",
        slug: "register",
        description: "Create your ZITADEL account",
      },
      {
        name: "Verify email",
        slug: "verify",
        description: "Verify your account with an email code",
      },
    ],
  },
];
