import {test} from "@playwright/test";
import {checkLogin, loginWithPassword} from "./login";

test("admin login", async ({page}) => {
    await loginWithPassword(page, "zitadel-admin@zitadel.localhost", "Password1.")
    await checkLogin(page, "ZITADEL Admin");
});
