import { test as base } from "@playwright/test";
import { loginScreenExpect, loginWithPassword } from "./login";
import { Config, ConfigReader } from "./config";

const test = base.extend<{ cfg: Config }>({
  cfg: async ({}, use) => {
    await use(new ConfigReader().config);
  },
});

test("admin login", async ({ page, cfg }) => {
  await loginWithPassword(page, cfg.zitadelAdminUser , "Password1!");
  await loginScreenExpect(page, "ZITADEL Admin");
});
