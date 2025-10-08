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
        name: "Loginname",
        slug: "loginname",
        description: "Start the loginflow with loginname",
      },
      {
        name: "Accounts",
        slug: "accounts",
        description: "List active and inactive sessions",
      },
    ],
  },
  {
    name: "Register",
    items: [
      {
        name: "Register",
        slug: "register",
        description: "Add a user with password or passkey",
      },
      {
        name: "IDP Register",
        slug: "idp",
        description: "Add a user from an external identity provider",
      },
    ],
  },
];
