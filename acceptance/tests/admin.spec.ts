import {test} from "@playwright/test";
import {loginWithPassword} from "./login";

test("admin login", async ({page}) => {
    await loginWithPassword(page, "zitadel-admin@zitadel.localhost", "Password1.")
    await page.getByRole("heading", {name: "Welcome ZITADEL Admin!"}).click();
});
