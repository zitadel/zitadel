export type Item = {
  name: string;
  slug: string;
  description?: string;
};

export enum ProviderSlug {
  GOOGLE = "google",
  GITHUB = "github",
}

export const demos: { name: string; items: Item[] }[] = [
  {
    name: "Login",
    items: [
      {
        name: "Loginname",
        slug: "loginname",
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
      {
        name: "Passkey Registration",
        slug: "passkey/add",
        description: "The page to add a users passkey device",
      },
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
        name: "IDP Register",
        slug: "register/idp",
        description: "Register with an Identity Provider",
      },
      {
        name: "Verify email",
        slug: "verify",
        description: "Verify your account with an email code",
      },
    ],
  },
];
